package interfaces

// a compiler needs some information on where to find things
type CompilerEnvironment interface {
	AllKnownProtosDir() string // shortcut to StoreDir+"/protos"
	StoreDir() string          // directory containing all known .proto and all compiled artefacts
	NewProtosDir() string      // directory containing .proto files which are meant to be compiled
	WorkDir() string
	CompilerOutDir() string // directory where to place artefacts from the compilers
}
