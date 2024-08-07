package protoretriever

import (
	"context"
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/utils"
	"io"
)

func ByFilename(ctx context.Context, filename, outdir string) error {
	req := &pb.ProtoFileRequest{ProtoFileName: filename}
	srv, err := pb.GetProtoRenderer2Client().GetProtoFile(ctx, req)
	if err != nil {
		return err
	}
	bs := utils.NewByteStreamReceiver(outdir)
	defer bs.Close()
	for {
		b, err := srv.Recv()
		xerr := bs.NewData(b)
		if xerr != nil {
			return xerr
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	bs.Close()
	return nil
}
