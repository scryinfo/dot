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
WIRES := certificate etcd_server etcd_client gindot sconfig simple db/redis_client \
	rpc/web_rpc/web_rpc_server \
	rpc/nobl/nobl_client rpc/nobl/nobl_server \
	rpc/bl/bl_client rpc/bl/bl_server rpc/bl/bl_server2

clean:
	cd rpc/proto && make clean
	cd rpc/web_rpc/client && rm -rf dist node_modules bun.lock
	find . -type f \( -executable -o -name "*.log" -o -name "wire_gen.go" \) -delete
upgrade:
	cd rpc/web_rpc/client && bun update --latest
format:
	cd rpc/web_rpc/client && bun run format
build: samples
	cd rpc/web_rpc/client && bun install && bun run build
rebuild: clean gen wire build

WIRES := certificate etcd_server etcd_client gindot sconfig simple db/redis_client \
	rpc/web_rpc/web_rpc_server \
	rpc/nobl/nobl_client rpc/nobl/nobl_server \
	rpc/bl/bl_client rpc/bl/bl_server rpc/bl/bl_server2

define wire_rule
$1/wire_gen.go: $1/$(notdir $1).go
	cd $1 && wire
endef
$(foreach wire,$(WIRES),$(eval $(call wire_rule,$(wire))))
wire:  $(addsuffix /wire_gen.go, $(WIRES))


define sample_rule
$1/$(notdir $1)$(EXE): $1/$(notdir $1).go $1/wire_gen.go
	cd $1 && ${go} build -o $(notdir $1) ./...
endef
$(foreach wire,$(WIRES),$(eval $(call sample_rule,$(wire))))
samples: $(foreach wire,$(WIRES),$(wire)/$(notdir $(wire))$(EXE))

gen:
	cd rpc/proto && make gen
