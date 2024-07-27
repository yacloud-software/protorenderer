package srv

type StandardCompilerEnvironment struct {
	workdir                  string
	new_protos_same_as_store bool
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
