package oidc_storage

import (
	"time"

	oidcapiv1 "github.com/scryinfo/dot/line/oidcdot/oidc_gen/oidcapi/v1"
	"github.com/zitadel/oidc/v4/pkg/oidc"
	"github.com/zitadel/oidc/v4/pkg/op"
)

type AuthRequest oidcapiv1.AuthRequest

var _ op.AuthRequest = (*AuthRequest)(nil)

func (p *AuthRequest) GetID() string {
	return p.Id
}

func (p *AuthRequest) GetACR() string {
	return "" // we won't handle acr in this example
}

func (p *AuthRequest) GetAMR() []string {
	// this example only uses password for authentication
	if p.DoneF {
		return []string{"pwd"}
	}
	return nil
}

func (p *AuthRequest) GetAudience() []string {
	return []string{p.ApplicationId} // this example will always just use the client_id as audience
}

func (p *AuthRequest) GetAuthTime() time.Time {
	return p.AuthTime.AsTime()
}

func (p *AuthRequest) GetClientID() string {
	return p.ApplicationId
}

func (p *AuthRequest) GetCodeChallenge() *oidc.CodeChallenge {
	return CodeChallengeToOIDC(p.CodeChallenge)
}

func (p *AuthRequest) GetNonce() string {
	return p.Nonce
}

func (p *AuthRequest) GetRedirectURI() string {
	return p.CallbackUri
}

func (p *AuthRequest) GetResponseType() oidc.ResponseType {
	return ResponseTypeEx.ToGoType(p.ResponseType)
}

func (p *AuthRequest) GetResponseMode() oidc.ResponseMode {
	return ResponseModeEx.ToGoType(p.ResponseMode)
}

func (p *AuthRequest) GetScopes() []string {
	return p.Scopes
}

func (p *AuthRequest) GetState() string {
	return p.TransferState
}

func (p *AuthRequest) GetSubject() string {
	return p.UserId
}

func (p *AuthRequest) Done() bool {
	return p.DoneF
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
func CodeChallengeToProto(challenge *oidc.CodeChallenge) *oidcapiv1.OIDCCodeChallenge {
	if challenge == nil {
		return nil
	}
	challengeMethod := oidc.CodeChallengeMethodPlain
	if challenge.Method == "S256" {
		challengeMethod = oidc.CodeChallengeMethodS256
	}
	return &oidcapiv1.OIDCCodeChallenge{
		Challenge: challenge.Challenge,
		Method:    string(challengeMethod),
	}
}
