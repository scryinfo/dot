module github.com/scryinfo/dot/dots/grpc

go 1.12

require (
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.4.0
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/improbable-eng/grpc-web v0.9.6
	github.com/json-iterator/go v1.1.7 // indirect
	github.com/mattn/go-isatty v0.0.10 // indirect
	github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223 // indirect
	github.com/pkg/errors v0.8.1
	github.com/rs/cors v1.6.0 // indirect
	github.com/scryinfo/dot v0.1.3
	github.com/scryinfo/dot/dots/gindot v0.0.0-20191026032307-4fe8cc8e04c9
	github.com/ugorji/go v1.1.7 // indirect
	go.uber.org/multierr v1.2.0 // indirect
	go.uber.org/zap v1.11.0
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859
	golang.org/x/sys v0.0.0-20191025090151-53bf42e6b339 // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/grpc v1.21.1
	gopkg.in/yaml.v2 v2.2.4 // indirect
)

replace (
	github.com/scryinfo/dot v0.1.3 => ../../
	github.com/scryinfo/dot/dots/gindot v0.0.0-20191026032307-4fe8cc8e04c9 => ../gindot
)
