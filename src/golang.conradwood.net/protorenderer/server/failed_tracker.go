package main

import (
	"golang.conradwood.net/protorenderer/compiler"
)

type failuretracker struct {
	failures []*failure_tracked
}
type failure_tracked struct {
	c        compiler.Compiler
	filename string
	message  string
}

func (f *failuretracker) AddFailed(c compiler.Compiler, filename, message string) {
	ft := &failure_tracked{
		c:        c,
		filename: filename,
		message:  message,
	}
	f.failures = append(f.failures, ft)
}
func (f *failuretracker) Failures() []*failure_tracked {
	return f.failures
}


















































































