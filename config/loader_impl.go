package config

import (
	"encoding/csv"
	"fmt"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"os"
)

type csvLoader struct {
}

func (cl *csvLoader) load(file string, processor func(row int, record []string) (interface{}, error), params []interface{}) error {
	cf, err := os.Open(file)
	if err != nil {
		return err
	}
	defer cf.Close()
	reader := csv.NewReader(transform.NewReader(cf, unicode.UTF8BOM.NewDecoder()))
	//reader := csv.NewReader(cf)
	reader.Comment = '#' // 可以设置读入文件中的注释符
	reader.Comma = ','   // 默认是逗号，也可以自己设置
	row := 0
	for {
		r, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if _, err = processor(row, r); err != nil {
			return err
		} else {
			row = row + 1
		}
	}
	return nil
}

type xlsxLoader struct {
}

func (xl *xlsxLoader) load(file string, processor func(row int, record []string) (interface{}, error), params []interface{}) error {
	if len(params) == 0 {
		return fmt.Errorf("required parameter sheet number: params[0]")
	}
	xf, err := excelize.OpenFile(file)
	if err != nil {
		return err
	}
	defer xf.Close()
	sheet := params[0].(int)
	sn := xf.GetSheetName(params[0].(int))
	if len(sn) == 0 {
		return fmt.Errorf("sheet[%d] not found in xlsx file %s", sheet, file)
	}
	rows, err := xf.GetRows(sn)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("read data from sheet[%d] of xlsx[%s] failed", sheet, file))
	}
	for rowIndex, row := range rows {
		if _, err := processor(rowIndex, row); err != nil {
			return err
		}
	}
	return nil
}
