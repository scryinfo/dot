// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"fmt"
	"os"

	"github.com/dgraph-io/badger/v4"
	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line"
	"github.com/scryinfo/dot/line/db/badgerdot"
	"github.com/scryinfo/dot/line/sconfig"
	"github.com/scryinfo/scryg/sutils/ssignal"
)

type Line struct {
	SConfig dot.SConfig
	Logger  *dot.LoggerType
	Badger  *badgerdot.BadgerDbDot
}

type LineConfig struct {
	Log      dot.LogConfig               `json:"log" toml:"log" yaml:"log" mapstructure:"log"`
	BadgerDb badgerdot.BadgerDbDotConfig `json:"badger_db" toml:"badger_db" yaml:"badger_db" mapstructure:"badger_db"`
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
	wire.FieldsOf(new(*LineConfig), "Log", "BadgerDb"),
	line.DbBadgerdotNewBadgerDot,
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
	db := line.Badger.Db()
	err = db.Update(func(txn *badger.Txn) error {
		key := []byte("test")
		v := []byte("test value")
		err2 := txn.Set(key, v)
		if err2 != nil {
			return err2
		}
		vGet, err2 := txn.Get(key)
		if err2 != nil {
			return err2
		}
		err2 = vGet.Value(func(val []byte) error {
			dot.Logger.Info().Msgf("get value: %s, set value: %s", string(val), string(v))
			return nil
		})

		return err2
	})
	if err != nil {
		dot.Logger.Error().Err(err).Msg("set failed")
	}

	ssignal.WaitCtrlC(func(s os.Signal) bool {
		return false
	})
	dot.Logger.Info().Msg("line exist")
}
