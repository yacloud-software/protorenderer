package main

import (
	"context"
	"flag"
	"fmt"
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/protorenderer"
	"path/filepath"
	//	ch "golang.conradwood.net/go-easyops/http"
	//	au "golang.conradwood.net/go-easyops/auth"
	ar "golang.conradwood.net/go-easyops/authremote"
	//	"golang.conradwood.net/go-easyops/tokens"
	"golang.conradwood.net/go-easyops/utils"
	//	"golang.conradwood.net/protorenderer/renderer"
	//	"io/ioutil"
	//	"net/http"
	"os"
	//	"path/filepath"
)

var (
	http_port   = flag.Int("http_port", 8081, "http port to listen on")
	protoClient pb.ProtoRendererServiceClient
	view        = flag.Bool("view", false, "view current proto docs")
	files       = flag.Bool("files", false, "download .proto, .pb.go, .class, .py, and nanopb files")
	sources     = flag.Bool("sources", false, "download .proto files")
	version     = flag.Bool("version", false, "get version")
	delete      = flag.Bool("delete", false, "delete files listed on command line")
	compile     = flag.Bool("compile", false, "compile files on the fly (but do not add to repository")
	outdir      = flag.String("outdir", "", "directory where to place compiled protos files")
	listflag    = flag.Bool("list", false, "list source files currently in repository")
	repoid      = flag.Uint64("repository_id", 0, "repository id of the proto being submitted. if not set, will look at deploy.yaml")
	pkgid       = flag.Uint64("package_id", 0, "package id to operate on")
	packages    = flag.Bool("packages", false, "get packages")
	get_zip     = flag.String("zip", "", "if non-nil get zip file for packagename")
)

func main() {
	flag.Parse()
	protoClient = pb.GetProtoRendererServiceClient()
	if *get_zip != "" {
		GetZip(*get_zip)
		os.Exit(0)
	}
	if *packages {
		Packages()
		os.Exit(0)
	}
	if *sources {
		Sources()
		os.Exit(0)
	}
	if *listflag {
		List()
		os.Exit(0)
	}
	if *compile {
		Compile()
		os.Exit(0)
	}
	if *delete {
		Delete()
		os.Exit(0)
	}
	if *files {
		get_files()
		os.Exit(0)
	}

	if *view {
		View()
	}
	files := flag.Args()
	if *version {
		doversion()
		os.Exit(0)
	}
	if len(files) > 0 {
		add(files)
		os.Exit(0)
	}
	fmt.Printf("Done\n")
}
func getrepoid() uint64 {
	return *repoid
}
func add(files []string) {
	ctx := getContext()

	for _, f := range files {
		fmt.Printf("Adding %s\n", f)
		name := f
		b, err := utils.ReadFile(f)
		utils.Bail("failed to read file", err)
		pf := &pb.AddProtoRequest{
			Name:         name,
			Content:      string(b),
			RepositoryID: getrepoid(),
		}
		ctx = getContext()
		v, err := protoClient.UpdateProto(ctx, pf)
		utils.Bail("failed to add proto", err)
		fmt.Printf("File is in Version: %d\n", v.Version)
	}
}

func doversion() {
	ctx := getContext()
	v, err := protoClient.GetVersion(ctx, &common.Void{})
	utils.Bail("failed to get version", err)
	fmt.Printf("Current Version: %d\n", v.Version)
	fmt.Printf("Compiling      : %v\n", v.Compiling)
	fmt.Printf("Next Version   : %d\n", v.NextVersion)
	fmt.Printf("ProtoVersion   : %d\n", v.ProtoVersion)
}
func Compile() {
	if len(flag.Args()) == 0 {
		fmt.Printf("missing filename(s)\n")
		os.Exit(10)
	}
	for _, fname := range flag.Args() {
		CompileFile(fname)
	}
}
func CompileFile(fname string) {
	fs, err := utils.ReadFile(fname)
	utils.Bail("failed to read file", err)
	ctx := getContext()
	req := &pb.CompileRequest{
		AddProtoRequest: &pb.AddProtoRequest{
			Name:    fname,
			Content: string(fs),
		},
	}
	res, err := protoClient.CompileFile(ctx, req)
	utils.Bail("failed to compile", err)
	if res.CompileError != "" {
		fmt.Printf("Failed to compile: %s\n", res.CompileError)
		os.Exit(10)
	}
	fmt.Printf("Compiled file \"%s\" (%d files returned)\n", res.SourceFilename, len(res.Files))
	for _, cf := range res.Files {
		if *outdir == "" {
			fmt.Printf("  %s [%d bytes]\n", cf.Filename, len(cf.Content))
			continue
		}
		// save file
		fname := fmt.Sprintf("%s/%s", *outdir, cf.Filename)
		os.MkdirAll(filepath.Dir(fname), 0777)
		err = utils.WriteFile(fname, cf.Content)
		utils.Bail("failed to store file", err)
	}
}

func getContext() context.Context {
	var ctx context.Context
	//	ctx = tokens.ContextWithToken()
	//	ctx = au.Context() // use env var
	ctx = ar.Context()
	return ctx
}

func List() {
	ctx := getContext()
	fl, err := protoClient.ListSourceFiles(ctx, &common.Void{})
	utils.Bail("failed to get list of files", err)
	for _, f := range fl.Files {
		fmt.Println(f)
	}
}
func Packages() {
	ctx := getContext()
	fl, err := protoClient.GetPackages(ctx, &common.Void{})
	utils.Bail("failed to get packages", err)
	for _, f := range fl.Packages {
		fmt.Printf("Package #%4s %s\n", f.ID, f.Name)
	}

}
