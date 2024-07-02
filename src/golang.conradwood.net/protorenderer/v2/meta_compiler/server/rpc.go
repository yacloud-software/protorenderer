package server

import (
	"context"
	"fmt"
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/protorenderer2"
)

func SubmitSource(ctx context.Context, req *pb.ProtocRequest) (*common.Void, error) {
	for _, pf := range req.ProtoFiles {
		fmt.Printf("meta compiler: %#v\n", pf)
	}
	return &common.Void{}, nil
}
