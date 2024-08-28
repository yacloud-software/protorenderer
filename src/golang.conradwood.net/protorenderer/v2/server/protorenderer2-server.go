package main

import (
	"flag"
	srv2 "golang.conradwood.net/protorenderer/v2/srv"
	"os"
)

func main() {
	flag.Parse()
	srv2.Start()
	os.Exit(0)

}
