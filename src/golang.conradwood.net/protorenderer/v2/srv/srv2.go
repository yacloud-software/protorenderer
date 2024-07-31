package srv

import (
	"context"
	"flag"
	"fmt"
	cma "golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/server"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/cmdline"
	"golang.conradwood.net/protorenderer/v2/compilers/java"
	"golang.conradwood.net/protorenderer/v2/interfaces"
	ms "golang.conradwood.net/protorenderer/v2/meta_compiler/server"
	"golang.conradwood.net/protorenderer/v2/metadata"
	"golang.conradwood.net/protorenderer/v2/store"
	"golang.conradwood.net/protorenderer/v2/versioninfo"
	"google.golang.org/grpc"
	"os"
	"time"
)

var (
	recompile_on_startup = flag.Bool("recompile_on_startup", false, "if true recompile store on startup")
	CompileEnv           *StandardCompilerEnvironment
	currentVersionInfo   *versioninfo.VersionInfo
)

func Start() {
	var err error
	fmt.Printf("Starting protorenderer-server (v2)\n")
	server.SetHealth(cma.Health_STARTING)

	hd, err := utils.HomeDir()
	utils.Bail("failed to get homedir", err)
	CompileEnv = &StandardCompilerEnvironment{workdir: hd + "/tmp/pr/v2"}
	metadata.MetaCache.SetEnv(CompileEnv)

	utils.RecreateSafely(CompileEnv.workdir + "/store")
	//scr := &StandardCompileResult{}
	mkdir(CompileEnv.AllKnownProtosDir())

	fmt.Printf("Creating workdir...\n")
	err = createWorkDir()
	utils.Bail("failed to create workdir", err)
	utils.RecreateSafely(CompileEnv.CompilerOutDir())
	mkdir(CompileEnv.NewProtosDir())

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
	server.SetHealth(cma.Health_STARTING)
	ctx := authremote.ContextWithTimeout(time.Duration(180) * time.Second)
	err := store.Retrieve(ctx, CompileEnv.StoreDir(), 0) // 0 == latest
	utils.Bail("failed to retrieve latest version", err)

	//	os.Exit(0)

	if false {
		java.Start(CompileEnv, CompileEnv.CompilerOutDir())
	}
	reloadVersionInfo(CompileEnv)
	b := *recompile_on_startup
	if !utils.FileExists(CompileEnv.StoreDir() + "/versioninfo.pbbin") {
		b = true
	}
	if b {
		err := RecompileStore(CompileEnv)
		utils.Bail("failed to recompile store", err)
	}
	server.SetHealth(cma.Health_READY)
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

func (pr *protoRenderer) GetVersionInfo(ctx context.Context, req *cma.Void) (*pb.VersionInfo, error) {
	return currentVersionInfo.ToProto(), nil
}

// store the versioninfo.pbbin
func saveVersionInfo() error {
	vi := currentVersionInfo
	ce := CompileEnv
	fname := ce.StoreDir() + "/versioninfo.pbbin"
	fmt.Printf("[recompilestore] Storing\n")
	b, err := utils.MarshalBytes(vi.ToProto())
	if err != nil {
		return err
	}
	err = utils.WriteFile(fname, b)
	if err != nil {
		return err
	}
	return nil
}

func reloadVersionInfo(ce interfaces.CompilerEnvironment) {
	vi := versioninfo.New()

	fname := ce.StoreDir() + "/versioninfo.pbbin"
	b, err := utils.ReadFile(fname)
	if err != nil {
		// cannot load, then make sure we write later
		vi.SetDirty() // must be stored later
		fmt.Printf("[recompilestore] Failed to load versioninfo: %s\n", err)
	} else {
		pbvi := &pb.VersionInfo{}
		err = utils.UnmarshalBytes(b, pbvi)
		if err != nil {
			fmt.Printf("failed to marshal versioninfo: %s!!!\n", err)
			panic("failed to marshal versioninfo")
		}
		vi = versioninfo.NewFromProto(pbvi)
		fmt.Printf("[recompilestore] Loaded versioninfo from %s\n", fname)
	}
	currentVersionInfo = vi

}
