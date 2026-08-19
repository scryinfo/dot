
.PHONY: clean gen build buf_dep upgrade format go_tools

clean:
	rm -rf ./kv_gen/connect/* ./kv_gen/gogrpc/* ./kv_ts/src/kv_gen/ kv_ts/dist/ kv_ts/node_modules/
gen:
	buf generate
	buf generate --template ./buf.gen_ts.yaml
	# protoc --proto_path=./ --go_opt=paths=source_relative --go_out=./kv_gen/gogrpc --go-grpc_opt=paths=source_relative --go-grpc_out=./kv_gen/gogrpc kv/v1/kv.proto

buf_dep:
	buf dep update
build_kv_ts:=kv_ts/dist/index.html
build: ${build_kv_ts}
	command go build ./...
${build_kv_ts}:
	cd kv_ts && bun run build
upgrade:
	cd kv_ts && bun update --latest
format:
	command go fmt ./...
	cd kv_ts && bun format
go_tools:
	command go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	command go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
