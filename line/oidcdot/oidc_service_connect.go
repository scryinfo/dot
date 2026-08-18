package oidcdot

import (
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/rpcdot"
)

type OidcServiceConnect struct {
	logger        *dot.LoggerType
	connectServer *rpcdot.ConnectServer
}

func NewOidcServiceConnect(config *OidcServiceConfig, connectServer *rpcdot.ConnectServer, logger *dot.LoggerType) (*OidcServiceConnect, error) {
	return &OidcServiceConnect{
		logger:        logger,
		connectServer: connectServer,
	}, nil
}
