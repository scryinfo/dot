// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"os"

	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/sconfig"
	"github.com/scryinfo/scryg/sutils/ssignal"
)

type App struct {
	SConfig *sconfig.SConfig
	Logger  *dot.LoggerType
}

type AppConfig struct {
	Log dot.LogConfig
}

func NewAppConfig(config *sconfig.SConfig) (*AppConfig, error) {
	return sconfig.NewAppConfig[AppConfig](config)
}

var AppSet = wire.NewSet(
	NewAppConfig,
	wire.Struct(new(App), "*"),
	sconfig.NewConfig,
	dot.InitLogger,
	wire.FieldsOf(new(*AppConfig), "Log"),
)

func main() {
	// dot.InitLogger(new(dot.TestLogConfig()))
	app, close, err := InitializeService()
	if err != nil {
		dot.Logger.Error().Err(err).Msg("initialize service failed")
		return
	}
	if close != nil {
		defer close()
	}

	dot.Logger.Info().Msg("dot ok")
	//second step ....
	_ = app

	ssignal.WaitCtrlC(func(s os.Signal) bool { //third wait for exit
		return false
	})
}
