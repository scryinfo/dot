package oidcdot

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/scryinfo/dot/dot"
	oidcapiv1 "github.com/scryinfo/dot/line/oidcdot/oidc_gen/oidcapi/v1"
	"github.com/scryinfo/dot/line/oidcdot/oidc_impl"
	"github.com/scryinfo/dot/line/rpcdot"
	"github.com/zitadel/oidc/v4/pkg/client/rp"
	httphelper "github.com/zitadel/oidc/v4/pkg/http"
	"github.com/zitadel/oidc/v4/pkg/oidc"
)

const (
	loginPath    = "/auth/login"
	logoutPath   = "/auth/logout"
	callbackPath = "/auth/callback"
)

type AuthConfig struct {
	// Hash Key： 32 or 64， for HMAC
	HashKey string `json:"hash_key" toml:"hash_key" yaml:"hash_key" mapstructure:"hash_key"`
	// 	Block Key： 16（AES-128）、24（AES-192）or 32 AES-256）
	EncryptKey   string `json:"encrypt_key" toml:"encrypt_key" yaml:"encrypt_key" mapstructure:"encrypt_key"`
	LoginPath    string `json:"login_path" toml:"login_path" yaml:"login_path" mapstructure:"login_path"`
	LogoutPath   string `json:"logout_path" toml:"logout_path" yaml:"logout_path" mapstructure:"logout_path"`
	CallbackPath string `json:"callback_path" toml:"callback_path" yaml:"callback_path" mapstructure:"callback_path"`
}

func NewAuthService(config *AuthConfig, mux *rpcdot.ConnectHttpServerMux, provider *OidcProvider, logger *dot.LoggerType) *AuthService {
	d := &AuthService{logger: logger, provider: provider}
	if config.HashKey == "" {
		config.HashKey = "t123456789012345678901234567890t"
	}
	if config.EncryptKey == "" {
		config.EncryptKey = "s123456789012345678901234567890s"
	}
	if config.LoginPath == "" {
		config.LoginPath = loginPath
	}
	if config.LogoutPath == "" {
		config.LogoutPath = logoutPath
	}
	if config.CallbackPath == "" {
		config.CallbackPath = callbackPath
	}
	{
		cookieHandler := httphelper.NewCookieHandler([]byte(config.HashKey), []byte(config.EncryptKey), httphelper.WithUnsecure())
		client := &http.Client{
			Timeout: time.Minute,
		}
		d.rpOptions = []rp.Option{
			rp.WithCookieHandler(cookieHandler),
			rp.WithVerifierOpts(rp.WithIssuedAtOffset(5 * time.Second)),
			rp.WithHTTPClient(client),
			rp.WithLogger(dot.MakeSlog(logger)),
			rp.WithSigningAlgsFromDiscovery(),
		}
	}

	mux.HandleFunc(config.LoginPath, d.Login)
	mux.HandleFunc(config.LogoutPath, d.Logout)
	mux.HandleFunc(config.CallbackPath, d.OidcCallback())

	// path, handle := oidcapiv1connect.NewAuthServiceHandler(d)
	// mux.Handle(path, handle)
	return d
}

type AuthService struct {
	logger    *dot.LoggerType
	provider  *OidcProvider
	rpOptions []rp.Option
}

// Callback implements [apiv1connect.AuthServiceHandler].
func (a *AuthService) OidcCallback() http.HandlerFunc {
	// res := &oidcapiv1.OidcCallbackResponse{
	// 	Resbase: &oidcapiv1.Resbase{},
	// }
	// err := oidc_impl.NewUuidV7Resbase(req.Msg.Reqbase, res.Resbase)
	// if err != nil {
	// 	a.logger.Error().Err(err).Send()
	// 	return nil, err
	// }
	// token, err := rp.CodeExchange[*oidc.IDTokenClaims](ctx, req.Msg.Code, a.provider.provider)
	// if err != nil {
	// 	a.logger.Error().AnErr("code exchange failed", err).Send()
	// 	return nil, err
	// }
	// userInfo, err := rp.Userinfo[*oidc.UserInfo](ctx, token.AccessToken, token.TokenType, token.IDTokenClaims.GetSubject(), a.provider.provider)
	// if err != nil {
	// 	a.logger.Error().AnErr("userinfo failed", err).Send()
	// 	return nil, err
	// }
	// res.AccessToken = token.AccessToken
	// res.IdToken = token.IDToken
	// res.UserId = userInfo.Subject
	// res.Email = userInfo.Email

	marshalUserinfo := func(w http.ResponseWriter, r *http.Request, tokens *oidc.Tokens[*oidc.IDTokenClaims], state string, rp rp.RelyingParty, info *oidc.UserInfo) {
		// fmt.Println("access token", tokens.AccessToken)
		// fmt.Println("refresh token", tokens.RefreshToken)
		// fmt.Println("id token", tokens.IDToken)
		data, err := json.Marshal(info)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(data)
	}
	return rp.CodeExchangeHandler(rp.UserinfoCallback(marshalUserinfo), a.provider.provider)
}

// Check implements [apiv1connect.AuthServiceHandler].
func (a *AuthService) Check(context.Context, *connect.Request[oidcapiv1.CheckRequest]) (*connect.Response[oidcapiv1.CheckResponse], error) {
	panic("unimplemented")
}

// Login implements [apiv1connect.AuthServiceHandler].
func (a *AuthService) Login(w http.ResponseWriter, req *http.Request) {
	// res := &oidcapiv1.LoginResponse{
	// 	Resbase: &oidcapiv1.Resbase{},
	// }
	// err := oidc_impl.NewUuidV7Resbase(req.Msg.Reqbase, res.Resbase)
	// if err != nil {
	// 	a.logger.Error().Err(err).Send()
	// 	return nil, err
	// }
	rp.AuthURLHandler(func() string {
		return oidc_impl.NewState()
	}, a.provider.provider)(w, req)
}

// Logout implements [apiv1connect.AuthServiceHandler].
func (a *AuthService) Logout(w http.ResponseWriter, re *http.Request) {
	panic("unimplemented")
}

// Reshresh implements [apiv1connect.AuthServiceHandler].
func (a *AuthService) Reshresh(context.Context, *connect.Request[oidcapiv1.ReshreshRequest]) (*connect.Response[oidcapiv1.ReshreshResponse], error) {
	panic("unimplemented")
}
