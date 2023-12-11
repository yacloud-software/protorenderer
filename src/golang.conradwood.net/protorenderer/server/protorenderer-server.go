package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"golang.conradwood.net/apis/common"
	ost "golang.conradwood.net/apis/objectstore"
	pb "golang.conradwood.net/apis/protorenderer"
	"sort"
	//	"golang.conradwood.net/go-easyops/auth"
	ar "golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/client"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/server"
	"golang.conradwood.net/go-easyops/utils"
	cm "golang.conradwood.net/protorenderer/common"
	"golang.conradwood.net/protorenderer/compiler"
	"golang.conradwood.net/protorenderer/db"
	fl "golang.conradwood.net/protorenderer/filelayouter"
	"golang.conradwood.net/protorenderer/meta"
	pc "golang.conradwood.net/protorenderer/protocache"
	"golang.conradwood.net/protorenderer/protoparser"
	"google.golang.org/grpc"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	VERSIONOBJECT  = "protorenderer_version"
	INDEX_FILENAME = "protorenderer_index_file"
)

var (
	compile_count       = 0 // how many times did we compile stuff
	osc                 ost.ObjectStoreClient
	uosLock             sync.Mutex
	start_update        = false
	topdir              = flag.String("workdir", "/tmp/protos", "all directories maintained will be relative to this")
	prefix_object_store = flag.String("prefix_object_store", "protorenderer-tmp", "a prefix to be used for objectstore put/get")
	read_object_store   = flag.Bool("read_object_store", true, "if false no objects will be retrieved from objectstore and the index will not be updated")
	update_object_store = flag.Bool("update_object_store", true, "if false index in objectstore will not be updated")
	compile_python      = flag.Bool("compile_python", true, "if true compile python...")
	compile_java        = flag.Bool("compile_java", true, "if true compile java...")
	compile_nano        = flag.Bool("compile_nanopb", true, "if true compile with nanopb")
	port                = flag.Int("port", 4102, "The grpc server port")
	recreate_version    = flag.Bool("initialise_version", false, "if true will start version counting from 0 again. this will reshuffle proto/service ids (that is: it will break deep html links). intented use is ONCE for the first time protorenderer-server is started")
	protocache          = pc.New()
	updateChan          = make(chan *updateinfo, 10000)
	names               = make(map[string]bool)
	current             *version
	nextVersion         *version
	completeVersion     *version
	compiling           = false
	debug               = flag.Bool("debug", false, "debug mode")
	dbproto             *db.DBDBProtoFile
)

type version struct {
	filelayouter       *fl.FileLayouter
	goCompiler         compiler.Compiler
	javaCompiler       compiler.Compiler
	pythonCompiler     compiler.Compiler
	nanopbCompiler     compiler.Compiler
	metaCompiler       *meta.MetaCompiler
	protocache_version int
	version            int
	failures           *failuretracker
}
type protoRenderer struct {
}
type updateinfo struct {
}

func main() {
	flag.Parse()
	server.SetHealth(server.STARTING)
	err := cm.RecreateSafely(TopDir())
	utils.Bail(fmt.Sprintf("Failed to recreate topdir (%s)", TopDir()), err)
	fmt.Printf("Starting ProtoRendererServiceServer...\n")
	fmt.Printf("Workdir: \"%s\"\n", TopDir())
	dbproto = db.DefaultDBDBProtoFile()

	lv := 0
	for {
		bs, err := client.Get(ar.Context(), VERSIONOBJECT)
		if err != nil {
			fmt.Printf("Failed to get version cache: %s\n", utils.ErrorString(err))
			if *recreate_version {
				err = client.PutWithID(ar.Context(), VERSIONOBJECT, []byte("1"))
				utils.Bail("failed to initialize version cache", err)
				time.Sleep(1 * time.Second)
				continue
			}
			time.Sleep(10 * time.Second)
			continue
		}
		v, err := strconv.Atoi(string(bs))
		if err != nil {
			fmt.Printf("Invalid integer in version cache: %s\n", err)
			time.Sleep(10 * time.Second)
			continue
		}
		lv = v
		break

	}

	fly := fl.New(protocache, fmt.Sprintf("%s/%d/", TopDir(), lv))
	cc := &compilerCallback{nfly: fly}
	nextVersion = &version{
		filelayouter:   fly,
		goCompiler:     compiler.NewGoCompiler(cc),
		javaCompiler:   compiler.NewJavaCompiler(cc),
		pythonCompiler: compiler.NewPythonCompiler(cc),
		nanopbCompiler: compiler.NewNanoPBCompiler(cc),
		metaCompiler:   meta.NewMetaCompiler(fly),
		version:        lv,
		failures:       &failuretracker{},
	}
	cc.metacompiler = nextVersion.metaCompiler
	current = nextVersion
	go updater()
	sd := server.NewServerDef()
	sd.SetPort(*port)
	e := new(protoRenderer)
	sd.SetRegister(server.Register(
		func(server *grpc.Server) error {
			pb.RegisterProtoRendererServiceServer(server, e)
			return nil
		},
	))
	go startup()
	err = server.ServerStartup(sd)
	utils.Bail("Unable to start server", err)
	os.Exit(0)

}

