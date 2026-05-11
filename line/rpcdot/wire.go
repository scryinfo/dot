package rpcdot

import (
	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	contextex "github.com/scryinfo/dot/line/context_ex"
)

var Newer = wire.NewSet(
	contextex.NewContextEx,
	dot.NewLogger,
	NewConnectHttpServerMux,
	NewHandlerMiddle,
	NewConnetServer,
	NewConnectServerEtcd,
	NewBothHttpServer,

	NewGrpcServer,
	NewHttpClientEx,
	NewGrpcClientEtcd,
	NewGrpcClientEx,
)
