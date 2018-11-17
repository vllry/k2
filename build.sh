 #!/bin/bash

 protoc -I ./ ./main.proto --go_out=plugins=grpc:.