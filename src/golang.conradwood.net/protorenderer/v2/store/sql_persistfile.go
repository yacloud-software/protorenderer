package store

import (
	"context"
	pb2 "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/protorenderer/db"
	"strings"
)

func GetFileID(ctx context.Context, filename string, repoid uint64) (uint64, error) {
	if repoid == 0 {
		return 0, errors.InvalidArgs(ctx, "missing repositoryid", "missing repositoryid for file \"%s\"", filename)
	}
	fname := strings.TrimPrefix(filename, "protos/")
	files, err := db.DefaultDBDBProtoFile().ByName(ctx, fname)
	if err != nil {
		return 0, err
	}
	if len(files) != 0 {
		return files[0].ID, nil
	}
	dd := &pb2.DBProtoFile{Name: fname, RepositoryID: repoid}
	_, err = db.DefaultDBDBProtoFile().Save(ctx, dd)
	if err != nil {
		return 0, err
	}
	return dd.ID, nil
}
