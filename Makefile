
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

.PHONY: clean upgrade format build

clean:
	rm -rf go.sum go.work.sum demo/go.sum node_modules bun.lock
	${go} clean
	cd demo && ${go} clean
	cd samples/rpc/proto && make clean
	cd samples/rpc/http/client && rm -rf dist node_modules bun.lock
tidy:
	${go} mod tidy
	cd demo && ${go} mod tidy
	cd demo/redis && ${go} mod tidy
	cd demo/redis/orm && ${go} mod tidy
	cd line/db/tools/gdao && ${go} mod tidy
	cd line/db/tools/gmodel && ${go} mod tidy
upgrade:
	${go} get -t -u ./... && ${go} mod tidy
	cd demo && ${go} get -t -u ./... && ${go} mod tidy
	cd demo/redis && ${go} get -t -u ./... && ${go} mod tidy
	cd demo/redis/orm && ${go} get -t -u ./... && ${go} mod tidy
	cd line/db/tools/gdao && ${go} get -t -u ./... && ${go} mod tidy
	cd line/db/tools/gmodel && ${go} get -t -u ./... && ${go} mod tidy
	cd samples/rpc/http/client && bun update --latest
format:
	${go} fmt ./...
	cd demo && ${go} fmt ./...
	cd demo/redis && ${go} fmt ./...
	cd demo/redis/orm && ${go} fmt ./...
	cd line/db/tools/gdao && ${go} fmt ./...
	cd line/db/tools/gmodel && ${go} fmt ./...
	cd samples/rpc/http/client && bun run format
build:
	bun install
	${go} build ./...
	cd demo && ${go} build ./...
	cd demo/redis && ${go} build ./...
	cd demo/redis/orm && ${go} build ./...
	cd line/db/tools/gdao && ${go} build ./...
	cd line/db/tools/gmodel && ${go} build ./...
	cd samples/rpc/http/client && bun run build
rebuild: clean gen wire build
go_fix:
	${go} fix ./...
	cd demo && ${go} fix ./...
	cd demo/redis && ${go} fix ./...
	cd demo/redis/orm && ${go} fix ./...
	cd line/db/tools/gdao && ${go} fix ./...
	cd line/db/tools/gmodel && ${go} fix ./...
wire:
	cd samples/certificate && wire
	cd samples/gindot && wire
	cd samples/sconfig && wire
	cd samples/simple && wire
	cd samples/db/redis_client && wire
	cd samples/rpc/http/server && wire
	cd samples/rpc/bl && wire
	cd samples/rpc/nobl/client && wire
	cd samples/rpc/nobl/server && wire
	cd samples/rpc/nobl/server2 && wire
gen:
	cd samples/rpc/proto && make gen

lint:
	${go} vet ./...
	govulncheck ./...
	staticcheck ./...
	nilaway ./...
	golangci-lint run --no-config --disable-all -E zerologlint ./...
lint_more:
	gosec -tags "Release" -quiet ./...
	revive ./...
	gocyclo ./...
	golangci-lint run ./...

go_tools:
	${go} install github.com/google/wire/cmd/wire@latest
	${go} install github.com/bufbuild/buf/cmd/buf@latest
	${go} install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	${go} install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
	${go} install github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto@latest
	${go} install github.com/mfridman/protoc-gen-go-json@latest
	${go} install github.com/fatih/gomodifytags@latest
	${go} install golang.org/x/tools/gopls@latest
	${go} install honnef.co/go/tools/cmd/staticcheck@latest
	${go} install github.com/cweill/gotests/gotests@latest
	${go} install github.com/josharian/impl@latest
	${go} install github.com/go-delve/delve/cmd/dlv@latest
	${go} install go.uber.org/nilaway/cmd/nilaway@latest
	${go} install golang.org/x/vuln/cmd/govulncheck@latest
	${go} install github.com/securego/gosec/v2/cmd/gosec@latest
	${go} install github.com/mgechev/revive@latest
	${go} install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	${go} install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
