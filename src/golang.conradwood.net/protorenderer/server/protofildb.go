package main

import (
	"context"
	pb "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/cache"
	"time"
)

var (
	pcache = cache.New("protofilecache", time.Duration(999)*time.Hour, 9999)
)

// for now, we're treating files with the same name the same, regardless of repo
func findOrUpdateProtoInDB(ctx context.Context, req *pb.AddProtoRequest) (*pb.DBProtoFile, error) {

	res, err := FindByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if res == nil {
		res = &pb.DBProtoFile{Name: req.Name, RepositoryID: req.RepositoryID}
		_, err = dbproto.Save(ctx, res)
		if err != nil {
			return nil, err
		}
		pcache.Put(req.Name, res)
	}

	if res.RepositoryID != req.RepositoryID && req.RepositoryID != 0 {
		res.RepositoryID = req.RepositoryID
		err = dbproto.Update(ctx, res)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func FindByName(ctx context.Context, name string) (*pb.DBProtoFile, error) {
	p := pcache.Get(name)
	if p != nil {
		return p.(*pb.DBProtoFile), nil
	}
	dbp, err := dbproto.ByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if len(dbp) != 0 {
		d := dbp[0]
		pcache.Put(name, d)
		return d, nil
	}
	return nil, nil
}








































