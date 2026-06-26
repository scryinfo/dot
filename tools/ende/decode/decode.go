// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"fmt"
	"os"

	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/sconfig"
	"github.com/scryinfo/dot/tools/ende"
	"github.com/scryinfo/scryg/sutils/ssignal"
)

type Line struct {
	SConfig *sconfig.SConfig
	Logger  *dot.LoggerType
	Decode  *ende.Decode
}

type LineConfig struct {
	Log    dot.LogConfig      `toml:"log" yaml:"log" json:"log" mapstructure:"log"`
	Decode *ende.DecodeConfig `toml:"decode" yaml:"decode" json:"decode" mapstructure:"decode"`
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
	sconfig.NewConfig,
	dot.NewLogger,
	wire.FieldsOf(new(*LineConfig), "Log", "Decode"),
	ende.NewDecode,
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
