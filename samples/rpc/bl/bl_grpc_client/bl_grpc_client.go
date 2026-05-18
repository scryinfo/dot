// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"context"
	"os"

	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/certificate"
	contextex "github.com/scryinfo/dot/line/context_ex"
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
	return sconfig.NewLineConfig[LineConfig](config)
}

var LineSet = wire.NewSet(
	wire.Struct(new(Line), "*"),
	wire.FieldsOf(new(*LineConfig), "Log", "GrpcClientEtcd", "EtcdClient"),
	NewLineConfig,
	sconfig.NewConfig,
	wire.Bind(new(dot.SConfig), new(*sconfig.SConfig)),
	dot.NewLogger,
	contextex.NewContextEx,
	rpcdot.NewGrpcClientEtcd,
	NewHiServiceClient,
	etcddot.NewClient,
	certificate.NewBaseCertificate,
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

	line.Logger.Info().Msg("dot ok")
	//second step ....
	res, err := line.HiServiceClient.Hi(context.Background(), &apiv1grpc.HiRequest{Name: "ttt"})
	if err != nil {
		line.Logger.Error().Err(err).Msg("hi failed")
		return
	}
	line.Logger.Info().Msg(res.Name)
	res, err = line.HiServiceClient.Hi(context.Background(), &apiv1grpc.HiRequest{Name: "ttt2"})
	if err != nil {
		line.Logger.Error().Err(err).Msg("hi failed")
		return
	}
	line.Logger.Info().Msg(res.Name)
	res, err = line.HiServiceClient.Hi(context.Background(), &apiv1grpc.HiRequest{Name: "ttt3"})
	if err != nil {
		line.Logger.Error().Err(err).Msg("hi failed")
		return
	}
	line.Logger.Info().Msg(res.Name)

	ssignal.WaitCtrlC(func(s os.Signal) bool { //third wait for exit
		return false
	})
}
