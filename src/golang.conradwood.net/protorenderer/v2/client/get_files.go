package main

import (
	"fmt"
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/utils"
	"io"
)

func GetPackageFiles() error {
	outdir := "/tmp/protomanager"
	ctx := authremote.Context()
	req := &pb.CompiledFilesRequest{Package: *get_package_files}
	srv, err := pb.GetProtoRenderer2Client().GetCompiledFiles(ctx, req)
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
	fmt.Printf("Retrieved %d files, saved to %s\n", bs.FileCount(), outdir)
	return nil
}
