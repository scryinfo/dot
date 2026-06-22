package denconfig

import (
	"crypto/ecdh"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/scryinfo/scryg/sutils/sfile"
)

func GenX25519(config *ConfigFile) error {
	{
		var err error
		config.PrivateFileName, err = filepath.Abs(config.PrivateFileName)
		if err != nil {
			slog.Error("", "", err)
			return err
		}
		config.PubFileName, err = filepath.Abs(config.PubFileName)
		if err != nil {
			slog.Error("", "", err)
			return err
		}
	}
	priv, err := ecdh.X25519().GenerateKey(nil)
	if err != nil {
		slog.Error("generate key : ", "", err)
		return err
	}
	{
		dd := filepath.Dir(config.PrivateFileName)
		if dd != "" {
			if !sfile.ExistDir(dd) {
				err = os.MkdirAll(dd, 0o755)
				if err != nil {
					slog.Error("", "", err)
					return err
				}
			}
		}
	}
	err = os.WriteFile(config.PrivateFileName, priv.Bytes(), 0o600)
	if err != nil {
		slog.Error("", "", err)
		return err
	}
	slog.Info(fmt.Sprintf("private key: %s\n len: %d\n %#v", config.PrivateFileName, len(priv.Bytes()), priv))
	{
		dd := filepath.Dir(config.PubFileName)
		if dd != "" {
			if !sfile.ExistDir(dd) {
				err = os.MkdirAll(dd, 0o755)
				if err != nil {
					slog.Error("", "", err)
					return err
				}
			}
		}
	}
	err = os.WriteFile(config.PubFileName, priv.PublicKey().Bytes(), 0o600)
	if err != nil {
		slog.Error("", "", err)
		return err
	}
	slog.Info(fmt.Sprintf("public key: %s\n len: %d\n%#v", config.PubFileName, len(priv.PublicKey().Bytes()), priv.PublicKey()))
	return nil
}

func Exit(err error) {
	code := 0
	if err != nil {
		code = 1
	}
	os.Exit(code)
}
