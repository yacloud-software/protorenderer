package main

import (
	"flag"
	"fmt"

	srv1 "golang.conradwood.net/protorenderer/v1/srv"
)

var (
	do_run = flag.Bool("obsolete_code", false, "if true run anyways")
)

func main() {
	flag.Parse()
	if !*do_run {
		fmt.Printf("This is obsolete code and shall not run anymore.\n")
		fmt.Printf("(also see -h)\n")
		panic("obsolete code")
	}
	srv1.Start()
}
