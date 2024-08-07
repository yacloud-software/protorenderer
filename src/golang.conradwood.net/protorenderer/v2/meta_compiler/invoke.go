package meta_compiler

import (
	"context"
	"fmt"
	google_protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/utils"
	pcmd "golang.conradwood.net/protorenderer/cmdline"
	cc "golang.conradwood.net/protorenderer/v2/compilers/common"
	"golang.conradwood.net/protorenderer/v2/helpers"
	"golang.conradwood.net/protorenderer/v2/interfaces"
	"sort"
	"sync"
	"time"
)

const (
	use_parallel = false
)

// TODO - instead of using gRPC use go-easyops IPC

type MetaCompiler struct {
	sync.Mutex
	id          string
	ce          interfaces.CompilerEnvironment
	files       []interfaces.ProtoFile
	cr          interfaces.CompileResult
	processed   map[string]bool
	descriptors map[string]*MessageDescriptor
}
type MessageDescriptor struct {
	ID        uint64
	descproto *google_protobuf.DescriptorProto
}

func New() *MetaCompiler {
	mc := &MetaCompiler{
		processed:   make(map[string]bool),
		id:          utils.RandomString(64),
		descriptors: make(map[string]*MessageDescriptor),
	}
	meta_compilers.Put(mc.id, mc)
	return mc
}

/*
the meta compilers works slightly different than the others. the protoc plugin is a small RPC stub, which then calls protorenderer. this function invokes protoc and protoc-gen-meta2 plugin, which then will call this process via gRPC
*/
func (gc *MetaCompiler) Compile(ctx context.Context, ce interfaces.CompilerEnvironment, files []interfaces.ProtoFile, outdir string, cr interfaces.CompileResult) error {
	port := pcmd.GetRPCPort()
	// keep the variables, we need them for the callback (protoc plugin calls this process via gRPC)
	gc.ce = ce
	gc.files = files
	gc.cr = cr
	helpers.Mkdir(ce.CompilerOutDir() + "/info")
	helpers.Mkdir(ce.StoreDir() + "/info")
	pcfname := cc.FindCompiler("protoc-gen-meta2")
	fmt.Printf("Using compiler: \"%s\"\n", pcfname)
	dir := ce.NewProtosDir()

	import_dirs := []string{
		dir,
		ce.AllKnownProtosDir(),
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
		fmt.Sprintf("--plugin=protoc-gen-meta2=%s", pcfname),
		"--meta2_out=/tmp", // has no output
		fmt.Sprintf("--meta2_opt=%s,%s,%d,%s", gc.id, sctx, port, cmdline.GetClientRegistryAddress()),
	}
	for _, id := range import_dirs {
		cmd = append(cmd, fmt.Sprintf("-I%s", id))
	}
	if use_parallel {
		wg := &sync.WaitGroup{}
		for _, xpf := range files {
			wg.Add(1)
			go func(pf interfaces.ProtoFile) {
				defer wg.Done()
				filename := pf.GetFilename()
				cmdfl := append(cmd, filename)

				out, err := l.SafelyExecuteWithDir(cmdfl, dir, nil)
				if err != nil {
					fmt.Printf("[metacompiler] protoc output: %s\n", out)
					fmt.Printf("[metacompiler] Failed to compile %s: %s\n", filename, err)
					cr.AddFailed(gc, pf, err, []byte(out))
				}
			}(xpf)
		}
		wg.Wait()
	} else {
		for _, pf := range files {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			filename := pf.GetFilename()
			cmdfl := append(cmd, filename)

			out, err := l.SafelyExecuteWithDir(cmdfl, dir, nil)
			if err != nil {
				fmt.Printf("[metacompiler] protoc output: %s\n", out)
				fmt.Printf("[metacompiler] Failed to compile %s: %s\n", filename, err)
				cr.AddFailed(gc, pf, err, []byte(out))
			}
		}
	}
	return nil
}
func (mc *MetaCompiler) ShortName() string {
	return "meta"
}

// this _only_ contains files that are being compiled at this run
func (mc *MetaCompiler) FileByName(name string) (interfaces.ProtoFile, error) {
	for _, pf := range mc.files {
		if pf.GetFilename() == name {
			return pf, nil
		}
	}
	for _, pf := range mc.files {
		fmt.Printf("Known file in meta compiler: \"%s\"\n", pf.GetFilename())
	}
	return nil, errors.Errorf("File \"%s\" not part of the meta compiler", name)
}

func (mc *MetaCompiler) WasFilenameSubmitted(filename string) bool {
	for _, pf := range mc.files {
		if pf.GetFilename() == filename {
			return true
		}
	}
	return false
}

func (mc *MetaCompiler) CompileResult() interfaces.CompileResult {
	return mc.cr
}

func (mc *MetaCompiler) CompilerEnvironment() interfaces.CompilerEnvironment {
	return mc.ce
}

func (mc *MetaCompiler) AddProcessed(filename string) {
	mc.Lock()
	defer mc.Unlock()
	mc.processed[filename] = true
}
func (mc *MetaCompiler) WasProcessed(filename string) bool {
	mc.Lock()
	defer mc.Unlock()
	_, b := mc.processed[filename]
	return b
}

func (md *MessageDescriptor) DescProto() *google_protobuf.DescriptorProto {
	return md.descproto
}
func (smc *MetaCompiler) GetMessageDescriptorByID(id uint64) *MessageDescriptor {
	for _, msg := range smc.descriptors {
		if msg.ID == id {
			return msg
		}
	}
	return nil
}
func (smc *MetaCompiler) GetMessageDescriptorByFQDN(fqdn string) *MessageDescriptor {
	msg, found := smc.descriptors[fqdn]
	if found {
		return msg
	}

	var names []string
	for k, _ := range smc.descriptors {
		names = append(names, k)
	}
	sort.Slice(names, func(i, j int) bool {
		return names[i] < names[j]
	})
	fmt.Printf("Known proto messages:\n")
	for _, n := range names {
		fmt.Printf(" \"%s\"\n", n)
	}
	fmt.Printf("Not found: \"%s\"\n", fqdn)

	return nil
}
func (mc *MetaCompiler) AddDescriptor(fqdn string, id uint64, msgtype *google_protobuf.DescriptorProto) {
	mc.descriptors[fqdn] = &MessageDescriptor{ID: id, descproto: msgtype}
}
