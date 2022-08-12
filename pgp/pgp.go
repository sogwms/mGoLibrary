package pgp

import (
	"bytes"
	"crypto"
	"io"
	"os"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

type EntityList = openpgp.EntityList
type Entity = openpgp.Entity
type KeyRing = openpgp.KeyRing
type MessageDetails = openpgp.MessageDetails

type Encrypter struct {
	buffer          *bytes.Buffer
	armorWriter     io.WriteCloser
	plaintextWriter io.WriteCloser
	to              EntityList
	signer          *Entity
}

// ReadArmoredKeyRing generates key entity from input
func ReadArmoredKeyRing(r io.Reader) (*Entity, error) {
	entityList, err := openpgp.ReadArmoredKeyRing(r)
	if err != nil {
		return nil, err
	}

	return entityList[0], nil
}

// ref @ReadArmoredKeyRing
func ReadArmoredKeyRingFromFile(filename string) (*Entity, error) {
	keyFile, _ := os.Open(filename)
	defer keyFile.Close()
	return ReadArmoredKeyRing(keyFile)
}

// ref @ReadArmoredKeyRing
func ReadArmoredKeyRingFromString(key string) (*Entity, error) {
	buffer := bytes.NewBufferString(key)
	return ReadArmoredKeyRing(buffer)
}

// decrypt armored message by keys
func DecryptArmoredMessage(r io.Reader, keys EntityList) (*MessageDetails, error) {
	armorDecoder, err := armor.Decode(r)
	if err != nil {
		return nil, err
	}
	messageDetail, err := openpgp.ReadMessage(armorDecoder.Body, keys, nil, nil)
	if err != nil {
		return nil, err
	}

	return messageDetail, nil
}

// ref @DecryptArmoredMessage
func DecryptArmoredMessageFromFile(filename string, keys EntityList) (*MessageDetails, error) {
	messageFile, _ := os.Open(filename)
	defer messageFile.Close()
	return DecryptArmoredMessage(messageFile, keys)
}

// ref @DecryptArmoredMessage
func DecryptArmoredMessageFromString(message string, keys EntityList) (*MessageDetails, error) {
	buffer := bytes.NewBufferString(message)
	return DecryptArmoredMessage(buffer, keys)
}

// NewEncrypter generates an encrypter which can be used to encrypt plaintext
// to: indicator of the receiver
// signer: indicator of the signer
func NewEncrypter(to EntityList, signer *Entity) (*Encrypter, error) {
	encryptor := new(Encrypter)
	encryptor.to = to
	encryptor.signer = signer
	if err := encryptor.Reset(); err != nil {
		return nil, err
	}

	return encryptor, nil
}

// DecodePrivateKey trys to decode the private key, if the key has been decoded, nothing will happen
func DecodePrivateKey(keyEntity *Entity, passphrase []byte) error {
	if keyEntity.PrivateKey != nil {
		if keyEntity.PrivateKey.Encrypted {
			passPhraseByte := []byte(passphrase)
			if err := keyEntity.PrivateKey.Decrypt(passPhraseByte); err != nil {
				return err
			}
			for _, subkey := range keyEntity.Subkeys {
				subkey.PrivateKey.Decrypt(passPhraseByte)
			}
		}
	}

	return nil
}

// Write will encrypt the input data
func (e *Encrypter) Write(p []byte) (n int, err error) {
	return e.plaintextWriter.Write(p)
}

// Buffer returns the buffer used to store the encrypted data. Used to fetch the output
func (e *Encrypter) Buffer() *bytes.Buffer {
	return e.buffer
}

// Close must be invoked after the contents have been written.
func (e *Encrypter) Close() {
	e.plaintextWriter.Close()
	e.armorWriter.Close()
}

func (e *Encrypter) Reset() error {
	to := e.to
	signer := e.signer

	packetConfig := &packet.Config{
		DefaultHash:            crypto.SHA256,
		DefaultCipher:          packet.CipherAES256,
		DefaultCompressionAlgo: packet.CompressionZIP,
		CompressionConfig: &packet.CompressionConfig{
			Level: 9,
		},
		RSABits: 4096,
	}
	buffer := bytes.NewBuffer(nil)
	armorWriter, err := armor.Encode(buffer, "PGP MESSAGE", nil)
	if err != nil {
		return err
	}

	receiver := to
	plaintextWriter, err := openpgp.Encrypt(armorWriter, receiver, signer, nil, packetConfig)
	if err != nil {
		return err
	}

	e.buffer = buffer
	e.armorWriter = armorWriter
	e.plaintextWriter = plaintextWriter

	return nil
}
