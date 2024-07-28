package srv

import (
	"context"
	"fmt"
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/v2/helpers"
	"golang.conradwood.net/protorenderer/v2/interfaces"
	"golang.conradwood.net/protorenderer/v2/meta_compiler"
	"golang.conradwood.net/protorenderer/v2/store"
	"golang.conradwood.net/protorenderer/v2/versioninfo"
	"time"
)

// recompiles the entire store, rebuilds versioninfo
func RecompileStore(ce interfaces.CompilerEnvironment) error {
	/*
		ce.(*StandardCompilerEnvironment).new_protos_same_as_store = true
		defer func() {
			ce.(*StandardCompilerEnvironment).new_protos_same_as_store = false
		}()
	*/

	fmt.Printf("[recompilestore] Start\n")
	vi := versioninfo.New()

	fname := ce.StoreDir() + "/versioninfo.pbbin"
	b, err := utils.ReadFile(fname)
	if err != nil {
		// cannot load, then make sure we write later
		vi.SetDirty() // must be stored later
		fmt.Printf("Failed to load versioninfo: %s\n", err)
	} else {
		pbvi := &pb.VersionInfo{}
		err = utils.UnmarshalBytes(b, pbvi)
		if err != nil {
			return err
		}
		vi = versioninfo.NewFromProto(pbvi)
	}

	fmt.Printf("[recompilestore] building missing meta (.info) files\n")

	must_save, err := build_version_info(ce, vi)
	if err != nil {
		return err
	}
	if must_save {
		vi.SetDirty()
	}
	// rebuild based on versioninfo
	pfs, err := helpers.FindProtoFiles(ce.AllKnownProtosDir())
	if err != nil {
		return err
	}

	fmt.Printf("[recompilestore] recompiling %d .proto files from %s\n", len(pfs), ce.AllKnownProtosDir())
	ctx := authremote.ContextWithTimeout(time.Duration(15) * time.Minute)
	err = compile_all_compilers(ctx, ce, vi.CompileResult(), pfs)
	if err != nil {
		return err
	}
	fmt.Printf("[recompilestore] Merging result to store\n")
	ce.(*StandardCompilerEnvironment).new_protos_same_as_store = false
	err = helpers.MergeCompilerEnvironment(ce, true)
	if err != nil {
		return err
	}
	if vi.IsDirty() {
		fmt.Printf("[recompilestore] Storing\n")
		b, err = utils.MarshalBytes(vi.ToProto())
		if err != nil {
			return err
		}
		err = utils.WriteFile(fname, b)
		if err != nil {
			return err
		}
		err = store.Store(ctx, ce.StoreDir()) // 0 == latest
		if err != nil {
			return err
		}
		fmt.Printf("Saved store\n")
	}
	fmt.Printf("[recompilestore] Completed\n")
	return nil
}

// get all meta info for all protos we have in our store.
// rebuild as/if/when necessary all meta information
func build_version_info(ce interfaces.CompilerEnvironment, vi *versioninfo.VersionInfo) (bool, error) {
	scr := vi.CompileResult()
	pfs, err := helpers.FindProtoFiles(ce.AllKnownProtosDir())
	if err != nil {
		return false, err
	}
	store_required := false
	ctx := context.Background()
	for _, pf := range pfs {
		mi, err := meta_compiler.GetMetaInfo(ctx, ce, pf.GetFilename())
		if err != nil {
			fmt.Printf("Failed to get meta info for %s: %s - compile triggered\n", pf.GetFilename(), utils.ErrorString(err))
			mc := meta_compiler.New()
			ctx = authremote.ContextWithTimeout(time.Duration(90) * time.Second)
			meta_needed := []interfaces.ProtoFile{pf}
			err = mc.Compile(ctx, ce, meta_needed, ce.CompilerOutDir()+"/meta", scr)
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
	return store_required, nil
}
