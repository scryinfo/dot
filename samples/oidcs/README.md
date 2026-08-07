# app 是一个应用，这里是一个web 应用

# app_server 是应用对应的服务端

# oidc_server 是oidc 协议的服务端，这最终存放用户等信息

# go tools

```bash
go := go
${go} install github.com/google/wire/cmd/wire@latest
${go} install github.com/bufbuild/buf/cmd/buf@latest
${go} install google.golang.org/protobuf/cmd/protoc-gen-go@latest
${go} install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
${go} install github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto@latest
${go} install github.com/mfridman/protoc-gen-go-json@latest
${go} install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

bun install -g @bufbuild/protoc-gen-es
```
