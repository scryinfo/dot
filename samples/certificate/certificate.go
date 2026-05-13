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
	Cert *certificate.Ecdsa
}

type LineConfig struct {
	Log dot.LogConfig
}

func NewAppConfig(config *sconfig.SConfig) (*LineConfig, error) {
	return sconfig.NewLiceConfig[LineConfig](config)
}

var LineSet = wire.NewSet(
	wire.Struct(new(Line), "*"),
	wire.FieldsOf(new(*LineConfig), "Log"),
	NewAppConfig,
	sconfig.NewConfig,
	dot.NewLogger,
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

	rootKey, err := certificate.MakeECDSAKey()
	if err != nil {
		return err
	}

	ca, err := cs.GenerateRoot(rootKey, "root.key", "root.cert", []string{"scry"}, []string{"scry"})
	if err != nil {
		return err
	}

	_, err = cs.GenerateLeaf(ca, rootKey, "server.key", "server.cert", []string{"scry"}, []string{"scry"})
	if err != nil {
		return err
	}

	_, err = cs.GenerateLeaf(ca, rootKey, "client.key", "client.cert", []string{"scry"}, []string{"scry"})
	if err != nil {
		return err
	}

	return nil

}
