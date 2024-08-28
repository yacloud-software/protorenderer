package protosubmitter

import (
	"flag"
	"fmt"
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/utils"
	"golang.yacloud.eu/yatools/gitrepo"
	"io"
	"os"
	//	"path/filepath"
	"strings"
	"time"
)

const (
	PROTO_COMPILE_RESULT = "/tmp/proto_compile_result"
)

var (
	debug_proto_submitter = flag.Bool("debug_protosubmitter", false, "true debug submit code")
)

type ProtoSubmitter interface {
	SubmitProtos(dir_or_files []string) error // files or dirs as provided
	CompileProtos(dir_or_files []string) error
}
type protoSubmitter struct {
	printer   func(txt string)
	git_cache map[string]*gitrepo.YAGitRepo // path->repoid
}

func New() ProtoSubmitter {
	return &protoSubmitter{git_cache: make(map[string]*gitrepo.YAGitRepo)}
}
func NewWithOutput(print func(txt string)) ProtoSubmitter {
	ps := New().(*protoSubmitter)
	ps.printer = print
	return ps
}

// CLI function, compile and submit directory or file to store
func (ps *protoSubmitter) SubmitProtos(dir_or_files []string) error {
	return ps.submit_protos_with_dir(dir_or_files, true)
}

// CLI function, compile but do not submit to store
func (ps *protoSubmitter) CompileProtos(dir_or_files []string) error {
	return ps.submit_protos_with_dir(dir_or_files, false)
}

/*
func (ps *protoSubmitter) find_git_dir() (string, error) {
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
*/
/*
resolve dirs to files and submit the final file
*/
func (ps *protoSubmitter) submit_protos_with_dir(proto_dir_or_files []string, save bool) error {
	// check if each file exists
	for _, df := range proto_dir_or_files {
		if !utils.FileExists(df) {
			return errors.Errorf("file \"%s\" does not exiset", df)
		}
	}
	curdir, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err)
	}

	var filenames []*rel_file_in_path
	for _, df := range proto_dir_or_files {
		if len(df) == 0 {
			return errors.Errorf("Empty file submitted")
		}
		if df[0] == '/' {
			return errors.Errorf("Absolute filename submitted (%s)", df)
		}
		d, err := IsDir(df)
		if err != nil {
			return err
		}
		// it's a file
		if !d {
			filenames = append(filenames, &rel_file_in_path{path: curdir, filename: df}) // relative file submitted
			continue
		}
		// it's a dir
		root := curdir + "/" + df
		err = utils.DirWalk(root, func(root, relfil string) error {
			if !strings.HasSuffix(relfil, ".proto") {
				return nil
			}
			filenames = append(filenames, &rel_file_in_path{path: root, filename: relfil})
			return nil
		},
		)
		if err != nil {
			return err
		}
	}
	ps.debugf("finding repos for %d files\n", len(filenames))
	err = ps.find_repos_for_files(filenames)
	if err != nil {
		return err
	}

	t := utils.Table{}
	t.AddHeaders("path", "filename", "gitrepoid", "gitrepopath")
	for _, f := range filenames {
		t.AddString(f.path)
		t.AddString(f.filename)
		t.AddUint64(f.repositoryid)
		t.AddString(f.gitrepopath)
		t.NewRow()
	}

	fmt.Println(t.ToPrettyString())

	fmap, err := ps.build_map_of_repoids(filenames)
	if err != nil {
		return err
	}
	for k, v := range fmap {
		for _, fname := range v {
			ps.debugf("Repo \"%d\": %s\n", k, fname.filename)
		}
		err := ps.submit_proto_files(v, save)
		if err != nil {
			return err
		}
	}
	return nil
	// return ps.submit_proto_files(filenames, save)
}

// fill in repositoryid for each filename
func (ps *protoSubmitter) find_repos_for_files(filenames []*rel_file_in_path) error {
	for _, f := range filenames {
		for path, gr := range ps.git_cache {
			if strings.HasPrefix(f.abs(), path) {
				f.repositoryid = gr.RepositoryID()
				f.gitrepopath = gr.Path()
				return nil
			}
		}
		ps.debugf("need Repo for \"%s\"\n", f.abs())
		gr, err := gitrepo.NewYAGitRepo(f.abs())
		if err != nil {
			return err
		}
		f.repositoryid = gr.RepositoryID()
		f.gitrepopath = gr.Path()
		ps.git_cache[gr.Path()] = gr
		ps.debugf("repo \"%s\" = %d\n", gr.Path(), gr.RepositoryID())
	}
	return nil
}

