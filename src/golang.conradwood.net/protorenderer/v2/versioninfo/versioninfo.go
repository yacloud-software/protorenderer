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
	return &VersionInfo{
		files:          make(map[string]*VersionFile),
		compile_result: &compileresult{},
		created:        time.Now(),
	}
}
func (vi *VersionInfo) SetDirty() {
	fmt.Printf("[versioninfo] marking as dirty\n")
	vi.isclean = false
}
func (vi *VersionInfo) IsDirty() bool {
	return !vi.isclean
}
func (vi *VersionInfo) ToProto() *pb.VersionInfo {
	res := &pb.VersionInfo{
		Created: uint32(vi.created.Unix()),
	}
	vi.Lock()
	defer vi.Unlock()
	for _, vf := range vi.files {
		pbvf := vf.ToProto()
		res.Files = append(res.Files, pbvf)
	}
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
		Filename: vf.filename,
		Result:   vf.lastCompileResult,
	}
	if res.Result == nil {
		res.Result = &pb.FileResult{Filename: res.Filename}
	}
	return res
}
func (vf *VersionFile) fromProto(pvf *pb.VersionFile) {
	vf.filename = pvf.Filename
	vf.lastCompileResult = pvf.Result
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
