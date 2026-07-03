package pebbleservice

import (
	kvv1 "github.com/scryinfo/dot/line/db/pebble_service/kv_gen/connect/kv/v1"
	"github.com/scryinfo/dot/line/rpcdot"
)

func NewPebbleGrpcClient(clientEx *rpcdot.GrpcClientEx) kvv1.KvServiceClient {
	return kvv1.NewKvServiceClient(clientEx.Client())
}
