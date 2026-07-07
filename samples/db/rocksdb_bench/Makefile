ifeq ($(OS),Windows_NT)
	EXE := .exe
	ifneq (,$(findstring MINGW,$(UNAME_S)))
	    go := ${shell which go}
	else
	    go := $(subst \,/,$(shell where.exe go))
	endif
	override VCPKG := $(abspath ./../../../vcpkg_installed)
	override ROCKSDB_INCLUDE :=${VCPKG}/x64-mingw-static/include
	override ROCKSDB_LIB :=${VCPKG}/x64-mingw-static/lib
	# -static -static-libgcc -static-libstdc++
	override CGO_LDFLAGS := -static -L${ROCKSDB_LIB} -lrocksdb -lstdc++ -lm -lz -lsnappy -lbz2 -llz4 -lzstd -lrpcrt4 -lshlwapi
	override CGO_CFLAGS :=-I${ROCKSDB_INCLUDE}
else
	go := ${shell which go}
	EXE :=
	CGO_LDFLAGS :=-L${ROCKSDB_LIB} -static -lrocksdb -lzstd -llz4 -lsnappy -lz -lbz2 -lstdc++ -lm -ldl -pthread
	CGO_CFLAGS :=-I${ROCKSDB_INCLUDE}
endif
$(info "go: ${go}")

go_rocksdb := CGO_CFLAGS="${CGO_CFLAGS}" CGO_LDFLAGS="${CGO_LDFLAGS}" ${go}

.PHONY: clean build rebuild

clean:
	${go} clean
	rm -rf go.sum rocksdb_bench$(EXE) temp/ logs/
rocksdb_bench$(EXE):
	CGO_ENABLED=1 ${go_rocksdb} test -ldflags="-extldflags=-Wl,--no-warnings" -c -o rocksdb_bench$(EXE) ./...
build: rocksdb_bench$(EXE)

rebuild: clean build

bench:
	${go_rocksdb} test -bench=. -run=^$$ -count=3 -benchmem -benchtime=3s
