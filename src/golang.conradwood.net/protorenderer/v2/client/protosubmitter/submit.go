package protosubmitter

import (
	"fmt"
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/utils"
	"golang.yacloud.eu/yatools/gitrepo"
	"io"
	"os"
	"strings"
	"time"
)

const (
	DEBUG                = false
	PROTO_COMPILE_RESULT = "/tmp/proto_compile_result"
)

// CLI function, compile and submit directory or file to store
func SubmitProtos(dir string) error {
	return submit_protos_with_dir(dir, true)
}

// CLI function, compile but do not submit to store
func CompileProtos(dir string) error {
	return submit_protos_with_dir(dir, false)
}

func SubmitProtosGit() error {
	dir, err := find_git_dir()
	if err != nil {
		return err
	}
	return submit_protos_with_dir(dir, true)
}
func CompileProtosGit() error {
	dir, err := find_git_dir()
	if err != nil {
		return err
	}
	return submit_protos_with_dir(dir, false)
}

func find_git_dir() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	gr, err := gitrepo.NewGitRepo(path)
	if err != nil {
		return "", err
	}
	proto_dir := gr.Path() + "/protos/"
	return proto_dir, nil
}

/*
given a directory, all .proto files in that directory will be submitted to protorenderer.
*/
func submit_protos_with_dir(proto_dir string, save bool) error {
	d, err := IsDir(proto_dir)
	if err != nil {
		return err
	}
	if !d {
		return submit_proto_filenames([]string{proto_dir}, save)
		//		return errors.Errorf("\"%s\" is not a directory", proto_dir)
	}

	var filenames []string // relative to proto_dir
	err = utils.DirWalk(proto_dir, func(root, relfil string) error {
		if !strings.HasSuffix(relfil, ".proto") {
			return nil
		}
		filenames = append(filenames, relfil)
		return nil
	},
	)
	if err != nil {
		return err
	}
	return submit_proto_files(proto_dir, filenames, save)
}

// given absolute filename(s), will submit those to protorenderer
// if filenames are part of different git repositories, it will make multiple calls, one per git repository
func submit_proto_filenames(abs_filenames []string, save bool) error {
	repo_files := make(map[string][]string) // git repository directory -> filenames
	for _, fname := range abs_filenames {
		if !utils.FileExists(fname) {
			return errors.Errorf("file \"%s\" does not exist", fname)
		}
		gr, err := gitrepo.NewGitRepo(fname)
		if err != nil {
			return err
		}
		path := gr.Path()
		nf := strings.TrimPrefix(fname, path)
		idx := strings.Index(nf, "/protos/")
		if idx != -1 {
			path = path + nf[:idx+8]
			nf = nf[idx+8:]
		}
		x := repo_files[path]
		x = append(x, nf)
		repo_files[path] = x
	}
	for k, v := range repo_files {
		for _, fname := range v {
			debugf("Repo \"%s\": %s\n", k, fname)
		}
		err := submit_proto_files(k, v, save)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
given a list of filenames, relative to proto_dir, will send those to protorenderer server.
The result of the compilation will be stored in PROTO_COMPILE_RESULT.
if 'save' is true, the files will be submitted in to protorenderers store.
currently proto_dir MUST be a git repository, or a directory within it

the filenames must match the golang convention. That is stripped from prefixing paths.

for example: VALID: "golang.conradwood.net/apis/common/common.proto"

for example: NOT VALID: "protos/golang.conradwood.net/apis/common/common.proto"
*/
func submit_proto_files(proto_dir string, filenames []string, save bool) error {
	utils.RecreateSafely(PROTO_COMPILE_RESULT)
	ctx := authremote.ContextWithTimeout(time.Duration(1800) * time.Second)

	// repoid
	repoid := uint32(0)
	gr, err := gitrepo.NewYAGitRepo(proto_dir)
	if err != nil {
		debugf("Not a YAGitRepo: \"%s\"\n", proto_dir)
	} else {
		repoid = uint32(gr.RepositoryID())
	}
	srv, err := pb.GetProtoRenderer2Client().Submit(ctx)
	if err != nil {
		return err
	}
	debugf("RepositoryID: %d\n", repoid)
	so := &pb.SubmitOption{Save: save}
	err = srv.Send(&pb.FileTransfer{SubmitOption: so})
	if err != nil {
		return err
	} //	utils.Bail("failed to stream options to server", err)
	bs := utils.NewByteStreamSender(func(key, filename string) error {
		// start new file
		err := srv.Send(&pb.FileTransfer{Filename: filename, RepositoryID: repoid})
		return err
	},
		// send contents
		func(b []byte) error {
			err := srv.Send(&pb.FileTransfer{Data: b})
			return err
		},
	)

	for _, fname := range filenames {
		ct, err := utils.ReadFile(proto_dir + "/" + fname)
		if err != nil {
			return err
		}
		debugf("Submitting %s (%d bytes)\n", "protos/"+fname, len(ct))
		fname = strings.TrimPrefix(fname, "protos/")
		err = bs.SendBytes("foo", fname, ct)
		if err != nil {
			return err
		}
	}
	//	utils.Bail("failed to read proto files", err)
	err = srv.Send(&pb.FileTransfer{TransferComplete: true}) // switching to recv mode now
	if err != nil {
		return err
	} //	utils.Bail("failed to send files", err)

	// receiving the results now...
	bsr := utils.NewByteStreamReceiver(PROTO_COMPILE_RESULT)

	for {
		recv, err := srv.Recv()
		if recv != nil {
			/*
				if recv.Filename != "" {
					debugf("Receiving: filename=%s, bytes=%d\n", recv.Filename, len(recv.Data))
				}
			*/
			if recv.Result != nil {
				debugf("Failure for %s: %v\n", recv.Result.Filename, recv.Result.Filename)
				for _, result := range recv.Result.CompileResults {
					debugf("    compiler: \"%s\" (success=%v)\n", result.CompilerName, result.Success)
					debugf("    error: %s\n", result.ErrorMessage)
					debugf("    output: %s\n", result.Output)
				}
			}
			if recv.Output != nil {
				for _, line := range recv.Output.Lines {
					fmt.Printf("server: \"%s\"\n", line)
				}
			}

		}
		xerr := bsr.NewData(recv)
		if xerr != nil {
			return xerr
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			} //			utils.Bail("failed to receive files", err)
		}
	}
	debugf("Done\n")
	return nil
}

func debugf(format string, args ...interface{}) {
	if !DEBUG {
		return
	}
	fmt.Printf(format, args...)
}
