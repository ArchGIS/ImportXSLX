package errs

import (
	"fmt"
	"time"
)

func (my FatalError) Error() string {
	return fmt.Sprintf("%v (when: %v)", my.What, my.When)
}

func NewFatal(what string) FatalError {
	return FatalError{When: time.Now(), What: what}
}

func NewFatalf(whatf string, args ...interface{}) FatalError {
	return FatalError{When: time.Now(), What: fmt.Sprintf(whatf, args...)}
}
