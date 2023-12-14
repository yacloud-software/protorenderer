package main

import (
	"fmt"
	pb "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/utils"
	"io"
	"os"
	"path/filepath"
)

func GetZip(pkg string) {
	fmt.Printf("Getting zip file for package \"%s\"\n", pkg)
	protoClient = pb.GetProtoRendererServiceClient()
	ctx := authremote.Context()
	srv, err := protoClient.GetFilesGoByPackageName(ctx, &pb.PackageName{PackageName: pkg})
	utils.Bail("failed to get zip", err)
	wd := "/tmp/x/"
	currentFilename := ""
	var curFile *os.File
	for {
		zf, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			utils.Bail("failed to get zip", err)
		}
		if zf.Filename != currentFilename {
			fmt.Printf("new file: %s\n", zf.Filename)
			if curFile != nil {
				curFile.Close()
				curFile = nil
			}
			fname := wd + zf.Filename
			cdir(fname)
			curFile, err = os.Create(fname)
			utils.Bail("failed to create file", err)
			currentFilename = zf.Filename
		}
		_, err = curFile.Write(zf.Payload)
		utils.Bail("failed to write", err)
	}
	if curFile != nil {
		curFile.Close()
		curFile = nil
	}
	fmt.Printf("files in %s\n", wd)
}
func cdir(fname string) {
	s := filepath.Dir(fname)
	os.MkdirAll(s, 0777)
}








































































































