package oidcdot

import (
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"time"

	"github.com/scryinfo/dot/dot"
	oidcconsts "github.com/scryinfo/dot/line/oidcdot/oidc_server/oidc_consts"
	"github.com/scryinfo/dot/line/oidcdot/oidc_server/oidc_storage"
	"github.com/scryinfo/dot/line/rpcdot"
	"github.com/zitadel/oidc/v4/pkg/op"
	"golang.org/x/text/language"
)

//go:embed oidc_web/dist/*
var webFS embed.FS

const loginFile = "oidc_web/dist/login.html"

type OidcServiceHttp struct {
	config               *OidcServiceConfig
	connectHttpServerMux *rpcdot.ConnectHttpServerMux
	store                *oidc_storage.StoragePebble2
	oidcProvider         op.OpenIDProvider
}

type OidcServiceConfig struct {
	OidcIssuer string `toml:"oidc_issuer" json:"oidc_issuer" yaml:"oidc_issuer" mapstructure:"oidc_issuer"` // "https://your-zitadel-instance.zitadel.cloud"
	Key        string `toml:"key" json:"key" yaml:"key" mapstructure:"key"`
	KeyId      string `toml:"key_id" json:"key_id" yaml:"key_id" mapstructure:"key_id"`
}

func init() {
	// if file, err := webFS.Open(loginFile); err == nil {
	// 	bs, err := io.ReadAll(file)
	// 	if err != nil {
	// 		dot.Logger.Error().Err(err).Send()
	// 	}
	// 	loginReader = bytes.NewReader(bs)
	// 	loginModTime = time.Now()
	// 	_ = file.Close()
	// }else {
	// 	dot.Logger.Error().Err(err).Send()
	// }
}

func NewOidcServiceHttp(config *OidcServiceConfig, mux *rpcdot.ConnectHttpServerMux, store *oidc_storage.StoragePebble2) (*OidcServiceHttp, error) {
	if len(config.Key) != 32 {
		err := fmt.Errorf("key must be 32 characters")
		dot.Logger.Error().Err(err).Send()
		return nil, err
	}

	d := &OidcServiceHttp{
		config:               config,
		connectHttpServerMux: mux,
		store:                store,
	}
	err := d.initOp()
	if err != nil {
		dot.Logger.Error().Err(err).Send()
		return nil, err
	}
	{
		staticFS, err := fs.Sub(webFS, "oidc_web/dist/assets")
		if err != nil {
			dot.Logger.Error().Err(err).Send()
			return nil, err
		}

		mux.Handle(
			"/assets/",
			http.StripPrefix("/assets/", http.FileServer(http.FS(staticFS))),
		)
		mux.HandleFunc("GET /login", d.LoginGet)
		mux.HandleFunc("POST /login", d.LoginPost)
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
	op, err := newOP(p.store, p.config.OidcIssuer, key, p.config.KeyId, dot.Slog)
	if err != nil {
		dot.Logger.Error().Err(err).Send()
		return err
	}
	p.oidcProvider = op
	return nil
}

func (p *OidcServiceHttp) LoginGet(w http.ResponseWriter, req *http.Request) {
	authRequestId := req.URL.Query().Get(oidcconsts.QueryAuthRequestID)
	if authRequestId == "" {
		http.Error(w, "authRequestId is required", http.StatusBadRequest)
		return
	}
	// set Content-Type header to text/html, otherwise the method http.ServeContent will check and set it, it take more time
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFileFS(w, req, webFS, loginFile)
	// http.ServeContent(w, req, "login.html", loginModTime, loginReader)
}
func (p *OidcServiceHttp) LoginPost(w http.ResponseWriter, req *http.Request) {
	authRequestId := req.URL.Query().Get(oidcconsts.QueryAuthRequestID)
	if authRequestId == "" {
		http.Error(w, "authRequestId is required", http.StatusBadRequest)
		return
	}
	username := req.FormValue(oidcconsts.UsernameParam)
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}
	password := req.FormValue(oidcconsts.PasswordParam)
	if password == "" {
		http.Error(w, "password is required", http.StatusBadRequest)
		return
	}

	authReq, err := p.store.AuthRequestById_(authRequestId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := p.store.UserByKeyname(req.Context(), username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ok, err := oidc_storage.PasswordHash.ComparePasswordAndHash(password, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "invalid password", http.StatusUnauthorized)
		return
	}
	authReq.DoneF = true
	authReq.UserId = user.Id
	err = p.store.SaveAuthRequest_(authReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = op.CreateAuthRequestCode(req.Context(), authReq, p.store, p.oidcProvider.Crypto())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	op.AuthResponse(authReq, p.oidcProvider, w, req)
}

func (p *OidcServiceHttp) Logout(w http.ResponseWriter, req *http.Request) {
	//todo
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
			op.WithCustomAuthEndpoint(op.NewEndpoint("auth")),
			// Pass our logger to the OP
			op.WithLogger(logger.WithGroup("op")),
		}, extraOptions...)...,
	)
	if err != nil {
		return nil, err
	}
	return handler, nil
}
