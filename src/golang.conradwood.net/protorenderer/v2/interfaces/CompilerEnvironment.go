package interfaces

// a compiler needs some information on where to find things
type CompilerEnvironment interface {
	AllKnownProtosDir() string // directory containing all .proto files relative to workdir
	ResultsDir() string        // directory containing all compiled things, e.g. java/classes and java/src files
	NewProtosDir() string      // directory containing .proto files which are meant to be compiled
	WorkDir() string
	CompilerOutDir() string // directory where to place artefacts from the compilers
}
