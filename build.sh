 #!/bin/bash
 set -e

go get -u github.com/golang/protobuf/proto github.com/golang/protobuf/protoc-gen-go
cd ./api/0.0
protoc -I ./ ./main.proto --go_out=plugins=grpc:.
cd ../..
go get
go build -v