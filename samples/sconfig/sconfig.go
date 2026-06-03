// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"fmt"
	"os"

	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line"
	"github.com/scryinfo/dot/line/sconfig"
	"github.com/scryinfo/scryg/sutils/ssignal"
)

type Line struct {
	SConfig dot.SConfig
	Logger  *dot.LoggerType
}

type LineConfig struct {
	Log dot.LogConfig `json:"log" toml:"log" yaml:"log" mapstructure:"log"`
}

func NewLineConfig(config *sconfig.SConfig) (*LineConfig, error) {
	lineConfig, err := sconfig.NewLineConfig[LineConfig](config)
	if err != nil {
		return nil, err
	}
	return sconfig.GenerateConfigWithArgs(config, lineConfig)
}

var LineSet = wire.NewSet(
	wire.Struct(new(Line), "*"),
	wire.FieldsOf(new(*LineConfig), "Log"),
	NewLineConfig,
	line.SconfigNewConfig,
	wire.Bind(new(dot.SConfig), new(*sconfig.SConfig)),
	dot.NewLogger,
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

	ssignal.WaitCtrlC(func(s os.Signal) bool { //third wait for exit
		return false
	})
}
