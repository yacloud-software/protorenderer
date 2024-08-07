package srv

import (
	"context"
	pb "golang.conradwood.net/apis/protorenderer"
	pb2 "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/cache"
	"time"
)

var (
	pcache = cache.New("protofilecache", time.Duration(999)*time.Hour, 9999)
)

// for now, we're treating files with the same name the same, regardless of repo
func findOrUpdateProtoInDB(ctx context.Context, req *pb.AddProtoRequest) (*pb2.DBProtoFile, error) {

	res, err := FindByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if res == nil {
		res = &pb2.DBProtoFile{Filename: req.Name, RepositoryID: req.RepositoryID}
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

func FindByName(ctx context.Context, name string) (*pb2.DBProtoFile, error) {
	p := pcache.Get(name)
	if p != nil {
		return p.(*pb2.DBProtoFile), nil
	}
	dbp, err := dbproto.ByFilename(ctx, name)
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
