package srv

import (
	"context"
	"fmt"
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/server"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/cmdline"
	"golang.conradwood.net/protorenderer/v2/compilers/java"
	ms "golang.conradwood.net/protorenderer/v2/meta_compiler/server"
	"golang.conradwood.net/protorenderer/v2/store"
	"google.golang.org/grpc"
	"os"
	"time"
)

var (
	CompileEnv *StandardCompilerEnvironment
)

func Start() {
	var err error

	server.SetHealth(common.Health_STARTING)

	CompileEnv = &StandardCompilerEnvironment{workdir: "/tmp/pr/v2"}
	//scr := &StandardCompileResult{}
	mkdir(CompileEnv.AllKnownProtosDir())

	fmt.Printf("Creating workdir...\n")
	err = createWorkDir()
	utils.Bail("failed to create workdir", err)
	utils.RecreateSafely(CompileEnv.CompilerOutDir())

	ctx := authremote.ContextWithTimeout(time.Duration(90) * time.Second)
	err = store.Retrieve(ctx, CompileEnv.AllKnownProtosDir(), 0) // 0 == latest
	utils.Bail("failed to retrieve latest version", err)

	//	os.Exit(0)

	if false {
		java.Start(CompileEnv, CompileEnv.CompilerOutDir())
	}
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

func (pr *protoRenderer) InternalMetaSubmit(ctx context.Context, req *pb.ProtocRequest) (*common.Void, error) {
	//fmt.Printf("[server] received request from protoc-gen-meta\n")
	return ms.InternalMetaSubmit(ctx, req)
}
func mkdir(dir string) {
	err := linux.CreateIfNotExists(dir, 0777)
	utils.Bail("failed to create dir", err)
	fmt.Printf("Created dir %s\n", dir)
}
