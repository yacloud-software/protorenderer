package main

import (
	"context"
	"flag"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/v2/compilers/golang"
	"golang.conradwood.net/protorenderer/v2/interfaces"
)

func main() {
	flag.Parse()
	_, pfs, err := find_protofiles()
	utils.Bail("failed to find protofiles", err)
	golang_compiler := golang.New()
	ctx := context.Background()
	golang_compiler.Compile(ctx, nil, pfs)

}

func find_protofiles() (string, []interfaces.ProtoFile, error) {
	dir, err := utils.FindFile("extra/protos")
	if err != nil {
		return dir, nil, err
	}
	var filenames []string
	utils.DirWalk(dir, func(root, relfil string) error {
		filenames = append(filenames, relfil)
		return nil
	},
	)
	var res []interfaces.ProtoFile
	return dir, res, nil
}
