//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
)

func InitializeService() *App {
	wire.Build(
		AppSet,
	)
	return nil
}
