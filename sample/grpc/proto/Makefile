
go := ${shell which go}

install:
	${go} install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.5
	${go} install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

gen:
	protoc --go_out=../go_out/hidot/ --go_opt=paths=source_relative --go-grpc_out=../go_out/hidot/ --go-grpc_opt=paths=source_relative hi.proto
