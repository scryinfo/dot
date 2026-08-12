package oidc_server

import (
	"time"

	"github.com/zitadel/oidc/v4/pkg/oidc"
)

type AuthRequest struct {
	ID            string
	CreationDate  time.Time
	ApplicationID string
	CallbackURI   string
	TransferState string
	Prompt        []string
	// UiLocales     []language.Tag
	LoginHint    string
	MaxAuthAge   *time.Duration
	UserID       string
	Scopes       []string
	ResponseType oidc.ResponseType
	ResponseMode oidc.ResponseMode
	Nonce        string
	// CodeChallenge *OIDCCodeChallenge

	done     bool
	authTime time.Time
}
