package meta_compiler

import (
	"context"
	"golang.conradwood.net/go-easyops/cache"
	"golang.conradwood.net/go-easyops/errors"
	"time"
)

var (
	meta_compilers = cache.New("meta_compiler_cache", time.Duration(30)*time.Minute, 100)
)

func GetMetaCompilerByID(ctx context.Context, id string) (*MetaCompiler, error) {
	omc := meta_compilers.Get(id)
	if omc == nil {
		return nil, errors.NotFound(ctx, "compiler not found")
	}

	return omc.(*MetaCompiler), nil
}
