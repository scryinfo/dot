ifeq ($(OS),Windows_NT)
	EXE := .exe
	ifneq (,$(findstring MINGW,$(UNAME_S)))
	    go := ${shell which go}
	else
	    go := $(subst \,/,$(shell where.exe go))
	endif
	EXE := .exe
else
	go := ${shell which go}
	EXE :=
endif
$(info "go: ${go}")

.PHONY: clean build rebuild

clean:
	${go} clean
	rm -rf go.sum connect_server_8090${EXE} connect_server_8092${EXE} grpc_server_8091${EXE} temp/ logs/ log*.log
build:
	${go_rocksdb} build ./...
	${go_rocksdb} test -c -o /dev/null ./...

rebuild: clean build

server:
	${go} build -o connect_server_8090${EXE} ../pebble_connect_server/...
	cp connect_server_8090${EXE} connect_server_8092${EXE}
	${go} build -o grpc_server_8091${EXE} ../pebble_grpc_server/...

run_servers:
	setsid ./connect_server_8090${EXE} -configpath=./
	setsid ./connect_server_8092${EXE} -configpath=./
	setsid ./grpc_server_8091${EXE} -configpath=./
bench:
	${go} test -bench=. -run=^$$ -count=3 -benchmem -benchtime=3s
