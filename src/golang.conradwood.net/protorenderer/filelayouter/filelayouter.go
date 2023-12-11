package filelayouter

/*
this will layout a bunch of .proto files on disk in directory
in such a way that import statements will work
suitable for running protoc across the directory

the 'detect' of paths should probably be pulled out and made a little smarter
*/

import (
	"context"
	"fmt"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/common"
	pc "golang.conradwood.net/protorenderer/protocache"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	FILELAYOUTER_PREFIX = "src"
)

var (
	filelock sync.Mutex
)

type FileLayouter struct {
	pc              *pc.ProtoCache
	toplevels       []string // top-level directories (detected)
	dir             string
	version         int // each time proto cache changes, version is incremented (updated by Save())
	current_version *trackedVersion
}

func New(c *pc.ProtoCache, dir string) *FileLayouter {
	res := &FileLayouter{pc: c}
	dir = strings.TrimSuffix(dir, "/")
	res.dir = dir
	return res
}

// returns the directory from which go compilers need to be run,
// e.g. "/tmp/protos/src"
func (f *FileLayouter) SrcDir() string {
	return f.TopDir() + FILELAYOUTER_PREFIX
}
func (f *FileLayouter) TopDir() string {
	res := f.dir
	if !strings.HasSuffix(res, "/") {
		res = res + "/"
	}
	return res
}

func (f *FileLayouter) Save(ctx context.Context) error {
	filelock.Lock()
	defer filelock.Unlock()
	new_version := f.current_version.Clone()
	dir := f.dir + "/" + FILELAYOUTER_PREFIX
	fmt.Printf("FILELAYOUTER - saving to dir \"%s\"\n", dir)
	err := common.RecreateSafely(f.dir)
	if err != nil {
		return err
	}

	// fill f.toplevels
	err = f.findTopLevels(ctx)

	for _, cf := range f.pc.Get(ctx) {
		tc := new_version.CurrentFile(cf)
		pf := cf.ProtoFile()
		path, err := f.resolvePath(ctx, pf)
		if err != nil {
			fmt.Printf("Cannot save proto \"%s\": %s\n", pf.Filename, err)
			continue
		}
		path = strings.Trim(path, "/")
		fname := filepath.Base(pf.Filename)
		fname = fmt.Sprintf("%s/%s/%s", dir, path, fname)
		//		fmt.Printf("Saving %s\n", fname)

		dir := filepath.Dir(fname)
		// debugging weird behaviour:
		/*
			if strings.Contains(dir, "googlecast") {
				fmt.Printf("fname=%s\ndir=%s\npath=%s\npf.filename=%s\n", fname, dir, path, pf.Filename)
				fmt.Printf("Creating dir \"%s\"\n", dir)
				panic("googlecast weirdness")
			}
		*/
		os.MkdirAll(dir, 0777)

		tc.filename = fname
		tc.prefix = filepath.Dir(path)
		tc.relativeDir = path

		err = utils.WriteFile(fname, []byte(pf.Content))
		if err != nil {
			return err
		}
		//		fmt.Printf("Saved \"%s\"\n", fname)
	}
	cl := len(new_version.Changes())
	fmt.Printf("Changes: %d\n", cl)
	if cl != 0 {
		f.version++
		new_version.version = f.version
		f.current_version = new_version
	}
	return nil
}
func (f *FileLayouter) CurrentVersionNumber() int {
	return f.version
}

// given a version number, returns all the protofiles that changed
func (f *FileLayouter) ChangedProtos(version int) []*TrackedChange {
	if f.current_version == nil {
		return nil
	}
	return f.current_version.ChangesSince(version)
}



































































