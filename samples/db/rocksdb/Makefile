ifeq ($(OS),Windows_NT)
	EXE := .exe
	override VCPKG := $(abspath ./../../../vcpkg_installed)
	override ROCKSDB_INCLUDE :=${VCPKG}/x64-mingw-static/include
	override ROCKSDB_LIB :=${VCPKG}/x64-mingw-static/lib
	# -static -static-libgcc -static-libstdc++
	override CGO_LDFLAGS := -static -L${ROCKSDB_LIB} -lrocksdb -lstdc++ -lm -lz -lsnappy -lbz2 -llz4 -lzstd -lrpcrt4 -lshlwapi
	override CGO_CFLAGS :=-I${ROCKSDB_INCLUDE}
else
	EXE :=
	CGO_LDFLAGS :=-L${ROCKSDB_LIB} -lrocksdb -lzstd -llz4 -lsnappy -lz -lbz2 -lstdc++ -lm -ldl -pthread
	CGO_CFLAGS :=-I${ROCKSDB_INCLUDE}
endif

go_rocksdb :=  CGO_CFLAGS="${CGO_CFLAGS}" CGO_LDFLAGS="${CGO_LDFLAGS}" command go

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
