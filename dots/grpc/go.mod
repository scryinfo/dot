module github.com/scryinfo/dot/dots/grpc

go 1.12

require (
	github.com/gin-gonic/gin v1.4.0
	github.com/improbable-eng/grpc-web v0.9.6
	github.com/pkg/errors v0.8.1
	github.com/scryinfo/dot v0.1.3
	github.com/scryinfo/dot/dots/gindot v0.0.0-20191121022614-959828ad21d4
	go.uber.org/zap v1.13.0
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859
	google.golang.org/grpc v1.21.1
)

replace (
	github.com/scryinfo/dot v0.1.3 => ../../
	github.com/scryinfo/dot/dots/gindot v0.0.0-20191121022614-959828ad21d4 => ../gindot
)
