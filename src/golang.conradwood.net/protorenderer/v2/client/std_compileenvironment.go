package main

type StandardCompilerEnvironment struct {
	workdir string
}

func (sce *StandardCompilerEnvironment) StoreDir() string {
	return sce.workdir + "/store"
}
func (sce *StandardCompilerEnvironment) AllKnownProtosDir() string {
	return sce.StoreDir() + "/protos"
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
