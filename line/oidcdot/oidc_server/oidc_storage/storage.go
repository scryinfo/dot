package oidc_storage

import (
	"context"
	"time"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	daobase "github.com/scryinfo/dot/line/db/dao/dao_base"
	"github.com/scryinfo/dot/line/db/pebble2dot"
	"github.com/zitadel/oidc/v4/pkg/oidc"
	"github.com/zitadel/oidc/v4/pkg/op"
)

var _ op.Storage = (*StoragePebble2)(nil)

type StoragerConfig struct {
}

type StoragePebble2 struct {
	db          *pebble2dot.Pebble2
	log         *dot.LoggerType
	signingKeys _SigningKeys

	authRequestDao     *AuthRequestDaoPebble2
	codeAuthRequestDao *CodeAuthRequestDaoPebble2
	oidcClientDao      *OidcClientDaoPebble2
	clientStatusDao    *ClientStatusDaoPebble2
	identityDao        *IdentityDaoPebble2
	refreshTokenDao    *RefreshTokenDaoPebble2
	tokenDao           *TokenDaoPebble2
	userDao            *UserDaoPebble2
	userIdentitiesDao  *UserIdentitiesDaoPebble2
}

// AuthRequestByCode implements [op.Storage].
func (s *StoragePebble2) AuthRequestByCode(context.Context, string) (op.AuthRequest, error) {
	panic("unimplemented")
}

// AuthRequestByID implements [op.Storage].
func (s *StoragePebble2) AuthRequestByID(context.Context, string) (op.AuthRequest, error) {
	panic("unimplemented")
}

// AuthorizeClientIDSecret implements [op.Storage].
func (s *StoragePebble2) AuthorizeClientIDSecret(ctx context.Context, clientID string, clientSecret string) error {
	panic("unimplemented")
}

// CreateAccessAndRefreshTokens implements [op.Storage].
func (s *StoragePebble2) CreateAccessAndRefreshTokens(ctx context.Context, request op.TokenRequest, currentRefreshToken string) (accessTokenID string, newRefreshToken string, expiration time.Time, err error) {
	panic("unimplemented")
}

// CreateAccessToken implements [op.Storage].
func (s *StoragePebble2) CreateAccessToken(context.Context, op.TokenRequest) (accessTokenID string, expiration time.Time, err error) {
	panic("unimplemented")
}

// CreateAuthRequest implements [op.Storage].
func (s *StoragePebble2) CreateAuthRequest(ctx context.Context, authRequest *oidc.AuthRequest, userId string) (op.AuthRequest, error) {
	panic("unimplemented")
}

// DeleteAuthRequest implements [op.Storage].
func (s *StoragePebble2) DeleteAuthRequest(context.Context, string) error {
	panic("unimplemented")
}

