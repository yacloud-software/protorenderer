package helpers

import (
	"golang.conradwood.net/go-easyops/utils"
	"path/filepath"
)

// write file, create parent directories as required
func WriteFileWithDir(filename string, content []byte) error {
	err := utils.WriteFile(filename, content)
	if err == nil {
		return nil
	}
	dir := filepath.Dir(filename)
	err = Mkdir(dir)
	if err != nil {
		return err
	}
	err = utils.WriteFile(filename, content)
	return err
}
