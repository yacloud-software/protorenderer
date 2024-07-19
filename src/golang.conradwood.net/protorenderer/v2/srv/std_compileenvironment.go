package srv

type StandardCompilerEnvironment struct {
	workdir string
}

func (sce *StandardCompilerEnvironment) AllKnownProtosDir() string {
	return sce.workdir + "/known_protos"
}
func (sce *StandardCompilerEnvironment) NewProtosDir() string {
	return sce.workdir + "/new_protos"
}
func (sce *StandardCompilerEnvironment) WorkDir() string {
	return sce.workdir
}

func (sce *StandardCompilerEnvironment) CompilerOutDir() string {
	return sce.workdir + "/compile_outdir"
}
