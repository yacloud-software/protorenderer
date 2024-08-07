package metadata

import (
	"fmt"
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/v2/interfaces"
	"strings"
	"sync"
)

type metaCache struct {
	sync.Mutex
	ce          interfaces.CompilerEnvironment
	metaEntries map[string]*metaEntry // keyed by .proto filename
	has_read    bool
}
type metaEntry struct {
	sync.Mutex
	filename      string // .info filename
	protofileinfo *pb.ProtoFileInfo
	mc            *metaCache
}

func New() *metaCache {
	return &metaCache{metaEntries: make(map[string]*metaEntry)}
}
func (mc *metaCache) Add(pfi *pb.ProtoFileInfo) {
	mc.Lock()
	defer mc.Unlock()
	fname := pfi.ProtoFile.Filename
	fmt.Printf("[meta] Adding %s\n", fname)
	mc.metaEntries[fname] = &metaEntry{filename: fname, protofileinfo: pfi, mc: mc}
}
func (mc *metaCache) Fork() interfaces.MetaCache {
	mc.Lock()
	defer mc.Unlock()
	res := &metaCache{
		ce:          mc.ce,
		has_read:    false,
		metaEntries: make(map[string]*metaEntry),
	}
	for k, v := range mc.metaEntries {
		res.metaEntries[k] = v
	}
	return res
}
func (mc *metaCache) SetEnv(ce interfaces.CompilerEnvironment) {
	mc.ce = ce
}

// read all from store if required
func (mc *metaCache) readAllIfNecessary() error {
	if mc.has_read {
		return nil
	}
	mc.Lock()
	defer mc.Unlock()
	infodir := mc.ce.StoreDir() + "/info"
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
		mc.metaEntries[proto_filename] = &metaEntry{mc: mc, filename: relfile}
		return nil
	})
	if err != nil {
		return err
	}
	mc.has_read = true

	return nil
}

func (mc *metaCache) AllWithDependencyOn(filename string, maxdepth uint32) ([]*pb.ProtoFileInfo, error) {
	err := mc.readAllIfNecessary()
	if err != nil {
		return nil, err
	}
	mc.Lock()
	defer mc.Unlock()
	res_m := make(map[string]*pb.ProtoFileInfo)
	err = mc.allWithDependencyOnRecursive(res_m, filename, maxdepth, 0)
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
			if imp.Filename == filename {
				_, fd := res[imp.Filename]
				if fd {
					continue
				}
				res[imp.Filename] = pfi
				err = mc.allWithDependencyOnRecursive(res, pfi.ProtoFile.Filename, maxdepth, cur_depth+1)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
func (me *metaEntry) GetProtoFileInfo() (*pb.ProtoFileInfo, error) {
	err := me.mc.readAllIfNecessary()
	if err != nil {
		return nil, err
	}
	me.Lock()
	defer me.Unlock()
	if me.protofileinfo != nil {
		return me.protofileinfo, nil
	}
	b, err := utils.ReadFile(me.mc.ce.StoreDir() + "/info/" + me.filename)
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

// get a protofileinfo for a file
func (mc *metaCache) ByProtoFile(pf interfaces.ProtoFile) *pb.ProtoFileInfo {
	return mc.ByFilename(pf.GetFilename())
}

// get a protofileinfo for a .proto  file
func (mc *metaCache) ByFilename(fname string) *pb.ProtoFileInfo {
	err := mc.readAllIfNecessary()
	if err != nil {
		fmt.Printf("[metadata] failed to read meta: %s\n", err)
		return nil
	}
	mc.Lock()
	me := mc.metaEntries[fname]
	mc.Unlock()
	if me == nil {
		return nil
	}
	pfi, err := me.GetProtoFileInfo()
	if err != nil {
		fmt.Printf("failed to get protofileinfo: %s\n", err)
		return nil
	}
	return pfi
}

func (mc *metaCache) ImportFrom(src interfaces.MetaCache) {
	cast, ok := src.(*metaCache)
	if !ok {
		panic("cannot currently do a generic interface import of metacaches")
	}
	cast.Lock()
	defer cast.Unlock()
	mc.Lock()
	defer mc.Unlock()
	for k, v := range cast.metaEntries {
		mc.metaEntries[k] = v
	}
}
