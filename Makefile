
ifeq ($(OS),Windows_NT)
	EXE := .exe
	VCPKG := $(subst \,/,$(abspath ./vcpkg_installed))
	ROCKSDB_INCLUDE :=${VCPKG}/x64-mingw-static/include
	ROCKSDB_LIB :=${VCPKG}/x64-mingw-static/lib
	CGO_LDFLAGS :=-L${ROCKSDB_LIB} -lrocksdb -lstdc++ -lm -lz -lsnappy -lbz2 -llz4 -lzstd -lrpcrt4 -lshlwapi
	CGO_CFLAGS :=-I${ROCKSDB_INCLUDE}
else
	EXE :=
	VCPKG :=""
	CGO_LDFLAGS :=-L${ROCKSDB_LIB} -lrocksdb -lzstd -llz4 -lsnappy -lz -lbz2 -lstdc++ -lm -ldl -pthread
	CGO_CFLAGS :=-I${ROCKSDB_INCLUDE}

endif

CGO_CFLAGS :=-I${ROCKSDB_INCLUDE}
go_rocksdb := CGO_CFLAGS="${CGO_CFLAGS}" CGO_LDFLAGS="${CGO_LDFLAGS}" command go

.PHONY: clean upgrade format build samples

ifeq ($(OS),Windows_NT)
VCPKG_INSTALLED:= ${VCPKG}/x64-mingw-static/lib/libz.a
${VCPKG_INSTALLED}:
	VCPKG_BUILD_TYPE=release vcpkg.exe install --triplet=x64-mingw-static
endif
clean:
	rm -rf go.sum go.work.sum demo/go.sum node_modules bun.lock
	command go clean
	cd demo && make clean
	cd samples && make clean
	cd line/db/pebble_service && make clean
	cd line/oidcdot/proto && make clean
clean_web:
	find . -name "node_modules" -type d -prune -exec rm -rf {} +
	find . -name "dist" -type d -prune -exec rm -rf {} +
tidy:
	command go mod tidy
	cd demo && make tidy
	cd line/db/tools/gdao && command go mod tidy
	cd line/db/tools/gmodel && command go mod tidy
	cd line/db/rocksdbdot && command go mod tidy
	cd samples && make tidy

upgrade:
	command go get -t -u ./... && command go mod tidy
	cd demo && make upgrade
	cd line/db/tools/gdao && command go get -t -u ./... && command go mod tidy
	cd line/db/tools/gmodel && command go get -t -u ./... && command go mod tidy
	cd line/db/rocksdbdot && command go get -t -u ./... && command go mod tidy
	cd samples && make upgrade
	cd line/db/pebble_service && make upgrade
	cd line/oidcdot/oicd_ts && bun update --latest
	cd line/db/pebble_service/kv_ts && bun update --latest

format:
	command go fmt ./...
	cd demo && make format
	cd line/db/tools/gdao && command go fmt ./...
	cd line/db/tools/gmodel && command go fmt ./...
	cd line/db/rocksdbdot && command go fmt ./...
	cd samples && make format
	cd line/db/pebble_service && make format
build: ${VCPKG_INSTALLED}
	bun install
	# command go build -ldflags="-s -w" ./...
	cd demo && make build
	cd line/db/tools/gdao && command go build ./...
	cd line/db/tools/gmodel && command go build ./...
	cd line/db/rocksdbdot && make build
	cd samples && make build
	cd line/db/pebble_service && make build
rebuild: clean gen wire build
test:
	command go test -tags="release" ./...

go_fix:
	command go fix ./...
	cd demo && make go_fix
	cd line/db/tools/gdao && command go fix ./...
	cd line/db/tools/gmodel && command go fix ./...
	cd line/db/rocksdbdot && command go fix ./...
wire:
	cd samples && make wire
samples:
	cd samples && make samples
gen:
	cd samples && make gen
	cd line/db/pebble_service && make gen
	cd line/oidcdot/proto && make gen

lint:
	command go vet ./...
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
	command go install github.com/google/wire/cmd/wire@latest
	command go install github.com/bufbuild/buf/cmd/buf@latest
	command go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	command go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
	command go install github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto@latest
	command go install github.com/mfridman/protoc-gen-go-json@latest
	command go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	command go install github.com/fatih/gomodifytags@latest
	command go install golang.org/x/tools/gopls@latest
	command go install honnef.co/go/tools/cmd/staticcheck@latest
	command go install github.com/cweill/gotests/gotests@latest
	command go install github.com/josharian/impl@latest
	command go install github.com/go-delve/delve/cmd/dlv@latest
	command go install go.uber.org/nilaway/cmd/nilaway@latest
	command go install golang.org/x/vuln/cmd/govulncheck@latest
	command go install github.com/securego/gosec/v2/cmd/gosec@latest
	command go install github.com/mgechev/revive@latest
	command go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	command go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	command go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	command go install github.com/nao1215/gup@latest
	command go install github.com/jmattheis/goverter/cmd/goverter@latest
bun_tools:
	bun install -g @bufbuild/protoc-gen-es
install_rocksdb:
	VCPKG_BUILD_TYPE=release vcpkg.exe install --triplet=x64-mingw-static
