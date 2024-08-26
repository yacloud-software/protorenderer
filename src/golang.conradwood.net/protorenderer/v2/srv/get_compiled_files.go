package srv

import (
	//	"context"
	"fmt"
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/utils"
	"strings"
)

func (e *protoRenderer) GetCompiledFiles(req *pb.CompiledFilesRequest, srv pb.ProtoRenderer2_GetCompiledFilesServer) error {
	package_name := req.Package
	if strings.Contains(package_name, "..") {
		return errors.Errorf("package name invalid \"%s\"", package_name)
	}
	fmt.Printf("Getting files for package \"%s\"\n", package_name)

	bss := utils.NewByteStreamSender(
		func(key, filename string) error {
			return srv.Send(&pb.FileTransfer{Filename: filename})
		},
		func(data []byte) error {
			return srv.Send(&pb.FileTransfer{Data: data})
		},
	)
	// 1. meta file:
	err := send_files_in_dir(bss, "info/"+package_name)
	if err != nil {
		return err
	}

	// now for each compiler...
	compilers := getCompilerArray()
	ctx := srv.Context()
	for _, compiler := range compilers {
		package_dirs, err := compiler.DirsForPackage(ctx, package_name)
		if err != nil {
			return err
		}
		for _, pd := range package_dirs {
			fmt.Printf("  Compiler \"%s\": %s\n", compiler.ShortName(), pd)
		}
	}
	return nil
}
func send_files_in_dir(bss *utils.ByteStreamSender, dir string) error {
	root := CompileEnv.StoreDir() + "/" + dir
	var fnames []string
	err := utils.DirWalk(root, func(r, rel string) error {
		fnames = append(fnames, rel)
		return nil
	})
	if err != nil {
		return err
	}
	for _, fname := range fnames {
		fmt.Printf("filename: %s\n", fname)
		b, err := utils.ReadFile(root + "/" + fname)
		dname := dir + "/" + fname
		err = bss.SendBytes(dname, dname, b)
		if err != nil {
			return err
		}
	}

	return nil
}
