package compiler

import (
	"fmt"
	"golang.conradwood.net/go-easyops/utils"
	"io/ioutil"
	"strings"
)

func AllJava(dir string) ([]string, error) {
	return AllFiles(dir, ".java")
}

func AllProtos(dir string) ([]string, error) {
	return AllFiles(dir, ".proto")
}

// suffix might be "" (empty string) to include all files in dir
func AllFiles(dir string, suffix string) ([]string, error) {
	dir = strings.TrimSuffix(dir, "/")
	files, err := addDir(dir, suffix)
	if err != nil {
		return nil, err
	}

	// strip prefix
	dir = dir + "/"
	for i, _ := range files {
		fn := files[i]
		if !strings.HasPrefix(fn, dir) {
			return nil, fmt.Errorf("Invalid filename \"%s\" - does not start with \"%s\"", fn, dir)
		}
		fn = strings.TrimPrefix(fn, dir)
		files[i] = fn
	}
	/*
		for _, f := range files {
			fmt.Printf("File: \"%s\"\n", f)
		}
	*/
	return files, nil
}

func addDir(dir string, suffix string) ([]string, error) {
	if !utils.FileExists(dir) {
		return make([]string, 0), nil
	}
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var res []string
	for _, fi := range fis {
		if !fi.IsDir() {
			if !strings.HasSuffix(fi.Name(), suffix) {
				continue
			}
			res = append(res, dir+"/"+fi.Name())
			continue
		}

		d, err := addDir(dir+"/"+fi.Name(), suffix)
		if err != nil {
			return nil, err
		}
		res = append(res, d...)

	}
	return res, nil
}
































































