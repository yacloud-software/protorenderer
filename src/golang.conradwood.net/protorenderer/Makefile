.PHONY:	client server protoc javaclient
all:	client server protoc javaclient
client:
	cd client && go install ${LDFLAGS} protorenderer-client.go `ls -1 *.go|grep -v protorenderer-client.go`
	cd v2/client && go install ${LDFLAGS} protorenderer2-client.go `ls -1 *.go|grep -v protorenderer2-client.go`
javaclient:
	cd v1/javaclient && go install ${LDFLAGS} protorenderer-java.go `ls -1 *.go|grep -v protorenderer-java.go`
server:
	cd server && go install ${LDFLAGS} protorenderer-server.go `ls -1 *.go|grep -v protorenderer-server.go`
protoc:
	cd v1/protoc && go install ${LDFLAGS} protoc-gen-meta.go `ls -1 *.go|grep -v protoc-gen-meta.go`
	cd v2/meta-compiler/protoc && go install ${LDFLAGS} protoc-gen-meta2.go `ls -1 *.go|grep -v protoc-gen-meta2.go`
test:
	cd compiler && go test
