package main

type StandardCompilerEnvironment struct {
	knownprotosdir string
	newprotosdir   string
	workdir        string
}

func (sce *StandardCompilerEnvironment) AllKnownProtosDir() string {
	return sce.knownprotosdir
}
func (sce *StandardCompilerEnvironment) NewProtosDir() string {
	return sce.newprotosdir
}
func (sce *StandardCompilerEnvironment) WorkDir() string {
	return sce.workdir
}
func (sce *StandardCompilerEnvironment) ResultsDir() string {
	return sce.workdir + "/store"
}
