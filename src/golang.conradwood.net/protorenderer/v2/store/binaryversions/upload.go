package binaryversions

import (
	"fmt"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/utils"
	pb "golang.yacloud.eu/apis/binaryversions"
	"io"
	"os"
	"path/filepath"
	//	"strings"
	"time"
)

const (
	PROTORENDERER_STORE_DIR_NAME = "protorenderer-store"
)

func Upload(dname string) error {
	c := pb.GetBinaryVersionsClient()
	ctx := authremote.ContextWithTimeout(time.Duration(3) * time.Minute) // long upload
	dir, err := c.MkOrGetDir(ctx, &pb.MkDirRequest{DirName: PROTORENDERER_STORE_DIR_NAME})
	if err != nil {
		return err
	}
	fmt.Printf("Dir: #%d (%s) in Realm #%d \"%s\"\n", dir.ID, dir.DirName, dir.Realm.ID, dir.Realm.Name)

	u := &Uploader{root: dname}
	//	ctx := authremote.Context()
	u.srv, err = pb.GetBinaryVersionsClient().UploadFiles(ctx)
	if err != nil {
		fmt.Printf("Upload start failed: %s\n", err)
		return err
	}
	err = u.srv.Send(&pb.UploadFile{Directory: dir})
	if err != nil {
		return err
	}
	err = utils.DirWalk(dname, u.Walker)
	//	err = u.UploadDir()
	if err != nil {
		return err
	}
	_, err = u.srv.CloseAndRecv()
	if err != nil {
		fmt.Printf("Closerecv failed\n")
		return err
	}
	return nil
}
func (u *Uploader) Walker(root string, rel string) error {
	fmt.Printf("File: %s,%s\n", root, rel)
	err := u.uploadFile(root, rel)
	return err
}
func (u *Uploader) UploadDir() error {
	dirname := u.root
	// now walk through all the files
	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		return u.uploadFile(path, info.Name())
	})
	if err != nil {
		return err
	}
	return nil
}

type Uploader struct {
	root            string
	context_created time.Time
	srv             pb.BinaryVersions_UploadFilesClient
}

func (u *Uploader) uploadFile(npath string, filename string) error {
	fname := npath + "/" + filename
	fmt.Printf("uploading \"%s\"\n", fname)
	buf := make([]byte, 8192)
	fd, err := os.Open(fmt.Sprintf("%s", fname))
	if err != nil {
		return fmt.Errorf("failed to open file: %s", err)
	}
	defer fd.Close()

	// send filename at least once
	uf := &pb.UploadFile{FileName: filename}
	err = u.srv.Send(uf)
	if err != nil {
		return err
	}

	for {
		n, err := fd.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		uf := &pb.UploadFile{FileName: filename, Data: buf[:n]}
		err = u.srv.Send(uf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Upload failed: %s\n", err)
			return err
		}
	}
	fmt.Printf("uploaded \"%s/%s\"\n", npath, filename)
	return nil
}
