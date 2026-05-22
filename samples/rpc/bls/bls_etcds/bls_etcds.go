// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"os"

	"golang.org/x/sync/errgroup"

	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	contextex "github.com/scryinfo/dot/line/context_ex"
	"github.com/scryinfo/dot/line/etcddot"
	"github.com/scryinfo/dot/line/sconfig"
	"github.com/scryinfo/scryg/sutils/ssignal"
)

// sdf
type Line struct {
	SConfig     dot.SConfig
	Logger      *dot.LoggerType
	EtcdServers []*etcddot.Server
}

type LineConfig struct {
	Log         dot.LogConfig          `json:"log" toml:"log" yaml:"log"`
	EtcdServers []etcddot.ServerConfig `json:"etcdServers" toml:"etcdServers" yaml:"etcdServers"`
}

func NewLineConfig(config *sconfig.SConfig) (*LineConfig, error) {
	lineConfig, err := sconfig.NewLineConfig[LineConfig](config)
	if err != nil {
		return nil, err
	}
	return sconfig.GenerateConfigWithArgs(config, lineConfig)
}

func NewEtcds(configs []etcddot.ServerConfig, ctxEx *contextex.ContextEx, logger *dot.LoggerType) ([]*etcddot.Server, func(), error) {
	etcdServers := make([]*etcddot.Server, len(configs))
	// the clear is a buildin func, so use "fClear"
	fClears := make([]func(), len(configs))
	clearsFun := func() {
		for i := len(fClears) - 1; i >= 0; i-- {
			c := fClears[i]
			if c != nil {
				c()
				fClears[i] = nil
			}
		}
	}
	group := errgroup.Group{}
	for i, _ := range configs {
		i := i
		config := &configs[i]
		group.Go(func() error {
			server, fClear, err := etcddot.NewServer(config, ctxEx, logger)
			fClears[i] = fClear
			etcdServers[i] = server
			if err != nil {
				return err
			}
			return nil
		})
	}
	err := group.Wait()
	return etcdServers, clearsFun, err
}

var LineSet = wire.NewSet(
	wire.Struct(new(Line), "*"),
	wire.FieldsOf(new(*LineConfig), "Log", "EtcdServers"),
	NewLineConfig,
	sconfig.NewConfig,
	wire.Bind(new(dot.SConfig), new(*sconfig.SConfig)),
	dot.NewLogger,
	contextex.NewContextEx,
	NewEtcds,
)

func main() {
	// dot.InitLogger(new(dot.TestLogConfig()))
	line, clear, err := InitializeService()
	if err != nil {
		if line != nil && line.Logger != nil {
			line.Logger.Error().Err(err).Msg("initialize service failed")
		}
		return
	}
	if clear != nil {
		defer clear()
	}

	line.Logger.Info().Msg("dot ok")
	//second step ....
	_ = line

	ssignal.WaitCtrlC(func(s os.Signal) bool { //third wait for exit
		return false
	})
}
