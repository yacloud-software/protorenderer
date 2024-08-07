package protosubmitter

import (
	"golang.conradwood.net/go-easyops/utils"
	"os"
)

func IsDir(name string) (bool, error) {
	fname, err := utils.FindFile(name)
	if err != nil {
		return false, err
	}
	// This returns an *os.FileInfo type
	fileInfo, err := os.Stat(fname)
	if err != nil {
		return false, err
	}

	// IsDir is short for fileInfo.Mode().IsDir()
	if fileInfo.IsDir() {
		return true, nil
	}
	return false, nil
}
