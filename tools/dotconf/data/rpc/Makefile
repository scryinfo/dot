
go := ${shell which go}

install:
	${go} install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.5
	${go} install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

gen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative config.proto
	protoc --js_out=import_style=commonjs:../../ui/src/rpc_face/ --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:../../ui/src/rpc_face/ config.proto
