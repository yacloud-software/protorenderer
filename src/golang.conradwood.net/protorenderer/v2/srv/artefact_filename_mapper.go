package srv

import (
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/protorenderer/v2/helpers"
	"golang.conradwood.net/protorenderer/v2/interfaces"
	"path/filepath"
	"strings"
)

/*
given a compiler name and a protofile will return the filenames in the store of the generated artefacts.
the filenames returned are relative to the storedir
this is probably not ideal as is. Problems we might run into:

- a new version of compiler places files in different directories. Remedy: record artefacts at compile-time
- might not be possible to find _all_ files for a given protofilename
- what means _all_ files anyway? (especially java is annoying like this)
*/
func get_artefact_filenames_for(ce interfaces.CompilerEnvironment, compilername, protofilename string) ([]string, error) {
	storedir := ce.StoreDir()
	if compilername == "golang" {
		base := filepath.Dir(protofilename)
		dir := storedir + "/" + compilername + "/" + base
		files, err := helpers.FindFiles(dir, ".go")
		if err != nil {
			return nil, err
		}
		var res []string
		for _, f := range files {
			res = append(res, compilername+"/"+base+"/"+f)
		}
		return res, nil
	} else if compilername == "java" {
		mi := ce.MetaCache().ByFilename(protofilename)
		if mi == nil {
			return nil, errors.Errorf("no metacache for \"%s\"", protofilename)
		}
		java_package := mi.PackageJava
		if java_package == "" {
			return nil, errors.Errorf("no javapackage for \"%s\"", protofilename)
		}
		java_dir := strings.ReplaceAll(java_package, ".", "/")
		//fmt.Printf("Java package: \"%s\" => %s\n", java_package, java_dir)
		base := java_dir
		files, err := helpers.FindFiles(storedir+"/"+compilername+"/src/"+java_dir, ".java")
		if err != nil {
			return nil, err
		}
		var res []string
		for _, f := range files {
			res = append(res, compilername+"/src/"+base+"/"+f)
		}

		files, err = helpers.FindFiles(storedir+"/"+compilername+"/classes/"+java_dir, ".class")
		if err != nil {
			return nil, err
		}

		for _, f := range files {
			res = append(res, compilername+"/classes/"+base+"/"+f)
		}

		return res, nil
	} else {
		return nil, errors.Errorf("unknown compiler \"%s\" whilst trying to retrieve artefacts for \"%s\"", compilername, protofilename)
	}
}
