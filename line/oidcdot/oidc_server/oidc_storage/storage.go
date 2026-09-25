package oidc_storage

import (
	"context"
	"errors"
	"time"

	"github.com/cockroachdb/pebble/v2"
	jose "github.com/go-jose/go-jose/v4"
	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/lib/kits"
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
func (s *StoragePebble2) AuthRequestByCode(ctx context.Context, code string) (op.AuthRequest, error) {
	codeAuth, err := s.codeAuthRequestDao.Find(daobase.IdType(code))
	if err != nil {
		dot.Logger.Error().Err(err).Send()
		return nil, err
	}
	auth, err := s.authRequestDao.Find(daobase.IdType(codeAuth.AuthRequestId))
	if err != nil {
		dot.Logger.Error().Err(err).Send()
		return nil, err
	}
	return &auth, nil
}

// AuthRequestByID implements [op.Storage].
func (s *StoragePebble2) AuthRequestByID(ctx context.Context, authRequestId string) (op.AuthRequest, error) {
	return s.AuthRequestById_(authRequestId)
}
func (s *StoragePebble2) AuthRequestById_(authRequestId string) (*AuthRequest, error) {
	auth, err := s.authRequestDao.Find(daobase.IdType(authRequestId))
	if err != nil {
		dot.Logger.Error().Err(err).Send()
	}
	return &auth, err
}

// AuthorizeClientIDSecret implements [op.Storage].
func (s *StoragePebble2) AuthorizeClientIDSecret(ctx context.Context, clientID string, clientSecret string) error {
	panic("unimplemented")
}

// CreateAccessAndRefreshTokens implements [op.Storage].
func (s *StoragePebble2) CreateAccessAndRefreshTokens(ctx context.Context, request op.TokenRequest, currentRefreshToken string) (accessTokenID string, newRefreshToken string, expiration time.Time, err error) {
	// generate tokens via token exchange flow if request is relevant
	if teReq, ok := request.(op.TokenExchangeRequest); ok {
		return s.exchangeRefreshToken(ctx, teReq)
	}

	// get the information depending on the request type / implementation
	applicationID, authTime, amr := getInfoFromRequest(request)

	// if currentRefreshToken is empty (Code Flow) we will have to create a new refresh token
	if currentRefreshToken == "" {
		refreshTokenID := kits.Ids.Uuid()
		accessToken, err := s.accessToken(applicationID, refreshTokenID, request.GetSubject(), request.GetAudience(), request.GetScopes())
		if err != nil {
			dot.Logger.Error().Err(err).Send()
			return "", "", time.Time{}, err
		}
		refreshToken, err := s.createRefreshToken(accessToken, amr, authTime)
		if err != nil {
			dot.Logger.Error().Err(err).Send()
			return "", "", time.Time{}, err
		}
		return accessToken.Id, refreshToken, accessToken.Expiration.AsTime(), nil
	}

	// if we get here, the currentRefreshToken was not empty, so the call is a refresh token request
	// we therefore will have to check the currentRefreshToken and renew the refresh token

	newRefreshToken = kits.Ids.Uuid()

	accessToken, err := s.accessToken(applicationID, newRefreshToken, request.GetSubject(), request.GetAudience(), request.GetScopes())
	if err != nil {
		dot.Logger.Error().Err(err).Send()
		return "", "", time.Time{}, err
	}

	if err := s.renewRefreshToken(currentRefreshToken, newRefreshToken, accessToken.Id); err != nil {
		dot.Logger.Error().Err(err).Send()
		return "", "", time.Time{}, err
	}

	return accessToken.Id, newRefreshToken, accessToken.Expiration.AsTime(), nil
}

// CreateAccessToken implements [op.Storage].
func (s *StoragePebble2) CreateAccessToken(ctx context.Context, request op.TokenRequest) (accessTokenID string, expiration time.Time, err error) {
	var applicationID string
	switch req := request.(type) {
	case *AuthRequest:
		// if authenticated for an app (auth code / implicit flow) we must save the client_id to the token
		applicationID = req.ApplicationId
	case op.TokenExchangeRequest:
		applicationID = req.GetClientID()
	}

	token, err := s.accessToken(applicationID, "", request.GetSubject(), request.GetAudience(), request.GetScopes())
	if err != nil {
		dot.Logger.Error().Err(err).Send()
		return "", time.Time{}, err
	}
	return token.Id, token.Expiration.AsTime(), nil
}

// CreateAuthRequest implements [op.Storage].
func (s *StoragePebble2) CreateAuthRequest(ctx context.Context, authReq *oidc.AuthRequest, userId string) (op.AuthRequest, error) {

	if len(authReq.Prompt) == 1 && authReq.Prompt[0] == "none" {
		dot.Logger.Error().Err(oidc.ErrLoginRequired()).Send()
		return nil, oidc.ErrLoginRequired()
	}

	// typically, you'll fill your storage / storage model with the information of the passed object
	request := s.authRequestToInternal(authReq, userId)
	err := s.authRequestDao.Add(request)
	if err != nil {
		dot.Logger.Error().Err(err).Send()
		return nil, err
	}
	// finally, return the request (which implements the AuthRequest interface of the OP
	return request, nil
}

// DeleteAuthRequest implements [op.Storage].
func (s *StoragePebble2) DeleteAuthRequest(ctx context.Context, id string) error {
	err := s.authRequestDao.RemoveBy(daobase.IdType(id))
	if err != nil {
		dot.Logger.Error().Err(err).Send()
		return err
	}
	m, err := s.codeAuthRequestDao.FindByAuthRequestId(id)
	if err != nil {
		dot.Logger.Error().Err(err).Send()
		return err
	}
	err = s.codeAuthRequestDao.Remove(m)
	if err != nil {
		dot.Logger.Error().Err(err).Send()
		return err
	}
	return nil
}

// GetClientByClientID implements [op.Storage].
func (s *StoragePebble2) GetClientByClientID(ctx context.Context, clientID string) (op.Client, error) {
	m, err := s.oidcClientDao.Find((daobase.IdType(clientID)))
	if err != nil {
		dot.Logger.Error().Err(err).Send()
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
func (s *StoragePebble2) SaveAuthCode(ctx context.Context, id string, code string) error {
	m := CodeAuthRequest{
		Code:          code,
		AuthRequestId: id,
	}
	err := s.codeAuthRequestDao.Add(&m)
	if err != nil {
		dot.Logger.Error().Err(err).Send()
		return err
	}
	return nil
}

// SetIntrospectionFromToken implements [op.Storage].
func (s *StoragePebble2) SetIntrospectionFromToken(ctx context.Context, userinfo *oidc.IntrospectionResponse, tokenID string, subject string, clientID string) error {
	token, err := s.tokenDao.Find(daobase.IdType(tokenID))
	if err != nil {
		if errors.Is(err, pebble.ErrNotFound) {
			userinfo.Active = false
			return nil
		}

		return err
	}
	user, err := s.userDao.Find(daobase.IdType(subject))
	if err != nil {
		return err
	}

	// if token.Revoked {
	// 	introspection.Active = false
	// 	return nil
	// }

	if time.Now().After(token.Expiration.AsTime()) {
		userinfo.Active = false
		return nil
	}

	userinfo.Active = true
	userinfo.Subject = subject
	userinfo.ClientID = clientID
	userinfo.Username = user.Username

	return nil
}

// SetUserinfoFromScopes implements [op.Storage].
func (s *StoragePebble2) SetUserinfoFromScopes(ctx context.Context, userinfo *oidc.UserInfo, userID string, clientID string, scopes []string) error {
	user, err := s.userDao.Find(daobase.IdType(userID))
	if err != nil {
		dot.Logger.Error().Err(err).Send()
		return err
	}

	userinfo.Subject = user.Id
	for _, scope := range scopes {
		switch scope {
		case oidc.ScopeProfile:
			userinfo.Name = user.Username
			userinfo.GivenName = user.FirstName
			userinfo.FamilyName = user.LastName
			// userinfo.Nickname = user.Nickname
			// userinfo.PreferredUsername = user.Username
			// userinfo.Picture = user.AvatarURL
			// userinfo.UpdatedAt = oidc.FromTime(user.UpdatedAt)

		case oidc.ScopeEmail:
			userinfo.Email = user.Email
			userinfo.EmailVerified = oidc.Bool(user.EmailVerified)

		case oidc.ScopePhone:
			userinfo.PhoneNumber = user.Phone
			userinfo.PhoneNumberVerified = oidc.Bool(user.PhoneVerified)

		case oidc.ScopeAddress:
			userinfo.Address = &oidc.UserInfoAddress{
				Formatted:     user.Address.Formatted,
				StreetAddress: user.Address.StreetAddress,
				Locality:      user.Address.Locality,
				Region:        user.Address.Region,
				PostalCode:    user.Address.PostalCode,
				Country:       user.Address.Country,
			}
		}
	}
	return nil
}

// SetUserinfoFromToken implements [op.Storage].
func (s *StoragePebble2) SetUserinfoFromToken(ctx context.Context, userinfo *oidc.UserInfo, tokenID string, subject string, origin string) error {
	user, err := s.userDao.Find(daobase.IdType(subject))
	if err != nil {
		return err
	}

	userinfo.Subject = subject
	userinfo.Name = user.Username
	// userinfo.GivenName = user.GivenName
	// userinfo.FamilyName = user.FamilyName
	// userinfo.NickName = user.NickName
	userinfo.PreferredUsername = user.Username

	userinfo.Email = user.Email
	userinfo.EmailVerified = oidc.Bool(user.EmailVerified)

	userinfo.PhoneNumber = user.Phone
	userinfo.PhoneNumberVerified = oidc.Bool(user.PhoneVerified)

	return nil
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
	NewAppSessionDaoPebble2,

	pebble2dot.NewPebble2,
)
