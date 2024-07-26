package meta_compiler

import (
	"context"
	"fmt"
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/v2/interfaces"
	"strings"
)

// get file meta info from disk, prefer compileroutput over known proto (searches both)
func GetMetaInfo(ctx context.Context, ce interfaces.CompilerEnvironment, filename string) (*pb.ProtoFileInfo, error) {
	info_name := strings.TrimSuffix(filename, ".proto")
	info_name = info_name + ".info"

	fname := ce.CompilerOutDir() + "/info/" + info_name
	if !utils.FileExists(fname) {
		ofname := fname
		//fmt.Printf("File \"%s\" does not exist\n", fname)
		fname = ce.StoreDir() + "/info/" + info_name
		if !utils.FileExists(fname) {
			fmt.Printf("File \"%s\" does not exist\n", ofname)
			fmt.Printf("File \"%s\" does not exist\n", fname)
			return nil, errors.NotFound(ctx, "no such file")
		}
	}
	b, err := utils.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	res := &pb.ProtoFileInfo{}
	err = utils.UnmarshalBytes(b, res)
	if err != nil {
		return nil, fmt.Errorf("Error in file \"%s\": %w", fname, err)
	}
	return res, nil
}
