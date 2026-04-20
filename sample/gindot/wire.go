//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
)

func InitializeService() (*App, func(), error) {
	wire.Build(
		AppSet,
	)
	return nil, nil, nil
}
