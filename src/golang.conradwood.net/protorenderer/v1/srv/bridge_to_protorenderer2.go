package srv

import (
	"flag"
	"fmt"
	pb "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/utils"
	"io"
	"sort"
	"sync"
	"time"
)

var (
	submit_to_pr2        = flag.Bool("submit_to_protorenderer2", true, "if true, submit to protorenderer2 as well")
	pr2_submit_chan      = make(chan *pb.AddProtoRequest, 5000)
	bridge_failures      = make(map[string]*pb.FailedBridgeFile)
	bridge_failures_lock sync.Mutex
)

func init() {
	go protorenderer2_submit_worker()
}

func submit_to_protorenderer2(req *pb.AddProtoRequest) {
	pr2_submit_chan <- req
}
func GetBridgeFailures() []*pb.FailedBridgeFile {
	bridge_failures_lock.Lock()
	defer bridge_failures_lock.Unlock()
	var res []*pb.FailedBridgeFile
	for _, v := range bridge_failures {
		res = append(res, v)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Filename < res[j].Filename
	})
	return res
}

func protorenderer2_submit_worker() {
	for {
		req := <-pr2_submit_chan
		err := submit_to_protorenderer2_werr(req)
		if err != nil {
			fmt.Printf("Error submitting to protorenderer2: %s\n", errors.ErrorStringWithStackTrace(err))
		}
		setFileError(req.Name, err)
	}
}
func submit_to_protorenderer2_werr(req *pb.AddProtoRequest) error {
	fmt.Printf("Submitting file \"%s\" to protorenderer2\n", req.Name)
	ctx := authremote.ContextWithTimeout(time.Duration(10) * time.Minute)
	repoid := uint32(req.RepositoryID)

	srv, err := protorenderer2.GetProtoRenderer2Client().Submit(ctx)
	if err != nil {
		return err
	}
	so := &protorenderer2.SubmitOption{Save: true}
	err = srv.Send(&protorenderer2.FileTransfer{SubmitOption: so})
	if err != nil {
		return err
	}
	bs := utils.NewByteStreamSender(func(key, filename string) error {
		// start new file
		err := srv.Send(&protorenderer2.FileTransfer{Filename: filename, RepositoryID: repoid})
		return err
	},
		// send contents
		func(b []byte) error {
			err := srv.Send(&protorenderer2.FileTransfer{Data: b})
			return err
		},
	)

	err = bs.SendBytes(req.Name, req.Name, []byte(req.Content))
	if err != nil {
		return err
	}
	err = srv.Send(&protorenderer2.FileTransfer{TransferComplete: true}) // switching to recv mode now
	if err != nil {
		return err
	}
	for {
		_, err := srv.Recv() // receive, but discard result
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}

	return nil
}

func setFileError(filename string, err error) {
	bridge_failures_lock.Lock()
	defer bridge_failures_lock.Unlock()
	if err == nil {
		delete(bridge_failures, filename)
	} else {
		bridge_failures[filename] = &pb.FailedBridgeFile{
			Occured:      uint32(time.Now().Unix()),
			Filename:     filename,
			ErrorMessage: utils.ErrorString(err),
		}
	}
}
