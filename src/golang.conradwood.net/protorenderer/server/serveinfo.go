package main

import (
	"context"
	"flag"
	"fmt"
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/protorenderer/meta"
	"path/filepath"
	"sort"
	"strings"
)

var (
	add_repo_to_flat = flag.Bool("add_repo_to_flat_package", true, "if true adds repositoryid to flat package")
)

func (p *protoRenderer) GetPackageByName(ctx context.Context, req *pb.PackageName) (*pb.Package, error) {
	pfr, err := p.FindPackageByName(ctx, req)
	if err != nil {
		return nil, err
	}
	if pfr.Exists {
		return pfr.Package, nil
	}
	return nil, errors.NotFound(ctx, "package \"%s\" not found", req.PackageName)
}
func (p *protoRenderer) FindPackageByName(ctx context.Context, req *pb.PackageName) (*pb.PackageFindResult, error) {
	res := &pb.PackageFindResult{Exists: false}
	ev := NeedVersion(ctx)
	if ev != nil {
		return nil, ev
	}
	result := completeVersion.metaCompiler.GetMostRecentResult()
	if result == nil {
		return nil, errors.Unavailable(ctx, "GetPackages (most recent result)")
	}
	var err error
	for _, pk := range result.Packages {
		if pk.FQDN == req.PackageName {
			res.Exists = true
			res.Package, err = p.GetPackageByID(ctx, &pb.ID{ID: pk.Proto.ID})
		}
	}
	return res, err
}
func (p *protoRenderer) GetPackages(ctx context.Context, req *common.Void) (*pb.FlatPackageList, error) {
	ev := NeedVersion(ctx)
	if ev != nil {
		return nil, ev
	}
	result := completeVersion.metaCompiler.GetMostRecentResult()
	if result == nil {
		return nil, errors.Unavailable(ctx, "GetPackages (most recent result)")
	}
	res := &pb.FlatPackageList{}
	for _, p := range result.Packages {
		pf := p.FQDN
		pf = filepath.Dir(pf)
		pp := &pb.FlatPackage{ID: p.Proto.ID, Name: p.Name, Prefix: pf, Filename: p.Filename}
		if pp.ID == "" {
			fmt.Printf("WARNING package %s - %s has no ID!!\n", pp.Prefix, pp.Name)
		}
		if *add_repo_to_flat {
			dbpf, err := FindByName(ctx, pp.Filename)
			if err == nil && dbpf != nil {
				pp.RepositoryID = dbpf.RepositoryID
			}
		}
		res.Packages = append(res.Packages, pp)
	}
	sort.Slice(res.Packages, func(i, j int) bool {
		return res.Packages[i].Name < res.Packages[j].Name
	})
	return res, nil
}
func (p *protoRenderer) GetPackageByID(ctx context.Context, req *pb.ID) (*pb.Package, error) {
	ev := NeedVersion(ctx)
	if ev != nil {
		return nil, ev
	}
	result := completeVersion.metaCompiler.GetMostRecentResult()
	if result == nil {
		return nil, errors.Unavailable(ctx, "GetPackageByID (most recent result)")
	}
	if req.ID == "" {
		return nil, errors.InvalidArgs(ctx, "invalid package id", "invalid package id \"%s\"", req.ID)
	}
	var rpkg *meta.Package
	for _, p := range result.Packages {
		if p.Proto.ID == req.ID {
			rpkg = p
			break
		}
	}
	if rpkg == nil {
		return nil, errors.InvalidArgs(ctx, "no such package id", "no such package id \"%s\"", req.ID)
	}
	res := &pb.Package{
		ID:     rpkg.Proto.ID,
		Name:   rpkg.Name,
		Prefix: rpkg.FQDN,
	}
	dbfile, err := FindByName(ctx, rpkg.Filename)
	repoid := uint64(0)
	if err != nil {
		return nil, err
	}
	if dbfile != nil {
		repoid = dbfile.RepositoryID
	}
	for _, s := range rpkg.Services {
		svc := &pb.Service{ID: s.ID, Name: s.Name, Comment: s.Comment, RepositoryID: repoid}
		res.Services = append(res.Services, svc)
		for _, r := range s.RPCs {
			rpc := &pb.RPC{
				ID:      r.ID,
				Name:    r.Name,
				Comment: r.Comment,
				Input:   messageToProto(r.Input),
				Output:  messageToProto(r.Output),
			}
			if strings.Contains(rpc.Comment, "DEPRECATED") {
				rpc.Deprecated = true
			}
			svc.RPCs = append(svc.RPCs, rpc)
		}
	}
	for _, m := range rpkg.Messages {
		msg := messageToProto(m)
		res.Messages = append(res.Messages, msg)
	}
	return res, nil
}

func messageToProto(m *meta.Message) *pb.Message {
	if m == nil {
		return nil
	}
	res := &pb.Message{
		ID:        m.ID,
		Name:      m.Name,
		Comment:   m.Comment,
		PackageID: m.Package.Proto.ID,
	}
	for _, f := range m.Fields {
		pf := &pb.Field{
			ID:       f.ID,
			Comment:  f.Comment,
			Name:     f.Name,
			Type:     fmt.Sprintf("%v", f.Type),
			Repeated: f.Repeated,
			Required: f.Required,
			Optional: f.Optional,
		}
		if f.Message != nil {
			pf.MessageID = f.Message.ID
			pf.MessageName = f.Message.Name
		}
		res.Fields = append(res.Fields, pf)
	}
	return res
}








































































