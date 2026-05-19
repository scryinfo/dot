// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"flag"
	"os"

	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/lib/kits"
	"github.com/scryinfo/dot/line/sconfig"
	"github.com/scryinfo/scryg/sutils/ssignal"
)

type Line struct {
	SConfig dot.SConfig
	Logger  *dot.LoggerType
}

type LineConfig struct {
	Log dot.LogConfig `json:"log" toml:"log" yaml:"log"`
}

func NewLineConfig(config *sconfig.SConfig) (*LineConfig, error) {
	return sconfig.NewLineConfig[LineConfig](config)
}

var LineSet = wire.NewSet(
	wire.Struct(new(Line), "*"),
	wire.FieldsOf(new(*LineConfig), "Log"),
	NewLineConfig,
	sconfig.NewConfig,
	wire.Bind(new(dot.SConfig), new(*sconfig.SConfig)),
	dot.NewLogger,
)

func main() {
	// dot.InitLogger(new(dot.TestLogConfig()))

	line, clear, err := InitializeService()
	if err != nil {
		dot.Logger.Error().Err(err).Msg("initialize service failed")
		return
	}
	{
		makeConfig := flag.Bool("MakeConfig", false, "make config file from the config struct")
		flag.Parse()
		if *makeConfig {
			var config LineConfig
			err := kits.Config.MakeConfig(line.SConfig, &config)
			if err != nil {
				line.Logger.Error().Err(err).Msg("make config failed")
			}
			line.Logger.Info().Msg("make config success")
		}
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
