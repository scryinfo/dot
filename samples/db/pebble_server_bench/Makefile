ifeq ($(OS),Windows_NT)
	EXE := .exe
	ifneq (,$(findstring MINGW,$(UNAME_S)))
	    go := ${shell which go}
	else
	    go := $(subst \,/,$(shell where.exe go))
	endif
	VCPKG := $(abspath ./../../../vcpkg_installed)
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

go_rocksdb := CGO_CFLAGS="${CGO_CFLAGS}" CGO_LDFLAGS="${CGO_LDFLAGS}" ${go}

.PHONY: clean build rebuild

clean:
	${go} clean
	rm -rf go.sum rocksdb$(EXE) temp/ logs/
build:
	${go_rocksdb} build ./...
	${go_rocksdb} test -c -o /dev/null ./...

rebuild: clean build

bench:
	${go_rocksdb} test -bench=. -run=^$$ -count=3 -benchmem -benchtime=3s
