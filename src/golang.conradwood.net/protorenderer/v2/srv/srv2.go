package srv

import (
	"context"
	"fmt"
	cma "golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/server"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/cmdline"
	"golang.conradwood.net/protorenderer/v2/common"
	"golang.conradwood.net/protorenderer/v2/compilers/java"
	"golang.conradwood.net/protorenderer/v2/helpers"
	"golang.conradwood.net/protorenderer/v2/interfaces"
	"golang.conradwood.net/protorenderer/v2/meta_compiler"
	ms "golang.conradwood.net/protorenderer/v2/meta_compiler/server"
	"golang.conradwood.net/protorenderer/v2/store"
	"golang.conradwood.net/protorenderer/v2/versioninfo"
	"google.golang.org/grpc"
	"os"
	"time"
)

var (
	CompileEnv *StandardCompilerEnvironment
)

func Start() {
	var err error

	server.SetHealth(cma.Health_STARTING)

	CompileEnv = &StandardCompilerEnvironment{workdir: "/tmp/pr/v2"}
	utils.RecreateSafely(CompileEnv.workdir + "/store")
	//scr := &StandardCompileResult{}
	mkdir(CompileEnv.AllKnownProtosDir())

	fmt.Printf("Creating workdir...\n")
	err = createWorkDir()
	utils.Bail("failed to create workdir", err)
	utils.RecreateSafely(CompileEnv.CompilerOutDir())
	mkdir(CompileEnv.NewProtosDir())

	ctx := authremote.ContextWithTimeout(time.Duration(90) * time.Second)
	err = store.Retrieve(ctx, CompileEnv.StoreDir(), 0) // 0 == latest
	utils.Bail("failed to retrieve latest version", err)

	//	os.Exit(0)

	if false {
		java.Start(CompileEnv, CompileEnv.CompilerOutDir())
	}

	sd := server.NewServerDef()
	sd.SetPort(cmdline.GetRPCPort())
	e := new(protoRenderer)
	sd.SetRegister(server.Register(
		func(server *grpc.Server) error {
			pb.RegisterProtoRenderer2Server(server, e)
			return nil
		},
	))
	sd.SetOnStartupCallback(server_started)
	err = server.ServerStartup(sd)
	utils.Bail("Unable to start server", err)
	os.Exit(0)

}

func server_started() {
	err := build_version_info()
	utils.Bail("failed to build versioninfo", err)
	server.SetHealth(cma.Health_READY)
}

// get all meta info for all protos we have in our store.
// rebuild as/if/when necessary all meta information
func build_version_info() error {
	vi := versioninfo.New()
	ce := CompileEnv
	pfs, err := helpers.FindProtoFiles(ce.AllKnownProtosDir())
	if err != nil {
		return err
	}
	store_required := false
	ctx := context.Background()
	for _, pf := range pfs {
		mi, err := meta_compiler.GetMetaInfo(ctx, ce, pf.GetFilename())
		if err != nil {
			fmt.Printf("Failed to get meta info for %s: %s - compile triggered\n", pf.GetFilename(), utils.ErrorString(err))
			scr := &common.StandardCompileResult{}
			rpc_port := cmdline.GetRPCPort()
			mc := meta_compiler.New()
			ctx = authremote.ContextWithTimeout(time.Duration(90) * time.Second)
			meta_needed := []interfaces.ProtoFile{pf}
			err = mc.Compile(ctx, rpc_port, ce, meta_needed, ce.CompilerOutDir()+"/meta", scr)
			if err != nil {
				fmt.Printf("Failed to compile meta info for %s: %s\n", pf.GetFilename(), utils.ErrorString(err))
				continue
			}
			mi, err = meta_compiler.GetMetaInfo(ctx, ce, pf.GetFilename())
			if err != nil {
				fmt.Printf("Failed to get meta info for %s: %s - compile unsuccessful\n", pf.GetFilename(), utils.ErrorString(err))
				continue
			}
			store_required = true
		}
		vi.GetOrAddFile(pf.GetFilename(), mi)
	}
	utils.Bail("failed to merge compiler env", helpers.MergeCompilerEnvironment(ce, true))
	if store_required {
		err = store.Store(ctx, CompileEnv.StoreDir()) // 0 == latest
		utils.Bail("failed to store", err)
	}
	return nil
}

type protoRenderer struct {
}

func (pr *protoRenderer) InternalMetaSubmit(ctx context.Context, req *pb.ProtocRequest) (*cma.Void, error) {
	//fmt.Printf("[server] received request from protoc-gen-meta\n")
	return ms.InternalMetaSubmit(ctx, req)
}
func mkdir(dir string) {
	err := linux.CreateIfNotExists(dir, 0777)
	utils.Bail("failed to create dir", err)
	fmt.Printf("Created dir %s\n", dir)
}
