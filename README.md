# go-grpc-server-shop
Simple gRPC server written in GO with simple shop like CRUD API

## Proto GO files generation
Requires <i>protoc</i> to be installed:
```shell
sudo apt install -y protobuf-compiler
```

Then install the GO specific proto generators: 
```shell
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
go install google.golang.org/protobuf/cmd/protoc-gen-go
```

and run
```shell
protoc -I proto/ --go_out=./proto --go_opt=paths=source_relative --go-grpc_out=./proto --go-grpc_opt=paths=source_relative proto/shop.proto
```
which generates GO code files into the proto directory.

## Local run and tests
```
go build
./go-grpc-server-shop

curl 127.0.0.1:8079/metrics
```
Create
```
grpcurl -d '{"name":"name-1", "price":45}' -cert hack/client-cert.pem -key hack/client-key.pem -cacert hack/ca-cert.pem localhost:8443 shop.v1.ShopService/Create
```
Get
```
grpcurl -d '{"id":"<ID>"}' -cert hack/client-cert.pem -key hack/client-key.pem -cacert hack/ca-cert.pem localhost:8443 shop.v1.ShopService/Get
```
Update
```
grpcurl -d '{"id":"<ID>", "name":"name-updated", "price":100.15}' -cert hack/client-cert.pem -key hack/client-key.pem -cacert hack/ca-cert.pem localhost:8443 shop.v1.ShopService/Update
```
Get all
```
grpcurl -d '{}' -cert hack/client-cert.pem -key hack/client-key.pem -cacert hack/ca-cert.pem localhost:8443 shop.v1.ShopService/GetAll
```
Remove
```
grpcurl -d '{"id":"<ID>"}' -cert hack/client-cert.pem -key hack/client-key.pem -cacert hack/ca-cert.pem localhost:8443 shop.v1.ShopService/Remove
```