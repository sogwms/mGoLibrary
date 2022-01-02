package main

import (
	"fmt"

	"github.com/sogwms/mGoLibrary/ico"
)

func main() {

	fmt.Println(ico.GetWebsiteIcoInBase64("https://shoat"))
	fmt.Println(ico.GetError())
}
