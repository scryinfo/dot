// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"os"

	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
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
	HiService         *connectimpl.HiService
	ConnectServer     *rpcdot.ConnectServer
	ConnectServerEtcd *rpcdot.ConnectServerEtcd
}

type LineConfig struct {
	Log               dot.LogConfig
	ConnectServerEtcd rpcdot.ConnectServerEtcdConfig
	ConnectServer     rpcdot.ConnectServerConfig
	EtcdClient        etcddot.ClientConfig
	HiService         connectimpl.HiServiceConfig
}

func NewLineConfig(config *sconfig.SConfig) (*LineConfig, error) {
	return sconfig.NewLineConfig[LineConfig](config)
}
func NewHandlerMiddle() rpcdot.HandlerMiddle {
	return nil
}

var LineSet = wire.NewSet(
	wire.Struct(new(Line), "*"),
	wire.FieldsOf(new(*LineConfig), "Log", "ConnectServerEtcd", "ConnectServer", "EtcdClient", "HiService"),
	NewLineConfig,
	sconfig.NewConfig,
	wire.Bind(new(dot.SConfig), new(*sconfig.SConfig)),
	dot.NewLogger,
	contextex.NewContextEx,
	NewHandlerMiddle,
	connectimpl.NewHiService,
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

	dot.Logger.Info().Msg("dot ok")
	//second step ....
	_ = line

	ssignal.WaitCtrlC(func(s os.Signal) bool { //third wait for exit
		return false
	})
}
