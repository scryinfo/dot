package oidc_storage

import (
	"fmt"

	"github.com/scryinfo/dot/dot"
	oidcapiv1 "github.com/scryinfo/dot/line/oidcdot/oidc_gen/oidcapi/v1"
	"github.com/zitadel/oidc/v4/pkg/oidc"
	"github.com/zitadel/oidc/v4/pkg/op"
)

type _ApplicationTypeEx struct{}

var ApplicationTypeEx = _ApplicationTypeEx{}
var _ EnumTo[oidcapiv1.ApplicationType, op.ApplicationType] = (*_ApplicationTypeEx)(nil)

func (m *_ApplicationTypeEx) ToGoType(appType oidcapiv1.ApplicationType) op.ApplicationType {
	switch appType {
	case oidcapiv1.ApplicationType_APPLICATION_TYPE_UNSPECIFIED:
		dot.Logger.Error().Msgf("unknown application type, alg=%d", appType)
		return op.ApplicationTypeWeb
	case oidcapiv1.ApplicationType_APPLICATION_TYPE_WEB:
		return op.ApplicationTypeWeb
	case oidcapiv1.ApplicationType_APPLICATION_TYPE_USER_AGENT:
		return op.ApplicationTypeUserAgent
	case oidcapiv1.ApplicationType_APPLICATION_TYPE_NATIVE:
		return op.ApplicationTypeNative
	default:
		dot.Logger.Error().Msgf("unknown application type, alg=%d", appType)
		return op.ApplicationTypeWeb
	}
}

func (m *_ApplicationTypeEx) ToProto(appType op.ApplicationType) oidcapiv1.ApplicationType {
	switch appType {
	case op.ApplicationTypeWeb:
		return oidcapiv1.ApplicationType_APPLICATION_TYPE_WEB
	case op.ApplicationTypeUserAgent:
		return oidcapiv1.ApplicationType_APPLICATION_TYPE_USER_AGENT
	case op.ApplicationTypeNative:
		return oidcapiv1.ApplicationType_APPLICATION_TYPE_NATIVE
	default:
		dot.Logger.Error().Msgf("unknown application type, =%s", appType)
		return oidcapiv1.ApplicationType_APPLICATION_TYPE_WEB
	}
}

// EnumsGo implements [EnumTo].
func (m *_ApplicationTypeEx) EnumsGo() []op.ApplicationType {
	return []op.ApplicationType{
		op.ApplicationTypeWeb,
		op.ApplicationTypeUserAgent,
		op.ApplicationTypeNative,
	}
}

// EnumsProto implements [EnumTo].
func (m *_ApplicationTypeEx) EnumsProto() []oidcapiv1.ApplicationType {
	return []oidcapiv1.ApplicationType{
		oidcapiv1.ApplicationType_APPLICATION_TYPE_WEB,
		oidcapiv1.ApplicationType_APPLICATION_TYPE_USER_AGENT,
		oidcapiv1.ApplicationType_APPLICATION_TYPE_NATIVE,
	}
}

type _AuthMethodEx struct{}

var AuthMethodEx = _AuthMethodEx{}
var _ EnumTo[oidcapiv1.AuthMethod, oidc.AuthMethod] = (*_AuthMethodEx)(nil)

func (m *_AuthMethodEx) ToGoType(p oidcapiv1.AuthMethod) oidc.AuthMethod {
	switch p {
	case oidcapiv1.AuthMethod_AUTH_METHOD_UNSPECIFIED:
		dot.Logger.Error().Msgf("unknown auth method, =%d", p)
		return oidc.AuthMethodNone
	case oidcapiv1.AuthMethod_AUTH_METHOD_BASIC:
		return oidc.AuthMethodBasic
	case oidcapiv1.AuthMethod_AUTH_METHOD_POST:
		return oidc.AuthMethodPost
	case oidcapiv1.AuthMethod_AUTH_METHOD_PRIVATE_KEY_JWT:
		return oidc.AuthMethodPrivateKeyJWT
	default:
		dot.Logger.Error().Msgf("unknown auth method, =%d", p)
		return oidc.AuthMethodNone
	}
}

func (m *_AuthMethodEx) ToProto(g oidc.AuthMethod) oidcapiv1.AuthMethod {
	switch g {
	case oidc.AuthMethodBasic:
		return oidcapiv1.AuthMethod_AUTH_METHOD_BASIC
	case oidc.AuthMethodPost:
		return oidcapiv1.AuthMethod_AUTH_METHOD_POST
	case oidc.AuthMethodPrivateKeyJWT:
		return oidcapiv1.AuthMethod_AUTH_METHOD_PRIVATE_KEY_JWT
	default:
		dot.Logger.Error().Msgf("unknown auth method, =%s", g)
		return oidcapiv1.AuthMethod_AUTH_METHOD_NONE
	}
}

