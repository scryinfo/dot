package oidcdot

import (
	"context"

	"connectrpc.com/connect"
	"github.com/scryinfo/dot/dot"
	oidcapiv1 "github.com/scryinfo/dot/line/oidcdot/oidc_gen/oidcapi/v1"
	"github.com/scryinfo/dot/line/oidcdot/oidc_gen/oidcapi/v1/oidcapiv1connect"
	"github.com/scryinfo/dot/line/rpcdot"
)

var _ oidcapiv1connect.AuthServiceHandler = (*AuthService)(nil)

func NewAuthService(mux *rpcdot.ConnectHttpServerMux, logger *dot.LoggerType) *AuthService {
	d := &AuthService{logger: logger}

	path, handle := oidcapiv1connect.NewAuthServiceHandler(d)
	mux.Handle(path, handle)
	return d
}

type AuthService struct {
	logger *dot.LoggerType
}

// Callback implements [apiv1connect.AuthServiceHandler].
func (a *AuthService) Callback(context.Context, *connect.Request[oidcapiv1.CallbackRequest]) (*connect.Response[oidcapiv1.CallbackResponse], error) {
	panic("unimplemented")
}

// Check implements [apiv1connect.AuthServiceHandler].
func (a *AuthService) Check(context.Context, *connect.Request[oidcapiv1.CheckRequest]) (*connect.Response[oidcapiv1.CheckResponse], error) {
	panic("unimplemented")
}

// Login implements [apiv1connect.AuthServiceHandler].
func (a *AuthService) Login(context.Context, *connect.Request[oidcapiv1.LoginRequest]) (*connect.Response[oidcapiv1.LoginResponse], error) {
	panic("unimplemented")
}

// Logout implements [apiv1connect.AuthServiceHandler].
func (a *AuthService) Logout(context.Context, *connect.Request[oidcapiv1.LogoutRequest]) (*connect.Response[oidcapiv1.LogoutResponse], error) {
	panic("unimplemented")
}

// Reshresh implements [apiv1connect.AuthServiceHandler].
func (a *AuthService) Reshresh(context.Context, *connect.Request[oidcapiv1.ReshreshRequest]) (*connect.Response[oidcapiv1.ReshreshResponse], error) {
	panic("unimplemented")
}
