ifeq ($(OS),Windows_NT)
	EXE := .exe
	VCPKG := $(abspath ./../../../vcpkg_installed)
	ROCKSDB_INCLUDE :=${VCPKG}/x64-mingw-static/include
	ROCKSDB_LIB :=${VCPKG}/x64-mingw-static/lib
	CGO_LDFLAGS :=-L${ROCKSDB_LIB} -lrocksdb -lstdc++ -lm -lz -lsnappy -lbz2 -llz4 -lzstd -lrpcrt4 -lshlwapi
	CGO_CFLAGS :=-I${ROCKSDB_INCLUDE}
else
	EXE :=
	CGO_LDFLAGS :=-L${ROCKSDB_LIB} -lrocksdb -lstdc++ -lm -lz -lsnappy -llz4 -lzstd -lbz2
	CGO_CFLAGS :=-I${ROCKSDB_INCLUDE}
endif

go_rocksdb := CGO_CFLAGS="${CGO_CFLAGS}" CGO_LDFLAGS="${CGO_LDFLAGS}" command go

.PHONY: clean build rebuild

clean:
	command go clean
	rm -f go.sum
build:
	${go_rocksdb} build ./...

rebuild: clean build
