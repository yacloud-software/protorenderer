package srv

import (
	"context"
	"fmt"
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/utils"
	"strings"
	"sync"
)

var (
	MetaCache = &metaCache{metaEntries: make(map[string]*metaEntry)}
)

type metaCache struct {
	sync.Mutex
	metaEntries map[string]*metaEntry // keyed by .proto filename
}
type metaEntry struct {
	sync.Mutex
	filename      string // .info filename
	protofileinfo *pb.ProtoFileInfo
}

func (pr *protoRenderer) GetReverseDependencies(ctx context.Context, req *pb.ReverseDependenciesRequest) (*pb.ReverseDependenciesResponse, error) {
	fmt.Printf("Getting dependencies for \"%s\"\n", req.Filename)
	err := MetaCache.readAllIfNecessary()
	if err != nil {
		return nil, err
	}
	res := &pb.ReverseDependenciesResponse{}
	with_deps, err := MetaCache.AllWithDependencyOn(req.Filename, req.MaxDepth)
	if err != nil {
		return nil, err
	}
	for _, meta := range with_deps {
		res.Filenames = append(res.Filenames, meta.ProtoFile.Name)
	}
	return res, nil
}

// read all from store if required
func (mc *metaCache) readAllIfNecessary() error {
	mc.Lock()
	defer mc.Unlock()
	infodir := CompileEnv.StoreDir() + "/info"
	err := utils.DirWalk(infodir, func(root, relfile string) error {
		if !strings.HasSuffix(relfile, ".info") {
			return nil
		}
		//		fmt.Printf("Info file: %s\n", relfile)
		proto_filename := strings.TrimSuffix(relfile, ".info")
		proto_filename = proto_filename + ".proto"
		_, found := mc.metaEntries[proto_filename]
		if found {
			return nil
		}
		mc.metaEntries[proto_filename] = &metaEntry{filename: relfile}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (mc *metaCache) AllWithDependencyOn(filename string, maxdepth uint32) ([]*pb.ProtoFileInfo, error) {
	mc.Lock()
	defer mc.Unlock()
	res_m := make(map[string]*pb.ProtoFileInfo)
	err := mc.allWithDependencyOnRecursive(res_m, filename, maxdepth, 0)
	if err != nil {
		return nil, err
	}
	var res []*pb.ProtoFileInfo
	for _, v := range res_m {
		res = append(res, v)
	}
	return res, nil
}

// adds dependent files to map (must be called with lock held)
func (mc *metaCache) allWithDependencyOnRecursive(res map[string]*pb.ProtoFileInfo, filename string, maxdepth, cur_depth uint32) error {
	if (maxdepth != 0) && (cur_depth > maxdepth) {
		return nil
	}
	fmt.Printf("Finding reverse deps for \"%s (depth=%d)\n", filename, cur_depth)
	for _, me := range mc.metaEntries {
		pfi, err := me.GetProtoFileInfo()
		if err != nil {
			return err
		}
		for _, imp := range pfi.Imports {
			if imp.Name == filename {
				_, fd := res[imp.Name]
				if fd {
					continue
				}
				res[imp.Name] = pfi
				err = mc.allWithDependencyOnRecursive(res, pfi.ProtoFile.Name, maxdepth, cur_depth+1)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
func (me *metaEntry) GetProtoFileInfo() (*pb.ProtoFileInfo, error) {
	me.Lock()
	defer me.Unlock()
	if me.protofileinfo != nil {
		return me.protofileinfo, nil
	}
	b, err := utils.ReadFile(CompileEnv.StoreDir() + "/info/" + me.filename)
	if err != nil {
		return nil, err
	}
	res := &pb.ProtoFileInfo{}
	err = utils.UnmarshalBytes(b, res)
	if err != nil {
		return nil, err
	}
	me.protofileinfo = res
	return res, nil

}
