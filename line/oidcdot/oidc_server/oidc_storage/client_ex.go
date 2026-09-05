package oidc_storage

import (
	"github.com/scryinfo/dot/dot"
	oidcapiv1 "github.com/scryinfo/dot/line/oidcdot/oidc_gen/oidcapi/v1"
	"github.com/zitadel/oidc/v4/pkg/oidc"
	"github.com/zitadel/oidc/v4/pkg/op"
)

type _ApplicationTypeEx int

var ApplicationTypeEx = _ApplicationTypeEx(0)
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

type _AuthMethodEx int

var AuthMethodEx = _AuthMethodEx(0)
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

type _AccessTokenTypeEx int

var AccessTokenTypeEx = _AccessTokenTypeEx(0)
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

type _ResponseTypeEx int

var ResponseTypeEx = _ResponseTypeEx(0)
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
