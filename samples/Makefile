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

.PHONY: clean clean_wire upgrade format build samples


WIRES := certificate/ecdsa certificate/ed25519 etcd_server etcd_client gindot sconfig sconfig_de simple \
	db/redis_client \
	rpc/web_rpc/web_rpc_server \
	rpc/nobl/nobl_connect_client rpc/nobl/nobl_grpc_client rpc/nobl/nobl_server \
	rpc/nobls/nobls_connect_client rpc/nobls/nobls_grpc_client rpc/nobls/nobls_server \
	rpc/bl/bl_grpc_client rpc/bl/bl_etcds rpc/bl/bl_server rpc/bl/bl_server2 rpc/bl/bl_server3 \
	rpc/bls/bls_grpc_client rpc/bls/bls_etcds rpc/bls/bls_server rpc/bls/bls_server2 rpc/bls/bls_server3 \
 	upgraders \
	../tools/ende/decode ../tools/ende/encode

clean:
	cd rpc/proto && make clean
	cd rpc/web_rpc/client && rm -rf dist node_modules bun.lock
	find . -type f \( -executable -o -name "*.log" -o -name "wire_gen.go" \) -delete
	cd db/rocksdb_bench && make clean
clean_wire:
	find . -type f \( -executable -o -name "wire_gen.go" \) -delete
tidy:
	cd db/rocksdb && ${go} mod tidy
	cd db/rocksdb_bench && ${go} mod tidy
upgrade:
	cd rpc/web_rpc/client && bun update --latest
	cd db/rocksdb && ${go} get -t ./... && ${go} mod tidy
	cd db/rocksdb_bench && ${go} get -t ./... && ${go} mod tidy
upgrade_latest:
	cd rpc/web_rpc/client && bun update --latest
	cd db/rocksdb && ${go} get -t -u ./... && ${go} mod tidy
	cd db/rocksdb_bench && ${go} get -t -u ./... && ${go} mod tidy
format:
	cd rpc/web_rpc/client && bun run format
	cd db/rocksdb && ${go} fmt ./...
	cd db/rocksdb_bench && ${go} fmt ./...

DIST_web_rpc_client := rpc/web_rpc/client/dist/index.html
$(DIST_web_rpc_client): $(shell find rpc/web_rpc/client/src -type f)
	cd rpc/web_rpc/client && bun install && bun run build
build: wire $(DIST_web_rpc_client)
	cd db/rocksdb_bench && make build

rebuild: clean gen build

clean_wire:
	find . -type f -name "wire_gen.go" -delete

define wire_rule
$1/wire_gen.go: $1/$(notdir $1).go
	cd $1 && wire
endef
$(foreach wire,$(WIRES),$(eval $(call wire_rule,$(wire))))
db/rocksdb/wire_gen.go: db/rocksdb/rocksdb.go
	cd db/rocksdb && wire
wire:  $(addsuffix /wire_gen.go, $(WIRES)) db/rocksdb/wire_gen.go

define sample_rule
$1/$(notdir $1)$(EXE): $1/$(notdir $1).go $1/wire_gen.go
	cd $1 && ${go} build -o $(notdir $1)$(EXE) ./...
endef
$(foreach wire,$(WIRES),$(eval $(call sample_rule,$(wire))))
db/rocksdb/rocksdb$(EXE): db/rocksdb/rocksdb.go db/rocksdb/wire_gen.go
	cd db/rocksdb && make build
samples: $(foreach wire,$(WIRES),$(wire)/$(notdir $(wire))$(EXE)) db/rocksdb/rocksdb$(EXE)

define generate_config_rule
$1/$(notdir $1)_gen.toml: $1/$(notdir $1).go $1/wire_gen.go
	cd $1 && ${go} run ./... -- -GenerateConfig
endef
$(foreach wire,$(WIRES),$(eval $(call generate_config_rule,$(wire))))
db/rocksdb/rocksdb_gen.toml: db/rocksdb/rocksdb.go db/rocksdb/wire_gen.go
	cd db/rocksdb && make generate_config
generate_config: $(foreach wire,$(WIRES),$(wire)/$(notdir $(wire))_gen.toml) db/rocksdb/rocksdb_gen.toml
clear_config:
	rm -f $(foreach wire,$(WIRES),$(wire)/$(notdir $(wire))_gen.toml) db/rocksdb/rocksdb_gen.toml

gen:
	cd rpc/proto && make gen
