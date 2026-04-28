
ifeq ($(OS),Windows_NT)
	ifneq (,$(findstring MINGW,$(UNAME_S)))
	    go := ${shell which go}
	else
	    go := $(subst \,/,$(shell where.exe go))
	endif
else
	go := ${shell which go}
endif
$(info "go: ${go}")

.PHONY: clean upgrade format build samples

clean:
	cd rpc/proto && make clean
	cd rpc/web_rpc/client && rm -rf dist node_modules bun.lock
	cd rpc/web_rpc/server && find . -maxdepth 1 -type f -executable -delete && rm -rf *.exe *.log
	cd certificate && find . -maxdepth 1 -type f -executable -delete && rm -rf *.exe *.log
	cd gindot && find . -maxdepth 1 -type f -executable -delete && rm -rf *.exe *.log
	cd rpc/bl/client && find . -maxdepth 1 -type f -executable -delete && rm -rf *.exe *.log
	cd rpc/bl/server && find . -maxdepth 1 -type f -executable -delete && rm -rf *.exe *.log
	cd rpc/bl/server2 && find . -maxdepth 1 -type f -executable -delete && rm -rf *.exe *.log
	cd rpc/nobl/client && find . -maxdepth 1 -type f -executable -delete && rm -rf *.exe *.log
	cd rpc/nobl/server && find . -maxdepth 1 -type f -executable -delete && rm -rf *.exe *.log
	cd sconfig/ && find . -maxdepth 1 -type f -executable -delete && rm -rf *.exe *.log
	cd simple && find . -maxdepth 1 -type f -executable -delete && rm -rf *.exe *.log
upgrade:
	cd rpc/web_rpc/client && bun update --latest
format:
	cd rpc/web_rpc/client && bun run format
build: samples
	cd rpc/web_rpc/client && bun run build
rebuild: clean gen wire build

wire:
	cd certificate && wire
	cd gindot && wire
	cd sconfig && wire
	cd simple && wire
	cd db/redis_client && wire
	cd rpc/web_rpc/server && wire
	cd rpc/nobl/client && wire
	cd rpc/nobl/server && wire
	cd rpc/bl/client && wire
	cd rpc/bl/server && wire
	cd rpc/bl/server2 && wire
samples:
	cd certificate && ${go} build -o certificate ./...
	cd gindot && ${go} build -o gindot ./...
	cd sconfig && ${go} build -o sconfig ./...
	cd simple && ${go} build -o simple ./...
	cd db/redis_client && ${go} build -o redis_client ./...
	cd rpc/web_rpc/server && ${go} build -o web_server ./...
	cd rpc/nobl/client && ${go} build -o nobl_client ./...
	cd rpc/nobl/server && ${go} build -o nobl_server ./...
	cd rpc/bl/client && ${go} build -o bl_client ./...
	cd rpc/bl/server && ${go} build -o bl_server ./...
	cd rpc/bl/server2 && ${go} build -o bl_server2 ./...
gen:
	cd rpc/proto && make gen
