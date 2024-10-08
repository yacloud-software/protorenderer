package srv

import (
	"flag"
	"fmt"
	"io"
	"sort"
	"sync"
	"time"

	pb "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/utils"
)

var (
	submit_to_pr2        = flag.Bool("submit_to_protorenderer2", false, "if true, submit to protorenderer2 as well")
	pr2_submit_chan      = make(chan []*pb.AddProtoRequest, 5000)
	bridge_failures      = make(map[string]*pb.FailedBridgeFile)
	bridge_failures_lock sync.Mutex
)

func init() {
	go protorenderer2_submit_worker()
}

func submit_to_protorenderer2(req []*pb.AddProtoRequest) {
	if len(req) == 0 {
		return
	}
	if *submit_to_pr2 {
		pr2_submit_chan <- req
	}
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
	}
}
func submit_to_protorenderer2_werr(reqs []*pb.AddProtoRequest) error {
	ctx := authremote.ContextWithTimeout(time.Duration(10) * time.Minute)
	//	repoid := uint32(req.RepositoryID)

	srv, err := protorenderer2.GetProtoRenderer2Client().Submit(ctx)
	if err != nil {
		for _, req := range reqs {
			setFileError(req, err)
		}
		return err
	}
	so := &protorenderer2.SubmitOption{Save: true}
	err = srv.Send(&protorenderer2.FileTransfer{SubmitOption: so})
	if err != nil {
		for _, req := range reqs {
			setFileError(req, err)
		}
		return err
	}
	bs := utils.NewByteStreamSender(func(key, filename string) error {
		repoid := uint32(0)
		for _, req := range reqs {
			if req.Name == filename {
				repoid = uint32(req.RepositoryID)
			}
		}
		if repoid == 0 {
			fmt.Printf("[bridge] - Warning, got no repoid for file %s\n", filename)
		}
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
	for _, apr := range reqs {
		fmt.Printf("Submitting file \"%s\" to protorenderer2\n", apr.Name)

		err = bs.SendBytes(apr.Name, apr.Name, []byte(apr.Content))
		if err != nil {
			return err
		}
	}
	err = srv.Send(&protorenderer2.FileTransfer{TransferComplete: true}) // switching to recv mode now
	if err != nil {
		return err
	}
	for {
		recv, err := srv.Recv() // receive, but discard content of received files
		if recv != nil {
			res := recv.Result
			if res != nil {
				for _, req := range reqs {
					if req.Name == res.Filename {
						setResult(req, res)
					}
				}

			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}

	return nil
}
func setResult(req *pb.AddProtoRequest, result *protorenderer2.FileResult) {
	allgood := true
	var failed_result *protorenderer2.CompileResult
	for _, cr := range result.CompileResults {
		if !cr.Success {
			allgood = false
			failed_result = cr
			break
		}
	}
	if allgood {
		setFileError(req, nil)
		return
	}
	ferr := fmt.Errorf("%s", failed_result.ErrorMessage)
	setFileError(req, ferr)

}
func setFileError(req *pb.AddProtoRequest, err error) {
	bridge_failures_lock.Lock()
	defer bridge_failures_lock.Unlock()
	filename := req.Name
	if err == nil {
		delete(bridge_failures, filename)
		fmt.Printf("File: %s: OK\n", filename)
	} else {
		fmt.Printf("File: %s: FAILED (%s)\n", filename, err)
		bridge_failures[filename] = &pb.FailedBridgeFile{
			Occured:      uint32(time.Now().Unix()),
			Filename:     filename,
			ErrorMessage: utils.ErrorString(err),
			RepositoryID: req.RepositoryID,
		}
	}
}
