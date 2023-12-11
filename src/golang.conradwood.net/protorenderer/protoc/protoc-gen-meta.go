package main

/**
* this is a protoc plugin
* it parses a single proto file and generates "protorenderer.Package|RPC|Service|Message" protos
* it relies on the protorenderer-server to issue IDs
 */
import (
	"fmt"
	"github.com/golang/protobuf/proto"
	//	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	//	"golang.conradwood.net/apis/create"
	pr "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/client"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/utils"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	debug = false
)

var (
	VerifyToken string
	sctx        string
	localport   int
)

func main() {
	if len(os.Args) > 1 && (os.Args[1]) == "-h" {
		print_help()
		os.Exit(0)
	}

	// oh we cannot use flag.Parse here..
	// this gives us a problem with finding the registry...
	data, err := ioutil.ReadAll(os.Stdin)
	utils.Bail("failed to read stdin", err)
	request := &plugin.CodeGeneratorRequest{}
	err = proto.Unmarshal(data, request)
	utils.Bail("failed to unmarshal", err)

	para := ""
	if request.Parameter != nil {
		para = *request.Parameter
	}
	sx := strings.Split(para, ",")
	if len(sx) != 4 {
		printf("Require 4 paras, got %d\n", len(sx))
		os.Exit(10)
	}
	VerifyToken = sx[0]
	sctx = sx[1]
	registry := sx[3]
	localport, err = strconv.Atoi(sx[2])
	if err != nil {
		printf("invalid port: %s\n", err)
		os.Exit(10)
	}
	cmdline.SetClientRegistryAddress(registry)
	response, err := generate_remote(request)
	if err != nil {
		printf("Failed to generate: %s\n", err)
		s := fmt.Sprintf("%s", err)
		response = &plugin.CodeGeneratorResponse{Error: &s}
	}

	data, err = proto.Marshal(response)
	utils.Bail("failed to marshal", err)
	_, err = os.Stdout.Write(data)
	utils.Bail("failed to write to stdout", err)
}
func generate_remote(req *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	client.GetSignatureFromAuth()
	con, err := client.ConnectWithIP(fmt.Sprintf("localhost:%d", localport))
	if err != nil {
		return nil, err
	}
	ps := pr.NewProtoRendererServiceClient(con)
	ctx, err := auth.RecreateContextWithTimeout(time.Duration(5)*time.Minute, []byte(sctx))
	if err != nil {
		printf("Failed to create context: %s\n", err)
		return nil, err
	}
	/*
		if auth.GetUser(ctx) == nil && auth.GetService(ctx) == nil {
			printf("Deserialised context has neither service nor user (%s)\n", sctx)
		}
	*/
	pr := &pr.ProtocRequest{VerifyToken: VerifyToken, ProtoFiles: req.ProtoFile}
	_, err = ps.SubmitSource(ctx, pr)
	if err != nil {
		printf("protoc-meta ERROR: %s\n", utils.ErrorString(err))
		return nil, err
	}
	response := &plugin.CodeGeneratorResponse{}
	return response, nil
}

func debugf(format string, args ...interface{}) {
	if !debug {
		return
	}
	fmt.Fprintf(os.Stderr, format, args...)
}
func printf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func print_help() {
	h := `protoc meta compiler

  Sourcecode: %s

  this plugin is a stub. It's only function is to take the input from the protoc compiler and forward it to protorender service, specifically the SubmitSource RPC. The result of the RPC is returned to protoc.
`
	fmt.Printf(h, cmdline.SourceCodePath())

}




















































































