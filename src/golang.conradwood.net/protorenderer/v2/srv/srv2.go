package srv

import (
	"context"
	"fmt"
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/server"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/cmdline"
	ms "golang.conradwood.net/protorenderer/v2/meta-compiler/server"
	"google.golang.org/grpc"
	"os"
)

var (
	CompileEnv *StandardCompilerEnvironment
)

func Start() {
	var err error
	server.SetHealth(common.Health_STARTING)

	//scr := &StandardCompileResult{}
	CompileEnv = &StandardCompilerEnvironment{workdir: "/tmp/pr/v2", knownprotosdir: "proto_files/protos", newprotosdir: "new_protos/protos"}
	mkdir(CompileEnv.WorkDir() + "/" + CompileEnv.AllKnownProtosDir())
	mkdir(CompileEnv.WorkDir() + "/" + CompileEnv.NewProtosDir())
	server.SetHealth(common.Health_READY)
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

func (pr *protoRenderer) SubmitSource(ctx context.Context, req *pb.ProtocRequest) (*common.Void, error) {
	return ms.SubmitSource(ctx, req)
}
func mkdir(dir string) {
	err := linux.CreateIfNotExists(dir, 0777)
	utils.Bail("failed to create dir", err)
	fmt.Printf("Created dir %s\n", dir)
}
