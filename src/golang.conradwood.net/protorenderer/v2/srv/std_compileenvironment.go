package srv

import (
	"fmt"
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/protorenderer/v2/interfaces"
	"strings"
)

type StandardCompilerEnvironment struct {
	workdir                  string
	new_protos_same_as_store bool
	mc                       interfaces.MetaCache
	server                   compile_serve_req
}

func (sce *StandardCompilerEnvironment) StoreDir() string {
	return sce.workdir + "/store"
}
func (sce *StandardCompilerEnvironment) AllKnownProtosDir() string {
	return sce.StoreDir() + "/protos"
}
func (sce *StandardCompilerEnvironment) NewProtosDir() string {
	if sce.new_protos_same_as_store {
		return sce.AllKnownProtosDir()
	}
	return sce.workdir + "/new_protos"
}
func (sce *StandardCompilerEnvironment) WorkDir() string {
	return sce.workdir
}

func (sce *StandardCompilerEnvironment) CompilerOutDir() string {
	return sce.workdir + "/compile_outdir"
}
func (sce *StandardCompilerEnvironment) MetaCache() interfaces.MetaCache {
	return sce.mc
}
func (sce *StandardCompilerEnvironment) Fork() *StandardCompilerEnvironment {
	res := &StandardCompilerEnvironment{
		workdir: sce.workdir,
		mc:      sce.mc.Fork(),
	}
	return res
}

func (sce *StandardCompilerEnvironment) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	srv := sce.server
	if srv != nil {
		s := fmt.Sprintf(format, args...)
		s = strings.TrimSuffix(s, "\n")
		lines := strings.Split(s, "\n")
		fts := &pb.FileTransferStdout{Lines: lines}
		srv.Send(&pb.FileTransfer{Output: fts})
	}
}
