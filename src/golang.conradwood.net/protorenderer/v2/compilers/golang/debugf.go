package golang

import (
	"fmt"
)

const (
	DEBUG_PREFIX = "[golang-compiler] "
)

func Debugf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	fmt.Print(DEBUG_PREFIX + s)
}
