package importer

import (
	"xl"
)

type Parser interface {
	Parse(*xl.Table) error
	Scheme() *ParseScheme
	CypherString(*xl.Table) (string, []error)
}

type ParseScheme struct {
	Name  string
	Cells []string
}

type Importer struct {
	parser Parser
	table  *xl.Table
}
