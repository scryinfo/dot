module github.com/scryinfo/dot/dots/grpc

go 1.12

require (
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.4.0
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/improbable-eng/grpc-web v0.9.6
	github.com/json-iterator/go v1.1.8 // indirect
	github.com/mattn/go-isatty v0.0.10 // indirect
	github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223 // indirect
	github.com/pkg/errors v0.8.1
	github.com/rs/cors v1.6.0 // indirect
	github.com/scryinfo/dot v0.1.3
	github.com/scryinfo/dot/dots/gindot v0.0.0-20191121022614-959828ad21d4
	github.com/ugorji/go v1.1.7 // indirect
	go.uber.org/atomic v1.5.1 // indirect
	go.uber.org/multierr v1.4.0 // indirect
	go.uber.org/zap v1.13.0
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859
	golang.org/x/sys v0.0.0-20191120155948-bd437916bb0e // indirect
	golang.org/x/text v0.3.2 // indirect
	golang.org/x/tools v0.0.0-20191121023328-35ba81b9fb22 // indirect
	google.golang.org/grpc v1.21.1
	gopkg.in/yaml.v2 v2.2.7 // indirect
)

replace (
	github.com/scryinfo/dot v0.1.3 => ../../
	github.com/scryinfo/dot/dots/gindot v0.0.0-20191121022614-959828ad21d4 => ../gindot
)
