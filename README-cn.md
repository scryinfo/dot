[中文](./README-cn.md)  
[EN](./README.md)

# dot

组件开发规范，主要有组件定义、组件依赖关系、组件生命周期、依赖注入、及常用的基础组件

- Dot: 组件，它没有类型或接口要求，都可以成为一个组件
- Line: dot的集合

- Injecter((wire)[github.com/google/wire]) ：组件依赖注入，基于wire

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

如何运行：

```bash
# 生成文件 wire.go
# 运行下面的命令
wire
# 配置配置文件 simple.toml
# 运行下面的命令
go run simple.go wire_gen.go
```

# 默认组件

## (Config)[./line/sconfig]

配置文件格式 json/toml/yaml, 命令行参数, 环境变量.

## (Log)[./line/slog]

高性能的日志组件，基于zerolog，可以配置谦容 slog

## (Certificate generated)[./line/certificate]

生成自定义证书. "sample/certificate" is an example.

## (db)[./line/db]

## (etcd)[./line/etcd]

etcd (client)[go.etcd.io/etcd/client/v3] and (server)[go.etcd.io/etcd/server/v3].

## (gindot)[./line/gindot]

dot for (gin)[github.com/gin-gonic/gin]

## (jsonrpc2)[./line/jsonrpc2]

## (rpcdot)[./line/rpcdot]

dot for (grpc)[google.golang.org/grpc] and (connect-rpc)[github.com/connectrpc/connect-go]  
dot HandlerMiddle: http auth 中间件，connect-rpc or grpc都可以使用
dot connect-rpc - HttpClientEx:
dot connect-rpc - ConnectHttpServerMux:
dot connect-rpc - ConnectServer:
dot connect-rpc - ConnectServerEtcd: etcd registry for connect-rpc
dot grpc - GrpcConnectEx: grpc connect 中间件
dot grpc - GrpcClientEtcd: grpc connect 带 etcd
dot grpc - grpc.Server: grpc 服务器

# [Code Style -- Go](https://github.com/scryinfo/scryg/blob/master/codestyle_go-cn.md)
