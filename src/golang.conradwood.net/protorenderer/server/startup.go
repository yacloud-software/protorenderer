package main

import (
	"fmt"
	ost "golang.conradwood.net/apis/objectstore"
	pb "golang.conradwood.net/apis/protorenderer"
	ar "golang.conradwood.net/go-easyops/authremote"
	"strings"
	"sync"
	"time"
)

const (
	STARTUP_READ_WORKERS   = 5
	STARTUP_SUBMIT_WORKERS = 5
)

var (
	read_object_chan     = make(chan *startup_read_msg)
	submit_object_chan   = make(chan *startup_submit_msg)
	startup_read_group   sync.WaitGroup
	startup_submit_group sync.WaitGroup
)

type startup_read_msg struct {
	line string
	exit bool
}
type startup_submit_msg struct {
	obj  *ost.Object
	line string
	r    *pb.AddProtoRequest
	exit bool
}

func startup() {

	if !*read_object_store {
		start_update = true
		return
	}

	for i := 0; i < STARTUP_READ_WORKERS; i++ {
		startup_read_group.Add(1)
		go startup_read_worker()
	}
	for i := 0; i < STARTUP_SUBMIT_WORKERS; i++ {
		startup_submit_group.Add(1)
		go startup_submit_worker()
	}

	time.Sleep(time.Duration(2) * time.Second)
	idxfilename := *prefix_object_store + INDEX_FILENAME
	b, err := osclient().Get(ar.Context(), &ost.GetRequest{ID: idxfilename})
	if err != nil {
		fmt.Printf("No index file \"%s\": %s\n", idxfilename, err)
		return
	}

	ls := strings.Split(string(b.Content), "\n")
	fmt.Printf("%d files in cache\n", len(ls))
	for _, line := range ls {
		read_object_chan <- &startup_read_msg{line: line}
	}

	for i := 0; i < STARTUP_READ_WORKERS; i++ {
		read_object_chan <- &startup_read_msg{exit: true}
	}
	for i := 0; i < STARTUP_SUBMIT_WORKERS; i++ {
		submit_object_chan <- &startup_submit_msg{exit: true}

	}
	startup_read_group.Wait()
	startup_submit_group.Wait()
	start_update = true
	ui := &updateinfo{}
	updateChan <- ui

}

func startup_read_worker() {
	for {
		o := <-read_object_chan
		if o.exit {
			break
		}
		b, err := osclient().Get(ar.Context(), &ost.GetRequest{ID: *prefix_object_store + o.line})
		if err != nil {
			fmt.Printf("Failed to get file %s: %s\n", o.line, err)
			continue
		}
		r := &pb.AddProtoRequest{Name: o.line, Content: string(b.Content)}
		if len(r.Name) == 0 || len(r.Content) == 0 {
			continue
		}

		submit_object_chan <- &startup_submit_msg{obj: b, line: o.line, r: r}

	}
	startup_read_group.Done()
}
func startup_submit_worker() {
	e := new(protoRenderer)
	for {
		o := <-submit_object_chan
		if o.exit {
			break
		}
		_, err := e.UpdateProto(ar.Context(), o.r)
		if err != nil {
			fmt.Printf("Failed to add proto %s: %s\n", o.line, err)
			continue
		}
		fmt.Printf("Added %s\n", o.line)

	}
	startup_submit_group.Done()
}



































































