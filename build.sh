 #!/bin/bash
 set -e

echo "Fetching protoc-gen-go"
go get -u github.com/golang/protobuf/proto github.com/golang/protobuf/protoc-gen-go
echo "Building protobufs"
cd ./api/0.0
protoc -I ./ ./*.proto --go_out=plugins=grpc:.
cd ../..
echo "Fetching dependencies"
dep ensure
echo "Building"
go build -v
