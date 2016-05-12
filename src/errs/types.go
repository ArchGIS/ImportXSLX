package errs

import (
	"time"
)

type FatalError struct {
	When time.Time
	What string
}

type RowError struct {
	Line  int
	Text  string
	Cells []string
}
