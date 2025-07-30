.PHONY: clean upgrade format build

go := ${shell which go}

clean:
	rm -f go.sum go.work.sum demo/go.sum
	${go} clean
	cd demo && ${go} clean
upgrade:
	${go} get -t -u all && ${go} mod tidy
format:
	${go} fmt ./...
build:
	${go} build ./...
