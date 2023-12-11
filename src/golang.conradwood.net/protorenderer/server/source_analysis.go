package main

import (
	"context"
	//	"flag"
	//	"fmt"
	"golang.conradwood.net/apis/common"
	pr "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/cache"
	//	"golang.conradwood.net/go-easyops/errors"
	"sync"
	"time"
)

var (
	packages  []*pr.Package
	msgids    = cache.NewResolvingCache("msgids", time.Duration(999999)*time.Hour, 9999999)
	msgidctr  = 1
	msgidlock sync.Mutex
)

func (p *protoRenderer) SubmitSource(ctx context.Context, req *pr.ProtocRequest) (*common.Void, error) {
	err := current.metaCompiler.SubmitSource(ctx, req)
	if err != nil {
		return nil, err
	}

	return &common.Void{}, nil
}

func (p *protoRenderer) GetIDForPackage(ctx context.Context, req *pr.PackageIDRequest) (*pr.ID, error) {
	return &pr.ID{ID: "package-id"}, nil
}
func (p *protoRenderer) GetIDForMessage(ctx context.Context, req *pr.MessageIDRequest) (*pr.ID, error) {
	return &pr.ID{ID: "package-id"}, nil
}























































