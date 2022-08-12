package converter

import (
	"github.com/sogwms/mGoLibrary/converter/def"
)

func Convert(f def.Reader, t def.Writer) error {
	data, err := f.GetData()
	if err != nil {
		return err
	}

	return t.Write(data)
}
