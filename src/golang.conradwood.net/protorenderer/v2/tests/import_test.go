package tests

import (
	"fmt"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/v2/client/protoretriever"
	"golang.conradwood.net/protorenderer/v2/client/protosubmitter"
	"golang.conradwood.net/protorenderer/v2/helpers"
	"path/filepath"
	"strings"
	"testing"
)

//TODO: test submission of individual files from different repositories

// test on-the-fly-compile for directories
func TestOnTheFlyCompileDirs(t *testing.T) {
	tests := []struct {
		dir      string
		expected map[string]int
	}{
		//		{"tests/04_test", map[string]int{"info": 1, "go": 1, "java": 1, "class": 6}}, // re-check to make sure previous tests do not influence it
		{"tests/01_test", map[string]int{"info": 1, "go": 1, "java": 1, "class": 6}},
		{"tests/02_test", map[string]int{"info": 12, "go": 12, "java": 66, "class": 279}},
		{"tests/01_test", map[string]int{"info": 1, "go": 1, "java": 1, "class": 6}}, // re-check to make sure
	}
	for _, ts := range tests {
		t.Run(ts.dir, func(tst *testing.T) {
			testcompile(tst, ts.dir, ts.dir, ts.expected)
		})
	}
}
func TestOnTheFlyCompileFiles(t *testing.T) {
	testcompile(t, "a", "tests/01_test/protos/golang.conradwood.net/apis/test/test.proto", map[string]int{"info": 1, "go": 1, "java": 1, "class": 6})
}

func TestSubmit(t *testing.T) {
	var proto_file string
	var proto_dir string

	proto_file = "golang.conradwood.net/apis/test/test.proto"
	proto_dir = "tests/01_test/protos/"
	test_submit_file(t, proto_dir, proto_file, map[string]int{"info": 1, "proto": 1, "go": 1, "java": 1, "class": 6})

	proto_file = "golang.conradwood.net/apis/test2/test2.proto"
	proto_dir = "tests/03_test/protos/"
	test_submit_file(t, proto_dir, proto_file, map[string]int{"info": 1, "proto": 1, "go": 1, "java": 1, "class": 6})
}

func TestPreviousResults(t *testing.T) {
	tdata := []struct {
		name         string
		java_package string
	}{
		{"works", ""},      //works
		{"fails", "3 sdf"}, // fails
	}
	for _, td := range tdata {
		t.Run(td.name, func(tlocal *testing.T) {
			test_prev_res(tlocal, td.java_package)
		})
	}
}

func test_prev_res(t *testing.T, javapackage string) {
	pfb := helpers.NewProtoFileBuilder("testdyn")
	if javapackage != "" {
		pfb.SetJavaPackage(javapackage)
	}
	b := pfb.Bytes()
	td := t.TempDir()
	write_fake_git_repo(td)
	fname := pfb.GetFilename()
	ffname := td + "/" + fname

	err := utils.WriteFileCreateDir(ffname, b)
	if err != nil {
		t.Fatalf("failed to write: %s", err)
		return
	}
	//	t.Logf("proto:\n%s", string(b))
	test_submit_file(t, td+"/", fname, nil)

}

func test_submit_file(t *testing.T, proto_dir, proto_file string, expected map[string]int) {
	var err error
	fname := proto_dir + proto_file
	fname, err = utils.FindFile(fname)
	if err != nil {
		t.Fatalf("file not found: %s", err)
	}
	err = protosubmitter.SubmitProtos(fname)
	if err != nil {
		t.Logf("failed: %s\n", errors.ErrorStringWithStackTrace(err))
		t.Fatalf("failed to submit file %s: %s", fname, err)
		return
	}

	// now retrieve the file again
	ctx := authremote.Context()
	tmpdir := t.TempDir()
	write_fake_git_repo(tmpdir)
	//	tmpdir = "/tmp/protorenderer_tests"
	err = protoretriever.ByFilename(ctx, proto_file, tmpdir)
	if err != nil {
		t.Fatalf("failed to retrieve files: %s\n", err)
		return
	}
	err = check_dir_against_expected(tmpdir, expected)
	if err != nil {
		t.Fatalf("result mismatch for %s: %s\n", proto_dir, err)
		return
	}
}

func testcompile(t *testing.T, testname, dir string, expected map[string]int) {
	run_test(t, testname, dir, expected, false)
}

// either compile or submit files in a directory and compare with result
func run_test(t *testing.T, testname, dir string, expected map[string]int, save bool) {
	fname, err := utils.FindFile(dir + "/protos")
	if err != nil {
		fname, err = utils.FindFile(dir)
	}
	if err != nil {
		t.Fatalf("failed to find dir: %s", err)
		return
	}
	if save {
		err = protosubmitter.SubmitProtos(fname)
	} else {
		err = protosubmitter.CompileProtos(fname)
	}
	if err != nil {
		t.Logf("failed: %s\n", errors.ErrorStringWithStackTrace(err))
		t.Fatalf("failed to submit dir %s: %s", dir, err)
		return
	}

	t.Logf("comparing result with expected...\n")
	err = check_dir_against_expected(protosubmitter.PROTO_COMPILE_RESULT, expected)
	if err != nil {
		t.Fatalf("Result mismatch for %s: %s\n", dir, err)
	}
}

// give a directory and a map of extensions and count it compares the two
// if map is nil or empty, it always returns true
func check_dir_against_expected(dir string, expected map[string]int) error {
	if expected == nil || len(expected) == 0 {
		return nil
	}
	filenames, err := helpers.FindFiles(dir)
	if err != nil {
		return err
	}
	// built map by extension
	exts := make(map[string]int)
	for _, f := range filenames {
		e := filepath.Ext(f)
		if e == ".goeasyops-dir" {
			continue
		}
		i := exts[e]
		i++
		exts[e] = i
	}
	failed := false
	errstr := ""
	for e_k, e_v := range expected {
		e_k = "." + e_k
		got := exts[e_k]
		if got != e_v {
			failed = true
			errstr = fmt.Sprintf("Extension %s, expected %d files, got %d files instead", e_k, e_v, got)
			failed_filenames, err := helpers.FindFiles(dir, e_k)
			if err == nil {
				fmt.Printf("Files with extension \"%s\":\n", e_k)
				for _, fname := range failed_filenames {
					fmt.Printf("  file: \"%s\"\n", fname)
				}
			}
		}
		//	t.Logf("%s Extension \"%s\": %d files, expected %d\n", s, e_k, got, e_v)
	}
	if failed {
		return fmt.Errorf("extension mismatch on dir %s: %s", dir, errstr)
	}
	for ext, ct := range exts {
		if expected[ext] != ct {
			xext := strings.TrimPrefix(ext, ".")
			if expected[xext] != ct {
				return errors.Errorf("found %d files with extension \"%s\", but expected \"%d\" files with that extension", ct, ext, expected[ext])
			}
		}

	}
	return nil

}
