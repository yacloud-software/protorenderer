package srv

import (
	"context"
	"fmt"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/protorenderer/v2/interfaces"
	"path/filepath"
	"strings"
)

type check_compiler struct {
}

func (cc *check_compiler) ShortName() string {
	return "packagecheck"
}
func (cc *check_compiler) Compile(context.Context, interfaces.CompilerEnvironment, []interfaces.ProtoFile, string, interfaces.CompileResult) error {
	return nil
}

func check_if_files_match_packages(ctx context.Context, ce interfaces.CompilerEnvironment, scr interfaces.CompileResult, pfs []interfaces.ProtoFile) error {
	fake_comp := &check_compiler{}
	for _, pf := range pfs {
		fname := pf.GetFilename()
		fmt.Printf("[packagecheck] File: %s\n", fname)
		mi := ce.MetaCache().ByFilename(fname)
		if mi == nil {
			return errors.Errorf("no meta for \"%s\"", fname)
		}
		if mi.Package == "" {
			scr.AddFailed(fake_comp, pf, errors.Errorf("file \"%s\" is missing mandatory option package", fname), nil)
			continue
		}
		if mi.PackageGo == "" {
			scr.AddFailed(fake_comp, pf, errors.Errorf("file \"%s\" is missing mandatory option go_package", fname), nil)
			continue
		}
		if mi.PackageJava == "" {
			scr.AddFailed(fake_comp, pf, errors.Errorf("file \"%s\" is missing mandatory option java_package", fname), nil)
			continue
		}
		dir := filepath.Dir(fname)
		dir = filepath.Base(dir)
		if dir != mi.Package {
			fmt.Printf("package: \"%s\", dir: \"%s\"\n", mi.Package, dir)
			scr.AddFailed(fake_comp, pf, errors.Errorf("package \"%s\" not compatible with filename \"%s\"", mi.Package, fname), nil)
			continue
		}

		if !strings.HasPrefix(fname, mi.PackageGo) {
			scr.AddFailed(fake_comp, pf, errors.Errorf("go package \"%s\" not compatible with filename \"%s\"", mi.PackageGo, fname), nil)
			continue
		}
		/*
					 java  package is really rubbish.
			golang.conradwood.net/apis/test/test.proto could be any of
			test.apis.net.conradwood.golang
			net.conradwood.golang.apis.test
		*/

		/*
			jpkg_to_file := java_package_name_reverse(mi.PackageJava)
			jpkg_to_file = strings.ReplaceAll(jpkg_to_file, ".", "/")
			if !strings.HasPrefix(fname, jpkg_to_file) {
				fmt.Printf("jpkg_to_file: %s\n", jpkg_to_file)
				scr.AddFailed(fake_comp, pf, errors.Errorf("java package \"%s\" not compatible with filename \"%s\"", mi.PackageJava, fname), nil)
			continue
			}
		*/

	}
	return nil
}

func java_package_name_reverse(pkg string) string {
	res := ""
	for _, element := range strings.Split(pkg, ".") {
		res = element + "." + res
	}
	return res
}
