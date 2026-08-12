package oidc_server

import (
	"context"
	"time"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/db/pebble2dot"
	"github.com/zitadel/oidc/v4/pkg/oidc"
	"github.com/zitadel/oidc/v4/pkg/op"
)

var _ op.Storage = (*Storager)(nil)

type StoragerConfig struct {
}

type Storager struct {
	db  *pebble2dot.Pebble2
	log *dot.LoggerType
}

// AuthRequestByCode implements [op.Storage].
func (s *Storager) AuthRequestByCode(context.Context, string) (op.AuthRequest, error) {
	panic("unimplemented")
}

// AuthRequestByID implements [op.Storage].
func (s *Storager) AuthRequestByID(context.Context, string) (op.AuthRequest, error) {
	panic("unimplemented")
}

// AuthorizeClientIDSecret implements [op.Storage].
func (s *Storager) AuthorizeClientIDSecret(ctx context.Context, clientID string, clientSecret string) error {
	panic("unimplemented")
}

// CreateAccessAndRefreshTokens implements [op.Storage].
func (s *Storager) CreateAccessAndRefreshTokens(ctx context.Context, request op.TokenRequest, currentRefreshToken string) (accessTokenID string, newRefreshToken string, expiration time.Time, err error) {
	panic("unimplemented")
}

// CreateAccessToken implements [op.Storage].
func (s *Storager) CreateAccessToken(context.Context, op.TokenRequest) (accessTokenID string, expiration time.Time, err error) {
	panic("unimplemented")
}

// CreateAuthRequest implements [op.Storage].
func (s *Storager) CreateAuthRequest(context.Context, *oidc.AuthRequest, string) (op.AuthRequest, error) {
	panic("unimplemented")
}

// DeleteAuthRequest implements [op.Storage].
func (s *Storager) DeleteAuthRequest(context.Context, string) error {
	panic("unimplemented")
}

// GetClientByClientID implements [op.Storage].
func (s *Storager) GetClientByClientID(ctx context.Context, clientID string) (op.Client, error) {
	panic("unimplemented")
}

// GetKeyByIDAndClientID implements [op.Storage].
func (s *Storager) GetKeyByIDAndClientID(ctx context.Context, keyID string, clientID string) (*jose.JSONWebKey, error) {
	panic("unimplemented")
}

// GetPrivateClaimsFromScopes implements [op.Storage].
func (s *Storager) GetPrivateClaimsFromScopes(ctx context.Context, userID string, clientID string, scopes []string) (map[string]any, error) {
	panic("unimplemented")
}

// GetRefreshTokenInfo implements [op.Storage].
func (s *Storager) GetRefreshTokenInfo(ctx context.Context, clientID string, token string) (userID string, tokenID string, err error) {
	panic("unimplemented")
}

// Health implements [op.Storage].
func (s *Storager) Health(context.Context) error {
	panic("unimplemented")
}

// KeySet implements [op.Storage].
func (s *Storager) KeySet(context.Context) ([]op.Key, error) {
	panic("unimplemented")
}

// RevokeToken implements [op.Storage].
func (s *Storager) RevokeToken(ctx context.Context, tokenOrTokenID string, userID string, clientID string) *oidc.Error {
	panic("unimplemented")
}

// SaveAuthCode implements [op.Storage].
func (s *Storager) SaveAuthCode(context.Context, string, string) error {
	panic("unimplemented")
}

// SetIntrospectionFromToken implements [op.Storage].
func (s *Storager) SetIntrospectionFromToken(ctx context.Context, userinfo *oidc.IntrospectionResponse, tokenID string, subject string, clientID string) error {
	panic("unimplemented")
}

// SetUserinfoFromScopes implements [op.Storage].
func (s *Storager) SetUserinfoFromScopes(ctx context.Context, userinfo *oidc.UserInfo, userID string, clientID string, scopes []string) error {
	panic("unimplemented")
}

// SetUserinfoFromToken implements [op.Storage].
func (s *Storager) SetUserinfoFromToken(ctx context.Context, userinfo *oidc.UserInfo, tokenID string, subject string, origin string) error {
	panic("unimplemented")
}

// SignatureAlgorithms implements [op.Storage].
func (s *Storager) SignatureAlgorithms(context.Context) ([]jose.SignatureAlgorithm, error) {
	panic("unimplemented")
}

// SigningKey implements [op.Storage].
func (s *Storager) SigningKey(context.Context) (op.SigningKey, error) {
	panic("unimplemented")
}

// TerminateSession implements [op.Storage].
func (s *Storager) TerminateSession(ctx context.Context, userID string, clientID string) error {
	panic("unimplemented")
}

// TokenRequestByRefreshToken implements [op.Storage].
func (s *Storager) TokenRequestByRefreshToken(ctx context.Context, refreshToken string) (op.RefreshTokenRequest, error) {
	panic("unimplemented")
}

// ValidateJWTProfileScopes implements [op.Storage].
func (s *Storager) ValidateJWTProfileScopes(ctx context.Context, userID string, scopes []string) ([]string, error) {
	panic("unimplemented")
}

func NewStorager(config *StoragerConfig, db *pebble2dot.Pebble2, logger *dot.LoggerType) (*Storager, error) {
	return &Storager{
		db:  db,
		log: logger,
	}, nil
}
