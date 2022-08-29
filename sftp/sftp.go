package sftp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Sftp struct {
	ssh  *ssh.Client  //ssh client
	sftp *sftp.Client //sftp client
}

func New(host, username string, key []byte, keyPassphrase string) (*Sftp, error) {
	sshConn, err := newSshConnection(host, username, key, keyPassphrase)
	if err != nil {
		return nil, err
	}
	sftpConn, err := newSftpConnection(sshConn)
	if err != nil {
		return nil, err
	}

	ret := new(Sftp)
	ret.ssh = sshConn
	ret.sftp = sftpConn

	return ret, nil
}

// param: filename or key string data
func ReadKey(param string) ([]byte, error) {
	// assume param is a filename
	info, err := os.Stat(param)
	if err == nil && !info.IsDir() {
		key, err := ioutil.ReadFile(param)
		if err != nil {
			return nil, err
		}
		return key, nil
	}

	// otherwise, key string data
	return []byte(param), nil
}

func ReadKeyFromFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// @ref ReadKey
func ReadKeyWithoutError(param string) []byte {
	ret, _ := ReadKey(param)
	return ret
}

// @ref ReadKeyFromFile
func ReadKeyFromFileWithoutError(filename string) []byte {
	ret, _ := ReadKeyFromFile(filename)
	return ret
}

func newSshConnection(host string, username string, key []byte, keyPassphrase string) (*ssh.Client, error) {

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if keyPassphrase != "" {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(keyPassphrase))
	}
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}
	config.KeyExchanges = append(config.KeyExchanges, "diffie-hellman-group-exchange-sha256")

	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func newSftpConnection(sshConn *ssh.Client) (*sftp.Client, error) {
	// open an SFTP session over an existing ssh connection.
	client, err := sftp.NewClient(sshConn)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (s *Sftp) ExecShellWithNewSession(cmd string) (string, error) {
	session, err := s.ssh.NewSession()
	if err != nil {
		return "", err
	}
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// src:local, dest:remote
func (s *Sftp) Upload(dest string, src io.Reader) (int64, error) {
	//  try to create the specific file on remote machine
	destFile, err := s.sftp.Create(dest)
	if err != nil {
		return 0, err
	}
	defer destFile.Close()

	// fianlly, write the data (streamly)
	bc, err := io.Copy(destFile, src)
	if err != nil {
		return 0, err
	}

	return bc, nil
}

// src:local, dest:remote
func (s *Sftp) UploadFromFile(dest, src string) error {
	// try to open the local file first
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	size, err := s.Upload(dest, srcFile)
	if err != nil {
		return err
	}

	// verify transfer size
	if localFileInfo, err := srcFile.Stat(); err == nil {
		if size != localFileInfo.Size() {
			return errors.New("unequal file size")
		}
	}

	return nil
}

func (s *Sftp) UploadFromString(dest string, data string) error {
	buffer := bytes.NewBufferString(data)

	_, err := s.Upload(dest, buffer)
	if err != nil {
		return err
	}

	return nil
}

// filename: remote filename to download
func (s *Sftp) Download(filename string) (io.ReadCloser, error) {
	srcFile, err := s.sftp.OpenFile(filename, os.O_RDONLY)
	if err != nil {
		return nil, err
	}

	return srcFile, nil
}

func (s *Sftp) DownloadByReader(filename string) (io.Reader, error) {
	rc, err := s.Download(filename)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(data)

	return buffer, nil
}

// dest: local filename, src: remote filename
func (s *Sftp) DownloadToLocal(dest string, src string) error {
	rc, err := s.Download(src)
	if err != nil {
		return err
	}
	defer rc.Close()

	descFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer descFile.Close()

	if _, err := descFile.ReadFrom(rc); err != nil {
		return err
	}

	return nil
}

func (s *Sftp) WalkDir(dir string) {
	w := s.sftp.Walk(dir)
	for w.Step() {
		if w.Err() != nil {
			continue
		}
		log.Println(w.Path())
	}
}

func (s *Sftp) Ls(dir string) {
	fileinfo, _ := s.sftp.ReadDir(dir)
	fmt.Println("name", "isDir", "Mode", "Size", "ModTime")
	for _, f := range fileinfo {
		fmt.Println(f.Name(), f.IsDir(), f.Mode(), f.Size(), f.ModTime())
	}
}

func (s *Sftp) Lstat(p string) (fs.FileInfo, error) {
	return s.sftp.Lstat(p)
}

func contains(str string, substrs []string) bool {
	ret := true
	for _, v := range substrs {
		ret = ret && strings.Contains(str, v)
		if !ret {
			break
		}
	}

	return ret
}

func (s *Sftp) CheckFileExistenceOnDir(dir string, substrs []string) string {
	fileinfo, _ := s.sftp.ReadDir(dir)
	for _, f := range fileinfo {
		filename := f.Name()
		if contains(filename, substrs) {
			return filename
		}
	}

	return ""
}

func (s *Sftp) GetFirstFilenameOnDir(dir string) string {
	fileinfo, _ := s.sftp.ReadDir(dir)
	if len(fileinfo) > 0 {
		return fileinfo[0].Name()
	}
	return ""
}

func (s *Sftp) GetAllFilesOnDir(dir string) ([]fs.FileInfo, error) {
	return s.sftp.ReadDir(dir)
}
