package ico

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"

	_ "github.com/biessek/golang-ico"
)

// 伪装浏览器访问
func httpRequest(url string) []byte {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}
	}

	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	request.Header.Add("Accept-encoding", "gzip, deflate, br")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36 Edg/96.0.1054.62")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return []byte{}
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte{}
	}

	return body
}

func getIcoByHttp(host string) []byte {
	return httpRequest(host + "/favicon.ico")
}

func convertIcoToPng(ico []byte) []byte {
	tmp := new(bytes.Buffer)
	tmp.Write(ico)
	img, _, err := image.Decode(tmp)
	if err != nil {
		fmt.Println(err)
	}
	buffer := new(bytes.Buffer)

	png.Encode(buffer, img)
	return buffer.Bytes()
}

func encodeImageByBase64(ico []byte) string {
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(ico)
}

func GetIcoBase64ByHttp(host string) string {
	return encodeImageByBase64(convertIcoToPng(getIcoByHttp(host)))
}
