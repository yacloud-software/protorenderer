package main

import (
	"flag"
	srv1 "golang.conradwood.net/protorenderer/v1/srv"
)

func main() {
	flag.Parse()
	srv1.Start()
}
