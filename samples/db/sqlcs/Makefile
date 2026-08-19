
.PHONY: clean gen

clean:
	rm -rf db/db.go db/models.go db/query.sql.go db/gen_to/*
gen:
	cd proto && make gen
	sqlc generate
	command go generate ./...
