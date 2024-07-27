package main

import (
	"bufio"
	"fmt"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/protorenderer/v2/store"
	"os"
)

func EditStore() error {
	CompileEnv := StandardCompilerEnvironment{workdir: "/tmp/protorenderer_edit_store"}
	ctx := authremote.Context()
	err := store.Retrieve(ctx, CompileEnv.StoreDir(), 0) // 0 == latest
	if err != nil {
		return err
	}
	// wait for input
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Store downloaded to %s, press return to save it again:", CompileEnv.StoreDir())
	_, err = reader.ReadString('\n')
	if err != nil {
		return err
	}
	fmt.Printf("Saving...\n")
	ctx = authremote.Context()
	err = store.Store(ctx, CompileEnv.StoreDir())
	if err != nil {
		return err
	}
	return nil
}
