package writer

import (
	"bytes"
	"io"
	"strconv"

	"github.com/sogwms/mGoLibrary/converter/def"

	"github.com/xuri/excelize/v2"
)

type xlsxWriter struct {
	f  *excelize.File
	to io.Writer
}

func NewXlsxWriter(w io.Writer) (def.Writer, error) {
	f := excelize.NewFile()

	ret := &xlsxWriter{
		f:  f,
		to: w,
	}

	return ret, nil
}

func (x *xlsxWriter) Write(data *def.FileData) error {
	rows := data.Rows

	activeSheet := x.f.GetSheetName(x.f.GetActiveSheetIndex())

	for i, row := range rows {
		for j, col := range row {
			axis := getAxis(j) + strconv.Itoa(i+1)
			err := x.f.SetCellValue(activeSheet, axis, col)
			if err != nil {
				return err
			}
		}
	}

	_, err := x.f.WriteTo(x.to)
	if err != nil {
		return err
	}

	return nil
}

// index: 0..
func getAxis(index int) string {
	s := strconv.FormatInt(int64(index), 26)
	baseChar := 'A'
	buf := bytes.NewBuffer(nil)
	for _, v := range s {
		if v >= 'a' {
			buf.WriteRune(v - 'a' + 10 + baseChar)
		} else {
			buf.WriteRune(v - '0' + baseChar)
		}
	}

	return buf.String()
}
