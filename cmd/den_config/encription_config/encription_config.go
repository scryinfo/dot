package main

import (
	"crypto/ecdh"
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"

	denconfig "github.com/scryinfo/dot/cmd/den_config"
	"github.com/scryinfo/dot/line/sconfig"
	"github.com/scryinfo/scryg/sutils/sfile"
)

func main() {
	config := denconfig.ParseFlags()
	err := config.ConfigFile()
	if err != nil {
		slog.Error("", "", err)
		denconfig.Exit(err)
		return
	}
	err = config.Abs()
	if err != nil {
		slog.Error("", "", err)
		denconfig.Exit(err)
		return
	}
	if !sfile.ExistFile(config.PubFileName) {
		err = denconfig.GenX25519(&config)
		if err != nil {
			slog.Error("", "", err)
			denconfig.Exit(err)
			return
		}
	}
	if !sfile.ExistFile(config.ConfigFileName) {
		err := fmt.Errorf("the config file dont exist: %s", config.ConfigFileName)
		slog.Error("", "", err)
		denconfig.Exit(err)
		return
	}
	plainConfig, err := os.ReadFile(config.ConfigFileName)
	if err != nil {
		slog.Error("", "", err)
		denconfig.Exit(err)
		return
	}
	var pub *ecdh.PublicKey
	{
		pubKey, err := os.ReadFile(config.PubFileName)
		if err != nil {
			slog.Error("", "", err)
			denconfig.Exit(err)
			return
		}
		pub, err = ecdh.X25519().NewPublicKey(pubKey)
		if err != nil {
			slog.Error("", "", err)
			denconfig.Exit(err)
			return
		}
	}
	chiperConfig, err := sconfig.EncriptionFile(plainConfig, pub)
	if err != nil {
		slog.Error("", "", err)
		denconfig.Exit(err)
		return
	}
	err = os.WriteFile(config.EnConfigFileName, []byte(base64.StdEncoding.EncodeToString(chiperConfig)), 0o0644)
	if err != nil {
		slog.Error("", "", err)
		denconfig.Exit(err)
		return
	}

}
