// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"fmt"
	"os"

	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line"
	"github.com/scryinfo/dot/line/etcddot"
	"github.com/scryinfo/dot/line/sconfig"
	"github.com/scryinfo/scryg/sutils/ssignal"
)

type Line struct {
	SConfig    *sconfig.SConfig
	Logger     *dot.LoggerType
	EtcdServer *etcddot.Server
}

type LineConfig struct {
	Log        dot.LogConfig        `json:"log" toml:"log" yaml:"log" mapstructure:"log"`
	EtcdServer etcddot.ServerConfig `json:"etcd_server" toml:"etcd_server" yaml:"etcd_server" mapstructure:"etcd_server"`
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
	dot.NewLogger,
	wire.FieldsOf(new(*LineConfig), "Log", "EtcdServer"),
	line.ContextexNewContextEx,
	line.EtcddotNewServer,
)

func main() {
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
	dot.Logger.Info().Msg("line run")

	_ = line

	ssignal.WaitCtrlC(func(s os.Signal) bool {
		return false
	})
	dot.Logger.Info().Msg("line exist")
}
