package helpers

import (
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/v2/interfaces"
)

// copy files from compiler output dir to store. if remove is true, remove all files in compiler out dir when done
func MergeCompilerEnvironment(source, target interfaces.CompilerEnvironment, remove bool) error {
	err := linux.CopyDir(source.CompilerOutDir(), target.StoreDir())
	if err != nil {
		return errors.Errorf("(1) failed to merge environments: %w", err)
	}
	err = linux.CopyDir(source.NewProtosDir(), target.StoreDir()+"/protos")
	if err != nil {
		return errors.Errorf("(2) failed to merge environments: %w", err)
	}
	if remove {
		err = utils.RecreateSafely(source.CompilerOutDir())
		if err != nil {
			return errors.Errorf("failed to merge environments: %w", err)
		}
	}
	target.MetaCache().ImportFrom(source.MetaCache()) // update the meta cache with new information
	return nil
}
