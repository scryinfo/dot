package oidcdot

import (
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/rpcdot"
)

type OidcServiceConfig struct {
}

type OidcServiceGrpc struct {
	logger     *dot.LoggerType
	grpcServer *rpcdot.GrpcServer
}

func NewOidcServiceGrpc(config *OidcServiceConfig, grpcServer *rpcdot.GrpcServer, logger *dot.LoggerType) (*OidcServiceGrpc, error) {
	return &OidcServiceGrpc{
		logger:     logger,
		grpcServer: grpcServer,
	}, nil
}
