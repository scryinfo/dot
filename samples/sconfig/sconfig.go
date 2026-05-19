// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"flag"
	"os"
	"path/filepath"

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
	Log dot.LogConfig
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
			name := filepath.Join(line.SConfig.ConfigPath(), line.SConfig.ConfigFile())
			ext := filepath.Ext(name)
			newName := name[:len(name)-len(ext)] + "_new" + ext
			kits.Config.WriteConfig(line.SConfig.ConfigFile(), &config, newName)
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
