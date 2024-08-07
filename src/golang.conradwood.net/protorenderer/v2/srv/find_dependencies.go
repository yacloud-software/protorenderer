package srv

import (
	"context"
	"fmt"
	pb "golang.conradwood.net/apis/protorenderer2"
	//	"golang.conradwood.net/go-easyops/utils"
	// "strings"
	// "sync"
)

func (pr *protoRenderer) GetReverseDependencies(ctx context.Context, req *pb.ReverseDependenciesRequest) (*pb.ReverseDependenciesResponse, error) {
	mc := CompileEnv.MetaCache()
	fmt.Printf("Getting dependencies for \"%s\"\n", req.Filename)

	res := &pb.ReverseDependenciesResponse{}
	with_deps, err := mc.AllWithDependencyOn(req.Filename, req.MaxDepth)
	if err != nil {
		return nil, err
	}
	for _, meta := range with_deps {
		res.Filenames = append(res.Filenames, meta.ProtoFile.Filename)
	}
	return res, nil
}
