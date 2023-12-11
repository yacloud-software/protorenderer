package main

import (
	"flag"
	"fmt"
	"golang.conradwood.net/apis/common"
	pr "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/utils"
	"strings"
)

func View() {
	as := flag.Args()
	filter := ""
	if len(as) != 0 {
		filter = strings.ToLower(as[0])
	}
	ps := pr.GetProtoRendererServiceClient()
	ctx := getContext()
	pkgs, err := ps.GetPackages(ctx, &common.Void{})
	utils.Bail("failed to get packages", err)
	for _, fpkg := range pkgs.Packages {
		if filter != "" && !strings.Contains(strings.ToLower(fpkg.Name), filter) {
			continue
		}
		ctx := getContext()
		pkg, err := ps.GetPackageByID(ctx, &pr.ID{ID: fpkg.ID})
		utils.Bail("Failed to get package", err)
		fmt.Printf("Package: #%s %s (prefix=%s)\n", pkg.ID, pkg.Name, pkg.Prefix)
		for _, svc := range pkg.Services {
			fmt.Printf("    Service: #%s %s (Repository %d)\n", svc.ID, svc.Name, svc.RepositoryID)
			for _, rpc := range svc.RPCs {
				ds := ""
				if rpc.Deprecated {
					ds = "DEPRECATED "
				}
				fmt.Printf("           RPC: %s#%s %s(%s) %s\n", ds, rpc.ID, rpc.Name, msgLine(rpc.Input), msgLine(rpc.Output))
			}
		}
	}
}
func msgLine(msg *pr.Message) string {
	if msg == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s[#%s]", msg.ID, msg.Name)
}








































































