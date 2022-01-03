package main

import (
	"fmt"

	"github.com/sogwms/mGoLibrary/ico"
)

func main() {

	fmt.Println(ico.GetWebsiteIcoInBase64("https://segmentfault.com/"))
	fmt.Println(ico.GetError())

	// url := "https://www.dingmos.com/usr/themes/Akina/images/d_logo.png"
	// fmt.Println(ico.GetIcoInBase64(url))
}
