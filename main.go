package main

import (
	"fmt"
	"m/ico"
)

func main() {
	fmt.Println(ico.GetIcoBase64ByHttp("https://tiddlywiki.com"))
}
