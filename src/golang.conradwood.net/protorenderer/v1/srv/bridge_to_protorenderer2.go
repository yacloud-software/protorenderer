package srv

import (
	"flag"
	"fmt"
	pb "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/utils"
	"time"
)

var (
	submit_to_pr2   = flag.Bool("submit_to_protorenderer2", true, "if true, submit to protorenderer2 as well")
	pr2_submit_chan = make(chan *pb.AddProtoRequest, 500)
)

func init() {
	go protorenderer2_submit_worker()
}

func submit_to_protorenderer2(req *pb.AddProtoRequest) {
	pr2_submit_chan <- req
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
	return nil
}
