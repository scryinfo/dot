module github.com/scryinfo/dot/dots/grpc

go 1.12

require (
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.4.0
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/improbable-eng/grpc-web v0.9.6
	github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223 // indirect
	github.com/pkg/errors v0.8.1
	github.com/rs/cors v1.6.0 // indirect
	github.com/scryinfo/dot v0.1.3-0.20190705064446-6614e45bf155
	github.com/scryinfo/dot/dots/gindot v0.0.0-20190705064650-8b2f44b376f8
	github.com/scryinfo/scryg v0.1.3-0.20190608053141-a292b801bfd6
	go.uber.org/zap v1.10.0
	golang.org/x/crypto v0.0.0-20190605123033-f99c8df09eb5 // indirect
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859
	golang.org/x/sys v0.0.0-20190624142023-c5567b49c5d0 // indirect
	golang.org/x/text v0.3.2 // indirect
	golang.org/x/tools v0.0.0-20190606050223-4d9ae51c2468 // indirect
	google.golang.org/grpc v1.21.1
)

replace (
	github.com/scryinfo/dot v0.1.3-0.20190625102047-666d44ee7d72 => ../../
	github.com/scryinfo/dot/dots/gindot v0.0.0-20190622091252-bab0929bd7e7 => ../gindot
)
