package java

import (
	"flag"
	"fmt"
	"golang.conradwood.net/go-easyops/linux"
	"strings"
	"sync"
	"time"
)

const (
	MAX_JAVA_FILES = 100
)

var (
	failed_java     = make(map[string]*JavaCompileError)
	max_workers     = flag.Int("java_max_workers", 2, "max number of parallel java processes")
	workerWaitGroup *sync.WaitGroup
	workerChan      = make(chan *jc, 50)
)

type JavaCompileError struct {
	jfg *JavaFileGroup
}
type JavaFile struct {
	Absolute string
	Relative string
}
type Java2Class struct {
	SourceDir string      // working directory for compiler
	JFiles    []*JavaFile // .java file name. SourceDir+Files[n] == absolute path
	Command   []string
}
type JavaFileGroup struct {
	jc     *Java2Class
	Files  []*JavaFile
	Result error // if nil all went well
}

func (j *Java2Class) Files() []*JavaFile {
	var res []*JavaFile
	for _, f := range j.JFiles {
		jce := failed_java[f.Relative]
		if jce != nil {
			fmt.Printf("marked as failed: skipping %s\n", f.Relative)
			continue
		}
		res = append(res, f)
	}
	return res
}

/*
we take ALL files and keep splitting into half the size groups to compile
until either:
* no group failed
* a group with a single file failed
*/
func (j *Java2Class) Compile() error {
	jfg := &JavaFileGroup{jc: j, Files: j.Files()}
	groups := []*JavaFileGroup{jfg}
	// limit to max size (makes no sense to waste cycle on sth we know won't work)
	repeat := true
	for repeat {
		repeat = false
		var ng []*JavaFileGroup
		for _, g := range groups {
			if len(g.Files) <= MAX_JAVA_FILES {
				ng = append(ng, g)
				continue
			}
			g1, g2 := g.Split()
			ng = append(ng, g1)
			ng = append(ng, g2)
			repeat = true
		}
		groups = ng
	}
	for {
		start_workers()
		for _, g := range groups {
			// g.compileGroup()
			workerChan <- &jc{abort: false, group: g}
		}
		stop_workers()

		// check for abort conditions:
		failure := false
		var singleFileFailure *JavaFileGroup
		for _, g := range groups {
			if g.Result != nil {
				failure = true
				if len(g.Files) < 2 {
					singleFileFailure = g
				}
			}
		}
		if singleFileFailure != nil {
			fmt.Printf("File failed to compile: %s (%s)\n", singleFileFailure.Files[0], singleFileFailure.Result)
			mark_as_failed(singleFileFailure)
			return singleFileFailure.Result
		}
		if !failure {
			return nil
		}
		// split each group and retry
		var ng []*JavaFileGroup
		for _, g := range groups {
			if len(g.Files) == 0 {
				continue
			}
			if g.Result == nil {
				continue
			}
			g1, g2 := g.Split()
			ng = append(ng, g1)
			ng = append(ng, g2)
		}
		groups = ng
	}
}
func mark_as_failed(e *JavaFileGroup) {
	fname := e.Files[0].Relative
	fmt.Printf("Mark as failed: %s\n", fname)
	failed_java[fname] = &JavaCompileError{jfg: e}

}
func (j *JavaFileGroup) Split() (*JavaFileGroup, *JavaFileGroup) {
	r1 := &JavaFileGroup{jc: j.jc}
	r2 := &JavaFileGroup{jc: j.jc}
	il := len(j.Files)
	if il == 0 {
		return r1, r2
	}
	if il == 1 {
		r1.Files = j.Files
		return r1, r2
	}
	half := il / 2
	r1.Files = j.Files[:half]
	r2.Files = j.Files[half:]
	return r1, r2

}

func (j *JavaFileGroup) compileGroup() {
	if len(j.Files) == 0 {
		return
	}
	fmt.Printf("Compiling %d java files...\n", len(j.Files))
	cwd := j.jc.SourceDir
	cmdandfile := j.jc.Command
	for _, j := range j.Files {
		cmdandfile = append(cmdandfile, j.Absolute)
	}
	l := linux.New()
	l.SetMaxRuntime(time.Duration(600) * time.Second)
	l.SetAllowConcurrency(true)
	out, err := l.SafelyExecuteWithDir(cmdandfile, cwd, nil)
	j.Result = err
	if err != nil {
		c := strings.Join(cmdandfile, " ")
		if len(c) > 3000 {
			c = c[:3000]
		}
		fmt.Printf("Whilst executing: %s\n", c)
		fmt.Printf("Failure compiling .java to .class: %s\n%s\n", out, err)
		if len(j.Files) == 1 {
			fmt.Printf("Broken file: %s\n", j.Files[0].Relative)
		}
	}

}

type jc struct {
	abort bool
	group *JavaFileGroup
}

func jcompiler() {
	for {
		j := <-workerChan
		if j.abort {
			break
		}
		j.group.compileGroup()
	}
	workerWaitGroup.Done()
}

func start_workers() {
	workerWaitGroup = &sync.WaitGroup{}
	for i := 0; i < *max_workers; i++ {
		workerWaitGroup.Add(1)
		go jcompiler()
	}
}
func stop_workers() {
	for i := 0; i < *max_workers; i++ {
		workerChan <- &jc{abort: true}
	}
	workerWaitGroup.Wait()
}

























































































