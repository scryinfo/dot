package pebbleservice

import (
	"connectrpc.com/connect"
	"github.com/scryinfo/dot/line/db/pebble_service/kv_gen/connect/kv/v1/kvv1connect"
	"github.com/scryinfo/dot/line/rpcdot"
)

func NewPebbleConnectClient(httpClientEx *rpcdot.HttpClientEx) kvv1connect.KvServiceClient {
	return kvv1connect.NewKvServiceClient(httpClientEx.Client(), httpClientEx.ServerAddress(), connect.WithGRPC())
}
