// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"connectrpc.com/connect"
	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line"
	pebbleservice "github.com/scryinfo/dot/line/db/pebble_service"
	kvv1 "github.com/scryinfo/dot/line/db/pebble_service/kv_gen/connect/kv/v1"
	"github.com/scryinfo/dot/line/db/pebble_service/kv_gen/connect/kv/v1/kvv1connect"
	"github.com/scryinfo/dot/line/rpcdot"
	"github.com/scryinfo/dot/line/sconfig"
	"github.com/scryinfo/scryg/sutils/ssignal"
)

type Line struct {
	SConfig      dot.SConfig
	Logger       *dot.LoggerType
	PebbleClient kvv1connect.KvServiceClient
}

type LineConfig struct {
	Log        dot.LogConfig           `json:"log" toml:"log" yaml:"log" mapstructure:"log"`
	HttpClient rpcdot.HttpClientConfig `json:"http_client" toml:"http_client" yaml:"http_client" mapstructure:"http_client"`
}

func NewLineConfig(config *sconfig.SConfig) (*LineConfig, error) {
	lineConfig, err := sconfig.NewLineConfig[LineConfig](config)
	if err != nil {
		return nil, err
	}
	return sconfig.GenerateConfigWithArgs(config, lineConfig)
}

var LineSet = wire.NewSet(
	NewLineConfig,
	wire.Struct(new(Line), "*"),
	line.SconfigNewConfig,
	wire.Bind(new(dot.SConfig), new(*sconfig.SConfig)),
	dot.NewLogger,
	wire.FieldsOf(new(*LineConfig), "Log", "HttpClient"),
	line.RpcdotNewHttpClientEx,
	pebbleservice.NewPebbleConnectClient,
	line.CertificateNewBaseCertificate,
)

func main() {
	line, clean, err := InitializeService()
	if err != nil {
		if line != nil {
			dot.Logger.Error().Err(err).Msg("initialize service failed")
		} else {
			dot.Logger.Info().Msg(err.Error())
			fmt.Printf("\n\n")
		}
		return
	}
	if clean != nil {
		defer clean()
	}
	dot.Logger.Info().Msg("line run")
	key := []byte("test")
	value := []byte("test value")
	_, err = line.PebbleClient.Set(context.Background(), &connect.Request[kvv1.SetRequest]{Msg: &kvv1.SetRequest{
		Key: key, Value: value, TtlSeconds: 0,
	}})
	if err != nil {
		dot.Logger.Error().Err(err).Send()
		return
	}

	res, err := line.PebbleClient.Get(context.Background(), &connect.Request[kvv1.GetRequest]{Msg: &kvv1.GetRequest{Key: key}})
	if err != nil {
		dot.Logger.Error().Err(err).Send()
		return
	}
	if !bytes.Equal(value, res.Msg.Value) {
		dot.Logger.Error().Msg("the get value is not eq value")
		return
	}

	ssignal.WaitCtrlC(func(s os.Signal) bool {
		return false
	})
	dot.Logger.Info().Msg("line exist")
}
