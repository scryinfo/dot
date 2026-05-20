ifeq ($(OS),Windows_NT)
	EXE := .exe
	ifneq (,$(findstring MINGW,$(UNAME_S)))
	    go := ${shell which go}
	else
	    go := $(subst \,/,$(shell where.exe go))
	endif
else
	go := ${shell which go}
	EXE :=
endif
$(info "go: ${go}")

.PHONY: clean upgrade format build samples
WIRES := certificate/ecdsa certificate/ed25519 etcd_server etcd_client gindot sconfig simple db/redis_client \
	rpc/web_rpc/web_rpc_server \
	rpc/nobl/nobl_connect_client rpc/nobl/nobl_grpc_client rpc/nobl/nobl_server \
	rpc/nobls/nobls_connect_client rpc/nobls/nobls_grpc_client rpc/nobls/nobls_server \
	rpc/bl/bl_grpc_client rpc/bl/bl_etcds rpc/bl/bl_server rpc/bl/bl_server2 rpc/bl/bl_server3 \
	rpc/bls/bls_grpc_client rpc/bls/bls_server rpc/bls/bls_server2 rpc/bls/bls_server3

clean:
	cd rpc/proto && make clean
	cd rpc/web_rpc/client && rm -rf dist node_modules bun.lock
	find . -type f \( -executable -o -name "*.log" -o -name "wire_gen.go" \) -delete
clean_wire:
	find . -type f \( -executable -o -name "wire_gen.go" \) -delete
upgrade:
	cd rpc/web_rpc/client && bun update --latest
format:
	cd rpc/web_rpc/client && bun run format

DIST_web_rpc_client := rpc/web_rpc/client/dist/index.html
$(DIST_web_rpc_client): $(shell find rpc/web_rpc/client/src -type f)
	cd rpc/web_rpc/client && bun install && bun run build
build: wire samples $(DIST_web_rpc_client)

rebuild: clean gen build

clean_wire:
	find . -type f -name "wire_gen.go" -delete

define wire_rule
$1/wire_gen.go: $1/$(notdir $1).go
	cd $1 && wire
endef
$(foreach wire,$(WIRES),$(eval $(call wire_rule,$(wire))))
wire:  $(addsuffix /wire_gen.go, $(WIRES))

define sample_rule
$1/$(notdir $1)$(EXE): $1/$(notdir $1).go $1/wire_gen.go
	cd $1 && ${go} build -o $(notdir $1)$(EXE) ./...
endef
$(foreach wire,$(WIRES),$(eval $(call sample_rule,$(wire))))
samples: $(foreach wire,$(WIRES),$(wire)/$(notdir $(wire))$(EXE))

gen:
	cd rpc/proto && make gen
