package protocache

import (
	"context"
	"fmt"
	pr "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/errors"
	"strings"
	"sync"
)

var (
	lock sync.Mutex
)

/*
* currently very simple in-memory proto cache
 */

// do not instantiate. use New()
type ProtoCache struct {
	version int // updated each time a file is added/removed
	protos  []*CachedFile
}

type CachedFile struct {
	version   int // incremented each time this file changes
	protoFile *pr.ProtoFile
}

func (c *CachedFile) Version() int {
	return c.version
}
func (c *CachedFile) ProtoFile() *pr.ProtoFile {
	return c.protoFile
}
func (c *CachedFile) Inc() {
	c.version++
}

func New() *ProtoCache {
	res := &ProtoCache{version: 1}
	return res
}
func (p *ProtoCache) Version() int {
	return p.version
}
func (p *ProtoCache) GetFile(ctx context.Context, filename string) (*pr.ProtoFile, error) {
	lf := 0
	var res *CachedFile
	for _, p := range p.protos {
		if p.protoFile.Filename != filename {
			continue
		}
		if p.version < lf {
			continue
		}
		res = p
		lf = p.version
	}
	if res == nil {
		return nil, fmt.Errorf("file \"%s\" not found", filename)
	}
	return res.protoFile, nil
}
func (p *ProtoCache) Delete(ctx context.Context, filename string) (int, error) {
	if filename == "" {
		return 0, fmt.Errorf("missing filename")
	}
	lock.Lock()
	defer lock.Unlock()
	var n []*CachedFile
	found := false
	for _, cf := range p.protos {
		if cf.protoFile.Filename == filename {
			found = true
			continue
		}
		n = append(n, cf)
	}
	if !found {
		return 0, errors.InvalidArgs(ctx, "file not found", "file \"%s\" not found", filename)
	}
	p.protos = n
	p.version++
	return p.version, nil
}

func (p *ProtoCache) AddOrUpdate(protofile *pr.ProtoFile) (int, error) {
	if protofile.Filename == "" {
		return 0, fmt.Errorf("Not adding protofile - missing name\n")
	}
	if !strings.Contains(protofile.Filename, "/") {
		return 0, fmt.Errorf("filename must contain a directory (at least one '/'). (not \"%s\")", protofile.Filename)
	}
	lock.Lock()
	defer lock.Unlock()
	for _, cf := range p.protos {
		pf := cf.protoFile
		if pf.GoPackage == protofile.GoPackage && pf.Filename == protofile.Filename {
			if pf.Content == protofile.Content {
				// same content - do nothing. Especially, do not increase version number
				return p.version, nil
			}
			cf.protoFile = protofile
			cf.Inc()
			p.version++
			return p.version, nil
		}
	}
	cf := &CachedFile{protoFile: protofile}
	p.protos = append(p.protos, cf)
	p.version++
	return p.version, nil
}
func (p *ProtoCache) Get(ctx context.Context) []*CachedFile {
	return p.protos
}






















































