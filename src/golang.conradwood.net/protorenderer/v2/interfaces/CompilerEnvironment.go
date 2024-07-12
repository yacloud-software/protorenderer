package interfaces

// a compiler needs some information on where to find things
type CompilerEnvironment interface {
	AllKnownProtosDir() string // directory containing all .proto files relative to workdir
	NewProtosDir() string      // directory containing .proto files which are meant to be compiled
	WorkDir() string
	CompilerOutDir() string // directory where to place artefacts from the compilers
}
