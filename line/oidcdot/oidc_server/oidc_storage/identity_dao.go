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

type IdentityDao struct {
	dao_pebble2.Daobase[Identity, *Identity]
}

func NewIdentityDao(db *pebble2dot.Pebble2, logger *dot.LoggerType) *IdentityDao {
	return &IdentityDao{
		Daobase: dao_pebble2.NewDaobase(db, logger, NewIdentityById),
	}
}

type Identity oidcapiv1.Identity

var _ daobase.Modal = (*Identity)(nil)

// func (o *Identity) UnmarshalJSON(data []byte) error {
// 	if o.Identity == nil {
// 		o.Identity = &oidcapiv1.Identity{}
// 	}
// 	return o.Identity.UnmarshalJSON(data)
// }

// func (o *Identity) MarshalJSON() ([]byte, error) {
// 	if o == nil || o.Identity == nil {
// 		return []byte("null"), nil
// 	}
// 	return o.Identity.MarshalJSON()
// }

func MakeIdentityByProto(p *oidcapiv1.Identity) *Identity {
	return (*Identity)(p)
}
func (m *Identity) IdentityToProto() *oidcapiv1.Identity {
	return (*oidcapiv1.Identity)(m)
}

// GetId implements [daobase.Modal].
// Subtle: this method shadows the method (*BanPlayers).GetId of BanPlayersM.BanPlayers.
func (m *Identity) GetId() daobase.IdType {
	return daobase.IdType(m.Id)
}

// SetId implements [daobase.Modal].
func (m *Identity) SetId(id daobase.IdType) {
	m.Id = string(id)
}

func NewIdentity() Identity {
	return Identity{
		Id: kits.Ids.NewXId(),
	}
}

func NewIdentityById(id daobase.IdType) Identity {
	return Identity{
		Id: string(id),
	}
}

// Key implements [daobase.Modal].
func (m *Identity) Key() []byte {
	return append(m.Prefix(), kits.StringToBytes(string(m.Id))...)
}

// Prefix implements [daobase.Modal].
func (m *Identity) Prefix() []byte {
	return daokeys.PrefixIdentity
}

// Value implements [daobase.Modal].
func (m *Identity) Value() ([]byte, error) {
	return daobase.Value(m)
}

// FromValue implements [daobase.Modal].
func (m *Identity) FromValue(bs []byte) error {
	return daobase.FromValue(m, bs)
}
