/*
   VersionInfo contains information about a specific version.
   A Version is a specific set of files, source (.proto) as well as artefacts.
   For each version we keep the following information, in addition to the artefacts:
   * which files compiled successfully with which compiler
   * which files failed to compile with which compiler
   * dependencies for each file (extracted by meta compiler)

   A Versioninfo only covers files in store, not new_protos

   A Version info may be serialised as part of binaryversions in the store.

   It serves two main purposes:

   1. Provide data for an algorithm to determine which existing .proto files need to be recompiled upon submission of a new .proto (its dependencies)
   2. Provide data for a renderer to indicate which files failed and why
*/

package versioninfo

import (
	"fmt"
	pb "golang.conradwood.net/apis/protorenderer2"
	//	"golang.conradwood.net/protorenderer/v2/common"
	"golang.conradwood.net/protorenderer/v2/interfaces"
	"sort"
	"sync"
	"time"
)

type VersionInfo struct {
	sync.Mutex
	files          map[string]*VersionFile
	compile_result *compileresult
	isclean        bool // true if it identical to the version uploaded to binary versions (and has not changed since)
	created        time.Time
}

func NewFromProto(vi *pb.VersionInfo) *VersionInfo {
	vin := New()
	vin.fromProto(vi)
	vin.isclean = true
	return vin
}
func New() *VersionInfo {
	vi := &VersionInfo{
		files:          make(map[string]*VersionFile),
		compile_result: &compileresult{},
		created:        time.Now(),
	}
	vi.compile_result.vi = vi
	return vi
}
func (vi *VersionInfo) SetDirty() {
	if vi.isclean == false {
		return
	}
	fmt.Printf("[versioninfo] marking as dirty\n")
	vi.isclean = false
}
func (vi *VersionInfo) IsDirty() bool {
	vi.update_version_files_from_compile_result()
	return !vi.isclean
}
func (vi *VersionInfo) ToProto() *pb.VersionInfo {
	vi.update_version_files_from_compile_result()
	res := &pb.VersionInfo{
		Created: uint32(vi.created.Unix()),
	}
	vi.Lock()
	defer vi.Unlock()
	for _, vf := range vi.files {
		pbvf := vf.ToProto()
		res.Files = append(res.Files, pbvf)
	}
	sort.Slice(res.Files, func(i, j int) bool {
		return res.Files[i].Filename < res.Files[j].Filename
	})
	return res
}
func (vi *VersionInfo) fromProto(pvi *pb.VersionInfo) {
	vi.Lock()
	defer vi.Unlock()
	vi.created = time.Unix(int64(pvi.Created), 0)
	for _, pvf := range pvi.Files {
		vf := &VersionFile{}
		vf.fromProto(pvf)
		vi.files[pvf.Filename] = vf
	}
}

type VersionFile struct {
	meta              *pb.ProtoFileInfo
	filename          string
	lastCompileResult *pb.FileResult // includes error messages for each compiler (as []pb.CompileFailure)
}

func (vf *VersionFile) ToProto() *pb.VersionFile {
	res := &pb.VersionFile{
		Filename:   vf.filename,
		FileResult: vf.lastCompileResult,
	}
	return res
}
func (vf *VersionFile) fromProto(pvf *pb.VersionFile) {
	vf.filename = pvf.Filename
	vf.lastCompileResult = pvf.FileResult
	if vf.lastCompileResult == nil {
		vf.lastCompileResult = &pb.FileResult{Filename: vf.filename}
	}
}
func (vi *VersionInfo) GetOrAddFile(file string, mi *pb.ProtoFileInfo) *VersionFile {
	vi.Lock()
	vf, ok := vi.files[file]
	if !ok {
		vf = &VersionFile{filename: file, meta: mi,
			lastCompileResult: &pb.FileResult{Filename: file},
		}
		vi.files[file] = vf
		fmt.Printf("Adding \"%s\"\n", file)
		vi.SetDirty()
	}

	vi.Unlock()
	return vf
}

// may return nil
func (vi *VersionInfo) GetVersionFile(file string) *VersionFile {
	vi.Lock()
	defer vi.Unlock()
	vf, ok := vi.files[file]
	if !ok {
		return nil
	}
	return vf
}

func (vi *VersionInfo) CompileResult() interfaces.CompileResult {
	return vi.compile_result
}

// this reads the compile result and updates the versionfiles. if necessary sets the dirty flag
func (vi *VersionInfo) update_version_files_from_compile_result() {
	vi.Lock()
	defer vi.Unlock()
	cr := vi.compile_result

	fmt.Printf("[versioninfo] updating version files from compile_result (%d files,%d results))\n", len(vi.files), len(cr.results))

	// Step #1 - see if we have proto files we did not have yet, add them if so
	for _, crf := range cr.results {
		file := crf.pf.GetFilename()
		_, f := vi.files[file]
		if f {
			continue
		}
		vf := &VersionFile{filename: file,
			lastCompileResult: &pb.FileResult{Filename: file},
		}
		vi.files[file] = vf
		fmt.Printf("Adding \"%s\"\n", file)
		vi.SetDirty()
	}

	// Step #2 - replace results with results from compile result
	for _, vf := range vi.files {
		change := vi.merge_result(vf, cr.getresultsforfile(vf.filename))
		if change != 0 {
			fmt.Printf("File compile result changed: %s\n", vf.filename)
			vi.SetDirty()
		}
	}
}

// return what sort of change it is
func (vi *VersionInfo) merge_result(vf *VersionFile, crfa []*compileresult_file) int {
	res := 0
	for _, crf := range crfa { // crf is the latest compile result for a given compiler
		new_cf := crf.toCompileResultProto()       // into pb.CompileResult
		cur_cf := vf.result_for_compiler(crf.comp) // get the pb.CompileResult for what is currently recorded in versioninfo
		fmt.Printf("[versioninfo] merging recorded (%v) and new (%v)\n", cur_cf, new_cf)
		// Check #1: recorded status is NIL, add latest result
		if cur_cf == nil {
			vf.lastCompileResult.CompileResults = append(vf.lastCompileResult.CompileResults, new_cf)
			res = 1
			continue
		}

		// at this point recorded status is a fail
		// Check #2: recorded status is fail, new status is OK
		if !crf.fail {
			vf.removeCompilerFailure(crf.comp)
			res = 1
			continue
		}
	}
	return res
}

// might return nil
func (vf *VersionFile) result_for_compiler(c interfaces.Compiler) *pb.CompileResult {
	if vf.lastCompileResult == nil {
		return nil
	}
	for _, cf := range vf.lastCompileResult.CompileResults {
		if cf.CompilerName == c.ShortName() {
			return cf
		}
	}
	return nil
}

// remove the entry for this compiler
func (vf *VersionFile) removeCompilerFailure(c interfaces.Compiler) {
	var n []*pb.CompileResult
	if vf.lastCompileResult == nil {
		return
	}
	for _, f := range vf.lastCompileResult.CompileResults {
		if f.CompilerName == c.ShortName() {
			continue
		}
		n = append(n, f)
	}
	vf.lastCompileResult.CompileResults = n
}

// true if compiler compiled successfully
func (vf *VersionFile) ResultForCompiler(compname string) *pb.CompileResult {
	return nil
}
