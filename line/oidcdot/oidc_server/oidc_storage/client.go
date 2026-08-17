package oidc_storage

import (
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/lib/kits"
	daobase "github.com/scryinfo/dot/line/db/dao/dao_base"
	daokeys "github.com/scryinfo/dot/line/db/dao/dao_keys"
	"github.com/scryinfo/dot/line/db/pebble2dot"
	"github.com/scryinfo/dot/line/db/pebble2dot/dao_pebble2"
	oidcapiv1 "github.com/scryinfo/dot/line/oidcdot/oidc_gen/oidcapi/v1"
)

type OidcClientDao struct {
	dao_pebble2.Daobase[OidcClient, *OidcClient]
}

func NewOidcClientDao(db *pebble2dot.Pebble2, logger *dot.LoggerType) *OidcClientDao {
	return &OidcClientDao{
		Daobase: dao_pebble2.NewDaobase(db, logger, NewOidcClientById),
	}
}

type OidcClient oidcapiv1.OidcClient

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
// Subtle: this method shadows the method (*BanPlayers).GetId of BanPlayersM.BanPlayers.
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
