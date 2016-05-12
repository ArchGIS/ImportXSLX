package importer

import (
	"fmt"
)

func (my ParseScheme) Find(name string) ParseSchemeCell {
	for _, cell := range my.Cells {
		if cell.Name == name {
			return cell
		}
	}

	panic(fmt.Sprintf("%s: not found cell named '%s'", my.Name, name))
}

func (my ParseScheme) IndexOf(name string) int {
	for index, cell := range my.Cells {
		if cell.Name == name {
			return index
		}
	}

	panic(fmt.Sprintf("%s: not found index of cell named '%s'", my.Name, name))
}
