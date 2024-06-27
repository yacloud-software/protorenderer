package store

import (
	"context"
	"fmt"
	//	ost "golang.conradwood.net/apis/objectstore"
	"golang.conradwood.net/go-easyops/client"
	"golang.conradwood.net/protorenderer/cmdline"
	"golang.conradwood.net/protorenderer/v2/helpers"
	"strings"
)

func Retrieve(ctx context.Context, dir string, version uint64) error {
	err := retrieve_from_object_store(ctx, dir, version)
	if err != nil {
		return err
	}
	return nil
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

	ls := strings.Split(string(b), "\n")
	fmt.Printf("%d files in cache\n", len(ls))
	for _, line := range ls {
		if line == "" {
			continue
		}
		fmt.Printf("Reading \"%s\"\n", line)

		filecontent, err := client.Get(ctx, cmdline.GetPrefixObjectStore()+line)
		if err != nil {
			fmt.Printf("Failed to get file %s: %s\n", line, err)
			return err
		}
		filename := dir + "/" + line
		fmt.Printf("Writing \"%s\"\n", filename)
		err = helpers.WriteFileWithDir(filename, filecontent)
		if err != nil {
			return err
		}
	}

	return nil

}
