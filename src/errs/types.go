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
	Texts []string
	Cells []string
}
