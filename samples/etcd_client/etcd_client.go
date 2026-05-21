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
	EtcdClient *etcddot.Client
}

type LineConfig struct {
	Log        dot.LogConfig        `json:"log" toml:"log" yaml:"log"`
	EtcdClient etcddot.ClientConfig `json:"etcdClient" toml:"etcdClient" yaml:"etcdClient"`
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
	wire.FieldsOf(new(*LineConfig), "Log", "EtcdClient"),
	contextex.NewContextEx,
	etcddot.NewClient,
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

	err = line.EtcdClient.Ping()
	if err != nil {
		dot.Logger.Error().Err(err).Msg("ping failed")
		return
	}
	dot.Logger.Info().Msg("ping success")

	ssignal.WaitCtrlC(func(s os.Signal) bool {
		return false
	})
	dot.Logger.Info().Msg("line exist")
}
