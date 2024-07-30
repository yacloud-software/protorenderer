package store

import (
	"context"
	"fmt"
	"golang.conradwood.net/protorenderer/v2/store/binaryversions"
	"time"
)

var (
	uploadTriggerChan = make(chan *uploadrequest, 10)
)

type uploadrequest struct {
	dir string
}

func init() {
	go upload_worker()
}
func Store(ctx context.Context, dir string) error {
	err := binaryversions.Upload(dir)
	return err
}

func TriggerUpload(dir string) {
	select {
	case uploadTriggerChan <- &uploadrequest{dir: dir}:
		//
	case <-time.After(time.Duration(1) * time.Second):
		// nope

	}
}
func upload_worker() {
	for {
		ur := <-uploadTriggerChan
		fmt.Printf("[uploadworker] uploading \"%s\"\n", ur.dir)
		err := binaryversions.Upload(ur.dir)
		if err != nil {
			fmt.Printf("[uploadworker] Failed to upload \"%s\": %s\n", ur.dir, err)
		} else {
			fmt.Printf("[uploadworker] uploaded \"%s\"\n", ur.dir)
		}
	}
}
