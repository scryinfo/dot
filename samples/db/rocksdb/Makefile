ifeq ($(OS),Windows_NT)
	EXE := .exe
	ifneq (,$(findstring MINGW,$(UNAME_S)))
	    go := ${shell which go}
	else
	    go := $(subst \,/,$(shell where.exe go))
	endif
	CGO_LDFLAGS :=-L${ROCKSDB_LIB} -lrocksdb -lstdc++ -lm -lz -lsnappy -llz4 -lzstd
else
	go := ${shell which go}
	EXE :=
	CGO_LDFLAGS :=-L${ROCKSDB_LIB} -lrocksdb -lzstd -llz4 -lsnappy -lz -lbz2 -lstdc++ -lm -ldl -pthread
endif
$(info "go: ${go}")

CGO_CFLAGS :=-I${ROCKSDB_INCLUDE}
go_rocksdb := CGO_CFLAGS="${CGO_CFLAGS}" CGO_LDFLAGS="${CGO_LDFLAGS}" ${go}

.PHONY: clean build rebuild

clean:
	${go} clean
	rm -f go.sum rocksdb$(EXE) rocksdb_gen.toml wire_gen.go
build:
	wire
	${go_rocksdb} build -o rocksdb$(EXE) ./...

rebuild: clean build

rocksdb_gen.toml: rocksdb.go wire_gen.go
	${go_rocksdb} run ./... -- -GenerateConfig
generate_config: rocksdb_gen.toml
