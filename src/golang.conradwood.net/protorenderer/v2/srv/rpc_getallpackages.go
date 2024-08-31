package srv

import (
	"context"
	"fmt"
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/protorenderer2"
)

func (e *protoRenderer) GetAllPackages(ctx context.Context, req *common.Void) (*pb.FlatPackageList, error) {
	fmt.Printf("Getting all packages...\n")
	res := &pb.FlatPackageList{}
	err := CompileEnv.MetaCache().AllPackages(func(pfi *pb.ProtoFileInfo) {
		fp := &pb.FlatPackage{
			ID:           "getallpackages_has_no_id",
			ShortName:    pfi.Package,
			FQDN:         pfi.PackageGo,
			RepositoryID: pfi.ProtoFile.RepositoryID,
		}
		res.Packages = append(res.Packages, fp)
		//
	})
	if err != nil {
		return nil, err
	}
	fmt.Printf("Got %d packages\n", len(res.Packages))
	return res, nil
}
