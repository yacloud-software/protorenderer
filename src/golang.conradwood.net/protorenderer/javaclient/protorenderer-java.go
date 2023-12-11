package main

import (
	"flag"
	"fmt"
	"golang.conradwood.net/apis/common"
	pr "golang.conradwood.net/apis/protorenderer"
	ar "golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/utils"
)

func main() {
	flag.Parse()
	ps := pr.GetProtoRendererServiceClient()
	ctx := ar.Context()
	pkgs, err := ps.GetPackages(ctx, &common.Void{})
	utils.Bail("failed to get packages", err)
	for _, pkg := range pkgs.Packages {
		fmt.Printf("Package: %s %s\n", pkg.Prefix, pkg.Name)
		ctx := ar.Context()
		fl, err := ps.GetFilesJavaClass(ctx, &pr.ID{ID: pkg.ID})
		if err != nil {
			fmt.Printf("Error retrieving java class file names for package. %s\n", err)
			continue
		}
		for _, f := range fl.Files {
			fmt.Printf("    %s  ", f)
			ff, err := ps.GetFile(ctx, &pr.FileRequest{Filename: f})
			if err != nil {
				fmt.Printf("Unable to get file: %s\n", err)
				continue
			}
			fmt.Printf(" size=%d Bytes\n", len(ff.Content))
		}
	}
	fmt.Printf("Done.\n")
}




















