# godino
A dinosaur gallery in Golang

## Protocol Buffers usage

A `.proto` file is used to declare an animal data model. In order to be able to compile it into Golang code, the following instructions need to be followed: https://github.com/golang/protobuf#installation

In a clean system where projects only use Go modules, the $GOPATH and $GOBIN variables may not be set. In that case, an Ubuntu system will keep Golang binaries in `/usr/local/go/bin`. So the binary `protoc-gen-go` containing the ProtoBuf Golang plugin will have to be installed either there (for system wide access) or in another folder and that added to the $PATH.