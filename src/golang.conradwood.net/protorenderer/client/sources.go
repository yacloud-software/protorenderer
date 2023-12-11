package main

import (
	"flag"
	"fmt"
	"golang.conradwood.net/apis/common"
	pr "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/utils"
	"os"
	"path/filepath"
	"strings"
)

func Sources() {
	as := flag.Args()
	filter := ""
	if len(as) != 0 {
		filter = strings.ToLower(as[0])
	}
	ps := pr.GetProtoRendererServiceClient()
	ctx := getContext()
	pkgs, err := ps.GetPackages(ctx, &common.Void{})
	utils.Bail("failed to get packages", err)
	for _, fpkg := range pkgs.Packages {
		if filter != "" && !strings.Contains(strings.ToLower(fpkg.Name), filter) {
			continue
		}
		ctx := getContext()
		pkgid := &pr.ID{ID: fpkg.ID}
		pkg, err := ps.GetPackageByID(ctx, pkgid)
		utils.Bail("Failed to get package", err)
		fmt.Printf("Package: #%s %s (prefix=%s) (fpkgname=%s)\n", pkg.ID, pkg.Name, pkg.Prefix, fpkg.Name)
		fl, err := ps.GetFilesProto(ctx, pkgid)
		utils.Bail("failed to get files", err)
		for i, filename := range fl.Files {
			file, err := ps.GetFile(ctx, &pr.FileRequest{PackageID: pkgid, Filename: filename})
			utils.Bail("failed to get file", err)
			fmt.Printf("    %02d File: %s (Repo %d)\n", i, filename, file.RepositoryID)
			fullpath := "/tmp/x/protos/" + filename
			dn := filepath.Dir(fullpath)
			os.MkdirAll(dn, 0777)
			utils.Bail("failed to write", utils.WriteFile(fullpath, file.Content))
			fmt.Printf("written to %s\n", fullpath)
		}
	}
}

































