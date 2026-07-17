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
	KILL_CMD = taskkill /F /IM $(1) /T
else
	go := ${shell which go}
	EXE :=
	CGO_LDFLAGS :=-L${ROCKSDB_LIB} -static -lrocksdb -lzstd -llz4 -lsnappy -lz -lbz2 -lstdc++ -lm -ldl -pthread
	CGO_CFLAGS :=-I${ROCKSDB_INCLUDE}
	KILL_CMD = pkill -f '^\./$(1).*' || true
endif
$(info "go: ${go}")
go_rocksdb := CGO_CFLAGS="${CGO_CFLAGS}" CGO_LDFLAGS="${CGO_LDFLAGS}" ${go}

.PHONY: clean build rebuild

clean:
	${go} clean
	rm -rf go.sum connect_server_8090${EXE} connect_server_8092${EXE} grpc_server_8091${EXE} pebble_server_bench${EXE} temp/ logs/ log*.log data/
build: pebble_server_bench${EXE}

pebble_server_bench${EXE}:
	${go_rocksdb} test -c -o pebble_server_bench${EXE} ./...
rebuild: clean build

servers:
	${go} build -ldflags="-s -w" -o connect_server_8090${EXE} ../pebble_connect_server/...
	cp -f connect_server_8090${EXE} connect_server_8092${EXE}
	${go} build -ldflags="-s -w" -o grpc_server_8091${EXE} ../pebble_grpc_server/...

run_servers:
	./connect_server_8090${EXE} -configpath=./ & \
	./connect_server_8092${EXE} -configpath=./ & \
	./grpc_server_8091${EXE} -configpath=./ & \
	sleep 1
stop_servers:
	$(call KILL_CMD,connect_server_8090$(EXE))
	$(call KILL_CMD,connect_server_8092$(EXE))
	$(call KILL_CMD,grpc_server_8091$(EXE))
bench:
	${go} test -bench=. -run=^$$ -count=3 -benchmem -benchtime=3s