// return a map with arrays of files for each repositoryid
func (ps *protoSubmitter) build_map_of_repoids(filenames []*rel_file_in_path) (map[uint64][]*rel_file_in_path, error) {
	res := make(map[uint64][]*rel_file_in_path)
	for _, rf := range filenames {
		res[rf.repositoryid] = append(res[rf.repositoryid], rf)
	}
	return res, nil
}

// given absolute filename(s), will submit those to protorenderer
// if filenames are part of different git repositories, it will make multiple calls, one per git repository
/*
func (ps *protoSubmitter) submit_proto_filenames(abs_filenames []string, save bool) error {
	repo_files := make(map[string][]string) // git repository directory -> filenames
	for _, fname := range abs_filenames {
		if len(fname) == 0 {
			return errors.Errorf("no filename provided")
		}
		if fname[0] != '/' {
			afname, err := filepath.Abs(fname)
			if err != nil {
				return errors.Errorf("filename \"%s\" is not absolute (%s)", fname, err)
			}
			fname = afname
		}
		if !utils.FileExists(fname) {
			return errors.Errorf("file \"%s\" does not exist", fname)
		}
		path := "http://nopath"
		gr, err := gitrepo.NewGitRepo(fname)
		if err == nil {
			path = gr.Path()
		}
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
			ps.debugf("Repo \"%s\": %s\n", k, fname)
		}
		err := ps.submit_proto_files(k, v, save)
		if err != nil {
			return err
		}
	}
	return nil
}
*/
/*
given a list of filenames, relative to proto_dir, will send those to protorenderer server.
The result of the compilation will be stored in PROTO_COMPILE_RESULT.
if 'save' is true, the files will be submitted in to protorenderers store.
currently proto_dir MUST be a git repository, or a directory within it

the filenames must match the golang convention. That is stripped from prefixing paths.

for example: VALID: "golang.conradwood.net/apis/common/common.proto"

for example: NOT VALID: "protos/golang.conradwood.net/apis/common/common.proto"

all files MUST have the same repositoryid
*/
func (ps *protoSubmitter) submit_proto_files(filenames []*rel_file_in_path, save bool) error {
	if len(filenames) == 0 {
		return nil
	}
	utils.RecreateSafely(PROTO_COMPILE_RESULT)
	ctx := authremote.ContextWithTimeout(time.Duration(1800) * time.Second)
	//	ps.debugf("proto-dir: %s\n", proto_dir)
	for _, f := range filenames {
		ps.debugf("submitting: \"%s\"\n", f.abs())
	}

	// repoid
	repoid := uint32(filenames[0].repositoryid)

	srv, err := pb.GetProtoRenderer2Client().Submit(ctx)
	if err != nil {
		return err
	}
	ps.debugf("RepositoryID: %d\n", repoid)
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
		ct, err := utils.ReadFile(fname.abs())
		if err != nil {
			return errors.Wrap(err)
		}
		ps.debugf("Submitting %s (%d bytes)\n", fname.abs(), len(ct))
		//		fname = strings.TrimPrefix(fname, "protos/")
		err = bs.SendBytes("foo", fname.filename, ct)
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
					ps.debugf("Receiving: filename=%s, bytes=%d\n", recv.Filename, len(recv.Data))
				}
			*/
			if recv.Result != nil {
				ps.debugf("Failure for %s: %v\n", recv.Result.Filename, recv.Result.Filename)
				for _, result := range recv.Result.CompileResults {
					ps.debugf("    compiler: \"%s\" (success=%v)\n", result.CompilerName, result.Success)
					ps.debugf("    error: %s\n", result.ErrorMessage)
					ps.debugf("    output: %s\n", result.Output)
				}
			}
			if recv.Output != nil {
				for _, line := range recv.Output.Lines {
					pf := ps.printer
					if pf != nil {
						pf(fmt.Sprintf("server: \"%s\"\n", line))
					}
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
	ps.debugf("Done\n")
	return nil
}

func (ps *protoSubmitter) debugf(format string, args ...interface{}) {
	if !*debug_proto_submitter {
		return
	}
	fmt.Printf(format, args...)
}

type rel_file_in_path struct {
	path         string // e.g. /home/cnw/devel/go/protorenderer/tests/02_test/protos
	filename     string // e.g. golang.conradwood.net/apis/common/common.proto
	repositoryid uint64 // the repo the file is in (or 0)
	gitrepopath  string // e.g. /home/cnw/devel/go/protorenderer (or "")
}

func (rfp *rel_file_in_path) abs() string {
	path := strings.TrimSuffix(rfp.path, "/")
	filename := strings.TrimPrefix(rfp.filename, "/")
	return path + "/" + filename
}
