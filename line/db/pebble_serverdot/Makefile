
ifeq ($(OS),Windows_NT)
	ifneq (,$(findstring MINGW,$(UNAME_S)))
	    go := ${shell which go}
	else
	    go := $(subst \,/,$(shell where.exe go))
	endif
else
	go := ${shell which go}
endif
$(info ${go})

clean:
	rm -rf ./kv_gen/connect/* ./kv_gen/gogrpc/* ./kv_ts/src/kv_gen/
gen:
	buf generate
	buf generate --template ./buf.gen_ts.yaml
	# protoc --proto_path=./ --go_opt=paths=source_relative --go_out=./kv_gen/gogrpc --go-grpc_opt=paths=source_relative --go-grpc_out=./kv_gen/gogrpc kv/v1/kv.proto

buf_dep:
	buf dep update
build:
	${go} build ./...

go_tools:
	${go} install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	${go} install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
