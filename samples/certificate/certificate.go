// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"os"

	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/certificate"
	"github.com/scryinfo/dot/line/sconfig"
	"github.com/scryinfo/scryg/sutils/ssignal"
)

type Line struct {
	Cert    *certificate.Ecdsa
	SConfig *sconfig.SConfig
	Logger  *dot.LoggerType
}

type LineConfig struct {
	Log dot.LogConfig
}

func NewAppConfig(config *sconfig.SConfig) (*LineConfig, error) {
	return sconfig.NewLiceConfig[LineConfig](config)
}

var LineSet = wire.NewSet(
	NewAppConfig,
	wire.Struct(new(Line), "*"),
	sconfig.NewConfig,
	dot.NewLogger,
	wire.FieldsOf(new(*LineConfig), "Log"),
	certificate.NewEcdsa,
)

func main() {
	line, clear, err := InitializeService()
	if err != nil {
		dot.Logger.Error().Err(err).Msg("initialize service failed")
		return
	}
	if clear != nil {
		defer clear()
	}

	err = makeSample(line.Cert)
	if err != nil {
		dot.Logger.Error().Err(err).Send()
	}
	ssignal.WaitCtrlC(func(s os.Signal) bool { //third wait for exit
		return false
	})
}

// Generate ca certificate, generate serve and client certificate under ca certificate
func makeSample(cs *certificate.Ecdsa) error {

	caPri, err := certificate.MakePriKey()
	if err != nil {
		return err
	}

	ca, err := cs.GenerateCaCertKey(caPri, "ca.key", "ca.pem", []string{"scry"}, []string{"scry"})
	if err != nil {
		return err
	}

	err = cs.GenerateCertKey(ca, caPri, "server.key", "server.pem", []string{"scry"}, []string{"scry"})
	if err != nil {
		return err
	}

	err = cs.GenerateCertKey(ca, caPri, "client.key", "client.pem", []string{"scry"}, []string{"scry"})
	if err != nil {
		return err
	}

	return nil

}
