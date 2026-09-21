package oidc_storage

import (
	"time"

	oidcapiv1 "github.com/scryinfo/dot/line/oidcdot/oidc_gen/oidcapi/v1"
	"github.com/zitadel/oidc/v4/pkg/oidc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *StoragePebble2) SaveAuthRequest_(authReq *AuthRequest) error {
	return s.authRequestDao.Add(authReq)
}

func (s *StoragePebble2) authRequestToInternal(authReq *oidc.AuthRequest, userID string) *AuthRequest {
	var codeChallenge *oidcapiv1.OIDCCodeChallenge
	if authReq.CodeChallenge != "" {
		codeChallenge = &oidcapiv1.OIDCCodeChallenge{
			Challenge: authReq.CodeChallenge,
			Method:    string(authReq.CodeChallengeMethod),
		}
	}

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
