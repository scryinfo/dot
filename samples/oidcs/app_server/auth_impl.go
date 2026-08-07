package main

import (
	"context"

	"connectrpc.com/connect"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/rpcdot"
	apiv1 "github.com/scryinfo/dot/samples/oidcs/go_out/connect/api/v1"
	"github.com/scryinfo/dot/samples/oidcs/go_out/connect/api/v1/apiv1connect"
)

var _ apiv1connect.AuthServiceHandler = (*AuthService)(nil)

func NewAuthService(mux *rpcdot.ConnectHttpServerMux, logger *dot.LoggerType) *AuthService {
	d := &AuthService{logger: logger}

	path, handle := apiv1connect.NewAuthServiceHandler(d)
	mux.Handle(path, handle)
	return d
}

type AuthService struct {
	logger *dot.LoggerType
}

// Callback implements [apiv1connect.AuthServiceHandler].
func (a *AuthService) Callback(context.Context, *connect.Request[apiv1.CallbackRequest]) (*connect.Response[apiv1.CallbackResponse], error) {
	panic("unimplemented")
}

// Check implements [apiv1connect.AuthServiceHandler].
func (a *AuthService) Check(context.Context, *connect.Request[apiv1.CheckRequest]) (*connect.Response[apiv1.CheckResponse], error) {
	panic("unimplemented")
}

// Login implements [apiv1connect.AuthServiceHandler].
func (a *AuthService) Login(context.Context, *connect.Request[apiv1.LoginRequest]) (*connect.Response[apiv1.LoginResponse], error) {
	panic("unimplemented")
}

// Logout implements [apiv1connect.AuthServiceHandler].
func (a *AuthService) Logout(context.Context, *connect.Request[apiv1.LogoutRequest]) (*connect.Response[apiv1.LogoutResponse], error) {
	panic("unimplemented")
}

// Reshresh implements [apiv1connect.AuthServiceHandler].
func (a *AuthService) Reshresh(context.Context, *connect.Request[apiv1.ReshreshRequest]) (*connect.Response[apiv1.ReshreshResponse], error) {
	panic("unimplemented")
}
