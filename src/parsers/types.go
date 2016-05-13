package parsers

import (
	"importer"
)

type Parser1 struct {
	scheme   *importer.ParseScheme
	epochs   map[string]string
	cultures map[string]string
}
