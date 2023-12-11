package filelayouter

import (
	"context"
	"fmt"
	pr "golang.conradwood.net/apis/protorenderer"
	"path/filepath"
	"strings"
)

/*
the layout of a directory containing .proto files is strictly defined.
Things to consider:
1) The "import" statements refer to files (not package/protos)
2) The java package must match the filename

Here we try to derive from the currently known protos what path a proto should be saved in
We look, for example, at import statements of other protos - if they refer to the proto
in question, we can deduce the path from that

if cannot determine by import path, we look at the java package to see if that matches
an existing top-level directory

it will return the path, e.g. "golang.conradwood.net/apis/common"
*/
func (f *FileLayouter) resolvePath(ctx context.Context, pf *pr.ProtoFile) (string, error) {
	if pf.Filename == "" {
		return "", fmt.Errorf("cannot resolve a path without filename")
	}
	if pf.JavaPackage == "" {
		return "", fmt.Errorf("Currently we need a java package defined for proto %s", pf.Filename)
	}

	var references []string // import statements referring to this protofile
	for _, cf := range f.pc.Get(ctx) {
		pfr := cf.ProtoFile()
		for _, i := range pfr.Imports {
			//				fmt.Printf("import statement in %s: %s\n", pfr.Filename, i)
			fn := filepath.Base(i)
			fname := filepath.Base(pf.Filename)
			if fn != fname {
				//			fmt.Printf("  %s != %s\n", fn, pf.Filename)
				continue
			}
			// same filename - is the import statements' toplevel dir
			// contained in "filename"? if so, that's the logical toplevel
			import_tld := tld(i)
			if strings.Contains(fname, import_tld) {
				if doDebug(pf) {
					fmt.Printf("import_tld = \"%s\", fname = \"%s\"\n", import_tld, fname)
				}
			}
			/*
				pd := filepath.Base(filepath.Dir(i))
				if pd != pf.GoPackage {
					fmt.Printf("  %s != %s\n", pd, pfr.GoPackage)
					continue
				}
			*/
			references = append(references, i)
		}
	}
	//	fmt.Printf("%d import statements for: %s\n", len(references), pf.Filename)
	s := ""
	mixed := false
	for _, r := range references {
		if s == "" {
			s = r
			continue
		}
		if s != r {
			mixed = true
			break
		}
	}
	if mixed {
		return "", fmt.Errorf("filelayouter - Unable to determine proto path.")
	}
	if s != "" {
		s = filepath.Dir(s)
		if doDebug(pf) {
			fmt.Printf("Proto (with at least one reference) %s in dir \"%s\"\n", pf.Filename, s)
		}
		f.addTopLevelFromDir(s)
		return s, nil
	}
	if doDebug(pf) {
		fmt.Printf("%s: no import statements to go by, using \"toplevels\"\n", pf.Filename)
		for _, r := range f.toplevels {
			fmt.Printf("Toplevel: \"%s\"\n", r)
		}
	}

	//	fmt.Printf("filename: \"%s\"\n", pf.Filename)
	for _, r := range f.toplevels {
		//		fmt.Printf("  tld: \"%s\"\n", r)
		pos := strings.Index(pf.Filename, r+"/")
		if pos == -1 {
			continue
		}
		s = pf.Filename[pos:]
		s = filepath.Dir(s)
		fmt.Printf("%s: returning toplevel %s\n", pf.Filename, s)
		return s, nil
	}

	res := filepath.Dir(pf.Filename)

	fmt.Printf("BAD (%s -> %s)\n", pf.Filename, res)
	return res, nil

}

func (f *FileLayouter) findTopLevels(ctx context.Context) error {
	for _, cf := range f.pc.Get(ctx) {
		pf := cf.ProtoFile()
		_, err := f.resolvePath(ctx, pf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *FileLayouter) addTopLevelFromDir(dir string) {
	s := tld(dir)
	for _, x := range f.toplevels {
		if x == s {
			return
		}
	}
	f.toplevels = append(f.toplevels, s)
	return
}

func tld(dir string) string {
	sx := strings.Split(dir, "/")
	if len(sx) == 0 {
		return dir
	}
	s := sx[0]
	return s
}

func doDebug(pf *pr.ProtoFile) bool {
	if strings.Contains(pf.Filename, "googlecast") {
		return true
	}
	return false
}












