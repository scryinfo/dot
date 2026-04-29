// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"os"

	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/rpcdot"
	"github.com/scryinfo/dot/line/sconfig"
	"github.com/scryinfo/dot/samples/rpc/go_impl/connectimpl"
	"github.com/scryinfo/scryg/sutils/ssignal"
)

type Line struct {
	// SConfig           *sconfig.SConfig
	// Logger            *dot.LoggerType
	HiService         *connectimpl.HiService
	ConnectHttpServer *rpcdot.ConnectServer
}

type LineConfig struct {
	Log           dot.LogConfig
	ConnectServer rpcdot.ConnectServerConfig
}

func NewLineConfig(config *sconfig.SConfig) (*LineConfig, error) {
	return sconfig.NewLiceConfig[LineConfig](config)
}

func NewHandlerMiddle() rpcdot.HandlerMiddle {
	return nil
}

var LineSet = wire.NewSet(
	wire.Struct(new(Line), "*"),
	wire.FieldsOf(new(*LineConfig), "Log", "ConnectServer"),
	NewLineConfig,
	sconfig.NewConfig,
	dot.NewLogger,
	rpcdot.NewConnetHttpServer,
	rpcdot.NewConnectHttpServerMux,
	NewHandlerMiddle,
	connectimpl.NewHiService,
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

	ssignal.WaitCtrlC(func(s os.Signal) bool { //third wait for exit
		return false
	})
}
