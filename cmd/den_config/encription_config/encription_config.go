package main

import (
	"fmt"
	"log/slog"

	denconfig "github.com/scryinfo/dot/cmd/den_config"
	"github.com/scryinfo/scryg/sutils/sfile"
)

func main() {
	config := denconfig.ParseFlags()
	err := config.Abs()
	if err != nil {
		slog.Error("", err)
		denconfig.Exit(err)
		return
	}
	if !sfile.ExistFile(config.PrivateFileName) {
		err = denconfig.GenEd25519(&config)
		if err != nil {
			denconfig.Exit(err)
			return
		}
	}
	if !sfile.ExistFile(config.ConfigFileName) {
		err := fmt.Errorf("the config file dont exist: %s", config.ConfigFileName)
		slog.Error("", err)
		denconfig.Exit(err)
		return
	}

}