// GetClientByClientID implements [op.Storage].
func (s *StoragePebble2) GetClientByClientID(ctx context.Context, clientID string) (op.Client, error) {
	m, err := s.oidcClientDao.Find((daobase.IdType(clientID)))
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// GetKeyByIDAndClientID implements [op.Storage].
func (s *StoragePebble2) GetKeyByIDAndClientID(ctx context.Context, keyID string, clientID string) (*jose.JSONWebKey, error) {
	panic("unimplemented")
}

// GetPrivateClaimsFromScopes implements [op.Storage].
func (s *StoragePebble2) GetPrivateClaimsFromScopes(ctx context.Context, userID string, clientID string, scopes []string) (map[string]any, error) {
	panic("unimplemented")
}

// GetRefreshTokenInfo implements [op.Storage].
func (s *StoragePebble2) GetRefreshTokenInfo(ctx context.Context, clientID string, token string) (userID string, tokenID string, err error) {
	panic("unimplemented")
}

// Health implements [op.Storage].
func (s *StoragePebble2) Health(context.Context) error {
	panic("unimplemented")
}

// KeySet implements [op.Storage].
func (s *StoragePebble2) KeySet(context.Context) ([]op.Key, error) {
	keys := make([]op.Key, 0, len(s.signingKeys.publicKeys))
	for _, key := range s.signingKeys.publicKeys {
		keys = append(keys, key)
	}
	return keys, nil
}

// RevokeToken implements [op.Storage].
func (s *StoragePebble2) RevokeToken(ctx context.Context, tokenOrTokenID string, userID string, clientID string) *oidc.Error {
	panic("unimplemented")
}

// SaveAuthCode implements [op.Storage].
func (s *StoragePebble2) SaveAuthCode(context.Context, string, string) error {
	panic("unimplemented")
}

// SetIntrospectionFromToken implements [op.Storage].
func (s *StoragePebble2) SetIntrospectionFromToken(ctx context.Context, userinfo *oidc.IntrospectionResponse, tokenID string, subject string, clientID string) error {
	panic("unimplemented")
}

// SetUserinfoFromScopes implements [op.Storage].
func (s *StoragePebble2) SetUserinfoFromScopes(ctx context.Context, userinfo *oidc.UserInfo, userID string, clientID string, scopes []string) error {
	panic("unimplemented")
}

// SetUserinfoFromToken implements [op.Storage].
func (s *StoragePebble2) SetUserinfoFromToken(ctx context.Context, userinfo *oidc.UserInfo, tokenID string, subject string, origin string) error {
	panic("unimplemented")
}

// SignatureAlgorithms implements [op.Storage].
func (s *StoragePebble2) SignatureAlgorithms(context.Context) ([]jose.SignatureAlgorithm, error) {
	algs := make([]jose.SignatureAlgorithm, 0, len(s.signingKeys.mapKeys))
	for alg, _ := range s.signingKeys.mapKeys {
		algs = append(algs, alg)
	}
	return algs, nil
}

// SigningKey implements [op.Storage].
func (s *StoragePebble2) SigningKey(ctx context.Context) (op.SigningKey, error) {
	// if authReq, ok := op.ClientFromContext(ctx); ok {
	//     clientID := authReq.GetClientID()

	//     // 2. 从你自己的数据库中查出该 clientID 注册的算法 (例如 ES256)
	//     if client, err := s.GetClientByID(ctx, clientID); err == nil {
	//         preferredAlg := client.IDTokenSignedResponseAlg // 拿到了算法！

	//         // 3. 匹配对应的私钥
	//         if key, exists := s.activeKeys[preferredAlg]; exists {
	//             return key, nil
	//         }
	//     }
	// }

	// 4. 如果 ctx 里没有 AuthRequest (比如非 Token 签发场景)，返回默认密钥兜底
	return s.signingKeys.defaultKey, nil
}

// TerminateSession implements [op.Storage].
func (s *StoragePebble2) TerminateSession(ctx context.Context, userID string, clientID string) error {
	panic("unimplemented")
}

// TokenRequestByRefreshToken implements [op.Storage].
func (s *StoragePebble2) TokenRequestByRefreshToken(ctx context.Context, refreshToken string) (op.RefreshTokenRequest, error) {
	panic("unimplemented")
}

// ValidateJWTProfileScopes implements [op.Storage].
func (s *StoragePebble2) ValidateJWTProfileScopes(ctx context.Context, userID string, scopes []string) ([]string, error) {
	panic("unimplemented")
}

func NewStoragePebble2(db *pebble2dot.Pebble2, logger *dot.LoggerType,
	authRequestDao *AuthRequestDaoPebble2, oidcClientDao *OidcClientDaoPebble2, clientStatusDao *ClientStatusDaoPebble2, codeAuthRequestDao *CodeAuthRequestDaoPebble2,
	identityDao *IdentityDaoPebble2, refreshTokenDao *RefreshTokenDaoPebble2,
	tokenDao *TokenDaoPebble2, userDao *UserDaoPebble2, userIdentitiesDao *UserIdentitiesDaoPebble2,
) (*StoragePebble2, error) {
	return &StoragePebble2{
		db:                 db,
		log:                logger,
		signingKeys:        NewSigningKeys(),
		authRequestDao:     authRequestDao,
		oidcClientDao:      oidcClientDao,
		clientStatusDao:    clientStatusDao,
		codeAuthRequestDao: codeAuthRequestDao,
		identityDao:        identityDao,
		refreshTokenDao:    refreshTokenDao,
		tokenDao:           tokenDao,
		userDao:            userDao,
		userIdentitiesDao:  userIdentitiesDao,
	}, nil
}

var Pebble2Set = wire.NewSet(
	NewStoragePebble2,
	NewAuthRequestDaoPebble2,
	NewOidcClientDaoPebble2,
	NewClientStatusDaoPebble2,
	NewCodeAuthRequestDaoPebble2,
	NewIdentityDaoPebble2,
	NewRefreshTokenDaoPebble2,
	NewTokenDaoPebble2,
	NewUserDaoPebble2,
	NewUserIdentitiesDaoPebble2,

	pebble2dot.NewPebble2,
)
