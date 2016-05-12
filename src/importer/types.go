package importer

import (
	"xl"

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
	table   *xl.Table
	scheme  ParseScheme
	header  *xlsx.Row
	rows    []*xlsx.Row
	epochs  map[string]string
	indexes map[string]int
}
