# dev tools
[protoc ](https://github.com/protocolbuffers/protobuf/releases) , move it to GOPATH/bin  
protoc-gen-go:  go install github.com/golang/protobuf/protoc-gen-go  
[protoc-gen-grpc-web](https://github.com/grpc/grpc-web/releases) rename protoc-gen-grpc-web after download, move it to GOPATH/bin

# go package
https://github.com/grpc/grpc-go  
github.com/improbable-eng/grpc-web,  grpc-web，ts-protoc-gen

# generate code
go code  
protoc --go_out=plugins=grpc:%out%/ hi.proto  
ts code  
protoc --js_out=import_style=commonjs:%out%/ --grpc-web_out=import_style=commonjs+dts,mode=grpcweb:%out%/ hi.proto  
js code  
protoc --js_out=import_style=commonjs:%out%/ --grpc-web_out=import_style=commonjs,mode=grpcweb:%out%/ hi.proto  
ts-protoc-gen
protoc --plugin="protoc-gen-ts" --js_out=import_style=commonjs,binary:%out%/ --ts_out=service=grpc-web:%out%/ hi.proto  

注：使用ts-protoc-gen生成的代码有一个bug ： exports is not defined
  需要在 “X_pb_service.js” 文件中加入如果代码(参见sample中的处理)：
```ts
export {HiDot, HiDotClient}
```