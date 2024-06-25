package java

import (
	"context"
	"fmt"
	//	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/protorenderer/v2/helpers"
	"golang.conradwood.net/protorenderer/v2/interfaces"
)

// compile all "known protos" to .java and .class files
// then copy them to storedir
func Start(ce interfaces.CompilerEnvironment, outdir string) error {
	srcdir := ce.WorkDir() + "/" + ce.AllKnownProtosDir()
	fmt.Printf("Compiling all java files from %s into %s...\n", srcdir, outdir)
	targetdir := outdir
	files, err := helpers.FindProtoFiles(srcdir)
	if err != nil {
		return err
	}
	fmt.Printf("%d files to compile\n", len(files))
	jc := New()
	ctx := context.Background()
	for _, pf := range files {
		pfs := []interfaces.ProtoFile{pf}
		err = jc.Compile(ctx, ce, pfs, targetdir, nil)
		if err != nil {
			//		return err
			fmt.Printf("Error compiling:%s\n", err)
		}
	}

	return nil
}
