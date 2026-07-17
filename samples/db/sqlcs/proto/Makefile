
ifeq ($(OS),Windows_NT)
	ifneq (,$(findstring MINGW,$(UNAME_S)))
	    go := ${shell which go}
	else
	    go := $(subst \,/,$(shell where.exe go))
	endif
else
	go := ${shell which go}
endif
$(info ${go})

clean:
	rm -rf ./data_gen/connect/*
gen:
	buf generate
buf_dep:
	buf dep update
build:
	${go} build ./...

go_tools:
	${go} install github.com/bufbuild/buf/cmd/buf@latest
