
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
	rm -f go.sum go.work.sum demo/go.sum
	${go} clean
	cd demo && ${go} clean
tidy:
	${go} mod tidy
	cd demo && ${go} mod tidy
	cd demo/redis && ${go} mod tidy
	cd demo/redis/orm && ${go} mod tidy
	cd dots/db/tools/gdao && ${go} mod tidy
	cd dots/db/tools/gmodel && ${go} mod tidy
upgrade:
	${go} get -t -u ./... && ${go} mod tidy
	cd demo && ${go} get -t -u ./... && ${go} mod tidy
	cd demo/redis && ${go} get -t -u ./... && ${go} mod tidy
	cd demo/redis/orm && ${go} get -t -u ./... && ${go} mod tidy
	cd dots/db/tools/gdao && ${go} get -t -u ./... && ${go} mod tidy
	cd dots/db/tools/gmodel && ${go} get -t -u ./... && ${go} mod tidy
format:
	${go} fmt ./...
	cd demo && ${go} fmt ./...
	cd demo/redis && ${go} fmt ./...
	cd demo/redis/orm && ${go} fmt ./...
	cd dots/db/tools/gdao && ${go} fmt ./...
	cd dots/db/tools/gmodel && ${go} fmt ./...
build:
	${go} build ./...
	cd demo && ${go} build ./...
	cd demo/redis && ${go} build ./...
	cd demo/redis/orm && ${go} build ./...
	cd dots/db/tools/gdao && ${go} build ./...
	cd dots/db/tools/gmodel && ${go} build ./...
