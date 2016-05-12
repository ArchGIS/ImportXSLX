package xl

import (
	"github.com/tealeg/xlsx"
)

func toString(cell *xlsx.Cell) string {
	stringValue, err := cell.String()
	if err != nil {
		// Не знаю, что здесь за ошибка может быть, лень
		// смотреть исходники функции xlsx.Cell.String(). Возможно она всегда nil
		// ошибку возвращает. Это можно проверить как-нибудь потом.
		panic(err)
	}

	return stringValue
}

func NewTable(xlsxFilePath string) (*Table, error) {
	file, err := xlsx.OpenFile(xlsxFilePath)
	if err != nil {
		return nil, err
	}

	xlsxSheet := file.Sheets[0]     // Работаем только с первым листом
	xlsxHeader := xlsxSheet.Rows[0] // Первый ряд - заголок
	xlsxRows := xlsxSheet.Rows[1:]  // Остальные - данные

	// Подготавливаем header для Table
	header := make(Header, len(xlsxHeader.Cells))
	for index, cell := range xlsxHeader.Cells {
		header[index] = toString(cell)
	}

	// Подготавливаем rows для Table
	rows := make([]Row, len(xlsxRows))
	for rowIndex, xlsxRow := range xlsxRows {
		/*
			row := Row{
				Header: &header,
				Cells:  make([]string, len(header)),
			}

			for cellIndex, cell := range xlsxRow.Cells {
				row.Cells[cellIndex] = toString(cell)
			}
		*/

		row := Row{
			Header: &header,
			Cells:  make(map[string]string, len(header)),
		}

		for cellIndex, cell := range xlsxRow.Cells {
			row.Cells[header[cellIndex]] = toString(cell)
		}

		rows[rowIndex] = row
	}

	return &Table{
		Header: header,
		Rows:   rows,
	}, nil
}
