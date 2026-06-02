// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"os"

	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line"
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
	HiService         *connectimpl.HiService
	ConnectServer     *rpcdot.ConnectServer
	ConnectServerEtcd *rpcdot.ConnectServerEtcd
}

type LineConfig struct {
	Log               dot.LogConfig                  `json:"log" toml:"log" yaml:"log"`
	ConnectServerEtcd rpcdot.ConnectServerEtcdConfig `json:"connect_server_etcd" toml:"connect_server_etcd" yaml:"connect_server_etcd"`
	ConnectServer     rpcdot.ConnectServerConfig     `json:"connect_server" toml:"connect_server" yaml:"connect_server"`
	EtcdClient        etcddot.ClientConfig           `json:"etcd_client" toml:"etcd_client" yaml:"etcd_client"`
	HiService         connectimpl.HiServiceConfig    `json:"hi_service" toml:"hi_service" yaml:"hi_service"`
}

func NewLineConfig(config *sconfig.SConfig) (*LineConfig, error) {
	lineConfig, err := sconfig.NewLineConfig[LineConfig](config)
	if err != nil {
		return nil, err
	}
	return sconfig.GenerateConfigWithArgs(config, lineConfig)
}
func NewHandlerMiddle() rpcdot.HandlerMiddle {
	return nil
}

var LineSet = wire.NewSet(
	wire.Struct(new(Line), "*"),
	wire.FieldsOf(new(*LineConfig), "Log", "ConnectServerEtcd", "ConnectServer", "EtcdClient", "HiService"),
	NewLineConfig,
	line.SconfigNewConfig,
	wire.Bind(new(dot.SConfig), new(*sconfig.SConfig)),
	dot.NewLogger,
	line.ContextexNewContextEx,
	NewHandlerMiddle,
	connectimpl.NewHiService,
	line.RpcdotNewConnectServerEtcd,
	line.RpcdotNewConnectHttpServerMux,
	line.RpcdotNewConnetServer,
	line.EtcddotNewClient,
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

	dot.Logger.Info().Msg("dot ok")
	//second step ....
	_ = line

	ssignal.WaitCtrlC(func(s os.Signal) bool { //third wait for exit
		return false
	})
}
