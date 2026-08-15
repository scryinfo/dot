package oidc_storage

import (
	"time"

	oidcapiv1 "github.com/scryinfo/dot/line/oidcdot/oidc_gen/oidcapi/v1"
	"github.com/zitadel/oidc/v4/pkg/oidc"
	"github.com/zitadel/oidc/v4/pkg/op"
)

type AuthRequest struct {
	*oidcapiv1.AuthRequest
}

var _ op.AuthRequest = (*AuthRequest)(nil)

func (p *AuthRequest) GetID() string {
	return p.AuthRequest.Id
}

func (p *AuthRequest) GetACR() string {
	return "" // we won't handle acr in this example
}

func (p *AuthRequest) GetAMR() []string {
	// this example only uses password for authentication
	if p.AuthRequest.Done {
		return []string{"pwd"}
	}
	return nil
}

func (p *AuthRequest) GetAudience() []string {
	return []string{p.AuthRequest.ApplicationId} // this example will always just use the client_id as audience
}

func (p *AuthRequest) GetAuthTime() time.Time {
	return p.AuthRequest.AuthTime.AsTime()
}

func (p *AuthRequest) GetClientID() string {
	return p.AuthRequest.ApplicationId
}

func (p *AuthRequest) GetCodeChallenge() *oidc.CodeChallenge {
	return CodeChallengeToOIDC(p.AuthRequest.CodeChallenge)
}

func (p *AuthRequest) GetNonce() string {
	return p.AuthRequest.Nonce
}

func (p *AuthRequest) GetRedirectURI() string {
	return p.AuthRequest.CallbackUri
}

func (p *AuthRequest) GetResponseType() oidc.ResponseType {
	return oidc.ResponseType(p.AuthRequest.ResponseType)
}

func (p *AuthRequest) GetResponseMode() oidc.ResponseMode {
	return oidc.ResponseMode(p.AuthRequest.ResponseMode)
}

func (p *AuthRequest) GetScopes() []string {
	return p.AuthRequest.Scopes
}

func (p *AuthRequest) GetState() string {
	return p.AuthRequest.TransferState
}

func (p *AuthRequest) GetSubject() string {
	return p.AuthRequest.UserId
}

func (p *AuthRequest) Done() bool {
	return p.AuthRequest.Done
}

func CodeChallengeToOIDC(challenge *oidcapiv1.OIDCCodeChallenge) *oidc.CodeChallenge {
	if challenge == nil {
		return nil
	}
	challengeMethod := oidc.CodeChallengeMethodPlain
	if challenge.Method == "S256" {
		challengeMethod = oidc.CodeChallengeMethodS256
	}
	return &oidc.CodeChallenge{
		Challenge: challenge.Challenge,
		Method:    challengeMethod,
	}
}
