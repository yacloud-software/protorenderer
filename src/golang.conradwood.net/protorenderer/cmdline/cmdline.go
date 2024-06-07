package cmdline

import (
	"flag"
)

var (
	port = flag.Int("port", 4102, "The grpc server port")

	java_compiler_bin   = flag.String("java_compiler", "/usr/bin/javac", "path to javac binary")
	java_release        = flag.String("java_release", "11", "flag --target [java_release] for javac: build for ssecific target version")
	java_use_std_protoc = flag.Bool("java_std_protoc", true, "if set use standard protoc compiler (the one with the OS rather than shipped in this repo")
	java_plugin_name    = flag.String("java_plugin_name", "protoc-gen-grpc-java-1.13.1-linux-x86_64.exe", "the name of the java gprc plugin in extra/compilers")
)

func GetRPCPort() int {
	return *port
}

func GetJavaCompilerBin() string {
	return *java_compiler_bin
}
func GetJavaRelease() string {
	return *java_release
}
func GetJavaStdProtoC() bool {
	return *java_use_std_protoc
}
func GetJavaPluginName() string {
	return *java_plugin_name
}
func GetCompilerVersion() string {
	return "original"
}
