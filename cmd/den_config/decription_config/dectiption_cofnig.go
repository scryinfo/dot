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
	if !sfile.ExistFile(config.PrivateFileName) || !sfile.ExistFile(config.PubFileName) {
		err := fmt.Errorf("the private or public file dont exist, private: %s, publick: %s", config.PrivateFileName, config.PubFileName)
		slog.Error("", err)
		denconfig.Exit(err)
		return
	}
	if !sfile.ExistFile(config.ConfigFileName) {
		err := fmt.Errorf("the config file dont exist: %s", config.ConfigFileName)
		slog.Error("", err)
		denconfig.Exit(err)
		return
	}

}
