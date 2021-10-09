.PHONY:	client server protoc javaclient
all:	client server protoc javaclient
client:
	cd client && go install protorenderer-client.go `ls -1 *.go|grep -v protorenderer-client.go`
javaclient:
	cd javaclient && go install protorenderer-java.go `ls -1 *.go|grep -v protorenderer-java.go`
server:
	cd server && go install protorenderer-server.go `ls -1 *.go|grep -v protorenderer-server.go`
protoc:
	cd protoc && go install protoc-gen-meta.go `ls -1 *.go|grep -v protoc-gen-meta.go`
test:
	cd compiler && go test
