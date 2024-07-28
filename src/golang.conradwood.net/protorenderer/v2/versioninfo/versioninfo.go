/*
   VersionInfo contains information about a specific version.
   A Version is a specific set of files, source (.proto) as well as artefacts.
   For each version we keep the following information, in addition to the artefacts:
   * which files compiled successfully with which compiler
   * which files failed to compile with which compiler
   * dependencies for each file (extracted by meta compiler)

   A Versioninfo only covers files in store, not new_protos

   A Version info may be serialised as part of binaryversions in the store.
*/

package versioninfo

import (
	"fmt"
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/protorenderer/v2/common"
	"sync"
)

type VersionInfo struct {
	sync.Mutex
	files          map[string]*VersionFile
	compile_result *common.StandardCompileResult
}

func New() *VersionInfo {
	return &VersionInfo{
		files:          make(map[string]*VersionFile),
		compile_result: &common.StandardCompileResult{},
	}
}

type VersionFile struct {
	meta     *pb.ProtoFileInfo
	filename string
}

func (vi *VersionInfo) GetOrAddFile(file string, mi *pb.ProtoFileInfo) *VersionFile {
	vi.Lock()
	vf, ok := vi.files[file]
	if !ok {
		vf = &VersionFile{filename: file, meta: mi}
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

func (vi *VersionInfo) CompileResult() *common.StandardCompileResult {
	return vi.compile_result
}
