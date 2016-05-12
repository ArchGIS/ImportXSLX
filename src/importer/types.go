package importer

import (
	"github.com/tealeg/xlsx"
)

type ParseSchemeCell struct {
	Name   string
	Parser func(string)
}

type ParseScheme struct {
	Name  string
	Cells []ParseSchemeCell
}

type Importer struct {
	scheme  ParseScheme
	header  *xlsx.Row
	rows    []*xlsx.Row
	epochs  map[string]string
	indexes map[string]int
}
