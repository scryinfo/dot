package oidcdot

import (
	"context"

	"connectrpc.com/connect"
	"github.com/scryinfo/dot/dot"
	oidcapiv1 "github.com/scryinfo/dot/line/oidcdot/oidc_gen/oidcapi/v1"
	"github.com/scryinfo/dot/line/oidcdot/oidc_gen/oidcapi/v1/oidcapiv1connect"
	"github.com/scryinfo/dot/line/oidcdot/oidc_impl"
	"github.com/scryinfo/dot/line/rpcdot"
	"github.com/zitadel/oidc/v4/pkg/client/rp"
	"github.com/zitadel/oidc/v4/pkg/oidc"
)

var _ oidcapiv1connect.AuthServiceHandler = (*AuthService)(nil)

func NewAuthService(mux *rpcdot.ConnectHttpServerMux, provider *OidcProvider, logger *dot.LoggerType) *AuthService {
	d := &AuthService{logger: logger, provider: provider}

	path, handle := oidcapiv1connect.NewAuthServiceHandler(d)
	mux.Handle(path, handle)
	return d
}

type AuthService struct {
	logger   *dot.LoggerType
	provider *OidcProvider
}

// Callback implements [apiv1connect.AuthServiceHandler].
func (a *AuthService) OidcCallback(ctx context.Context, req *connect.Request[oidcapiv1.OidcCallbackRequest]) (*connect.Response[oidcapiv1.OidcCallbackResponse], error) {
	res := &oidcapiv1.OidcCallbackResponse{
		Resbase: &oidcapiv1.Resbase{},
	}
	err := oidc_impl.NewUuidV7Resbase(req.Msg.Reqbase, res.Resbase)
	if err != nil {
		a.logger.Error().Err(err).Send()
		return nil, err
	}
	token, err := rp.CodeExchange[*oidc.IDTokenClaims](ctx, req.Msg.Code, a.provider.provider)
	if err != nil {
		a.logger.Error().AnErr("code exchange failed", err).Send()
		return nil, err
	}
	userInfo, err := rp.Userinfo[*oidc.UserInfo](ctx, token.AccessToken, token.TokenType, token.IDTokenClaims.GetSubject(), a.provider.provider)
	if err != nil {
		a.logger.Error().AnErr("userinfo failed", err).Send()
		return nil, err
	}
	res.AccessToken = token.AccessToken
	res.IdToken = token.IDToken
	res.UserId = userInfo.Subject
	res.Email = userInfo.Email
	return connect.NewResponse(res), nil
}

// Check implements [apiv1connect.AuthServiceHandler].
func (a *AuthService) Check(context.Context, *connect.Request[oidcapiv1.CheckRequest]) (*connect.Response[oidcapiv1.CheckResponse], error) {
	panic("unimplemented")
}

// Login implements [apiv1connect.AuthServiceHandler].
func (a *AuthService) Login(ctx context.Context, req *connect.Request[oidcapiv1.LoginRequest]) (*connect.Response[oidcapiv1.LoginResponse], error) {
	res := &oidcapiv1.LoginResponse{
		Resbase: &oidcapiv1.Resbase{},
	}
	err := oidc_impl.NewUuidV7Resbase(req.Msg.Reqbase, res.Resbase)
	if err != nil {
		a.logger.Error().Err(err).Send()
		return nil, err
	}
	res.State = oidc_impl.NewState()
	res.RedirectUrl = rp.AuthURL(res.State, a.provider.provider)
	return connect.NewResponse(res), nil
}

// Logout implements [apiv1connect.AuthServiceHandler].
func (a *AuthService) Logout(context.Context, *connect.Request[oidcapiv1.LogoutRequest]) (*connect.Response[oidcapiv1.LogoutResponse], error) {
	panic("unimplemented")
}

// Reshresh implements [apiv1connect.AuthServiceHandler].
func (a *AuthService) Reshresh(context.Context, *connect.Request[oidcapiv1.ReshreshRequest]) (*connect.Response[oidcapiv1.ReshreshResponse], error) {
	panic("unimplemented")
}
