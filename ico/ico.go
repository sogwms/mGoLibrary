package ico

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	_ "github.com/biessek/golang-ico"
)

var lastError error

// 伪装浏览器访问
func httpRequest(url string) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}

	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	/*避免不必要的解压缩*/
	// request.Header.Add("Accept-encoding", "gzip, deflate, br")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36 Edg/96.0.1054.62")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return []byte{}, err
	}
	if response.StatusCode != 200 {
		return []byte{}, errors.New(response.Status)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func getIcoByHttp(url string) ([]byte, error) {
	return httpRequest(url)
}

func convertIcoToPng(ico []byte, e error) ([]byte, error) {
	if e != nil {
		return nil, e
	}

	tmp := new(bytes.Buffer)
	tmp.Write(ico)
	img, _, err := image.Decode(tmp)
	if err != nil {
		return nil, err
	}
	buffer := new(bytes.Buffer)

	png.Encode(buffer, img)
	return buffer.Bytes(), nil
}

func encodeImageByBase64(ico []byte, e error) (string, error) {
	if e != nil {
		return "", e
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(ico), nil
}

func GetIcoInBase64(url string) string {
	img, err := encodeImageByBase64(convertIcoToPng(getIcoByHttp(url)))

	if err != nil {
		fmt.Println(err)
		lastError = err
	}

	return img
}

func getWebsiteIconUrl(host string) string {
	resp, err := http.Get(host)
	if err != nil {
		return ""
	}
	data, _ := io.ReadAll(resp.Body)

	re := regexp.MustCompile(`<link\s*rel\s*=\s*"\s*(?:shortcut)?\s*icon"\s*href\s*=\s*"(\S*)"\s*/?>`)

	matched := re.FindSubmatch(data)
	if matched != nil {
		addr := string(matched[1])
		if addr[0] != 'h' {
			return host + addr
		} else {
			return addr
		}
	}

	return host + "/favicon.ico"
}

func GetWebsiteIcoInBase64(host string) string {

	hostAddr := ""
	if host[0] >= '0' && host[0] <= '9' {
		hostAddr = "http://" + host
	} else {
		parts := strings.Split(host, "/")
		hostAddr = parts[0] + "//" + parts[2]
	}

	fmt.Println("host: ", hostAddr)
	url := getWebsiteIconUrl(hostAddr)
	fmt.Println("icon-url:", url)

	img, err := encodeImageByBase64(convertIcoToPng(getIcoByHttp(url)))

	if err != nil {
		fmt.Println(err)
		lastError = err
	}

	return img
}

func GetError() error {
	return lastError
}
