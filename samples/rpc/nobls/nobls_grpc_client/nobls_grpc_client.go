// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line"
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
	Log        dot.LogConfig           `json:"log" toml:"log" yaml:"log" mapstructure:"log"`
	GrpcClient rpcdot.GrpcClientConfig `json:"grpc_client" toml:"grpc_client" yaml:"grpc_client" mapstructure:"grpc_client"`
}

func NewLineConfig(config *sconfig.SConfig) (*LineConfig, error) {
	lineConfig, err := sconfig.NewLineConfig[LineConfig](config)
	if err != nil {
		return nil, err
	}
	return sconfig.GenerateConfigWithArgs(config, lineConfig)
}
func NewHiServiceClient(clientEx *rpcdot.GrpcClientEx) apiv1grpc.HiServiceClient {
	return apiv1grpc.NewHiServiceClient(clientEx.Client())
}

var LineSet = wire.NewSet(
	wire.Struct(new(Line), "*"),
	wire.FieldsOf(new(*LineConfig), "Log", "GrpcClient"),
	NewLineConfig,
	line.SconfigNewConfig,
	wire.Bind(new(dot.SConfig), new(*sconfig.SConfig)),
	dot.NewLogger,
	NewHiServiceClient,
	line.RpcdotNewGrpcClientEx,
	line.CertificateNewBaseCertificate,
)

func main() {
	// dot.InitLogger(new(dot.TestLogConfig()))
	line, clean, err := InitializeService()
	if err != nil {
		if line != nil && line.Logger != nil {
			line.Logger.Error().Err(err).Msg("initialize service failed")
		} else {
			fmt.Printf("%s\n", err.Error())
		}
		return
	}
	if clean != nil {
		defer clean()
	}

	dot.Logger.Info().Msg("dot ok")
	//second step ....
	_ = line

	res, err := line.HiServiceClient.Hi(context.Background(), &apiv1grpc.HiRequest{Name: "test"})
	if err != nil {
		dot.Logger.Error().Err(err).Msg("hi failed")
		return
	}

	dot.Logger.Info().Msgf("hi response: %v", res.Name)

	ssignal.WaitCtrlC(func(s os.Signal) bool { //third wait for exit
		return false
	})
}
