
clean:
	rm -rf ../go_out/connect/* ../go_out/gorpc/* ../http/client/src/api/*
gen:
	buf generate
	protoc --proto_path=./ --go_opt=paths=source_relative --go_out=../go_out/gogrpc --go-grpc_opt=paths=source_relative --go-grpc_out=../go_out/gogrpc api/v1/hi.proto

go_tools:
	command go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	command go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
