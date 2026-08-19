
clean:
	rm -rf ./data_gen/connect/*
gen:
	buf generate
buf_dep:
	buf dep update
build:
	command go build ./...

go_tools:
	command go install github.com/bufbuild/buf/cmd/buf@latest
