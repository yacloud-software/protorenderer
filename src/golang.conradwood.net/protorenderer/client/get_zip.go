package main

import (
	"fmt"
	pb "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/utils"
	"io"
)

func GetZip(pkg string) {
	fmt.Printf("Getting zip file for package \"%s\"\n", pkg)
	protoClient = pb.GetProtoRendererServiceClient()
	ctx := authremote.Context()
	srv, err := protoClient.GetFilesGoByPackageName(ctx, &pb.PackageName{PackageName: pkg})
	utils.Bail("failed to get zip", err)
	res := make([]byte, 0)
	for {
		zf, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			utils.Bail("failed to get zip", err)
		}
		res = append(res, zf.Payload...)
	}
	err = utils.WriteFile("/tmp/proto.zip", res)
	utils.Bail("failed to write file", err)
	fmt.Printf("zipfile in /tmp/proto.zip (length=%d bytes)\n", len(res))
}
