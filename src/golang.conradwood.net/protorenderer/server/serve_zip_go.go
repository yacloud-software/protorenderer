package main

import (
	"fmt"
	pb "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/protorenderer/compiler"
	"golang.conradwood.net/protorenderer/meta"
)

func (e *protoRenderer) GetFilesGoByPackageName(req *pb.PackageName, srv pb.ProtoRendererService_GetFilesGoByPackageNameServer) error {
	ctx := srv.Context()
	ev := NeedVersion(ctx)
	if ev != nil {
		return ev
	}
	result := completeVersion.metaCompiler.GetMostRecentResult()
	if result == nil {
		return errors.Unavailable(ctx, "GetPackages (most recent result)")
	}
	var packages []*meta.Package
	for _, p := range result.Packages {
		pf := p.FQDN
		if pf != req.PackageName {
			continue
		}
		packages = append(packages, p)
		fmt.Printf("pf==%s\n", pf)
	}

	// get files from the gocompiler
	filetype := ".go"
	file_compiler := completeVersion.goCompiler

	// "packages" is now the list of packages we want..
	var files []compiler.File
	for _, pkg := range packages {
		rf, err := file_compiler.Files(ctx, pkg.Proto, filetype)
		if err != nil {
			return err
		}
		for _, f := range rf {
			files = append(files, f)
			fmt.Printf("Adding file %s\n", f.GetFilename())
		}
	}

	// "files" now contains the files we want
	fz := &zipcopier{srv: srv}
	for _, f := range files {
		buf, err := f.GetContent()
		if err != nil {
			fmt.Printf("Failed to get content for \"%s\": %s\n", f.GetFilename(), err)
			return err
		}
		_, err = fz.Write(f.GetFilename(), buf)
		if err != nil {
			fmt.Printf("Failed to write file %s: %s\n", f.GetFilename(), err)
			return err
		}
		fmt.Printf("File sent: %s\n", f.GetFilename())
	}
	return nil
}

type zipcopier struct {
	srv pb.ProtoRendererService_GetFilesGoByPackageNameServer
}

func (z *zipcopier) Write(filename string, buf []byte) (int, error) {
	start := 0
	for {
		to_send := len(buf)
		if to_send > 8192 {
			to_send = 8192
		}
		if to_send+start > len(buf) {
			to_send = len(buf) - start
		}
		if to_send == 0 {
			break
		}
		fmt.Printf("Sending %d bytes (offset %d)...\n", to_send, start)
		zf := &pb.FileStream{
			Filename: filename,
			Payload:  buf[start : start+to_send],
		}
		err := z.srv.Send(zf)
		if err != nil {
			return 0, err
		}
		start = start + to_send
	}
	return len(buf), nil
}












































