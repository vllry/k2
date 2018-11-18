 #!/bin/bash
 set -e

go get -u github.com/golang/protobuf/proto github.com/golang/protobuf/protoc-gen-go
protoc -I ./ ./main.proto --go_out=plugins=grpc:.
go get
go build -v