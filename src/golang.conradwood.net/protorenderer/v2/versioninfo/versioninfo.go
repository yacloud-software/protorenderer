package versioninfo

import (
	"fmt"
	pb "golang.conradwood.net/apis/protorenderer2"
	"sync"
)

type VersionInfo struct {
	sync.Mutex
	files map[string]*VersionFile
}

func New() *VersionInfo {
	return &VersionInfo{
		files: make(map[string]*VersionFile),
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
