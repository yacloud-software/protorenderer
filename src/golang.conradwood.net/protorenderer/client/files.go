package main

import (
	"flag"
	"fmt"
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/protorenderer"
	ar "golang.conradwood.net/go-easyops/authremote"
	"os"
	"path/filepath"
	"strings"
	//	"golang.conradwood.net/go-easyops/tokens"
	"golang.conradwood.net/go-easyops/utils"
)

func get_files() {
	save := false
	if *outdir == "" {
		fmt.Printf("Forgot to specify outdir. Not saving files...\n")
	} else {
		save = true
		os.MkdirAll(*outdir, 0777)
	}
	protoClient = pb.GetProtoRendererServiceClient()
	ctx := ar.Context()
	pkgs, err := protoClient.GetPackages(ctx, &common.Void{})
	utils.Bail("failed to get packages", err)
	for _, pkg := range pkgs.Packages {
		if !packageMatches(pkg) {
			continue
		}
		pkgid := &pb.ID{ID: pkg.ID}
		gofiles, err := protoClient.GetFilesGO(ctx, pkgid)
		utils.Bail("failed to get go files", err)

		javafiles, err := protoClient.GetFilesJavaClass(ctx, pkgid)
		utils.Bail("failed to get java files", err)

		protofiles, err := protoClient.GetFilesProto(ctx, pkgid)
		utils.Bail("failed to get proto files", err)

		nanofiles, err := protoClient.GetFilesNanoPB(ctx, pkgid)
		utils.Bail("failed to get proto files", err)

		fmt.Printf("Name: %s (%d go files, %d java files, %d protos)\n", pkg.Name, len(gofiles.Files), len(javafiles.Files), len(protofiles.Files))
		for _, f := range gofiles.Files {
			ctx = ar.Context()
			fr := &pb.FileRequest{Compiler: pb.CompilerType_GOLANG, PackageID: pkgid, Filename: f}
			file, err := protoClient.GetFile(ctx, fr)
			utils.Bail("failed to get file", err)
			fmt.Printf("  Go File  : %s (%d bytes) (repo=%d)\n", f, len(file.Content), file.RepositoryID)
			if save {
				write(*outdir+"/go/"+f, file.Content)
			}
		}
		for _, f := range nanofiles.Files {
			ctx = ar.Context()
			fr := &pb.FileRequest{Compiler: pb.CompilerType_NANOPB, PackageID: pkgid, Filename: f}
			file, err := protoClient.GetFile(ctx, fr)
			utils.Bail("failed to get file", err)
			fmt.Printf("  Go File  : %s (%d bytes)\n", f, len(file.Content))
			if save {
				write(*outdir+"/nanopb/"+f, file.Content)
			}
		}

		for _, f := range javafiles.Files {
			fmt.Printf("  Java File: %s\n", f)
			ctx = ar.Context()
			file, err := protoClient.GetFile(ctx, &pb.FileRequest{PackageID: pkgid, Filename: f})
			utils.Bail("failed to get file", err)
			if save {
				write(*outdir+"/java/"+f, file.Content)
			}
		}

		for _, f := range protofiles.Files {
			fmt.Printf(" Proto File: %s\n", f)
			ctx = ar.Context()
			file, err := protoClient.GetFile(ctx, &pb.FileRequest{PackageID: pkgid, Filename: f})
			utils.Bail("failed to get file", err)
			if save {
				write(*outdir+"/protos/"+f, file.Content)
			}
		}

	}
}

func write(s string, b []byte) {
	p := filepath.Dir(s)
	os.MkdirAll(p, 0777)
	err := utils.WriteFile(s, b)
	utils.Bail("failed to write file", err)
}

func packageMatches(fp *pb.FlatPackage) bool {
	if len(flag.Args()) == 0 {
		return true
	}
	sname := strings.ToLower(fp.Name)
	spkg := strings.ToLower(fp.Prefix)
	for _, arg := range flag.Args() {
		a := strings.ToLower(arg)
		if strings.Contains(sname, a) {
			return true
		}
		if strings.Contains(spkg, a) {
			return true
		}
	}
	return false
}






















































































