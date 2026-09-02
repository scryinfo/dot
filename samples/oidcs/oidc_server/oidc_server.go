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
	"github.com/scryinfo/dot/line/oidcdot"
	"github.com/scryinfo/dot/line/oidcdot/oidc_server/oidc_storage"
	"github.com/scryinfo/dot/line/rpcdot"
	"github.com/scryinfo/dot/line/sconfig"
	"github.com/scryinfo/scryg/sutils/ssignal"
)

type Line struct {
	SConfig         dot.SConfig
	Logger          *dot.LoggerType
	ConnectServer   *rpcdot.ConnectServer
	OidcServiceHttp *oidcdot.OidcServiceHttp
}

type LineConfig struct {
	Log           dot.LogConfig              `json:"log" toml:"log" yaml:"log" mapstructure:"log"`
	ConnectServer rpcdot.ConnectServerConfig `json:"connect_server" toml:"connect_server" yaml:"connect_server" mapstructure:"connect_server"`
	OidcService   oidcdot.OidcServiceConfig  `json:"oidc_service" toml:"oidc_service" yaml:"oidc_service" mapstructure:"oidc_service"`
	Pebble2       pebble2dot.Pebble2Config   `json:"pebble2" toml:"pebble2" yaml:"pebble2" mapstructure:"pebble2"`
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
	wire.FieldsOf(new(*LineConfig), "Log", "ConnectServer", "OidcService", "Pebble2"),
	line.RpcdotNewConnetServer,
	line.RpcdotNewConnectHttpServerMux,
	line.RpcdotNewHandlerMiddle,

	oidcdot.NewOidcServiceHttp,
	oidc_storage.Pebble2Set,
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

func makeTestData(line *Line) {

}
