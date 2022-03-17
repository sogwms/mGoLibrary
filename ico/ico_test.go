package ico

import (
	"testing"
	// "github.com/sogwms/mGoLibrary/ico"
)

var websiteCases = [...]string{
	"https://www.iodraw.com/",
	"https://www.jenkins.io/zh/",
	"https://wordpress.org/",
	"https://riscv.org/",
	"https://khalilstemmler.com/articles/software-design-architecture/full-stack-software-design/",
	"https://www.code-nav.cn/",
}

func TestIco(t *testing.T) {

	t.Log("Done")

	for _, v := range websiteCases {
		GetWebsiteIcoInBase64(v)
		if GetError() != nil {
			t.Fatal(GetError())
			ClearError()
		}
	}

}
