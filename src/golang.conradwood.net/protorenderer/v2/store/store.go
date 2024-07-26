package store

import (
	"context"
	"golang.conradwood.net/protorenderer/v2/store/binaryversions"
)

func Store(ctx context.Context, dir string) error {
	err := binaryversions.Upload(dir)
	return err
}
