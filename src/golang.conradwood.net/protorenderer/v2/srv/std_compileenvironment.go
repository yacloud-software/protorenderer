package srv

type StandardCompilerEnvironment struct {
	knownprotosdir string
	workdir        string
}

func (sce *StandardCompilerEnvironment) AllKnownProtosDir() string {
	return sce.knownprotosdir
}
func (sce *StandardCompilerEnvironment) NewProtosDir() string {
	return sce.workdir + "/new_protos"
}
func (sce *StandardCompilerEnvironment) WorkDir() string {
	return sce.workdir
}
func (sce *StandardCompilerEnvironment) ResultsDir() string {
	return sce.workdir + "/store"
}
func (sce *StandardCompilerEnvironment) CompilerOutDir() string {
	return sce.workdir + "/compile_outdir"
}
