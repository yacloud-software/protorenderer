package binaryversions

import (
	"context"
	"fmt"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/utils"
	pb "golang.yacloud.eu/apis/binaryversions"
	"io"
	"os"
	"path/filepath"
)

// version 0 == default
func Download(ctx context.Context, realm, destination string, version uint64) error {
	if !*use_store {
		return nil
	}
	fmt.Printf("[store] Downloading \"%s\" to \"%s\"\n", PROTORENDERER_STORE_DIR_NAME, destination)
	c := pb.GetBinaryVersionsClient()
	dir, err := c.MkOrGetDir(ctx, &pb.MkDirRequest{DirName: PROTORENDERER_STORE_DIR_NAME, Realm: &pb.Realm{Name: realm}})
	if err != nil {
		fmt.Printf("failed to mkdir %s: %s\n", PROTORENDERER_STORE_DIR_NAME, err)
		return errors.Wrap(err)
	}

	dvl, err := c.DirVersions(ctx, &pb.DirVersionRequest{DirName: PROTORENDERER_STORE_DIR_NAME})
	if err != nil {
		return errors.Wrap(err)
	}
	if len(dvl.Version) == 0 {
		// empty
		fmt.Printf("No previous version of \"%s\"\n", PROTORENDERER_STORE_DIR_NAME)
		err = utils.RecreateSafely(destination)
		return errors.Wrap(err)
	}

	ov := &pb.OpenDirVersionRequest{Directory: dir, Version: version}
	v, err := c.OpenDirVersion(ctx, ov)
	if err != nil {
		return errors.Wrap(err)
	}
	dfr := &pb.DownloadFileRequest{DirVersion: v}
	srv, err := c.DownloadFiles(ctx, dfr)
	if err != nil {
		return errors.Wrap(err)
	}
	var wr *DownloadWriter
	for {
		fdata, err := srv.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err)
		}
		if fdata.FileName == "" {
			if wr == nil {
				return fmt.Errorf("received data before filename")
			}
			err = wr.Write(fdata)
			if err != nil {
				return errors.Wrap(err)
			}
			continue
		}
		if wr != nil {
			wr.Close()
			wr = nil
		}
		fname := fmt.Sprintf("%s/%s", destination, fdata.FileName)
		dname := filepath.Dir(fname)
		os.MkdirAll(dname, 0777)
		if *debug {
			fmt.Printf("[store] Writing to disk (%s)\n", fdata.FileName)
		}
		fd, err := os.Create(fname)
		if err != nil {
			return errors.Wrap(err)
		}
		wr = &DownloadWriter{filename: fname, fd: fd}
		//	fmt.Printf("fdata: %s\n", fdata)
	}
	return nil
}

type DownloadWriter struct {
	filename string
	fd       *os.File
}

func (dw *DownloadWriter) Close() error {
	if dw.fd == nil {
		return nil
	}
	err := dw.fd.Close()
	dw.fd = nil
	return errors.Wrap(err)
}
func (dw *DownloadWriter) Write(fdata *pb.DownloadFile) error {
	_, err := dw.fd.Write(fdata.Data)
	return errors.Wrap(err)
}
