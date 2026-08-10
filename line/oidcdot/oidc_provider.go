package oidcdot

import (
	"context"

	"github.com/scryinfo/dot/dot"
	"github.com/zitadel/oidc/v4/pkg/client/rp"
)

type OidcProviderConfig struct {
	OidcIssuer   string `toml:"oidc_issuer" json:"oidc_issuer" yaml:"oidc_issuer" mapstructure:"oidc_issuer"`         // "https://your-zitadel-instance.zitadel.cloud"
	ClientID     string `toml:"client_id" json:"client_id" yaml:"client_id" mapstructure:"client_id"`                 // "your-client-id"
	ClientSecret string `toml:"client_secret" json:"client_secret" yaml:"client_secret" mapstructure:"client_secret"` // "your-client-secret"
	RedirectURL  string `toml:"redirect_url" json:"redirect_url" yaml:"redirect_url" mapstructure:"redirect_url"`     // "http://localhost:8080/callback"
}

type OidcProvider struct {
	config   *OidcProviderConfig
	logger   *dot.LoggerType
	provider rp.RelyingParty
}

func NewOidcProvider(config *OidcProviderConfig, logger *dot.LoggerType) (*OidcProvider, error) {
	d := &OidcProvider{
		config: config,
		logger: logger,
	}
	err := d.newOIDCRelyingParty()
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (s *OidcProvider) newOIDCRelyingParty() error {
	ctx := context.Background()

	config := s.config
	rpOpts := []rp.Option{}
	provider, err := rp.NewRelyingPartyOIDC(ctx, config.OidcIssuer, config.ClientID, config.ClientSecret, config.RedirectURL, []string{}, rpOpts...)
	if err != nil {
		s.logger.Error().AnErr("init oidc client failed", err).Send()
		return err
	}
	s.provider = provider
	return nil
}
