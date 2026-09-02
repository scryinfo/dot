package oidcdot

import (
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/oidcdot/oidc_server/oidc_storage"
	"github.com/scryinfo/dot/line/rpcdot"
)

type OidcServiceHttp struct {
	logger               *dot.LoggerType
	connectHttpServerMux *rpcdot.ConnectHttpServerMux
	store                *oidc_storage.StoragePebble2
}

type OidcServiceConfig struct {
	OidcIssuer string `toml:"oidc_issuer" json:"oidc_issuer" yaml:"oidc_issuer" mapstructure:"oidc_issuer"` // "https://your-zitadel-instance.zitadel.cloud"

}

func NewOidcServiceHttp(config *OidcServiceConfig, mux *rpcdot.ConnectHttpServerMux, store *oidc_storage.StoragePebble2, logger *dot.LoggerType) (*OidcServiceHttp, error) {
	d := &OidcServiceHttp{
		logger:               logger,
		connectHttpServerMux: mux,
		store:                store,
	}

	return d, nil
}
