package interfaces

// a compiler needs some information on where to find things
type CompilerEnvironment interface {
	AllKnownProtosDir() string // directory containing all .proto files relative to workdir

}
