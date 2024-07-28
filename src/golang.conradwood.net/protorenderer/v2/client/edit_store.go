package main

import (
	"bufio"
	"fmt"
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/v2/store"
	"os"
	"time"
)

func EditStore() error {
	CompileEnv := StandardCompilerEnvironment{workdir: "/tmp/protorenderer_edit_store"}
	ctx := authremote.ContextWithTimeout(time.Duration(180) * time.Second)
	err := store.Retrieve(ctx, CompileEnv.StoreDir(), *version) // 0 == latest
	if err != nil {
		return err
	}

	vi := &pb.VersionInfo{}
	b, err := utils.ReadFile(CompileEnv.StoreDir() + "/versioninfo.pbbin")
	if err != nil {
		fmt.Printf("versioninfo not found: %s\n", err)
	} else {
		err = utils.UnmarshalBytes(b, vi)
		if err != nil {
			fmt.Printf("failed to unmarshal versioninfo: %s\n", err)
		}
	}
	b, err = utils.MarshalYaml(vi)
	if err != nil {
		fmt.Printf("failed to marshal versioninfo: %s\n", err)
	} else {
		err = utils.WriteFile("/tmp/versioninfo.yaml", b)
		if err != nil {
			fmt.Printf("failed to write versioninfo yaml: %s\n", err)
		} else {
			fmt.Printf("Saved /tmp/versioninfo.yaml\n")
		}
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
