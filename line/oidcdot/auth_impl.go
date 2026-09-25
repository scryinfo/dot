package oidcdot

import (
	"context"
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/scryinfo/dot/dot"
	oidcapiv1 "github.com/scryinfo/dot/line/oidcdot/oidc_gen/oidcapi/v1"
	"github.com/scryinfo/dot/line/oidcdot/oidc_impl"
	oidcrp "github.com/scryinfo/dot/line/oidcdot/oidc_rp"
	"github.com/scryinfo/dot/line/oidcdot/oidc_server/oidc_storage"
	"github.com/scryinfo/dot/line/rpcdot"
	"github.com/zitadel/oidc/v4/pkg/client/rp"
	httphelper "github.com/zitadel/oidc/v4/pkg/http"
	"github.com/zitadel/oidc/v4/pkg/oidc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	loginPath    = "/auth2/login"
	logoutPath   = "/auth2/logout"
	callbackPath = "/auth2/callback"
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

//go:embed oidc_templates/*
var goTemplateFs embed.FS

const rpLoginFile = "oidc_templates/rp_login.gohtml"

var (
	tmplRpLogin *template.Template
)

func init() {
	tmpl, err := template.ParseFS(goTemplateFs, rpLoginFile)
	if err != nil {
		panic(err)
	}
	tmplRpLogin = tmpl
}

func NewAuthService(config *AuthConfig, mux *rpcdot.ConnectHttpServerMux, provider *OidcProvider, dao *oidc_storage.AppSessionDaoPebble2, logger *dot.LoggerType) (*AuthService, error) {
	d := &AuthService{logger: logger, provider: provider, appSessionHelper: oidcrp.NewAppSessionHelper(dao)}
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
			rp.WithLogger(dot.Slog),
			rp.WithSigningAlgsFromDiscovery(),
		}
	}

	{
		staticFS, err := fs.Sub(webFS, "oidc_web/dist/assets")
		if err != nil {
			dot.Logger.Error().Err(err).Send()
			return nil, err
		}

		mux.Handle(
			"/auth2/assets/",
			http.StripPrefix("/auth2/assets/", http.FileServer(http.FS(staticFS))),
		)
	}

	mux.HandleFunc(config.LoginPath, d.Login)
	mux.HandleFunc(config.LogoutPath, d.Logout)
	mux.HandleFunc(config.CallbackPath, d.OidcCallback())

	// path, handle := oidcapiv1connect.NewAuthServiceHandler(d)
	// mux.Handle(path, handle)
	return d, nil
}

type AuthService struct {
	logger           *dot.LoggerType
	provider         *OidcProvider
	rpOptions        []rp.Option
	appSessionHelper *oidcrp.AppSessionHelper
}

// Callback implements [apiv1connect.AuthServiceHandler].
func (p *AuthService) OidcCallback() http.HandlerFunc {
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

	marshalUserinfo := func(w http.ResponseWriter, req *http.Request, tokens *oidc.Tokens[*oidc.IDTokenClaims], state string, rp rp.RelyingParty, info *oidc.UserInfo) {
		session, err := p.appSessionHelper.GetWebId(w, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if session.RedirectUri != "" {
			http.Redirect(w, req, session.RedirectUri, http.StatusFound)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = tmplRpLogin.Execute(w, map[string]string{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		http.ServeFileFS(w, req, goTemplateFs, rpLoginFile)
	}
	return rp.CodeExchangeHandler(rp.UserinfoCallback(marshalUserinfo), p.provider.provider)
}

// Check implements [apiv1connect.AuthServiceHandler].
func (p *AuthService) Check(context.Context, *connect.Request[oidcapiv1.CheckRequest]) (*connect.Response[oidcapiv1.CheckResponse], error) {
	panic("unimplemented")
}

func (p *AuthService) Login(w http.ResponseWriter, req *http.Request) {
	// res := &oidcapiv1.LoginResponse{
	// 	Resbase: &oidcapiv1.Resbase{},
	// }
	// err := oidc_impl.NewUuidV7Resbase(req.Msg.Reqbase, res.Resbase)
	// if err != nil {
	// 	a.logger.Error().Err(err).Send()
	// 	return nil, err
	// }
	session, err := p.appSessionHelper.GetWebId(w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	redirectUri := p.appSessionHelper.RedirectUri(req)
	if redirectUri != "" {
		session.RedirectUri = redirectUri
	}
	session.CreationDate = timestamppb.Now()
	err = p.appSessionHelper.SaveSession(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rp.AuthURLHandler(func() string {
		return oidc_impl.NewState()
	}, p.provider.provider)(w, req)
}

func (p *AuthService) Logout(w http.ResponseWriter, re *http.Request) {
	panic("unimplemented")
}

func (p *AuthService) Reshresh(context.Context, *connect.Request[oidcapiv1.ReshreshRequest]) (*connect.Response[oidcapiv1.ReshreshResponse], error) {
	panic("unimplemented")
}
