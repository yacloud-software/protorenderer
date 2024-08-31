package interfaces

import (
	pb "golang.conradwood.net/apis/protorenderer2"
)

// a compiler needs some information on where to find things
type CompilerEnvironment interface {
	AllKnownProtosDir() string // shortcut to StoreDir+"/protos"
	StoreDir() string          // directory containing all known .proto and all compiled artefacts
	NewProtosDir() string      // directory containing .proto files which are meant to be compiled
	WorkDir() string
	CompilerOutDir() string                    // directory where to place artefacts from the compilers
	MetaCache() MetaCache                      // get the meta cache for this compiler run
	Printf(format string, args ...interface{}) // print for this compiler-environment. quite possibly sent output to client
}

type MetaCache interface {
	Fork() MetaCache
	ImportFrom(MetaCache) // add all files from another metacache
	Add(*pb.ProtoFileInfo)
	ByProtoFile(pf ProtoFile) *pb.ProtoFileInfo
	ByFilename(filename string) *pb.ProtoFileInfo // by .proto file
	AllWithDependencyOn(filename string, maxdepth uint32) ([]*pb.ProtoFileInfo, error)
	AllPackages(f func(pfi *pb.ProtoFileInfo)) error
}
