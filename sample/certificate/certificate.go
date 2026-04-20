// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"os"

	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/certificate"
	"github.com/scryinfo/scryg/sutils/ssignal"
)

type App struct {
	Cert *certificate.Ecdsa
}

var AppSet = wire.NewSet(
	wire.Struct(new(App), "*"),
	certificate.NewEcdsa,
)

func main() {
	dot.InitLogger(new(dot.TestLogConfig()))
	app := InitializeService()
	//second step ....

	err := makeSample(app.Cert)
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
