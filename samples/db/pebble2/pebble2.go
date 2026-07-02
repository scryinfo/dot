// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/cockroachdb/pebble/v2"
	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line"
	"github.com/scryinfo/dot/line/db/pebble2dot"
	"github.com/scryinfo/dot/line/sconfig"
	"github.com/scryinfo/scryg/sutils/ssignal"
)

type Line struct {
	SConfig dot.SConfig
	Logger  *dot.LoggerType
	Pebble  *pebble2dot.Pebble2
}

type LineConfig struct {
	Log     dot.LogConfig            `json:"log" toml:"log" yaml:"log" mapstructure:"log"`
	Pebble2 pebble2dot.Pebble2Config `json:"pebble2" toml:"pebble2" yaml:"pebble2" mapstructure:"pebble2"`
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
	wire.FieldsOf(new(*LineConfig), "Log", "Pebble2"),
	line.DbPebble2dotNewPebble2,
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

	key := []byte("test")
	value := []byte("test value")
	err2 := line.Pebble.Db().Set(key, value, nil)
	if err2 != nil {
		dot.Logger.Error().Err(err2).Msg("set failed")
	}

	val, closer, err := line.Pebble.Db().Get(key)
	if err == pebble.ErrNotFound {
		dot.Logger.Info().Msg("not found")
	} else if err != nil {
		dot.Logger.Error().Err(err).Msg("get failed")
	}

	defer closer.Close()
	if bytes.Equal(val, value) {
		dot.Logger.Info().Msg("get success")
	} else {
		dot.Logger.Error().Msg("get failed")
	}

	ssignal.WaitCtrlC(func(s os.Signal) bool {
		return false
	})
	dot.Logger.Info().Msg("line exist")
}