// return Toplevel dir, with trailing /
func TopDir() string {
	var err error
	res := *topdir
	if !strings.HasPrefix(res, "/") {
		res, err = filepath.Abs(res)
		if err != nil {
			panic(fmt.Sprintf("Unable to get absolute dir for '%s':%s", res, err))
		}
	}
	if res != "" && !strings.HasSuffix(res, "/") {
		res = res + "/"
	}
	return res
}

/************************************
* grpc functions
************************************/
func (e *protoRenderer) ListSourceFiles(ctx context.Context, req *common.Void) (*pb.FilenameList, error) {
	res := &pb.FilenameList{}
	for k, v := range names {
		if !v {
			continue
		}
		res.Files = append(res.Files, k)
	}
	sort.Slice(res.Files, func(i, j int) bool {
		return res.Files[i] < res.Files[j]
	})
	return res, nil
}
func (e *protoRenderer) DeleteFile(ctx context.Context, req *pb.DeleteRequest) (*common.Void, error) {
	if !names[req.Name] {
		return nil, errors.InvalidArgs(ctx, "file not found", "file \"%s\" not fond", req.Name)
	}
	names[req.Name] = false
	if updateObjectStore() {
		pir := &ost.PutWithIDRequest{ID: *prefix_object_store + req.Name, Expiry: 1}
		_, err := osclient().PutWithID(ctx, pir)
		if err != nil {
			return nil, err
		}
		err = rewriteIndexFile(ctx)
		if err != nil {
			return nil, err
		}
	}
	pv, err := protocache.Delete(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	nv := 0
	if pv != current.protocache_version {
		nv = nextVersion.version
	}
	fmt.Printf("NextVersion: %d, ProtocacheVersion: %d\n", nv, pv)
	ui := &updateinfo{}
	updateChan <- ui
	return &common.Void{}, nil
}

func (e *protoRenderer) GetVersion(ctx context.Context, req *common.Void) (*pb.Version, error) {
	ev := NeedVersion(ctx)
	if ev != nil {
		return nil, ev
	}
	return &pb.Version{
		Version:      uint64(completeVersion.version),
		Compiling:    compiling,
		NextVersion:  uint64(current.version),
		ProtoVersion: uint64(completeVersion.protocache_version),
	}, nil
}
func (e *protoRenderer) UpdateProto(ctx context.Context, req *pb.AddProtoRequest) (*pb.AddProtoResponse, error) {
	if len(req.Content) == 0 {
		return nil, errors.InvalidArgs(ctx, "Content of proto file must be >0", "Content of proto file (%s) must be >0", req.Name)
	}
	if len(req.Name) == 0 {
		return nil, errors.InvalidArgs(ctx, "length of name of proto file must be >0", "length of name of proto file must be >0")
	}
	fmt.Printf("New Proto submitted %s (%d) bytes:\n", req.Name, len(req.Content))
	// normalise - if starts with '/protos/' remove that bit
	rn := req.Name
	rn = strings.TrimPrefix(rn, "/")
	rn = strings.TrimPrefix(rn, "protos/")
	req.Name = rn

	pp, err := protoparser.Parse(req.Content)
	if err != nil {
		fmt.Printf("Failed to parse: %s\n", err)
		return nil, err
	}
	fmt.Printf("  Go   Package: %s\n", pp.GoPackage)
	fmt.Printf("  Java Package: %s\n", pp.JavaPackage)
	if len(pp.JavaPackage) < 2 {
		pp.JavaPackage = "no.javapackage.proto." + pp.GoPackage
		if *compile_java {
			return nil, errors.InvalidArgs(ctx, "Invalid java package", "Invalid java package \"%s\" in file %s", pp.JavaPackage, req.Name)
		}
	}

	names[req.Name] = true
	if updateObjectStore() && start_update {
		pir := &ost.PutWithIDRequest{ID: *prefix_object_store + req.Name, Content: []byte(req.Content)}
		_, err := osclient().PutWithID(ctx, pir)
		if err != nil {
			return nil, err
		}
		err = rewriteIndexFile(ctx)
		if err != nil {
			return nil, err
		}
		_, err = findOrUpdateProtoInDB(ctx, req)
		if err != nil {
			// ignoring the error during migration
			fmt.Printf("failed to update db: %s\n", utils.ErrorString(err))
		}
	}
	add := pp.Protofile()
	// maintain the repositoryid in the cache
	pfdb, err := dbproto.ByName(context.Background(), req.Name)
	if err == nil && len(pfdb) != 0 {
		add.RepositoryID = pfdb[0].RepositoryID
	}
	add.Filename = req.Name
	nv := current.version
	pv, err := protocache.AddOrUpdate(add)
	if err != nil {
		fmt.Printf("failed to add protofile: %s\n", err)
		return nil, err
	}
	if pv != current.protocache_version {
		nv = nextVersion.version
	}
	if start_update {
		ui := &updateinfo{}
		updateChan <- ui
	}
	res := &pb.AddProtoResponse{
		ProtoVersion: uint64(pv),
		Version:      uint64(nv),
		GoPackage:    pp.GoPackage,
		JavaPackage:  pp.JavaPackage,
		Imports:      pp.Imports,
	}
	return res, nil
}

func (e *protoRenderer) MiniParser(ctx context.Context, req *pb.AddProtoRequest) (*pb.ProtoFile, error) {
	pp, err := protoparser.Parse(req.Content)
	if err != nil {
		return nil, err
	}
	res := pp.Protofile()
	return res, nil

}

// the thing that actually compiles
func updater() {
	lastone := 0
	fmt.Printf("Updater started...\n")
	for {
		compiling = false
		_ = <-updateChan
		compiling = true
		v := protocache.Version()
		if v == lastone {
			fmt.Printf("Nothing to do: new version: %d, lastone= %d\n", v, lastone)
			continue
		}
		fmt.Printf("Compiling...\n")
		current = nextVersion
		current.protocache_version = v
		nv := current.version + 1
		nfly := fl.New(protocache, fmt.Sprintf("%s/%d/", TopDir(), nv))
		cc := &compilerCallback{nfly: nfly}
		nextVersion = &version{
			filelayouter:   nfly,
			goCompiler:     compiler.NewGoCompiler(cc),
			javaCompiler:   compiler.NewJavaCompiler(cc),
			pythonCompiler: compiler.NewPythonCompiler(cc),
			nanopbCompiler: compiler.NewNanoPBCompiler(cc),
			metaCompiler:   meta.NewMetaCompiler(nfly),
			version:        nv,
			failures:       &failuretracker{},
		}
		cc.metacompiler = nextVersion.metaCompiler
		var err error
		fmt.Printf("********************** Creating new version ***************************\n")
		ctx := ar.Context()
		current.filelayouter.Save(ctx)
		if err != nil {
			fmt.Printf("Error saving: %s\n", err)
			continue
		}

		err = current.metaCompiler.Compile(*port)
		if err != nil {
			fmt.Printf("Error compiling metadata: %s\n", err)
		}

		err = current.goCompiler.Compile(current.failures)
		if err != nil {
			fmt.Printf("Error compiling go: %s\n", err)
		}

		if *compile_nano {
			err = current.nanopbCompiler.Compile(current.failures)
			if err != nil {
				fmt.Printf("Error compiling nanopb: %s\n", err)
			}
		}

		if *compile_java {
			err = current.javaCompiler.Compile(current.failures)
			if err != nil {
				fmt.Printf("Error compiling java: %s\n", err)
			}
		}
		if *compile_python {
			err = current.pythonCompiler.Compile(current.failures)
			if err != nil {
				fmt.Printf("Error compiling python: %s\n", err)
			}
		}
		lastone = v

		completeVersion = current
		completeVersion.protocache_version = v

		ctx = ar.Context()
		bs := fmt.Sprintf("%d", completeVersion.version)
		if updateObjectStore() {
			client.PutWithID(ctx, VERSIONOBJECT, []byte(bs))
		}
		compile_count++
		if compile_count == 1 {
			server.SetHealth(server.READY)
		}
		fmt.Printf("********************** Created version %d ***************************\n", completeVersion.version)
	}
}

func updateObjectStore() bool {
	if !*read_object_store {
		return false
	}
	return *update_object_store
}

func osclient() ost.ObjectStoreClient {
	if osc != nil {
		return osc
	}
	osc = ost.NewObjectStoreClient(client.Connect("objectstore.ObjectStore"))
	return osc
}

// rewrite index file based on 'names'
func rewriteIndexFile(ctx context.Context) error {
	uosLock.Lock()
	defer uosLock.Unlock()
	var bf bytes.Buffer
	for k, _ := range names {
		if names[k] {
			bf.WriteString(fmt.Sprintf("%s\n", k))
		}
	}
	pir := &ost.PutWithIDRequest{ID: *prefix_object_store + INDEX_FILENAME, Content: bf.Bytes()}
	_, err := osclient().PutWithID(ctx, pir)
	return err
}

func (e *protoRenderer) FindServiceByName(ctx context.Context, req *pb.FindServiceByNameRequest) (*pb.ServiceList, error) {
	if completeVersion == nil || completeVersion.metaCompiler == nil {
		return nil, errors.FailedPrecondition(ctx, "version not yet available")
	}
	pkgs := completeVersion.metaCompiler.Packages()
	if pkgs == nil {
		return nil, errors.FailedPrecondition(ctx, "packages not yet available")
	}
	res := &pb.ServiceList{}
	for _, pkg := range pkgs {
		for _, svc := range pkg.Services {
			//			fmt.Printf("[%s] Service: %s / %s\n", req.Name, pkg.Name, svc.Name)
			add := false
			if svc.Name == req.Name {
				add = true
			}
			fp := fmt.Sprintf("%s.%s", pkg.Name, svc.Name)
			if fp == req.Name {
				add = true
			}
			if add {
				asvc := &pb.ServiceResponse{
					Service: &pb.Service{
						ID:      svc.ID,
						Name:    svc.Name,
						Comment: svc.Comment,
					},
					Package:     pkg.Proto,
					PackageName: pkg.Name,
					PackageFQDN: pkg.FQDN,
				}
				res.Services = append(res.Services, asvc)
			}
		}
	}
	return res, nil
}
func (e *protoRenderer) FindServiceByID(ctx context.Context, req *pb.ID) (*pb.ServiceResponse, error) {
	ev := NeedVersion(ctx)
	if ev != nil {
		return nil, ev
	}
	result := completeVersion.metaCompiler.GetMostRecentResult()
	if result == nil {
		return nil, errors.Unavailable(ctx, "GetPackages (most recent result)")
	}
	for _, pkg := range result.Packages {
		for _, svc := range pkg.Services {
			if svc.ID == req.ID {
				asvc := &pb.ServiceResponse{
					Service: &pb.Service{
						ID:        svc.ID,
						Name:      svc.Name,
						Comment:   svc.Comment,
						PackageID: pkg.Proto.ID,
					},
					Package:     pkg.Proto,
					PackageName: pkg.Name,
					PackageFQDN: pkg.FQDN,
				}
				return asvc, nil
			}
		}
	}
	return nil, errors.NotFound(ctx, "service not found (id=%s)", req.ID)

}
func (e *protoRenderer) GetFailedFiles(ctx context.Context, req *common.Void) (*pb.FailedFilesList, error) {
	if completeVersion == nil {
		return nil, errors.Unavailable(ctx, "GetFailedFiles")
	}
	result := completeVersion.metaCompiler.GetMostRecentResult()
	if result == nil {
		return nil, errors.Unavailable(ctx, "GetFailedFiles")
	}
	res := &pb.FailedFilesList{}
	rt := completeVersion.failures
	for _, f := range rt.Failures() {
		ff := &pb.FailedFile{
			Filename: f.filename,
			Message:  f.message,
			Compiler: f.c.Name(),
		}
		res.Files = append(res.Files, ff)
	}
	return res, nil
}
























































