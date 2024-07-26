package store

import (
	"context"
	"fmt"
	//	ost "golang.conradwood.net/apis/objectstore"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/client"
	"golang.conradwood.net/protorenderer/cmdline"
	"golang.conradwood.net/protorenderer/v2/helpers"
	"golang.conradwood.net/protorenderer/v2/store/binaryversions"
	"strings"
	"sync"
)

const (
	RETRIEVE_WORKERS = 30
)

func Retrieve(ctx context.Context, dir string, version uint64) error {

	berr := binaryversions.Download("protorenderer", dir, dir)
	err := retrieve_from_object_store(ctx, dir, version)
	if err != nil {
		return err
	}

	return berr

}

func retrieve_from_object_store(ctx context.Context, dir string, version uint64) error {
	bs, err := client.Get(ctx, cmdline.VERSIONOBJECT)
	if err != nil {
		return err
	}

	fmt.Printf("Version: %#v\n", bs)

	idxfilename := cmdline.GetPrefixObjectStore() + cmdline.INDEX_FILENAME
	b, err := client.Get(ctx, idxfilename)
	if err != nil {
		fmt.Printf("No index file \"%s\": %s\n", idxfilename, err)
		return err
	}

	wg := &sync.WaitGroup{}
	ch := make(chan *retrieve_request)
	for i := 0; i < RETRIEVE_WORKERS; i++ {
		wg.Add(1)
		go retrieve_worker(ch, wg)
	}

	ls := strings.Split(string(b), "\n")
	fmt.Printf("%d files in cache\n", len(ls))
	for _, line := range ls {
		if line == "" {
			continue
		}
		rr := &retrieve_request{dir: dir, line: line}
		ch <- rr
	}

	rr := &retrieve_request{exit: true}
	for i := 0; i < RETRIEVE_WORKERS; i++ {
		ch <- rr
	}
	wg.Wait()
	return nil

}

type retrieve_request struct {
	dir  string
	line string
	exit bool
}

func retrieve_worker(ch chan *retrieve_request, wg *sync.WaitGroup) {
	for {
		rr := <-ch
		if rr.exit {
			break
		}
		line := rr.line
		dir := rr.dir
		fmt.Printf("Reading \"%s\"\n", line)
		ctx := authremote.Context()
		filecontent, err := client.Get(ctx, cmdline.GetPrefixObjectStore()+line)
		if err != nil {
			fmt.Printf("Failed to get file %s: %s\n", line, err)
			continue
		}
		filename := dir + "/" + line
		fmt.Printf("Writing \"%s\"\n", filename)
		err = helpers.WriteFileWithDir(filename, filecontent)
		if err != nil {
			continue
		}
	}
	wg.Done()
}
