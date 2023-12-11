package common

import (
	"flag"
	"fmt"
	"golang.conradwood.net/go-easyops/utils"
	"os"
)

var (
	compiler_versions = flag.String("compiler_versions", "original", "refers to the directory under extra/ to look for compilers")
)

func GetCompilerVersion() string {
	return *compiler_versions
}

// if dir does not exist: create it and create a file ".filelayouter"
// if dir DOES exist, check for existence of a file  ".filelayouter", if so, recreate
// otherwise error
func RecreateSafely(dirname string) error {
	var err error
	fname := dirname + "/.protorenderer"
	if utils.FileExists(dirname) {
		if utils.FileExists(fname) {
			err = os.RemoveAll(dirname)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("not recreating dir \"%s\", it exists already, but .filelayouter does not.", dirname)
		}
	}
	if utils.FileExists(dirname) {
		return fmt.Errorf("Attempt to delete dir \"%s\" failed", dirname)
	}
	err = os.MkdirAll(dirname, 0777)
	if err != nil {
		return err
	}
	err = utils.WriteFile(fname, make([]byte, 0))
	return err
}











































































