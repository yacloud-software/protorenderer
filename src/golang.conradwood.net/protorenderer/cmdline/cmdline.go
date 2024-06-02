package cmdline

import (
	"flag"
)

var (
	port = flag.Int("port", 4102, "The grpc server port")
)

func GetRPCPort() int {
	return *port
}
