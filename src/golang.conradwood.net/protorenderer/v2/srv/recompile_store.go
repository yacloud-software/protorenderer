package srv

import (
	"context"
	"fmt"
	//	pb "golang.conradwood.net/apis/protorenderer2"
	"flag"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/v2/helpers"
	"golang.conradwood.net/protorenderer/v2/interfaces"
	"golang.conradwood.net/protorenderer/v2/meta_compiler"
	"golang.conradwood.net/protorenderer/v2/store"
	"golang.conradwood.net/protorenderer/v2/versioninfo"
	"time"
)

var (
	recompile_ignore_existing_versioninfo = flag.Bool("recompile_rebuild_from_scratch", false, "if true, ignore versioninfo etc, and rebuild from scratch")
)

// recompiles the entire store, rebuilds versioninfo
func RecompileStore(ce interfaces.CompilerEnvironment) error {
	/*
		ce.(*StandardCompilerEnvironment).new_protos_same_as_store = true
		defer func() {
			ce.(*StandardCompilerEnvironment).new_protos_same_as_store = false
		}()
	*/
	fname := ce.StoreDir() + "/versioninfo.pbbin"
	fmt.Printf("[recompilestore] Start\n")
	vi := currentVersionInfo
	if *recompile_ignore_existing_versioninfo {
		vi = versioninfo.New()
	}
	fmt.Printf("[recompilestore] building missing meta (.info) files\n")

	must_save, err := build_version_info(ce, vi)
	if err != nil {
		return errors.Wrap(err)
	}
	if must_save {
		vi.SetDirty()
	}
	// rebuild based on versioninfo
	pfs, err := helpers.FindProtoFiles(ce.AllKnownProtosDir())
	if err != nil {
		return errors.Wrap(err)
	}

	fmt.Printf("[recompilestore] total of %d protofiles found in store\n", len(pfs))
	// remove protofiles which are marked as "successful"
	var n_pfs []interfaces.ProtoFile
	for _, pf := range pfs {
		//		vf := vi.GetVersionFile(pf.GetFilename())
		//		if vf == nil { // TODO: || vf.HasAtLeastOneCompilerFailed() {
		n_pfs = append(n_pfs, pf)
		//		}
	}
	pfs = n_pfs
	fmt.Printf("[recompilestore] recompiling %d protofiles\n", len(pfs))

	fmt.Printf("[recompilestore] recompiling %d .proto files from %s\n", len(pfs), ce.AllKnownProtosDir())
	ctx := authremote.ContextWithTimeout(time.Duration(15) * time.Minute)
	err = compile_all_compilers(ctx, ce, vi.CompileResult(), pfs)
	if err != nil {
		return errors.Wrap(err)
	}
	fmt.Printf("[recompilestore] Merging result to store\n")
	ce.(*StandardCompilerEnvironment).new_protos_same_as_store = false
	err = helpers.MergeCompilerEnvironment(ce, ce, true)
	if err != nil {
		return errors.Wrap(err)
	}
	if vi.IsDirty() {
		err := saveVersionInfo()
		if err != nil {
			return errors.Wrap(err)
		}
		err = store.Store(ctx, ce.StoreDir()) // 0 == latest
		if err != nil {
			return errors.Wrap(err)
		}
		fmt.Printf("Saved store\n")
	}
	vpb := vi.ToProto()
	fname = "/tmp/versioninfo.yaml"
	err = utils.WriteYaml(fname, vpb)
	if err != nil {
		fmt.Printf("[recompilestore] write failure: %s", err)
	} else {
		fmt.Printf("[recompilestore] versioninfo written as debug output to %s\n", fname)
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
	utils.Bail("failed to merge compiler env", helpers.MergeCompilerEnvironment(ce, ce, true))
	return store_required, nil
}