// EnumsGo implements [EnumTo].
func (m *_AuthMethodEx) EnumsGo() []oidc.AuthMethod {
	return []oidc.AuthMethod{
		oidc.AuthMethodBasic,
		oidc.AuthMethodPost,
		oidc.AuthMethodPrivateKeyJWT,
	}
}

// EnumsProto implements [EnumTo].
func (m *_AuthMethodEx) EnumsProto() []oidcapiv1.AuthMethod {
	return []oidcapiv1.AuthMethod{
		oidcapiv1.AuthMethod_AUTH_METHOD_BASIC,
		oidcapiv1.AuthMethod_AUTH_METHOD_POST,
		oidcapiv1.AuthMethod_AUTH_METHOD_PRIVATE_KEY_JWT,
	}
}

type _AccessTokenTypeEx struct{}

var AccessTokenTypeEx = _AccessTokenTypeEx{}
var _ EnumTo[oidcapiv1.AccessTokenType, op.AccessTokenType] = (*_AccessTokenTypeEx)(nil)

func (m *_AccessTokenTypeEx) ToGoType(p oidcapiv1.AccessTokenType) op.AccessTokenType {
	switch p {
	case oidcapiv1.AccessTokenType_ACCESS_TOKEN_TYPE_UNSPECIFIED:
		dot.Logger.Error().Msgf("unknown access token type, =%d", p)
		return op.AccessTokenTypeBearer
	case oidcapiv1.AccessTokenType_ACCESS_TOKEN_TYPE_BEARER:
		return op.AccessTokenTypeBearer
	case oidcapiv1.AccessTokenType_ACCESS_TOKEN_TYPE_JWT:
		return op.AccessTokenTypeJWT
	default:
		dot.Logger.Error().Msgf("unknown access token type, =%d", p)
		return op.AccessTokenTypeBearer
	}
}

func (m *_AccessTokenTypeEx) ToProto(g op.AccessTokenType) oidcapiv1.AccessTokenType {
	switch g {
	case op.AccessTokenTypeBearer:
		return oidcapiv1.AccessTokenType_ACCESS_TOKEN_TYPE_BEARER
	case op.AccessTokenTypeJWT:
		return oidcapiv1.AccessTokenType_ACCESS_TOKEN_TYPE_JWT
	default:
		dot.Logger.Error().Msgf("unknown access token type, =%s", g)
		return oidcapiv1.AccessTokenType_ACCESS_TOKEN_TYPE_UNSPECIFIED
	}
}

// EnumsGo implements [EnumTo].
func (m *_AccessTokenTypeEx) EnumsGo() []op.AccessTokenType {
	return []op.AccessTokenType{
		op.AccessTokenTypeBearer,
		op.AccessTokenTypeJWT,
	}
}

// EnumsProto implements [EnumTo].
func (m *_AccessTokenTypeEx) EnumsProto() []oidcapiv1.AccessTokenType {
	return []oidcapiv1.AccessTokenType{
		oidcapiv1.AccessTokenType_ACCESS_TOKEN_TYPE_BEARER,
		oidcapiv1.AccessTokenType_ACCESS_TOKEN_TYPE_JWT,
	}
}

type _ResponseTypeEx struct{}

var ResponseTypeEx = _ResponseTypeEx{}
var _ EnumTo[oidcapiv1.ResponseType, oidc.ResponseType] = (*_ResponseTypeEx)(nil)

func (m *_ResponseTypeEx) ToGoType(p oidcapiv1.ResponseType) oidc.ResponseType {
	switch p {
	case oidcapiv1.ResponseType_RESPONSE_TYPE_UNSPECIFIED:
		dot.Logger.Error().Msgf("unknown response type, =%d", p)
		return oidc.ResponseTypeCode
	case oidcapiv1.ResponseType_RESPONSE_TYPE_CODE:
		return oidc.ResponseTypeCode
	case oidcapiv1.ResponseType_RESPONSE_TYPE_ID_TOKEN:
		return oidc.ResponseTypeIDToken
	case oidcapiv1.ResponseType_RESPONSE_TYPE_ID_TOKEN_ONLY:
		return oidc.ResponseTypeIDTokenOnly
	default:
		dot.Logger.Error().Msgf("unknown response type, =%d", p)
		return oidc.ResponseTypeCode
	}
}

