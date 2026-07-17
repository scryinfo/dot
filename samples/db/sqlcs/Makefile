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
clean:
	rm -rf db/db.go db/models.go db/query.sql.go db/gen_to/*
gen:
	cd proto && make gen
	sqlc generate
	${go} generate ./...
