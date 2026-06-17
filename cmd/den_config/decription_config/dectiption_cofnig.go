package main

import (
	"crypto/ecdh"
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"

	denconfig "github.com/scryinfo/dot/cmd/den_config"
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
	if !sfile.ExistFile(config.PrivateFileName) {
		err := fmt.Errorf("the private file dont exist, private: %s", config.PrivateFileName)
		slog.Error("", "", err)
		denconfig.Exit(err)
		return
	}
	if !sfile.ExistFile(config.ConfigFileName) {
		err := fmt.Errorf("the config file dont exist: %s", config.ConfigFileName)
		slog.Error("", "", err)
		denconfig.Exit(err)
		return
	}

	var priv *ecdh.PrivateKey
	{
		pubKey, err := os.ReadFile(config.PrivateFileName)
		if err != nil {
			slog.Error("", "", err)
			denconfig.Exit(err)
			return
		}
		priv, err = ecdh.X25519().NewPrivateKey(pubKey)
		if err != nil {
			slog.Error("", "", err)
			denconfig.Exit(err)
			return
		}
	}
	var chiperConfig []byte
	{
		chiperConfig, err = os.ReadFile(config.ConfigFileName)
		if err != nil {
			slog.Error("", "", err)
			denconfig.Exit(err)
			return
		}
		chiperConfig, err = base64.StdEncoding.DecodeString(string(chiperConfig))
		if err != nil {
			slog.Error("", "", err)
			denconfig.Exit(err)
			return
		}
	}
	plainConfig, err := denconfig.DecriptionFile(chiperConfig, priv)
	if err != nil {
		slog.Error("", "", err)
		denconfig.Exit(err)
		return
	}

	err = os.WriteFile(config.DeConfigFileName, plainConfig, 0o0644)
	if err != nil {
		slog.Error("", "", err)
		denconfig.Exit(err)
		return
	}

}
