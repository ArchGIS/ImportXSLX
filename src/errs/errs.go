package errs

import (
	"fmt"
)

func CellInfoNotFound(cellName string) error {
	return fmt.Errorf("Не найдена информация по полю %s", cellName)
}
