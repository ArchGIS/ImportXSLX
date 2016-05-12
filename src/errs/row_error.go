package errs

import (
	"fmt"
	"strings"
	"xl"
)

func (my *RowError) PushError(text string) {
	my.Texts = append(my.Texts, text)
}

func NewRow(row xl.Row, line int) *RowError {
	cells := make([]string, 0, len(row.Cells))
	for _, cellName := range *row.Header {
		cells = append(cells, row.Cells[cellName])
	}

	return &RowError{
		Line:  line + 1, // Учитываем header row
		Texts: make([]string, 0, len(row.Cells)),
		Cells: cells,
	}
}

func (my RowError) Error() string {
	if len(my.Texts) > 0 {
		return fmt.Sprintf(
			`{"line":%d,"texts":[%s],"cells":[%s]}`,
			my.Line,
			`"`+strings.Join(my.Texts, `","`)+`"`,
			`"`+strings.Join(my.Cells, `","`)+`"`,
		)
	} else {
		return ""
	}
}

/*
func NewRow(row xl.Row, line int, text string) *RowError {
	cells := make([]string, len(row.Cells))
	for _, cellName := range *row.Header {
		cells = append(cells, row.Cells[cellName])
	}

	return &RowError{
		Line:  line + 1, // Учитываем header row
		Text:  text,
		Cells: cells,
	}
}

func (my RowError) Error() string {
	return fmt.Sprintf(
		`{"line":%d,"text":"%s",cells:[%s]}`,
		my.Line,
		my.Text,
		`"`+strings.Join(my.Cells, `","`)+`"`,
	)
}
*/
