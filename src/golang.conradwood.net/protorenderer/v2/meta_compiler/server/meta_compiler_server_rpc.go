package server

import (
	"context"
	"fmt"
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/protorenderer2"
	mcomp "golang.conradwood.net/protorenderer/v2/meta_compiler"
)

type ServerMetaCompiler struct {
	mc *mcomp.MetaCompiler
}

// called by the protoc plugin
func InternalMetaSubmit(ctx context.Context, req *pb.ProtocRequest) (*common.Void, error) {
	mc, err := mcomp.GetMetaCompilerByID(ctx, req.MetaCompilerID)
	if err != nil {
		fmt.Printf("meta compiler server invoked with an invalid meta compiler id (%s)\n", err)
		return nil, err
	}
	icm := &ServerMetaCompiler{mc: mc}
	for _, pf := range req.ProtoFiles {
		fmt.Printf("meta compiler: %#v\n", pf)
	}

	for _, pf := range req.ProtoFiles {
		fmt.Printf("Protofile: %s\n", *pf.Name)
		err := icm.handle_protofile(ctx, pf)
		if err != nil {
			return nil, err
		}
	}

	return &common.Void{}, nil
}
