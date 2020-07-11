module github.com/scryinfo/dot/dots/grpc

go 1.14

require (
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/gin-gonic/gin v1.5.0
	github.com/golang/mock v1.4.3
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/improbable-eng/grpc-web v0.12.0
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/pkg/errors v0.9.1
	github.com/rs/cors v1.7.0 // indirect
	github.com/scryinfo/dot v0.1.5-0.20200620045039-00dee0543f89
	github.com/scryinfo/dot/dots/gindot v0.0.0-20200520093457-f8a16513551b
	github.com/stretchr/testify v1.4.0
	go.etcd.io/etcd v0.0.0-20200520232829-54ba9589114f
	go.uber.org/zap v1.14.0
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a
	google.golang.org/grpc v1.28.0
)

replace (
	github.com/scryinfo/dot => ../../
	github.com/scryinfo/dot/dots/gindot v0.0.0-20200311030916-18de37ac25e4 => ../gindot
)
