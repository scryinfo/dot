// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"context"
	"os"

	"connectrpc.com/connect"
	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/rpcdot"
	"github.com/scryinfo/dot/line/sconfig"
	apiv1 "github.com/scryinfo/dot/samples/rpc/go_out/connect/api/v1"
	"github.com/scryinfo/dot/samples/rpc/go_out/connect/api/v1/apiv1connect"
	"github.com/scryinfo/scryg/sutils/ssignal"
)

type Line struct {
	SConfig         *sconfig.SConfig
	Logger          *dot.LoggerType
	HiServiceClient apiv1connect.HiServiceClient
}

type LineConfig struct {
	Log        dot.LogConfig           `json:"log" toml:"log" yaml:"log"`
	HttpClient rpcdot.HttpClientConfig `json:"httpClient" toml:"httpClient" yaml:"httpClient"`
}

func NewLineConfig(config *sconfig.SConfig) (*LineConfig, error) {
	return sconfig.NewLiceConfig[LineConfig](config)
}
func NewHiServiceClient(httpClientEx *rpcdot.HttpClientEx) apiv1connect.HiServiceClient {
	return apiv1connect.NewHiServiceClient(httpClientEx.Client(), httpClientEx.ServerAddress(), connect.WithGRPC())
}

var LineSet = wire.NewSet(
	wire.Struct(new(Line), "*"),
	wire.FieldsOf(new(*LineConfig), "Log", "HttpClient"),
	NewLineConfig,
	sconfig.NewConfig,
	dot.NewLogger,
	NewHiServiceClient,
	rpcdot.NewHttpClientEx,
)

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
