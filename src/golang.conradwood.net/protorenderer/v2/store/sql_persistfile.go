package store

import (
	"context"
	pb2 "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/protorenderer/db"
	"strings"
)

// one may set a repoid, but once it is set (that is: not-zero), it may not be changed
func GetOrCreateFile(ctx context.Context, filename string, repoid uint64) (*pb2.DBProtoFile, error) {
	if repoid == 0 {
		return nil, errors.InvalidArgs(ctx, "missing repositoryid", "missing repositoryid for file \"%s\"", filename)
	}
	fname := strings.TrimPrefix(filename, "protos/")
	files, err := db.DefaultDBDBProtoFile().ByName(ctx, fname)
	if err != nil {
		return nil, err
	}
	var res *pb2.DBProtoFile
	if len(files) != 0 {
		res = files[0]
	} else {
		res = &pb2.DBProtoFile{Name: fname, RepositoryID: repoid}
		_, err = db.DefaultDBDBProtoFile().Save(ctx, res)
		if err != nil {
			return nil, err
		}
	}

	if res.RepositoryID == 0 {
		// update repositoryid of database
		res.RepositoryID = repoid
		err = db.DefaultDBDBProtoFile().Update(ctx, res)
		if err != nil {
			return nil, err
		}
	} else {
		if repoid != res.RepositoryID {
			return nil, errors.InvalidArgs(ctx, "repoid mismatch", "mismatch of repository id. file \"%s\" has repository id %d previously, but changed to %d", fname, res.RepositoryID, repoid)
		}
	}
	return res, nil
}
func FileByName(ctx context.Context, filename string) (*pb2.DBProtoFile, error) {
	fname := strings.TrimPrefix(filename, "protos/")
	files, err := db.DefaultDBDBProtoFile().ByName(ctx, fname)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, errors.NotFound(ctx, "file does not exist", "file %s does not exist", filename)
	}
	return files[0], nil
}

// update the package option in database
func UpdatePackageInFile(ctx context.Context, pf *pb2.DBProtoFile, packagename string) error {
	pf.Package = packagename
	err := db.DefaultDBDBProtoFile().Update(ctx, pf)
	return err
}
