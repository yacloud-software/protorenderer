package srv

import (
	"context"
	"fmt"
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/v2/helpers"
	"golang.conradwood.net/protorenderer/v2/interfaces"
	"strings"
)

type protofilesender interface {
}

func (pr *protoRenderer) GetProtoFile(req *pb.ProtoFileRequest, srv pb.ProtoRenderer2_GetProtoFileServer) error {
	ce := CompileEnv
	fname := req.ProtoFileName
	if len(fname) == 0 || strings.Contains(fname, " ") || strings.Contains(fname, "..") {
		fmt.Printf("Invalid filename: \"%s\"\n", fname)
		return errors.InvalidArgs(srv.Context(), "invalid filename", "invalid filename: \"%s\"", fname)
	}
	fmt.Printf("sending files for \"%s\"\n", fname)

	ctx := srv.Context()
	bs := utils.NewByteStreamSender(
		func(key, filename string) error {
			return srv.Send(&pb.FileTransfer{Filename: filename})
		},
		func(b []byte) error {
			return srv.Send(&pb.FileTransfer{Data: b})
		},
	)
	err := get_proto_file(ctx, bs, fname)
	if err != nil {
		return err
	}
	err = get_info_file(ctx, bs, fname)
	if err != nil {
		return err
	}
	for _, compname := range []string{"golang", "java"} {
		ctx = srv.Context()
		err = get_compiler_file(ctx, bs, ce, compname, fname)
		if err != nil {
			return err
		}
	}
	fmt.Printf("All Files sent for \"%s\"\n", fname)
	return nil
}

func get_proto_file(ctx context.Context, bs *utils.ByteStreamSender, protofilename string) error {
	ce := CompileEnv
	absname := fmt.Sprintf("%s/protos/%s", ce.StoreDir(), protofilename)
	key := "protos/" + protofilename
	//	fmt.Printf("Filename: \"%s\"\n", absname)
	b, err := utils.ReadFile(absname)
	if err != nil {
		return err
	}
	//	fmt.Printf("File size: %d\n", len(b))
	err = bs.SendBytes(key, key, b)
	if err != nil {
		return err
	}
	return nil
}
func get_info_file(ctx context.Context, bs *utils.ByteStreamSender, protofilename string) error {
	ce := CompileEnv
	infofilename := helpers.ChangeExt(protofilename, ".info")
	absname := fmt.Sprintf("%s/info/%s", ce.StoreDir(), infofilename)
	key := "info/" + infofilename
	//fmt.Printf("Filename: \"%s\"\n", absname)
	b, err := utils.ReadFile(absname)
	if err != nil {
		return err
	}
	//	fmt.Printf("File size: %d\n", len(b))
	err = bs.SendBytes(key, key, b)
	if err != nil {
		return err
	}
	return nil
}

func get_compiler_file(ctx context.Context, bs *utils.ByteStreamSender, ce interfaces.CompilerEnvironment, compname, protofilename string) error {
	fnames, err := get_artefact_filenames_for(ce, compname, protofilename)
	if err != nil {
		return err
	}
	fmt.Printf("Artefacts for compiler \"%s\":\n", compname)
	for _, fname := range fnames {
		key := fname
		fmt.Printf("    Sending \"%s\"\n", fname)
		b, err := utils.ReadFile(ce.StoreDir() + "/" + fname)
		if err != nil {
			return err
		}
		//	fmt.Printf("File size: %d\n", len(b))
		err = bs.SendBytes(key, key, b)
		if err != nil {
			return err
		}
	}
	return nil
}
