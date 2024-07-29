package golang

import (
	"flag"
	"fmt"
)

const (
	DEBUG_PREFIX = "[golang-compiler] "
)

var (
	debug = flag.Bool("debug_golang_compiler", false, "debug golang compiler output")
)

func Debugf(format string, args ...interface{}) {
	if !*debug {
		return
	}
	s := fmt.Sprintf(format, args...)
	fmt.Print(DEBUG_PREFIX + s)
}
