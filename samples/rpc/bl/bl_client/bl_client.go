// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"os"

	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/etcddot"
	"github.com/scryinfo/dot/line/rpcdot"
	"github.com/scryinfo/dot/line/sconfig"
	apiv1grpc "github.com/scryinfo/dot/samples/rpc/go_out/gogrpc/api/v1"
	"github.com/scryinfo/scryg/sutils/ssignal"
)

type Line struct {
	SConfig         *sconfig.SConfig
	Logger          *dot.LoggerType
	HiServiceClient apiv1grpc.HiServiceClient
}

type LineConfig struct {
	Log            dot.LogConfig
	GrpcClientEtcd rpcdot.GrpcClientEtcdConfig
	EtcdClient     etcddot.ClientConfig
}

func NewLineConfig(config *sconfig.SConfig) (*LineConfig, error) {
	return sconfig.NewLiceConfig[LineConfig](config)
}

var LineSet = wire.NewSet(
	wire.Struct(new(Line), "*"),
	wire.FieldsOf(new(*LineConfig), "Log", "GrpcClientEtcd", "EtcdClient"),
	NewLineConfig,
	sconfig.NewConfig,
	dot.NewLogger,
	rpcdot.NewGrpcClientEtcd,
	NewHiServiceClient,
	etcddot.NewClient,
)

func NewHiServiceClient(clientEtcd *rpcdot.GrpcClientEtcd) apiv1grpc.HiServiceClient {
	return apiv1grpc.NewHiServiceClient(clientEtcd.Client())
}

func main() {
	// dot.InitLogger(new(dot.TestLogConfig()))
	line, clear, err := InitializeService()
	if err != nil {
		dot.Logger.Error().Err(err).Msg("initialize service failed")
		return
	}
	if clear != nil {
		defer clear()
	}

	dot.Logger.Info().Msg("dot ok")
	//second step ....
	_ = line

	ssignal.WaitCtrlC(func(s os.Signal) bool { //third wait for exit
		return false
	})
}
