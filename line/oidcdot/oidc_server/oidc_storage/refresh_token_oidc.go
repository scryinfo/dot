package oidc_storage

import "github.com/zitadel/oidc/v4/pkg/op"

var _ op.TokenRequest = (*RefreshToken)(nil)

// GetAudience implements [op.TokenRequest].
func (m *RefreshToken) GetAudience() []string {
	return m.Audience
}

// GetScopes implements [op.TokenRequest].
func (m *RefreshToken) GetScopes() []string {
	return m.Scopes
}

// GetSubject implements [op.TokenRequest].
func (m *RefreshToken) GetSubject() string {
	return m.UserId
}
