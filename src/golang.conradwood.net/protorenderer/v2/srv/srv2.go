package srv

import (
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/server"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/cmdline"
	"google.golang.org/grpc"
	"os"
)

func Start() {
	var err error
	server.SetHealth(common.Health_STARTING)
	sd := server.NewServerDef()
	sd.SetPort(cmdline.GetRPCPort())
	e := new(protoRenderer)
	sd.SetRegister(server.Register(
		func(server *grpc.Server) error {
			pb.RegisterProtoRenderer2Server(server, e)
			return nil
		},
	))

	err = server.ServerStartup(sd)
	utils.Bail("Unable to start server", err)
	os.Exit(0)

}

type protoRenderer struct {
}
