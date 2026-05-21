// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"os"

	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	contextex "github.com/scryinfo/dot/line/context_ex"
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
	Log        dot.LogConfig        `json:"log" toml:"log" yaml:"log"`
	EtcdServer etcddot.ServerConfig `json:"etcdServer" toml:"etcdServer" yaml:"etcdServer"`
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
	wire.FieldsOf(new(*LineConfig), "Log", "EtcdServer"),
	contextex.NewContextEx,
	etcddot.NewServer,
)

func main() {
	line, clean, err := InitializeService()
	if err != nil {
		dot.Logger.Error().Err(err).Msg("initialize service failed")
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
