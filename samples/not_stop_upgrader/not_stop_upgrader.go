// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"fmt"

	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line"
	upgrader "github.com/scryinfo/dot/line/not_stop_upgrade"
	"github.com/scryinfo/dot/line/rpcdot"
	"github.com/scryinfo/dot/line/sconfig"
	"github.com/scryinfo/dot/samples/rpc/go_impl/connectimpl"
)

type Line struct {
	// SConfig           *sconfig.SConfig
	Logger              *dot.LoggerType
	HiService           *connectimpl.HiService
	ConnectServer       *rpcdot.ConnectServer
	Upgrader            *upgrader.UpgraderListener
	ConnectServerConfig *rpcdot.ConnectServerConfig
}

type LineConfig struct {
	Log              dot.LogConfig                   `json:"log" toml:"log" yaml:"log" mapstructure:"log"`
	ConnectServer    rpcdot.ConnectServerConfig      `json:"connect_server" toml:"connect_server" yaml:"connect_server" mapstructure:"connect_server"`
	HiService        connectimpl.HiServiceConfig     `json:"hi_service" toml:"hi_service" yaml:"hi_service" mapstructure:"hi_service"`
	UpgraderListener upgrader.UpgraderListenerConfig `json:"upgrader_listener" toml:"upgrader_listener" yaml:"upgrader_listener" mapstructure:"upgrader_listener"`
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
	wire.FieldsOf(new(*LineConfig), "Log", "ConnectServer", "HiService", "UpgraderListener"),
	NewLineConfig,
	line.SconfigNewConfig,
	wire.Bind(new(dot.SConfig), new(*sconfig.SConfig)),
	dot.NewLogger,
	line.RpcdotNewConnetServer,
	line.RpcdotNewConnectHttpServerMux,
	line.RpcdotNewHandlerMiddle,
	connectimpl.NewHiService,
	upgrader.NewUpgraderListener,
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
	if line.ConnectServerConfig != nil && !line.ConnectServerConfig.AutoRun {
		err := line.ConnectServer.StartWithListener(line.Upgrader.Listener)
		if err != nil {
			line.Logger.Error().Err(err).Msg("start connect server failed")
		}
	}

	if line.Upgrader != nil && line.Upgrader.WaitUpgrader != nil {
		line.Upgrader.WaitUpgrader()
	} else {
		line.Logger.Error().Msg("upgrader wait func is nil")
	}

	dot.Logger.Info().Msg("line exist")
}