func (m *_ResponseTypeEx) ToProto(g oidc.ResponseType) oidcapiv1.ResponseType {
	switch g {
	case oidc.ResponseTypeCode:
		return oidcapiv1.ResponseType_RESPONSE_TYPE_CODE
	case oidc.ResponseTypeIDToken:
		return oidcapiv1.ResponseType_RESPONSE_TYPE_ID_TOKEN
	case oidc.ResponseTypeIDTokenOnly:
		return oidcapiv1.ResponseType_RESPONSE_TYPE_ID_TOKEN_ONLY
	default:
		dot.Logger.Error().Msgf("unknown response type, =%s", g)
		return oidcapiv1.ResponseType_RESPONSE_TYPE_CODE
	}
}

// EnumsGo implements [EnumTo].
func (m *_ResponseTypeEx) EnumsGo() []oidc.ResponseType {
	return []oidc.ResponseType{
		oidc.ResponseTypeCode,
		oidc.ResponseTypeIDToken,
		oidc.ResponseTypeIDTokenOnly,
	}
}

// EnumsProto implements [EnumTo].
func (m *_ResponseTypeEx) EnumsProto() []oidcapiv1.ResponseType {
	return []oidcapiv1.ResponseType{
		oidcapiv1.ResponseType_RESPONSE_TYPE_CODE,
		oidcapiv1.ResponseType_RESPONSE_TYPE_ID_TOKEN,
		oidcapiv1.ResponseType_RESPONSE_TYPE_ID_TOKEN_ONLY,
	}
}

// GrantType
type _GrantTypeEx struct{}

var GrantTypeEx = _GrantTypeEx{}
var _ EnumTo[oidcapiv1.GrantType, oidc.GrantType] = (*_GrantTypeEx)(nil)

func (m *_GrantTypeEx) ToGoType(p oidcapiv1.GrantType) oidc.GrantType {
	switch p {
	case oidcapiv1.GrantType_GRANT_TYPE_UNSPECIFIED:
		dot.Logger.Error().Msgf("unknown grant type, =%d", p)
		panic(fmt.Errorf("unknown grant type, =%d", p))
	case oidcapiv1.GrantType_GRANT_TYPE_CODE:
		return oidc.GrantTypeCode
	case oidcapiv1.GrantType_GRANT_TYPE_REFRESH_TOKEN:
		return oidc.GrantTypeRefreshToken
	case oidcapiv1.GrantType_GRANT_TYPE_CLIENT_CREDENTIALS:
		return oidc.GrantTypeClientCredentials
	case oidcapiv1.GrantType_GRANT_TYPE_BEARER:
		return oidc.GrantTypeBearer
	case oidcapiv1.GrantType_GRANT_TYPE_TOKEN_EXCHANGE:
		return oidc.GrantTypeTokenExchange
	case oidcapiv1.GrantType_GRANT_TYPE_IMPLICIT:
		return oidc.GrantTypeImplicit
	case oidcapiv1.GrantType_GRANT_TYPE_DEVICE_CODE:
		return oidc.GrantTypeDeviceCode
	default:
		dot.Logger.Error().Msgf("unknown grant type, =%d", p)
		panic(fmt.Errorf("unknown grant type, =%d", p))
	}
}

func (m *_GrantTypeEx) ToProto(g oidc.GrantType) oidcapiv1.GrantType {
	switch g {
	case oidc.GrantTypeCode:
		return oidcapiv1.GrantType_GRANT_TYPE_CODE
	case oidc.GrantTypeRefreshToken:
		return oidcapiv1.GrantType_GRANT_TYPE_REFRESH_TOKEN
	case oidc.GrantTypeClientCredentials:
		return oidcapiv1.GrantType_GRANT_TYPE_CLIENT_CREDENTIALS
	case oidc.GrantTypeBearer:
		return oidcapiv1.GrantType_GRANT_TYPE_BEARER
	case oidc.GrantTypeTokenExchange:
		return oidcapiv1.GrantType_GRANT_TYPE_TOKEN_EXCHANGE
	case oidc.GrantTypeImplicit:
		return oidcapiv1.GrantType_GRANT_TYPE_IMPLICIT
	case oidc.GrantTypeDeviceCode:
		return oidcapiv1.GrantType_GRANT_TYPE_DEVICE_CODE
	default:
		dot.Logger.Error().Msgf("unknown grant type, =%s", g)
		return oidcapiv1.GrantType_GRANT_TYPE_UNSPECIFIED
	}
}

