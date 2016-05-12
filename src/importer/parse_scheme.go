package importer

import (
	"errs"
)

func (my ParseScheme) Contains(name string) bool {
	return my.Find(name) != ""
}

func (my ParseScheme) Find(name string) string {
	for _, cell := range my.Cells {
		if cell == name {
			return cell
		}
	}

	return ""
}

func (my ParseScheme) IndexOf(name string) int {
	for index, cell := range my.Cells {
		if cell == name {
			return index
		}
	}

	panic(errs.NewFatalf("%s: not found index of cell named '%s'", my.Name, name))
}
