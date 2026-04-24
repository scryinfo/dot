
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
	rm -rf ../go_out/connect/* ../go_out/gorpc/* ../http/client/src/api/*
gen:
	buf generate
	protoc --proto_path=./ --go_opt=paths=source_relative --go_out=../go_out/gogrpc --go-grpc_opt=paths=source_relative --go-grpc_out=../go_out/gogrpc api/v1/hi.proto

go_tools:
	${go} install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	${go} install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
