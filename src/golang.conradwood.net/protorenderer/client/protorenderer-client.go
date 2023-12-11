package main

import (
	"context"
	"flag"
	"fmt"
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/protorenderer"
	"path/filepath"
	"strings"
	//	ch "golang.conradwood.net/go-easyops/http"
	"golang.conradwood.net/go-easyops/cmdline"
	//	au "golang.conradwood.net/go-easyops/auth"
	ar "golang.conradwood.net/go-easyops/authremote"
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
	files       = flag.Bool("files", false, "download .proto, .pb.go, .class, .py, and nanopb files (if extra arguments are given on the commandline, limit downloads to those packages matching remaining non-option arguments")
	sources     = flag.Bool("sources", false, "download .proto files (remaining args, if present, filter the output)")
	compilers   = flag.String("compilers", "", "specify compilers. if empty (Default): autodetect")
	version     = flag.Bool("version", false, "get version")
	delete      = flag.Bool("delete", false, "delete files listed on command line")
	compile     = flag.Bool("compile", false, "compile files on the fly (but do not add to repository")
	show_failed = flag.Bool("failed", false, "show failed files")
	outdir      = flag.String("outdir", "", "directory where to place compiled protos files")
	listflag    = flag.Bool("list", false, "list source files currently in repository")
	repoid      = flag.Uint64("repository_id", 0, "repository id of the proto being submitted. if not set, will look at deploy.yaml")
	pkgid       = flag.Uint64("package_id", 0, "package id to operate on")
	packages    = flag.Bool("packages", false, "get packages")
	get_zip     = flag.String("zip", "", "if non-nil get zip file for packagename")
	find_pkg    = flag.String("find_package", "", "if non-nil find package by name")
	debug       = flag.Bool("debug", false, "debug mode")
)

func main() {
	flag.Parse()
	protoClient = pb.GetProtoRendererServiceClient()
	if *show_failed {
		utils.Bail("failed to show failed files", showFailed())
		os.Exit(0)
	}
	if *find_pkg != "" {
		utils.Bail("failed to find package", FindPkg(*find_pkg))
		os.Exit(0)
	}
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
		os.Exit(0)
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
	fmt.Printf("Registry: %s\n", cmdline.GetRegistryAddress())
	for _, fname := range flag.Args() {
		CompileFile(fname)
	}
}
func CompileFile(fname string) {
	fs, err := utils.ReadFile(fname)
	utils.Bail("failed to read file", err)
	ctx := getContext()
	req := &pb.CompileRequest{
		Compilers: GetSpecifiedCompilers(),
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
			fmt.Printf("  %s [%d bytes] (not saved, no outdir)\n", cf.Filename, len(cf.Content))
			continue
		}
		// save file
		cfname := strings.TrimPrefix(cf.Filename, "protos/")
		fname := fmt.Sprintf("%s/%s", *outdir, cfname)
		os.MkdirAll(filepath.Dir(fname), 0777)
		err = utils.WriteFile(fname, cf.Content)
		utils.Bail("failed to store file", err)
		fmt.Printf("Saved to \"%s\"\n", fname)
	}
}

func getContext() context.Context {
	var ctx context.Context
	//	ctx = authremote.Context()
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
	t := utils.Table{}
	t.AddHeaders("ID", "Prefix", "Name", "RepositoryID")
	for _, f := range fl.Packages {
		t.AddString(f.ID)
		t.AddString(f.Prefix)
		t.AddString(f.Name)
		t.AddUint64(f.RepositoryID)
		t.NewRow()
	}
	fmt.Println(t.ToPrettyString())
}

// return the compilers, either from commandline or autodetect
func GetSpecifiedCompilers() []pb.CompilerType {
	if *compilers == "" {
		return []pb.CompilerType{
			pb.CompilerType_GOLANG,
		}
	}
	var res []pb.CompilerType
	sx := strings.Split(*compilers, ",")
	for _, comp_name := range sx {
		comp_name = strings.ToLower(strings.Trim(comp_name, " "))
		found := false
		for cdefnum, cdef := range pb.CompilerType_name {
			cdef = strings.ToLower(cdef)
			if cdef == comp_name {
				found = true
				res = append(res, pb.CompilerType(cdefnum))
				break
			}
		}
		if !found {
			fmt.Printf("\"%s\" is not a supported compiler.\n", comp_name)
			for _, cdef := range pb.CompilerType_name {
				fmt.Printf("\"%s\"\n", cdef)
			}
			os.Exit(10)
		}

	}
	return res
}

func FindPkg(pkgname string) error {
	fmt.Printf("Finding package \"%s\"\n", pkgname)
	ctx := ar.Context()
	pkg := &pb.PackageName{
		PackageName: pkgname,
	}
	res, err := protoClient.GetPackageByName(ctx, pkg)
	if err != nil {
		return err
	}
	fmt.Printf("Name: %s\n", res.Name)
	fmt.Printf("Prefix: %s\n", res.Prefix)
	svc := "none"
	cmt := "n/a"
	if len(res.Services) > 0 {
		svc = res.Services[0].Name
		cmt = res.Services[0].Comment
	}
	fmt.Printf("Service: %s\n", svc)
	fmt.Printf("Comment: %s\n", cmt)
	return nil
}

func showFailed() error {
	ctx := getContext()
	res, err := protoClient.GetFailedFiles(ctx, &common.Void{})
	if err != nil {
		return err
	}
	t := &utils.Table{}
	t.AddHeaders("filename", "compiler", "message")
	for _, f := range res.Files {
		t.AddString(f.Filename)
		t.AddString(f.Compiler)
		t.AddString(f.Message)
		t.NewRow()
	}
	fmt.Println(t.ToPrettyString())
	return nil
}








































