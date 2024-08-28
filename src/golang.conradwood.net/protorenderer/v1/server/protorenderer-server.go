package main

import (
	"flag"
	srv1 "golang.conradwood.net/protorenderer/v1/srv"
	srv2 "golang.conradwood.net/protorenderer/v2/srv"
	"os"
)

var (
	v2 = flag.Bool("v2", false, "if true, use v2 renderer")
)

func main() {
	flag.Parse()
	if *v2 {
		srv2.Start()
		os.Exit(0)
	}
	srv1.Start()
}
