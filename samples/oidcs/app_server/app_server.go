// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line"
	"github.com/scryinfo/dot/line/db/pebble2dot"
	"github.com/scryinfo/dot/line/oidcdot"
	"github.com/scryinfo/dot/line/rpcdot"
	"github.com/scryinfo/dot/line/sconfig"
	"github.com/scryinfo/scryg/sutils/ssignal"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Line struct {
	// SConfig           *sconfig.SConfig
	Logger        *dot.LoggerType
	AuthService   *oidcdot.AuthService
	ConnectServer *rpcdot.ConnectServer
}

type LineConfig struct {
	Log           dot.LogConfig              `json:"log" toml:"log" yaml:"log" mapstructure:"log"`
	ConnectServer rpcdot.ConnectServerConfig `json:"connect_server" toml:"connect_server" yaml:"connect_server" mapstructure:"connect_server"`
	OidcProvider  oidcdot.OidcProviderConfig `json:"oidc_provider" toml:"oidc_provider" yaml:"oidc_provider" mapstructure:"oidc_provider"`
	AuthConfig    oidcdot.AuthConfig         `json:"auth_config" toml:"auth_config" yaml:"auth_config" mapstructure:"auth_config"`
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
	wire.Struct(new(Line), "*"),
	wire.FieldsOf(new(*LineConfig), "Log", "ConnectServer", "OidcProvider", "AuthConfig", "Pebble2"),
	NewLineConfig,
	line.SconfigNewConfig,
	wire.Bind(new(dot.SConfig), new(*sconfig.SConfig)),
	dot.NewLogger,
	line.RpcdotNewConnetServer,
	line.RpcdotNewConnectHttpServerMux,
	line.RpcdotNewHandlerMiddle,
	oidcdot.NewAuthService,
	oidcdot.NewOidcProvider,
	oidcdot.OidcPebble2Set,
)

func main() {
	// cleanTracer := initTracer()
	// defer cleanTracer()

	// dot.InitLogger(new(dot.TestLogConfig()))
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

	dot.Logger.Info().Msg("dot ok")
	//second step ....
	_ = line

	ssignal.WaitCtrlC(func(s os.Signal) bool { //third wait for exit
		return false
	})
}

func initTracer() func() {
	exporter, err := stdouttrace.New(
		stdouttrace.WithWriter(os.Stdout),
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		panic(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
	)

	otel.SetTracerProvider(tp)

	return func() {
		_ = tp.Shutdown(context.Background())
	}
}
