package helpers

import (
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/v2/interfaces"
)

// copy files from compiler output dir to store. if remove is true, remove all files in compiler out dir when done
func MergeCompilerEnvironment(ce interfaces.CompilerEnvironment, remove bool) error {
	err := linux.CopyDir(ce.CompilerOutDir(), ce.StoreDir())
	if err != nil {
		return errors.Errorf("failed to merge environments: %w", err)
	}
	if remove {
		err = utils.RecreateSafely(ce.CompilerOutDir())
		if err != nil {
			return errors.Errorf("failed to merge environments: %w", err)
		}
	}
	return nil
}
