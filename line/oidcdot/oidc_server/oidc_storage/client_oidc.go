package oidc_storage

import (
	"net/url"
	"slices"
	"time"

	"github.com/scryinfo/dot/dot"
	oidcconsts "github.com/scryinfo/dot/line/oidcdot/oidc_server/oidc_consts"
	"github.com/zitadel/oidc/v4/pkg/oidc"
	"github.com/zitadel/oidc/v4/pkg/op"
)

// implements [op.Client]
var _ op.Client = (*OidcClient)(nil)

// AccessTokenType implements [op.Client].
func (m *OidcClient) AccessTokenType() op.AccessTokenType {
	return AccessTokenTypeEx.ToGoType(m.AccessTokenTypeF)
}

// ApplicationType implements [op.Client].
func (m *OidcClient) ApplicationType() op.ApplicationType {
	return ApplicationTypeEx.ToGoType(m.ApplicationTypeF)
}

// AuthMethod implements [op.Client].
func (m *OidcClient) AuthMethod() oidc.AuthMethod {
	return AuthMethodEx.ToGoType(m.AuthMethodF)
}

// ClockSkew implements [op.Client].
func (m *OidcClient) ClockSkew() time.Duration {
	return m.ClockSkewF.AsDuration()
}

// DevMode implements [op.Client].
func (m *OidcClient) DevMode() bool {
	return m.DevModeF
}

// GetID implements [op.Client].
func (m *OidcClient) GetID() string {
	return m.Id
}

// GrantTypes implements [op.Client].
func (m *OidcClient) GrantTypes() []oidc.GrantType {
	re := make([]oidc.GrantType, 0, len(m.GrantTypesF))
	for _, gt := range m.GrantTypesF {
		re = append(re, GrantTypeEx.ToGoType(gt))
	}
	return re
}

// IDTokenLifetime implements [op.Client].
func (m *OidcClient) IDTokenLifetime() time.Duration {
	return m.IdTokenLifetimeF.AsDuration()
}

// IDTokenUserinfoClaimsAssertion implements [op.Client].
func (m *OidcClient) IDTokenUserinfoClaimsAssertion() bool {
	return m.IdTokenUserinfoClaimsAssertionF
}

// IsScopeAllowed implements [op.Client].
func (m *OidcClient) IsScopeAllowed(scope string) bool {
	return slices.Contains(_scopes, scope)
}

// LoginURL implements [op.Client].
func (m *OidcClient) LoginURL(authRequestID string) string {
	baseURL := m.LoginUrlF
	if baseURL == "" {
		baseURL = oidcconsts.LoginEndpoint
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		dot.Logger.Error().AnErr("failed to parse login URL", err).Send()
		return oidcconsts.LoginEndpoint + "?" + oidcconsts.QueryAuthRequestID + "=" + authRequestID
	}
	q := u.Query()
	q.Set(oidcconsts.QueryAuthRequestID, authRequestID)
	u.RawQuery = q.Encode()
	return u.String()
}

// PostLogoutRedirectURIs implements [op.Client].
func (m *OidcClient) PostLogoutRedirectURIs() []string {
	return m.PostLogoutRedirectUriGlobsF
}

// RedirectURIs implements [op.Client].
func (m *OidcClient) RedirectURIs() []string {
	return m.RedirectUrisF
}

// ResponseTypes implements [op.Client].
func (m *OidcClient) ResponseTypes() []oidc.ResponseType {
	re := make([]oidc.ResponseType, 0, len(m.ResponseTypesF))
	for _, rt := range m.ResponseTypesF {
		re = append(re, ResponseTypeEx.ToGoType(rt))
	}
	return re
}

// RestrictAdditionalAccessTokenScopes implements [op.Client].
func (m *OidcClient) RestrictAdditionalAccessTokenScopes() func(scopes []string) []string {
	return func(scopes []string) []string { return scopes }
}

// RestrictAdditionalIdTokenScopes implements [op.Client].
func (m *OidcClient) RestrictAdditionalIdTokenScopes() func(scopes []string) []string {
	return func(scopes []string) []string { return scopes }
}
