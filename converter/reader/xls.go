package reader

import (
	"io"

	"github.com/sogwms/mGoLibrary/converter/def"

	"github.com/extrame/xls"
)

type xlsReader struct {
	wb       *xls.WorkBook
	filename string
}

func NewXlsReader(reader io.ReadSeeker, charset string) (def.Reader, error) {
	wb, err := xls.OpenReader(reader, charset)
	if err != nil {
		return nil, err
	}

	ret := &xlsReader{
		wb: wb,
	}

	return ret, nil
}

func (r *xlsReader) GetRows() (res [][]string, err error) {
	// TODO: check active sheet
	sheet := r.wb.GetSheet(0)
	for i := 0; i < int(sheet.MaxRow)+1; i++ {
		rawRow := sheet.Row(i)
		row := make([]string, 0)
		for j := 0; j < rawRow.LastCol(); j++ {
			col := rawRow.Col(j)
			row = append(row, col)
		}
		res = append(res, row)
	}

	return
}

func (r *xlsReader) GetData() (fd *def.FileData, err error) {
	rows, err := r.GetRows()
	if err != nil {
		return
	}

	fd = new(def.FileData)
	fd.Rows = rows
	fd.Metadata = nil

	return
}
