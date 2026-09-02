package oidcdot

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/oidcdot/oidc_server/oidc_storage"
	"github.com/scryinfo/dot/line/rpcdot"
	"github.com/zitadel/oidc/v4/pkg/op"
	"golang.org/x/text/language"
)

type OidcServiceHttp struct {
	config               *OidcServiceConfig
	logger               *dot.LoggerType
	connectHttpServerMux *rpcdot.ConnectHttpServerMux
	store                *oidc_storage.StoragePebble2
	oidcProvider         op.OpenIDProvider
}

type OidcServiceConfig struct {
	OidcIssuer string `toml:"oidc_issuer" json:"oidc_issuer" yaml:"oidc_issuer" mapstructure:"oidc_issuer"` // "https://your-zitadel-instance.zitadel.cloud"
	Key        string `toml:"key" json:"key" yaml:"key" mapstructure:"key"`
	KeyId      string `toml:"key_id" json:"key_id" yaml:"key_id" mapstructure:"key_id"`
}

func NewOidcServiceHttp(config *OidcServiceConfig, mux *rpcdot.ConnectHttpServerMux, store *oidc_storage.StoragePebble2, logger *dot.LoggerType) (*OidcServiceHttp, error) {
	if len(config.Key) != 32 {
		err := fmt.Errorf("key must be 32 characters")
		logger.Error().Err(err).Send()
		return nil, err
	}

	d := &OidcServiceHttp{
		config:               config,
		logger:               logger,
		connectHttpServerMux: mux,
		store:                store,
	}
	err := d.initOp()
	if err != nil {
		return nil, err
	}
	mux.Handle("/", d.oidcProvider)

	return d, nil
}

func (p *OidcServiceHttp) initOp() error {
	if p.oidcProvider != nil {
		return nil
	}
	var key [32]byte
	copy(key[:], p.config.Key)
	logger := dot.MakeSlog(p.logger)
	op, err := newOP(p.store, p.config.OidcIssuer, key, p.config.KeyId, logger)
	if err != nil {
		p.logger.Error().Err(err).Send()
		return err
	}
	p.oidcProvider = op
	return nil
}

// newOP will create an OpenID Provider for localhost on a specified port
// and a predefined default logout uri
// it will enable all options (see descriptions)
func newOP(
	storage op.Storage,
	issuer string,
	key [32]byte, // encryption key
	keyId string,
	logger *slog.Logger,
	extraOptions ...op.Option,
) (op.OpenIDProvider, error) {
	config := &op.Config{
		CryptoKey:   key,
		CryptoKeyId: keyId,

		// will be used if the end_session endpoint is called without a post_logout_redirect_uri
		// DefaultLogoutRedirectURI: pathLoggedOut,

		// enables code_challenge_method S256 for PKCE (and therefore PKCE in general)
		CodeMethodS256: true,

		// enables additional client_id/client_secret authentication by form post (not only HTTP Basic Auth)
		AuthMethodPost: true,

		// enables additional authentication by using private_key_jwt
		AuthMethodPrivateKeyJWT: true,

		// enables refresh_token grant use
		GrantTypeRefreshToken: true,

		// enables use of the `request` Object parameter
		RequestObjectSupported: true,

		// this example has only static texts (in English), so we'll set the here accordingly
		SupportedUILocales: []language.Tag{language.English},

		DeviceAuthorization: op.DeviceAuthorizationConfig{
			Lifetime:     5 * time.Minute,
			PollInterval: 5 * time.Second,
			UserFormPath: "/device",
			UserCode:     op.UserCodeBase20,
		},
	}
	handler, err := op.NewProvider(config, storage, op.StaticIssuer(issuer),
		append([]op.Option{
			//we must explicitly allow the use of the http issuer
			op.WithAllowInsecure(),
			// as an example on how to customize an endpoint this will change the authorization_endpoint from /authorize to /auth
			// op.WithCustomAuthEndpoint(op.NewEndpoint("auth")),
			// Pass our logger to the OP
			op.WithLogger(logger.WithGroup("op")),
		}, extraOptions...)...,
	)
	if err != nil {
		return nil, err
	}
	return handler, nil
}
