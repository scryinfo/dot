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

# grpc stream
[see](https://grpc.io/docs/languages/go/basics/)
## server stream
the client sends a request to the server and gets a stream to read a sequence of messages back. The client reads from the returned stream until there are no more messages. As you can see in our example, you specify a server-side streaming method by placing the stream keyword before the response type.
## client stream
 the client writes a sequence of messages and sends them to the server, again using a provided stream. Once the client has finished writing the messages, it waits for the server to read them all and return its response. You specify a client-side streaming method by placing the stream keyword before the request type
## bidirectional stream
both sides send a sequence of messages using a read-write stream. The two streams operate independently, so clients and servers can read and write in whatever order they like: for example, the server could wait to receive all the client messages before writing its responses, or it could alternately read a message then write a message, or some other combination of reads and writes. The order of messages in each stream is preserved. You specify this type of method by placing the stream keyword before both the request and the response

## in web browser
| | |  |
| :----: | :----: |  :----:|
|             | browser http| websocket|
|server stream| 独立请求| 只能请求一次|
|client stream| 独立请求| 多次|
|bidirectional| 独立请求| 多次|
