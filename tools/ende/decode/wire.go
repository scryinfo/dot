//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
)

func InitializeService() (*Line, func(), error) {
	wire.Build(
		LineSet,
	)
	return nil, nil, nil
}
