package srv

import (
	"flag"
	"fmt"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/utils"
)

var (
	copy_protos = flag.Bool("copy_extra_protos", false, "if true copy text protos from extras directory to workdir on startup")
)

// take all the .proto files we know and copy them  to our workdir
func createWorkDir() error {
	if !*copy_protos {
		return nil
	}
	// copy protos from extra/test_protos
	test_proto_dir, err := utils.FindFile("extra/test_protos/previous_protos/protos")
	if err != nil {
		return err
	}

	proto_dir := CompileEnv.WorkDir() + "/" + CompileEnv.AllKnownProtosDir()
	fmt.Printf("Copying \"%s\" -> \"%s\"...\n", test_proto_dir, proto_dir)
	err = linux.CopyDir(test_proto_dir, proto_dir)
	utils.Bail("failed to copy new protos", err)

	return nil
}
