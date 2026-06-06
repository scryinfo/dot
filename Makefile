
ifeq ($(OS),Windows_NT)
	EXE := .exe
	ifneq (,$(findstring MINGW,$(UNAME_S)))
	    go := ${shell which go}
	else
	    go := $(subst \,/,$(shell where.exe go))
	endif
	VCPKG := $(abspath ./vcpkg_installed)
	ROCKSDB_INCLUDE :=${VCPKG}/x64-mingw-static/include
	ROCKSDB_LIB :=${VCPKG}/x64-mingw-static/lib
	CGO_LDFLAGS :=-L${ROCKSDB_LIB} -lrocksdb -lstdc++ -lm -lz -lsnappy -lbz2 -llz4 -lzstd -lrpcrt4 -lshlwapi
	CGO_CFLAGS :=-I${ROCKSDB_INCLUDE}
else
	go := ${shell which go}
	EXE :=
	CGO_LDFLAGS :=-L${ROCKSDB_LIB} -lrocksdb -lzstd -llz4 -lsnappy -lz -lbz2 -lstdc++ -lm -ldl -pthread
	CGO_CFLAGS :=-I${ROCKSDB_INCLUDE}
endif
$(info "go: ${go}")

CGO_CFLAGS :=-I${ROCKSDB_INCLUDE}
go_rocksdb := CGO_CFLAGS="${CGO_CFLAGS}" CGO_LDFLAGS="${CGO_LDFLAGS}" ${go}

.PHONY: clean upgrade format build samples

clean:
	rm -rf go.sum go.work.sum demo/go.sum node_modules bun.lock
	${go} clean
	cd demo && ${go} clean
	cd samples && make clean
tidy:
	${go} mod tidy
	cd demo && ${go} mod tidy
	cd demo/redis && ${go} mod tidy
	cd demo/redis/orm && ${go} mod tidy
	cd line/db/tools/gdao && ${go} mod tidy
	cd line/db/tools/gmodel && ${go} mod tidy
	cd line/db/rocksdbdot && ${go} mod tidy
	cd samples && make tidy
upgrade:
	${go} get -t -u ./... && ${go} mod tidy
	cd demo && ${go} get -t -u ./... && ${go} mod tidy
	cd demo/redis && ${go} get -t -u ./... && ${go} mod tidy
	cd demo/redis/orm && ${go} get -t -u ./... && ${go} mod tidy
	cd line/db/tools/gdao && ${go} get -t -u ./... && ${go} mod tidy
	cd line/db/tools/gmodel && ${go} get -t -u ./... && ${go} mod tidy
	cd line/db/rocksdbdot && ${go} get -t -u ./... && ${go} mod tidy
	cd samples && make upgrade
format:
	${go} fmt ./...
	cd demo && ${go} fmt ./...
	cd demo/redis && ${go} fmt ./...
	cd demo/redis/orm && ${go} fmt ./...
	cd line/db/tools/gdao && ${go} fmt ./...
	cd line/db/tools/gmodel && ${go} fmt ./...
	cd line/db/rocksdbdot && ${go} fmt ./...
	cd samples && make format
build:
	bun install
	${go} build ./...
	cd demo && ${go} build ./...
	cd demo/redis && ${go} build ./...
	cd demo/redis/orm && ${go} build ./...
	cd line/db/tools/gdao && ${go} build ./...
	cd line/db/tools/gmodel && ${go} build ./...
	cd line/db/rocksdbdot && make build
	cd samples && make build
rebuild: clean gen wire build
go_fix:
	${go} fix ./...
	cd demo && ${go} fix ./...
	cd demo/redis && ${go} fix ./...
	cd demo/redis/orm && ${go} fix ./...
	cd line/db/tools/gdao && ${go} fix ./...
	cd line/db/tools/gmodel && ${go} fix ./...
	cd line/db/rocksdbdot && ${go} fix ./...
wire:
	cd samples && make wire
samples:
	cd samples && make samples
gen:
	cd samples && make gen

lint:
	${go} vet ./...
	govulncheck ./...
	staticcheck ./...
	nilaway ./...
	golangci-lint run --no-config --disable-all -E zerologlint ./...
lint_more:
	gosec -tags "Release" -quiet ./...
	revive ./...
	gocyclo ./...
	golangci-lint run ./...

go_tools:
	${go} install github.com/google/wire/cmd/wire@latest
	${go} install github.com/bufbuild/buf/cmd/buf@latest
	${go} install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	${go} install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
	${go} install github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto@latest
	${go} install github.com/mfridman/protoc-gen-go-json@latest
	${go} install github.com/fatih/gomodifytags@latest
	${go} install golang.org/x/tools/gopls@latest
	${go} install honnef.co/go/tools/cmd/staticcheck@latest
	${go} install github.com/cweill/gotests/gotests@latest
	${go} install github.com/josharian/impl@latest
	${go} install github.com/go-delve/delve/cmd/dlv@latest
	${go} install go.uber.org/nilaway/cmd/nilaway@latest
	${go} install golang.org/x/vuln/cmd/govulncheck@latest
	${go} install github.com/securego/gosec/v2/cmd/gosec@latest
	${go} install github.com/mgechev/revive@latest
	${go} install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	${go} install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

install_rocksdb:
	VCPKG_BUILD_TYPE=release d:/lang/vcpkg/vcpkg.exe install --triplet=x64-mingw-static
