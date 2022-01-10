package main

import (
	"fmt"

	"github.com/sogwms/mGoLibrary/ico"
)

func main() {

	fmt.Println(ico.GetWebsiteIcoInBase64("https://wordpress.org/"))
	// fmt.Println(ico.GetWebsiteIcoInBase64("https://riscv.org/"))
	// fmt.Println(ico.GetWebsiteIcoInBase64("https://khalilstemmler.com/articles/software-design-architecture/full-stack-software-design/"))
	// fmt.Println(ico.GetWebsiteIcoInBase64("https://www.code-nav.cn/"))
	fmt.Println(ico.GetError())

	// url := "https://www.dingmos.com/usr/themes/Akina/images/d_logo.png"
	// fmt.Println(ico.GetIcoInBase64(url))

	// resp, err := soup.Get("https://www.code-nav.cn/")
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(-1)
	// }

	// doc := soup.HTMLParse(resp)
	// links := doc.Find("head").FindAll("link")
	// for _, link := range links {
	// 	if strings.Contains(link.Attrs()["rel"], "icon") {
	// 		fmt.Println(link.Attrs()["href"])
	// 	}
	// }
}
