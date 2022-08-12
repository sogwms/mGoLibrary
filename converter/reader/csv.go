package reader

import (
	"encoding/csv"
	"io"

	"github.com/sogwms/mGoLibrary/converter/def"
)

type csvReader struct {
	r *csv.Reader
}

func NewCsvReader(reader io.Reader) (def.Reader, error) {

	r := csv.NewReader(reader)

	ret := &csvReader{
		r: r,
	}

	return ret, nil
}

func (r *csvReader) GetData() (fd *def.FileData, err error) {
	rows, err := r.r.ReadAll()
	if err != nil {
		return nil, err
	}

	fd = new(def.FileData)
	fd.Rows = rows
	fd.Metadata = nil

	return
}
