ifeq ($(OS),Windows_NT)
	EXE := .exe
	override VCPKG := $(abspath ./../../../vcpkg_installed)
	override ROCKSDB_INCLUDE :=${VCPKG}/x64-windows/include
	override ROCKSDB_LIB :=${VCPKG}/x64-windows/lib
	# -static -static-libgcc -static-libstdc++
	override CGO_LDFLAGS := -static -L${ROCKSDB_LIB} -lrocksdb -lstdc++ -lm -lz -lsnappy -lbz2 -llz4 -lzstd -lrpcrt4 -lshlwapi
	override CGO_CFLAGS :=-I${ROCKSDB_INCLUDE}
	vcpath := $(shell "/c/Program Files (x86)/Microsoft Visual Studio/Installer/vswhere.exe" -latest -property installationPath)
	vcpath := $(subst \,/,$(vcpath))
	# vcvars64 := $(shell cmd.exe /c call msvc.bat)
	go_rocksdb := ${vcvars64} CGO_CFLAGS="${CGO_CFLAGS}" CGO_LDFLAGS="${CGO_LDFLAGS}" command go

else
	EXE :=
	CGO_LDFLAGS :=-L${ROCKSDB_LIB} -lrocksdb -lzstd -llz4 -lsnappy -lz -lbz2 -lstdc++ -lm -ldl -pthread
	CGO_CFLAGS :=-I${ROCKSDB_INCLUDE}
	go_rocksdb :=  CGO_CFLAGS="${CGO_CFLAGS}" CGO_LDFLAGS="${CGO_LDFLAGS}" command go
endif
$(info "tttsdfd, ${ttt}\n")
$(info "${vcpath}/VC/Auxiliary/Build/vcvars64.bat")
echo "test  :"

.PHONY: clean build rebuild

clean:
	command go clean
	rm -f go.sum rocksdb$(EXE) rocksdb_gen.toml wire_gen.go
rocksdb.exe: wire_gen.go
	cmd.exe //c call msvc.bat "${CGO_CFLAGS}" "${CGO_LDFLAGS}" rocksdb$(EXE)
rocksdb: wire_gen.go
	${go_rocksdb} build -o rocksdb$(EXE) ./...
wire_gen.go:
	wire
build: rocksdb$(EXE)

rebuild: clean build

rocksdb_gen.toml: rocksdb.go wire_gen.go
	${go_rocksdb} run ./... -- -GenerateConfig
generate_config: rocksdb_gen.toml