// EnumsGo implements [EnumTo].
func (m *_GrantTypeEx) EnumsGo() []oidc.GrantType {
	return []oidc.GrantType{
		oidc.GrantTypeCode,
		oidc.GrantTypeRefreshToken,
		oidc.GrantTypeClientCredentials,
		oidc.GrantTypeBearer,
		oidc.GrantTypeTokenExchange,
		oidc.GrantTypeImplicit,
		oidc.GrantTypeDeviceCode,
	}
}

// EnumsProto implements [EnumTo].
func (m *_GrantTypeEx) EnumsProto() []oidcapiv1.GrantType {
	return []oidcapiv1.GrantType{
		oidcapiv1.GrantType_GRANT_TYPE_CODE,
		oidcapiv1.GrantType_GRANT_TYPE_REFRESH_TOKEN,
		oidcapiv1.GrantType_GRANT_TYPE_CLIENT_CREDENTIALS,
		oidcapiv1.GrantType_GRANT_TYPE_BEARER,
		oidcapiv1.GrantType_GRANT_TYPE_TOKEN_EXCHANGE,
		oidcapiv1.GrantType_GRANT_TYPE_IMPLICIT,
		oidcapiv1.GrantType_GRANT_TYPE_DEVICE_CODE,
	}
}

// Scope
type _ScopeEx struct{}

var ScopeEx = _ScopeEx{}
var _ EnumTo[oidcapiv1.Scope, string] = (*_ScopeEx)(nil)

func (m *_ScopeEx) ToGoType(p oidcapiv1.Scope) string {
	switch p {
	case oidcapiv1.Scope_SCOPE_UNSPECIFIED:
		dot.Logger.Error().Msgf("unknown scope, =%d", p)
		panic(fmt.Errorf("unknown scope, =%d", p))
	case oidcapiv1.Scope_SCOPE_OPENID:
		return oidc.ScopeOpenID
	case oidcapiv1.Scope_SCOPE_PROFILE:
		return oidc.ScopeProfile
	case oidcapiv1.Scope_SCOPE_EMAIL:
		return oidc.ScopeEmail
	case oidcapiv1.Scope_SCOPE_ADDRESS:
		return oidc.ScopeAddress
	case oidcapiv1.Scope_SCOPE_PHONE:
		return oidc.ScopePhone
	case oidcapiv1.Scope_SCOPE_OFFLINE_ACCESS:
		return oidc.ScopeOfflineAccess
	default:
		dot.Logger.Error().Msgf("unknown scope, =%d", p)
		panic(fmt.Errorf("unknown scope, =%d", p))
	}
}

func (m *_ScopeEx) ToProto(scope string) oidcapiv1.Scope {
	switch scope {
	case oidc.ScopeOpenID:
		return oidcapiv1.Scope_SCOPE_OPENID
	case oidc.ScopeProfile:
		return oidcapiv1.Scope_SCOPE_PROFILE
	case oidc.ScopeEmail:
		return oidcapiv1.Scope_SCOPE_EMAIL
	case oidc.ScopeAddress:
		return oidcapiv1.Scope_SCOPE_ADDRESS
	case oidc.ScopePhone:
		return oidcapiv1.Scope_SCOPE_PHONE
	case oidc.ScopeOfflineAccess:
		return oidcapiv1.Scope_SCOPE_OFFLINE_ACCESS
	default:
		dot.Logger.Error().Msgf("unknown scope, =%s", scope)
		return oidcapiv1.Scope_SCOPE_UNSPECIFIED
	}
}

func (m *_ScopeEx) EnumsGo() []string {
	return []string{
		oidc.ScopeOpenID,
		oidc.ScopeProfile,
		oidc.ScopeEmail,
		oidc.ScopeAddress,
		oidc.ScopePhone,
		oidc.ScopeOfflineAccess,
	}
}
func (m *_ScopeEx) EnumsProto() []oidcapiv1.Scope {
	return []oidcapiv1.Scope{
		oidcapiv1.Scope_SCOPE_OPENID,
		oidcapiv1.Scope_SCOPE_PROFILE,
		oidcapiv1.Scope_SCOPE_EMAIL,
		oidcapiv1.Scope_SCOPE_ADDRESS,
		oidcapiv1.Scope_SCOPE_PHONE,
		oidcapiv1.Scope_SCOPE_OFFLINE_ACCESS,
	}
}
