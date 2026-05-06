[中文](./README-cn.md)  
[EN](./README.md)

# dot

Component development specification, including component definition, component dependencies, component life cycle, dependency injection, and common basic components

- Dot: A component which has no type or interface requirements, anything can be a component
- Line: the collect of dot components

- Injecter((wire)[github.com/google/wire]) ：It is component dependency injection, base on wire

[the simple samp1e](./samples/simple)

```go
// file name: simple.go
// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"os"

	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/sconfig"
	"github.com/scryinfo/scryg/sutils/ssignal"
)

type Line struct {
	SConfig *sconfig.SConfig
	Logger  *dot.LoggerType
}

type LineConfig struct {
	Log dot.LogConfig
}

func NewLineConfig(config *sconfig.SConfig) (*LineConfig, error) {
	return sconfig.NewLiceConfig[LineConfig](config)
}

var LineSet = wire.NewSet(
	NewLineConfig,
	wire.Struct(new(Line), "*"),
	sconfig.NewConfig,
	dot.NewLogger,
	wire.FieldsOf(new(*LineConfig), "Log"),
)

func main() {
	line, clean, err := InitializeService()
	if err != nil {
		dot.Logger.Error().Err(err).Msg("initialize service failed")
		return
	}
	if clean != nil {
		defer clean()
	}
	dot.Logger.Info().Msg("line run")

	_ = line

	ssignal.WaitCtrlC(func(s os.Signal) bool {
		return false
	})
	dot.Logger.Info().Msg("line exist")
}

```

how to run it：

```bash
# make the file wire.go
# run the command
wire
# config the simple.toml file
# run the command
go run simple.go wire_gen.go
```

# Default components

## (Config)[./line/sconfig]

the json/toml/yaml format, command line, and environment variables.

## (Log)[./line/slog]

High performance logs based on zerolog.

## (Certificate generated)[./line/certificate]

Generate root and sub certificates. "sample/certificate" is an example.

## (db)[./line/db]

## (etcd)[./line/etcd]

etcd (client)[go.etcd.io/etcd/client/v3] and (server)[go.etcd.io/etcd/server/v3].

## (gindot)[./line/gindot]

dot for (gin)[github.com/gin-gonic/gin]

## (jsonrpc2)[./line/jsonrpc2]

## (rpcdot)[./line/rpcdot]

dot for (grpc)[google.golang.org/grpc] and (connect-rpc)[github.com/connectrpc/connect-go]  
dot HandlerMiddle: auth middleware for connect-rpc or grpc  
dot connect-rpc - HttpClientEx:
dot connect-rpc - ConnectHttpServerMux:
dot connect-rpc - ConnectServer:
dot connect-rpc - ConnectServerEtcd: etcd registry for connect-rpc
dot grpc - GrpcConnectEx: grpc connect middleware
dot grpc - GrpcClientEtcd: grpc connect with etcd
dot grpc - grpc.Server: grpc server

# [Code Style -- Go](https://github.com/scryinfo/scryg/blob/master/codestyle_go.md)
