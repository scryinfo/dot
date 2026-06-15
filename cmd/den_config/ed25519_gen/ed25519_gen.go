package main

import (
	"log/slog"

	denconfig "github.com/scryinfo/dot/cmd/den_config"
)

func main() {
	config := denconfig.ParseFlags()
	err := config.Abs()
	if err != nil {
		slog.Error("", err)
	} else {
		err = denconfig.GenEd25519(&config)
	}

	denconfig.Exit(err)
}
