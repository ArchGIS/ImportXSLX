package importer

import (
	"errs"
	"xl"
)

func New(xslxFilePath string, parser Parser) (*Importer, error) {
	table, err := xl.NewTable(xslxFilePath)
	if err != nil {
		return nil, err
	}

	return &Importer{
		table:  table,
		parser: parser,
	}, nil
}

func (my *Importer) ValidateHeader() []error {
	validationErrs := []error{}
	scheme := my.parser.Scheme()

	for _, cellName := range my.table.Header {
		if !scheme.Contains(cellName) {
			validationErrs = append(validationErrs, errs.CellInfoNotFound(cellName))
		}
	}

	return validationErrs
}

func (my *Importer) Parse() error {
	return my.parser.Parse(my.table)
}

func (my *Importer) CypherString() (string, []error) {
	return my.parser.CypherString(my.table)
}
