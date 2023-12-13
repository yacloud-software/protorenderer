package compiler

import (
	"fmt"
	"golang.conradwood.net/go-easyops/utils"
)

type StdFile struct {
	// exposed filename
	Filename string
	// filename on disk
	ctFilename string
	version    int
}

func (s *StdFile) GetVersion() int {
	return s.version
}
func (s *StdFile) GetContent() ([]byte, error) {
	if s.ctFilename == "" {
		return nil, fmt.Errorf("no location for file \"%s\" found", s.Filename)
	}
	return utils.ReadFile(s.ctFilename)
}

func (s *StdFile) GetFilename() string {
	return s.Filename
}



































































































