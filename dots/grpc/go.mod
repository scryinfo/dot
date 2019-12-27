module github.com/scryinfo/dot/dots/grpc

go 1.12

require (
	github.com/gin-gonic/gin v1.5.0
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/improbable-eng/grpc-web v0.9.6
	github.com/pkg/errors v0.8.1
	github.com/rs/cors v1.7.0 // indirect
	github.com/scryinfo/dot v0.1.3
	github.com/scryinfo/dot/dots/gindot v0.0.0-20191121022614-959828ad21d4
	go.uber.org/zap v1.13.0
	golang.org/x/net v0.0.0-20190923162816-aa69164e4478
	google.golang.org/grpc v1.21.1
)

replace (
	github.com/scryinfo/dot v0.1.3 => ../../
	github.com/scryinfo/dot/dots/gindot v0.0.0-20191121022614-959828ad21d4 => ../gindot
)
