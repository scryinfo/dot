package oidc_storage

import (
	"context"
	"fmt"
	"time"

	"github.com/scryinfo/dot/lib/kits"
	daobase "github.com/scryinfo/dot/line/db/dao/dao_base"
	"github.com/zitadel/oidc/v4/pkg/oidc"
	"github.com/zitadel/oidc/v4/pkg/op"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *StoragePebble2) SaveAuthRequest_(authReq *AuthRequest) error {
	return s.authRequestDao.Add(authReq)
}

func (s *StoragePebble2) authRequestToInternal(authReq *oidc.AuthRequest, userID string) *AuthRequest {
	codeChallenge := CodeChallengeToProto(&oidc.CodeChallenge{
		Challenge: authReq.CodeChallenge,
		Method:    authReq.CodeChallengeMethod,
	})

	return &AuthRequest{
		Id:            NewAuthRequestId(),
		CreationDate:  timestamppb.New(time.Now()),
		ApplicationId: authReq.ClientID,
		CallbackUri:   authReq.RedirectURI,
		TransferState: authReq.State,
		Prompt:        PromptToInternal(authReq.Prompt),
		// UiLocales:     authReq.UILocales,
		LoginHint:     authReq.LoginHint,
		MaxAuthAge:    MaxAgeToInternal(authReq.MaxAge),
		UserId:        userID,
		Scopes:        authReq.Scopes,
		ResponseType:  ResponseTypeEx.ToProto(authReq.ResponseType),
		ResponseMode:  ResponseModeEx.ToProto(authReq.ResponseMode),
		Nonce:         authReq.Nonce,
		CodeChallenge: codeChallenge,
	}
}

func (s *StoragePebble2) accessToken(applicationID, refreshTokenID, subject string, audience, scopes []string) (*Token, error) {
	token := NewToken()
	token.ApplicationId = applicationID
	token.RefreshTokenId = refreshTokenID
	token.Subject = subject
	token.Audience = audience
	token.Expiration = timestamppb.New(time.Now().Add(5 * time.Minute))
	token.Scopes = scopes
	s.tokenDao.Add(&token)
	return &token, nil
}

// createRefreshToken will store a refresh_token in-memory based on the provided information
func (s *StoragePebble2) createRefreshToken(accessToken *Token, amr []string, authTime time.Time) (string, error) {
	token := NewRefreshToken()
	token.Token = accessToken.RefreshTokenId
	token.AuthTime = timestamppb.New(authTime)
	token.Amr = amr
	token.ApplicationId = accessToken.ApplicationId
	token.UserId = accessToken.Subject
	token.Audience = accessToken.Audience
	token.Expiration = timestamppb.New(time.Now().Add(5 * time.Hour))
	token.Scopes = accessToken.Scopes
	token.AccessToken = accessToken.Id
	s.refreshTokenDao.Add(&token)
	return token.Token, nil
}

// renewRefreshToken checks the provided refresh_token and creates a new one based on the current
//
// [Refresh Token Rotation] is implemented.
//
// [Refresh Token Rotation]: https://www.rfc-editor.org/rfc/rfc6819#section-5.2.2.3
func (s *StoragePebble2) renewRefreshToken(currentRefreshToken, newRefreshToken, newAccessToken string) error {

	refreshToken, err := s.refreshTokenDao.Find(daobase.IdType(currentRefreshToken))
	if err != nil {
		return err
	}
	err = s.refreshTokenDao.Remove(&refreshToken)
	if err != nil {
		return err
	}

	if refreshToken.Expiration.AsTime().Before(time.Now()) {
		return fmt.Errorf("expired refresh token")
	}

	// creates a new refresh token based on the current one
	refreshToken.Token = newRefreshToken
	refreshToken.Id = ""
	refreshToken.Expiration = timestamppb.New(time.Now().Add(5 * time.Hour))
	refreshToken.AccessToken = newAccessToken
	return s.refreshTokenDao.Add(&refreshToken)
}
func (s *StoragePebble2) exchangeRefreshToken(ctx context.Context, request op.TokenExchangeRequest) (accessTokenID string, newRefreshToken string, expiration time.Time, err error) {
	applicationID := request.GetClientID()
	authTime := request.GetAuthTime()

	refreshTokenID := kits.Ids.Uuid()
	accessToken, err := s.accessToken(applicationID, refreshTokenID, request.GetSubject(), request.GetAudience(), request.GetScopes())
	if err != nil {
		return "", "", time.Time{}, err
	}

	refreshToken, err := s.createRefreshToken(accessToken, nil, authTime)
	if err != nil {
		return "", "", time.Time{}, err
	}

	return accessToken.Id, refreshToken, accessToken.Expiration.AsTime(), nil
}
func getInfoFromRequest(req op.TokenRequest) (clientID string, authTime time.Time, amr []string) {
	authReq, ok := req.(*AuthRequest) // Code Flow (with scope offline_access)
	if ok {
		return authReq.ApplicationId, authReq.GetAuthTime(), authReq.GetAMR()
	}
	refreshReq, ok := req.(*RefreshToken) // Refresh Token Request
	if ok {
		return refreshReq.ApplicationId, refreshReq.AuthTime.AsTime(), refreshReq.Amr
	}
	return "", time.Time{}, nil
}
