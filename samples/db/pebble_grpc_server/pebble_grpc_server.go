// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"fmt"
	"os"

	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line"
	"github.com/scryinfo/dot/line/db/pebble2dot"
	pebbleservice "github.com/scryinfo/dot/line/db/pebble_service"
	"github.com/scryinfo/dot/line/rpcdot"
	"github.com/scryinfo/dot/line/sconfig"
	"github.com/scryinfo/scryg/sutils/ssignal"
)

type Line struct {
	SConfig       dot.SConfig
	Logger        *dot.LoggerType
	Pebble        *pebble2dot.Pebble2
	PebbleService *pebbleservice.PebbleGrpcService
}

type LineConfig struct {
	Log        dot.LogConfig            `json:"log" toml:"log" yaml:"log" mapstructure:"log"`
	Pebble2    pebble2dot.Pebble2Config `json:"pebble2" toml:"pebble2" yaml:"pebble2" mapstructure:"pebble2"`
	GrpcServer rpcdot.GrpcServerConfig  `json:"grpc_server" toml:"grpc_server" yaml:"grpc_server" mapstructure:"grpc_server"`
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
	wire.FieldsOf(new(*LineConfig), "Log", "Pebble2", "GrpcServer"),
	line.DbPebble2dotNewPebble2,
	line.DbPebbleServiceNewPebbleGrpcService,
	line.RpcdotNewGrpcServer,
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

	_ = line

	ssignal.WaitCtrlC(func(s os.Signal) bool {
		return false
	})
	dot.Logger.Info().Msg("line exist")
}
