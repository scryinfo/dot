// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"flag"
	"os"

	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/lib/kits"
	contextex "github.com/scryinfo/dot/line/context_ex"
	"github.com/scryinfo/dot/line/etcddot"
	"github.com/scryinfo/dot/line/rpcdot"
	"github.com/scryinfo/dot/line/sconfig"
	"github.com/scryinfo/dot/samples/rpc/go_impl/connectimpl"
	"github.com/scryinfo/scryg/sutils/ssignal"
)

// sdf
type Line struct {
	SConfig           dot.SConfig
	Logger            *dot.LoggerType
	EtcdServer        *etcddot.Server
	HiService         *connectimpl.HiService
	ConnectServer     *rpcdot.ConnectServer
	ConnectServerEtcd *rpcdot.ConnectServerEtcd
}

type LineConfig struct {
	Log               dot.LogConfig                  `json:"log" toml:"log" yaml:"log"`
	EtcdServer        etcddot.ServerConfig           `json:"etcdServer" toml:"etcdServer" yaml:"etcdServer"`
	ConnectServerEtcd rpcdot.ConnectServerEtcdConfig `json:"connectServerEtcd" toml:"connectServerEtcd" yaml:"connectServerEtcd"`
	ConnectServer     rpcdot.ConnectServerConfig     `json:"connectServer" toml:"connectServer" yaml:"connectServer"`
	EtcdClient        etcddot.ClientConfig           `json:"etcdClient" toml:"etcdClient" yaml:"etcdClient"`
	HiService         connectimpl.HiServiceConfig    `json:"hiService" toml:"hiService" yaml:"hiService"`
}

func NewLineConfig(config *sconfig.SConfig) (*LineConfig, error) {
	return sconfig.NewLineConfig[LineConfig](config)
}
func NewHandlerMiddle() rpcdot.HandlerMiddle {
	return nil
}

var LineSet = wire.NewSet(
	wire.Struct(new(Line), "*"),
	wire.FieldsOf(new(*LineConfig), "Log", "EtcdServer", "ConnectServerEtcd", "ConnectServer", "EtcdClient", "HiService"),
	NewLineConfig,
	sconfig.NewConfig,
	wire.Bind(new(dot.SConfig), new(*sconfig.SConfig)),
	dot.NewLogger,
	contextex.NewContextEx,
	NewHandlerMiddle,
	connectimpl.NewHiService,
	etcddot.NewServer,
	rpcdot.NewConnectServerEtcd,
	rpcdot.NewConnectHttpServerMux,
	rpcdot.NewConnetServer,
	etcddot.NewClient,
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
	{
		makeConfig := flag.Bool("MakeConfig", false, "make config file from the config struct")
		flag.Parse()
		if *makeConfig {
			var config LineConfig
			err := kits.Config.MakeConfig(line.SConfig, &config)
			if err != nil {
				line.Logger.Error().Err(err).Msg("make config failed")
			}
			line.Logger.Info().Msg("make config success")
			return
		}
	}

	dot.Logger.Info().Msg("dot ok")
	//second step ....
	_ = line

	ssignal.WaitCtrlC(func(s os.Signal) bool { //third wait for exit
		return false
	})
}
