// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"context"
	"os"

	"connectrpc.com/connect"
	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line"
	"github.com/scryinfo/dot/line/rpcdot"
	"github.com/scryinfo/dot/line/sconfig"
	apiv1 "github.com/scryinfo/dot/samples/rpc/go_out/connect/api/v1"
	"github.com/scryinfo/dot/samples/rpc/go_out/connect/api/v1/apiv1connect"
	"github.com/scryinfo/scryg/sutils/ssignal"
)

type Line struct {
	SConfig         dot.SConfig
	Logger          *dot.LoggerType
	HiServiceClient apiv1connect.HiServiceClient
}

type LineConfig struct {
	Log        dot.LogConfig           `json:"log" toml:"log" yaml:"log"`
	HttpClient rpcdot.HttpClientConfig `json:"http_client" toml:"http_client" yaml:"http_client"`
}

func NewLineConfig(config *sconfig.SConfig) (*LineConfig, error) {
	lineConfig, err := sconfig.NewLineConfig[LineConfig](config)
	if err != nil {
		return nil, err
	}
	return sconfig.GenerateConfigWithArgs(config, lineConfig)
}
func NewHiServiceClient(httpClientEx *rpcdot.HttpClientEx) apiv1connect.HiServiceClient {
	// return apiv1connect.NewHiServiceClient(httpClientEx.Client(), httpClientEx.ServerAddress(), connect.WithGRPC())
	return apiv1connect.NewHiServiceClient(httpClientEx.Client(), httpClientEx.ServerAddress())
}

var LineSet = wire.NewSet(
	wire.Struct(new(Line), "*"),
	wire.FieldsOf(new(*LineConfig), "Log", "HttpClient"),
	NewLineConfig,
	line.SconfigNewConfig,
	wire.Bind(new(dot.SConfig), new(*sconfig.SConfig)),
	dot.NewLogger,
	NewHiServiceClient,
	line.RpcdotNewHttpClientEx,
	line.CertificateNewBaseCertificate,
)

func main() {
	// dot.InitLogger(new(dot.TestLogConfig()))
	line, clear, err := InitializeService()
	if err != nil {
		if line != nil && line.Logger != nil {
			line.Logger.Error().Err(err).Msg("initialize service failed")
		}
		return
	}
	if clear != nil {
		defer clear()
	}

	dot.Logger.Info().Msg("dot ok")
	//second step ....
	_ = line

	res, err := line.HiServiceClient.Hi(context.Background(), &connect.Request[apiv1.HiRequest]{
		Msg: &apiv1.HiRequest{
			Name: "test",
		},
	})
	if err != nil {
		dot.Logger.Error().Err(err).Msg("hi failed")
		return
	}

	dot.Logger.Info().Msgf("hi response: %v", res.Msg)

	ssignal.WaitCtrlC(func(s os.Signal) bool { //third wait for exit
		return false
	})
}
