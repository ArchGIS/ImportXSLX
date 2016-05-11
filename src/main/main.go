package main

import (
	"fmt"

	"github.com/tealeg/xlsx"
)

func number(data string) {

}

type ParseScheme map[string]func(string)

type Importer struct {
	scheme map[string]func(string)
	header *xlsx.Row
	rows   []*xlsx.Row
}

func NewImporter(xslxFilePath string, scheme ParseScheme) (*Importer, error) {
	file, err := xlsx.OpenFile(xslxFilePath)
	if err != nil {
		return nil, err
	}

	sheet := file.Sheets[0]
	header := sheet.Rows[0]
	rows := sheet.Rows[1:]

	return &Importer{
		scheme: scheme,
		header: header,
		rows:   rows,
	}, nil
}

func (my *Importer) ValidateHeader() []error {
	errs := []error{}

	for _, cell := range my.header.Cells {
		cellName, err := cell.String()
		if err != nil {
			panic(err)
		}

		if _, ok := my.scheme[cellName]; ok {
			fmt.Printf("%s -- ok\n", cellName)
		} else {
			errs = append(
				errs, fmt.Errorf("Не найдена информация по полю %s\n", cellName),
			)
		}
	}

	return errs
}

var scheme1 = ParseScheme{
	"Номер":                    number,
	"Название":                 number,
	"Эпоха":                    number,
	"Описание":                 number,
	"Библиографические ссылки": number,
	"Страницы":                 number,
}

func main() {
	importer, err := NewImporter("test.xlsx", scheme1)
	if err != nil {
		panic(err)
	}
	errs := importer.ValidateHeader()
	if len(errs) != 0 {
		println("ERRORS:")
		for err := range errs {
			fmt.Printf("%+v\n", err)
		}
	}

	/*
		for _, sheet := range file.Sheets {
			for _, row := range sheet.Rows {
				for _, cell := range row.Cells {
					s, _ := cell.String()
					fmt.Printf("%s\n", s)
				}
			}
		}*/
}
