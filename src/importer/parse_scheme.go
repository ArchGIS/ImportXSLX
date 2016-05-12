package importer

import (
	"errs"
)

func (my ParseScheme) Contains(name string) bool {
	cell := my.Find(name)
	return cell.Name != ""
}

func (my ParseScheme) Find(name string) ParseSchemeCell {
	for _, cell := range my.Cells {
		if cell.Name == name {
			return cell
		}
	}

	return ParseSchemeCell{}
}

func (my ParseScheme) IndexOf(name string) int {
	for index, cell := range my.Cells {
		if cell.Name == name {
			return index
		}
	}

	panic(errs.NewFatalf("%s: not found index of cell named '%s'", my.Name, name))
}
