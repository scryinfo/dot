ifeq ($(OS),Windows_NT)
	EXE := .exe
	VCPKG := $(abspath ./../../../vcpkg_installed)
	ROCKSDB_INCLUDE :=${VCPKG}/x64-mingw-static/include
	ROCKSDB_LIB :=${VCPKG}/x64-mingw-static/lib
	# -static -static-libgcc -static-libstdc++
	#CGO_LDFLAGS := -static -L${ROCKSDB_LIB} -lrocksdb -lstdc++ -lm -lz -lsnappy -lbz2 -llz4 -lzstd -lrpcrt4 -lshlwapi
	CGO_LDFLAGS := -static -L${ROCKSDB_LIB} -lrocksdb -lstdc++ -lm -lz -lsnappy -llz4 -lzstd -lbz2 -lrpcrt4 -lshlwapi
	CGO_CFLAGS :=-I${ROCKSDB_INCLUDE}
	go_rocksdb := CGO_CFLAGS="${CGO_CFLAGS}" CGO_LDFLAGS="${CGO_LDFLAGS}" command go

else
	EXE :=
	CGO_LDFLAGS :=-L${ROCKSDB_LIB} -lrocksdb -lzstd -llz4 -lsnappy -lz -lbz2 -lstdc++ -lm -ldl -pthread
	CGO_CFLAGS :=-I${ROCKSDB_INCLUDE}
	go_rocksdb :=  CGO_CFLAGS="${CGO_CFLAGS}" CGO_LDFLAGS="${CGO_LDFLAGS}" command go
endif
$(info "")

.PHONY: clean build rebuild

clean:
	command go clean
	rm -f go.sum rocksdb$(EXE) rocksdb_gen.toml wire_gen.go
rocksdb$(EXE): wire_gen.go
	${go_rocksdb} build -o rocksdb$(EXE) ./...
wire_gen.go:
	wire
build: rocksdb$(EXE)

rebuild: clean build

rocksdb_gen.toml: rocksdb.go wire_gen.go
	${go_rocksdb} run ./... -- -GenerateConfig
generate_config: rocksdb_gen.toml
