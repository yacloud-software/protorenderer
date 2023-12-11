package main

import (
	"flag"
	"fmt"
	ost "golang.conradwood.net/apis/objectstore"
	pb "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/client"
	"golang.conradwood.net/go-easyops/utils"
	//	"sort"
	//	"strings"
)

const (
	VERSIONOBJECT  = "protorenderer_version"
	INDEX_FILENAME = "protorenderer_index_file"
)

var (
	osc                 ost.ObjectStoreClient
	prefix_object_store = flag.String("prefix_object_store", "protorenderer-tmp", "a prefix to be used for objectstore put/get")
)

func Delete() {
	ctx := getContext()
	s := ""
	deli := ""
	for _, file := range flag.Args() {
		s = s + deli + file
		deli = ", "
		_, err := protoClient.DeleteFile(ctx, &pb.DeleteRequest{Name: file})
		utils.Bail("failed to delete file", err)
	}

	/*
		b, err := osclient().Get(ctx, &ost.GetRequest{ID: *prefix_object_store + INDEX_FILENAME})
		utils.Bail("failed to get index", err)
		ls := strings.Split(string(b.Content), "\n")
		fmt.Printf("%d files in cache\n", len(ls))
		ns := ""
		sort.Slice(ls, func(i, j int) bool {
			return ls[i] < ls[j]
		})
		for _, line := range ls {
			if line == "" {
				continue
			}
			if ListedFile(line) {
				fmt.Printf("[DELETED] \"%s\"\n", line)
				continue
			}
			fmt.Printf("[KEPT   ] %s\n", line)
			ns = ns + line + "\n"
		}
		pir := &ost.PutWithIDRequest{ID: *prefix_object_store + INDEX_FILENAME, Content: []byte(ns)}
		_, err = osclient().PutWithID(ctx, pir)
		utils.Bail("failed to update index", err)
	*/
	fmt.Printf("Deleted %s, done\n", s)
}
func ListedFile(filename string) bool {
	for _, f := range flag.Args() {
		if f == filename {
			return true
		}
	}
	return false
}
func osclient() ost.ObjectStoreClient {
	if osc != nil {
		return osc
	}
	osc = ost.NewObjectStoreClient(client.Connect("objectstore.ObjectStore"))
	return nil
}


























































