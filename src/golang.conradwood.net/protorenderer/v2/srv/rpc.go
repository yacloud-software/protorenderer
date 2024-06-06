package srv

import (
	"context"
	pb1 "golang.conradwood.net/apis/protorenderer"
	//	pb "golang.conradwood.net/apis/protorenderer2"
	//	"golang.conradwood.net/go-easyops/utils"
	//	"io"
	"sync"
)

var (
	compile_lock sync.Mutex
)

func (pr *protoRenderer) UpdateProto(ctx context.Context, req *pb1.AddProtoRequest) (*pb1.AddProtoResponse, error) {
	return nil, nil
}
