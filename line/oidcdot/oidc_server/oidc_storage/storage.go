package oidc_storage

import (
	"context"
	"time"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/db/pebble2dot"
	"github.com/zitadel/oidc/v4/pkg/oidc"
	"github.com/zitadel/oidc/v4/pkg/op"
)

var _ op.Storage = (*Storage)(nil)

type StoragerConfig struct {
}

type Storage struct {
	db  *pebble2dot.Pebble2
	log *dot.LoggerType
}

// AuthRequestByCode implements [op.Storage].
func (s *Storage) AuthRequestByCode(context.Context, string) (op.AuthRequest, error) {
	panic("unimplemented")
}

// AuthRequestByID implements [op.Storage].
func (s *Storage) AuthRequestByID(context.Context, string) (op.AuthRequest, error) {
	panic("unimplemented")
}

// AuthorizeClientIDSecret implements [op.Storage].
func (s *Storage) AuthorizeClientIDSecret(ctx context.Context, clientID string, clientSecret string) error {
	panic("unimplemented")
}

// CreateAccessAndRefreshTokens implements [op.Storage].
func (s *Storage) CreateAccessAndRefreshTokens(ctx context.Context, request op.TokenRequest, currentRefreshToken string) (accessTokenID string, newRefreshToken string, expiration time.Time, err error) {
	panic("unimplemented")
}

// CreateAccessToken implements [op.Storage].
func (s *Storage) CreateAccessToken(context.Context, op.TokenRequest) (accessTokenID string, expiration time.Time, err error) {
	panic("unimplemented")
}

// CreateAuthRequest implements [op.Storage].
func (s *Storage) CreateAuthRequest(context.Context, *oidc.AuthRequest, string) (op.AuthRequest, error) {
	panic("unimplemented")
}

// DeleteAuthRequest implements [op.Storage].
func (s *Storage) DeleteAuthRequest(context.Context, string) error {
	panic("unimplemented")
}

// GetClientByClientID implements [op.Storage].
func (s *Storage) GetClientByClientID(ctx context.Context, clientID string) (op.Client, error) {
	panic("unimplemented")
}

// GetKeyByIDAndClientID implements [op.Storage].
func (s *Storage) GetKeyByIDAndClientID(ctx context.Context, keyID string, clientID string) (*jose.JSONWebKey, error) {
	panic("unimplemented")
}

// GetPrivateClaimsFromScopes implements [op.Storage].
func (s *Storage) GetPrivateClaimsFromScopes(ctx context.Context, userID string, clientID string, scopes []string) (map[string]any, error) {
	panic("unimplemented")
}

// GetRefreshTokenInfo implements [op.Storage].
func (s *Storage) GetRefreshTokenInfo(ctx context.Context, clientID string, token string) (userID string, tokenID string, err error) {
	panic("unimplemented")
}

// Health implements [op.Storage].
func (s *Storage) Health(context.Context) error {
	panic("unimplemented")
}

// KeySet implements [op.Storage].
func (s *Storage) KeySet(context.Context) ([]op.Key, error) {
	panic("unimplemented")
}

// RevokeToken implements [op.Storage].
func (s *Storage) RevokeToken(ctx context.Context, tokenOrTokenID string, userID string, clientID string) *oidc.Error {
	panic("unimplemented")
}

// SaveAuthCode implements [op.Storage].
func (s *Storage) SaveAuthCode(context.Context, string, string) error {
	panic("unimplemented")
}

// SetIntrospectionFromToken implements [op.Storage].
func (s *Storage) SetIntrospectionFromToken(ctx context.Context, userinfo *oidc.IntrospectionResponse, tokenID string, subject string, clientID string) error {
	panic("unimplemented")
}

// SetUserinfoFromScopes implements [op.Storage].
func (s *Storage) SetUserinfoFromScopes(ctx context.Context, userinfo *oidc.UserInfo, userID string, clientID string, scopes []string) error {
	panic("unimplemented")
}

// SetUserinfoFromToken implements [op.Storage].
func (s *Storage) SetUserinfoFromToken(ctx context.Context, userinfo *oidc.UserInfo, tokenID string, subject string, origin string) error {
	panic("unimplemented")
}

// SignatureAlgorithms implements [op.Storage].
func (s *Storage) SignatureAlgorithms(context.Context) ([]jose.SignatureAlgorithm, error) {
	panic("unimplemented")
}

// SigningKey implements [op.Storage].
func (s *Storage) SigningKey(context.Context) (op.SigningKey, error) {
	panic("unimplemented")
}

// TerminateSession implements [op.Storage].
func (s *Storage) TerminateSession(ctx context.Context, userID string, clientID string) error {
	panic("unimplemented")
}

// TokenRequestByRefreshToken implements [op.Storage].
func (s *Storage) TokenRequestByRefreshToken(ctx context.Context, refreshToken string) (op.RefreshTokenRequest, error) {
	panic("unimplemented")
}

// ValidateJWTProfileScopes implements [op.Storage].
func (s *Storage) ValidateJWTProfileScopes(ctx context.Context, userID string, scopes []string) ([]string, error) {
	panic("unimplemented")
}

func NewStorager(config *StoragerConfig, db *pebble2dot.Pebble2, logger *dot.LoggerType) (*Storage, error) {
	return &Storage{
		db:  db,
		log: logger,
	}, nil
}
