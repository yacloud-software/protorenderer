package migrate

import (
	"fmt"
	"golang.conradwood.net/apis/gitserver"
	pb "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/errors"
	// "golang.yacloud.eu/yatools/gitrepo"
)

func Fix(file *pb.FailedBridgeFile) error {
	// url for gitrepo
	ctx := authremote.Context()
	gitrepo, err := gitserver.GetGIT2Client().RepoByID(ctx, &gitserver.ByIDRequest{ID: file.RepositoryID})
	if err != nil {
		return err
	}
	url := ""
	if len(gitrepo.URLs) == 0 {
		return errors.Errorf("repository %d has no urls", file.RepositoryID)
	}
	fmt.Printf("URL: %s\n", url)

	return nil
}
