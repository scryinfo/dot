package oidc_storage

import (
	"time"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/lib/kits"
	"github.com/scryinfo/dot/line/db/badgerdot"
	"github.com/scryinfo/dot/line/db/badgerdot/dao_badger"
	daobase "github.com/scryinfo/dot/line/db/dao/dao_base"
	daokeys "github.com/scryinfo/dot/line/db/dao/dao_keys"
	"github.com/scryinfo/dot/line/db/pebble2dot"
	"github.com/scryinfo/dot/line/db/pebble2dot/dao_pebble2"
	oidcapiv1 "github.com/scryinfo/dot/line/oidcdot/oidc_gen/oidcapi/v1"
	"github.com/zitadel/oidc/v4/pkg/oidc"
	"github.com/zitadel/oidc/v4/pkg/op"
)

type OidcClientDaoPebble2 struct {
	dao_pebble2.Daobase[OidcClient, *OidcClient]
}

type OidcClientDaoBadger struct {
	dao_badger.Daobase[OidcClient, *OidcClient]
}

func NewOidcClientDaoPebble2(db *pebble2dot.Pebble2, logger *dot.LoggerType) *OidcClientDaoPebble2 {
	return &OidcClientDaoPebble2{
		Daobase: dao_pebble2.NewDaobase(db, logger, NewOidcClientById),
	}
}

func NewOidcClientDaoBadger(db *badgerdot.BadgerDbDot, logger *dot.LoggerType) *OidcClientDaoBadger {
	return &OidcClientDaoBadger{
		Daobase: dao_badger.NewDaobase(db, logger, NewOidcClientById),
	}
}

type OidcClient oidcapiv1.OidcClient

// Expire implements [daobase.Modal].
func (m *OidcClient) Expire() bool {
	return daobase.ModalExpire(m.ExpireTs)
}

// GetExpireTs implements [daobase.Modal].
func (m *OidcClient) GetExpireTs() uint64 {
	return m.ExpireTs
}

// SetExpireTs implements [daobase.Modal].
func (m *OidcClient) SetExpireTs(ts uint64) {
	m.ExpireTs = ts
}

var _ daobase.Modal = (*OidcClient)(nil)

// func (o *OidcClient) UnmarshalJSON(data []byte) error {
// 	if o.OidcClient == nil {
// 		o.OidcClient = &oidcapiv1.OidcClient{}
// 	}
// 	return o.OidcClient.UnmarshalJSON(data)
// }

// func (o *OidcClient) MarshalJSON() ([]byte, error) {
// 	if o == nil || o.OidcClient == nil {
// 		return []byte("null"), nil
// 	}
// 	return o.OidcClient.MarshalJSON()
// }

func MakeOidcClientByProto(p *oidcapiv1.OidcClient) *OidcClient {
	return (*OidcClient)(p)
}
func (m *OidcClient) OidcClientToProto() *oidcapiv1.OidcClient {
	return (*oidcapiv1.OidcClient)(m)
}

// GetId implements [daobase.Modal].
func (m *OidcClient) GetId() daobase.IdType {
	return daobase.IdType(m.Id)
}

// SetId implements [daobase.Modal].
func (m *OidcClient) SetId(id daobase.IdType) {
	m.Id = string(id)
}

func NewOidcClient() OidcClient {
	return OidcClient{
		Id: kits.Ids.NewXId(),
	}
}

func NewOidcClientById(id daobase.IdType) OidcClient {
	return OidcClient{
		Id: string(id),
	}
}

// Key implements [daobase.Modal].
func (m *OidcClient) Key() []byte {
	return append(m.Prefix(), kits.StringToBytes(string(m.Id))...)
}

// Prefix implements [daobase.Modal].
func (m *OidcClient) Prefix() []byte {
	return daokeys.PrefixOidcClient
}

// Value implements [daobase.Modal].
func (m *OidcClient) Value() ([]byte, error) {
	return daobase.Value(m)
}

// FromValue implements [daobase.Modal].
func (m *OidcClient) FromValue(bs []byte) error {
	return daobase.FromValue(m, bs)
}

// implements [op.Client]
var _ op.Client = (*OidcClient)(nil)

// AccessTokenType implements [op.Client].
func (m *OidcClient) AccessTokenType() op.AccessTokenType {
	return op.AccessTokenType(m.AccessTokenTypeF)
}

// ApplicationType implements [op.Client].
func (m *OidcClient) ApplicationType() op.ApplicationType {
	panic("unimplemented")
}

// AuthMethod implements [op.Client].
func (m *OidcClient) AuthMethod() oidc.AuthMethod {
	panic("unimplemented")
}

// ClockSkew implements [op.Client].
func (m *OidcClient) ClockSkew() time.Duration {
	panic("unimplemented")
}

// DevMode implements [op.Client].
func (m *OidcClient) DevMode() bool {
	panic("unimplemented")
}

// GetID implements [op.Client].
func (m *OidcClient) GetID() string {
	panic("unimplemented")
}

// GrantTypes implements [op.Client].
func (m *OidcClient) GrantTypes() []oidc.GrantType {
	panic("unimplemented")
}

// IDTokenLifetime implements [op.Client].
func (m *OidcClient) IDTokenLifetime() time.Duration {
	panic("unimplemented")
}

// IDTokenUserinfoClaimsAssertion implements [op.Client].
func (m *OidcClient) IDTokenUserinfoClaimsAssertion() bool {
	panic("unimplemented")
}

// IsScopeAllowed implements [op.Client].
func (m *OidcClient) IsScopeAllowed(scope string) bool {
	panic("unimplemented")
}

// LoginURL implements [op.Client].
func (m *OidcClient) LoginURL(string) string {
	panic("unimplemented")
}

// PostLogoutRedirectURIs implements [op.Client].
func (m *OidcClient) PostLogoutRedirectURIs() []string {
	panic("unimplemented")
}

// RedirectURIs implements [op.Client].
func (m *OidcClient) RedirectURIs() []string {
	panic("unimplemented")
}

// ResponseTypes implements [op.Client].
func (m *OidcClient) ResponseTypes() []oidc.ResponseType {
	panic("unimplemented")
}

// RestrictAdditionalAccessTokenScopes implements [op.Client].
func (m *OidcClient) RestrictAdditionalAccessTokenScopes() func(scopes []string) []string {
	panic("unimplemented")
}

// RestrictAdditionalIdTokenScopes implements [op.Client].
func (m *OidcClient) RestrictAdditionalIdTokenScopes() func(scopes []string) []string {
	panic("unimplemented")
}
