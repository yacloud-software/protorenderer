package meta

/*
* compile meta information for documentation
 */
import (
	"context"
	"flag"
	"fmt"
	pr "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/filelayouter"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	debug             = flag.Bool("debug_meta", false, "debug mode of meta compiler")
	accept_all_tokens = flag.Bool("accept_all_tokens_by_protoc", false, "do not use in production. disables verification between protoc and this service")
	//mostRecent        *Result
)

type MetaCompiler struct {
	fl                  *filelayouter.FileLayouter
	cache               *MetaCache
	lastVersion         int
	verifytoken         string
	result              *Result
	processingProtofile *filelayouter.TrackedChange
}

// currently this thing holds it all in cache
// eventually it probably needs to database it
type MetaCache struct {
	packages []*pr.Package
}

func NewMetaCompiler(f *filelayouter.FileLayouter) *MetaCompiler {
	res := &MetaCompiler{fl: f, cache: &MetaCache{}}
	return res
}
func (m *MetaCompiler) Error() error {
	return nil
}

// since we are submitting to ourselves, we need to know our local tcp/ip port
func (m *MetaCompiler) Compile(myport int) error {
	port := myport

	fmt.Printf("Meta compiling...\n")
	tcs := m.fl.ChangedProtos(m.lastVersion)
	srcdir := m.fl.SrcDir()
	fmt.Printf("Srcdir: \"%s\"\n", srcdir)
	m.verifytoken = utils.RandomString(64)
	pcfname := findCompiler()
	fmt.Printf("Using: %s\n", pcfname)

	ctx := authremote.Context()
	if ctx == nil {
		panic("No context!")
	}
	l := linux.New()
	l.SetMaxRuntime(time.Duration(600) * time.Second)
	sctx, err := auth.SerialiseContextToString(ctx)
	if err != nil {
		fmt.Printf("Meta-Compiler: Unable to serialise context: %s\n", err)
		return err
	}
	cmd := []string{
		cmdline.GetYACloudDir() + "/ctools/dev/go/current/protoc/protoc",
		fmt.Sprintf("--plugin=protoc-gen-meta=%s", pcfname),
		"--meta_out=/tmp", // has no output
		fmt.Sprintf("--meta_opt=%s,%s,%d,%s", m.verifytoken, sctx, port, cmdline.GetClientRegistryAddress()),
	}
	for _, tc := range tcs {
		//		pf := tc.Protofile()
		fname := strings.TrimPrefix(tc.Filename(), srcdir)
		fname = strings.TrimPrefix(fname, "/")
		fmt.Printf(" meta: %s\n", fname)
		if tc.Filename() == "" {
			continue
		}
		cmdfl := append(cmd, fname)
		m.processingProtofile = tc
		out, err := l.SafelyExecuteWithDir(cmdfl, srcdir, nil)
		if err != nil {
			fmt.Printf("protoc output: %s\n", out)
			fmt.Printf("Failed to compile: %s\n", err)
			continue
		}
		if out != "" {
			fmt.Printf("protoc output: %s\n", out)
		}
	}
	//mostRecent = m.result
	fmt.Printf("meta compiler done\n")
	return nil
}
func debugf(format string, args ...interface{}) {
	if !*debug {
		return
	}
	fmt.Fprintf(os.Stderr, format, args...)
}

// this function is called by protoc - it is the implementation of the meta compiler
func (m *MetaCompiler) SubmitSource(ctx context.Context, req *pr.ProtocRequest) error {
	if !*accept_all_tokens {
		if req.VerifyToken != m.verifytoken {
			fmt.Printf("Invoked with an invalid token: (%s)\n", req.VerifyToken)
			return errors.AccessDenied(ctx, "invalid token from protoc")
		}
	}
	/*
		for _, pf := range req.ProtoFiles {
			fmt.Printf("Protoc Request: %s\n", *pf.Name)
		}
	*/
	err := m.generate(req)
	if err != nil {
		return err
	}
	return nil
}
func findCompiler() string {
	check := []string{
		"dist/linux/amd64/protoc-gen-meta",
		"linux/amd64/protoc-gen-meta",
		"/home/cnw/go/bin/protoc-gen-meta",
		"/home/cnw/devel/go/protorenderer/dist/linux/amd64/protoc-gen-meta",
	}
	for _, c := range check {
		if !utils.FileExists(c) {
			continue
		}
		if c[0] == '/' {
			return c
		}
		cs, err := filepath.Abs(c)
		if err != nil {
			fmt.Printf("Unable to absolutise \"%s\": %s\n", cs, err)
			return c
		}
		return cs
	}
	fmt.Printf("protoc-gen-meta not found\n")
	panic("protoc-gen-meta not found")
}
func (m *MetaCompiler) Packages() []*Package {
	if m.result != nil {
		return m.result.Packages
	}
	return nil

}

// most recently parsed result - get package
func (m *MetaCompiler) PackageByID(pkgid string) *Package {
	for _, r := range m.result.Packages {
		if r.Proto.ID == pkgid {
			return r
		}
	}
	return nil
}




























































































