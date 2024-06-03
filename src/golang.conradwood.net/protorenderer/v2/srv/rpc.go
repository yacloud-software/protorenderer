package srv

import (
	"context"
	pb1 "golang.conradwood.net/apis/protorenderer"
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/utils"
	"io"
)

func (pr *protoRenderer) UpdateProto(ctx context.Context, req *pb1.AddProtoRequest) (*pb1.AddProtoResponse, error) {
	return nil, nil
}
func (pr *protoRenderer) Compile(srv pb.ProtoRenderer2_CompileServer) error {
	bsr := utils.NewByteStreamReceiver(CompileEnv.WorkDir() + "/" + CompileEnv.NewProtosDir())
	for {
		rcv, err := srv.Recv()
		if rcv != nil {
			err = bsr.NewData(rcv)
			if err != nil {
				return err
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	err := bsr.Close()
	if err != nil {
		return err
	}
	return nil
}
