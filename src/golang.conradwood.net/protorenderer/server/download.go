package main

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	h2g "golang.conradwood.net/apis/h2gproxy"
	pb "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/protorenderer/compiler"
	"io"
	"path/filepath"
	"strings"
)

func (p *protoRenderer) StreamHTTP(req *h2g.StreamRequest, srv pb.ProtoRendererService_StreamHTTPServer) error {
	ctx := srv.Context()
	fmt.Printf("Downloading %s\n", req.Path)
	files, err := filesForTar(ctx, req.Path)
	if err != nil {
		return err
	}
	fmt.Printf("Returning tar with %d files\n", len(files))
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	for _, f := range files {
		fmt.Printf("Adding: %s\n", f.GetFilename())
		ct, err := f.GetContent()
		if err != nil {
			fmt.Printf("cannot get content: %s\n", err)
			return err
		}
		hdr := &tar.Header{
			Name: f.GetFilename(),
			Mode: 0600,
			Size: int64(len(ct)),
		}
		err = tw.WriteHeader(hdr)
		if err != nil {
			fmt.Printf("Failed to write header: %s\n", err)
			return err
		}
		_, err = tw.Write(ct)
		if err != nil {
			fmt.Printf("Failed to write content: %s\n", err)
			return err
		}
	}
	err = tw.Close()
	if err != nil {
		return err
	}
	vfname := filepath.Base(req.Path)
	size := uint64(buf.Len())
	err = srv.Send(&h2g.StreamDataResponse{Response: &h2g.StreamResponse{
		Filename: vfname,
		Size:     size,
		MimeType: "application/tar",
	}})
	sbuf := make([]byte, 32768) // max size of buf
	read_size := 0
	for {
		rs, err := buf.Read(sbuf)
		read_size = read_size + rs
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		err = srv.Send(&h2g.StreamDataResponse{Data: sbuf[:rs]})
		if err != nil {
			return err
		}
	}

	return nil
}

type pfile struct {
	Filename string
	Content  []byte
}

func (p *pfile) GetFilename() string {
	return p.Filename
}
func (p *pfile) GetVersion() int {
	return 0
}
func (p *pfile) GetContent() ([]byte, error) {
	return p.Content, nil
}
func filesForTar(ctx context.Context, path string) ([]compiler.File, error) {
	if completeVersion == nil {
		return nil, errors.Unavailable(ctx, "not available atm - try later")
	}

	var res []compiler.File

	/********** handle protos ***********/
	if strings.Contains(path, "protos.tar") {
		result := completeVersion.metaCompiler.GetMostRecentResult()
		if result == nil {
			return nil, errors.Unavailable(ctx, "GetPackages")
		}

		for _, p := range result.Packages {
			for _, pf := range p.Protofiles {
				pl, err := protocache.GetFile(ctx, pf.Filename)
				if err != nil {
					return nil, err
				}
				f := &pfile{
					Filename: pf.Filename,
					Content:  []byte(pl.Content),
				}
				res = append(res, f)
			}
		}
		return res, nil
	}

	/********** std stuff ***********/
	filetype := ".go"
	comp := completeVersion.goCompiler
	if strings.Contains(path, "python") {
		filetype = ".py"
		comp = completeVersion.pythonCompiler
	} else if strings.Contains(path, "java") {
		filetype = ".class"
		comp = completeVersion.javaCompiler
	}

	result := completeVersion.metaCompiler.GetMostRecentResult()
	if result == nil {
		return nil, errors.Unavailable(ctx, "GetPackages")
	}

	for _, pkg := range result.Packages {
		rf, err := comp.Files(ctx, pkg.Proto, filetype)
		if err != nil {
			return nil, err
		}
		for _, xf := range rf {
			f, err := comp.GetFile(ctx, xf.GetFilename())
			if err != nil {
				return nil, err
			}
			res = append(res, f)
		}
	}
	return res, nil
}














































