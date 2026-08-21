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

type UserIdentitiesDao struct {
	daobase.Dao[UserIdentities, *UserIdentities]
}

func NewUserIdentitiesDao(db *pebble2dot.Pebble2, logger *dot.LoggerType) *UserIdentitiesDao {
	return &UserIdentitiesDao{
		Dao: dao_pebble2.NewPointDaobase(db, logger, NewUserIdentitiesById),
	}
}

type UserIdentities oidcapiv1.UserIdentities

var _ daobase.Modal = (*UserIdentities)(nil)

// func (o *UserIdentities) UnmarshalJSON(data []byte) error {
// 	if o.UserIdentities == nil {
// 		o.UserIdentities = &oidcapiv1.UserIdentities{}
// 	}
// 	return o.UserIdentities.UnmarshalJSON(data)
// }

// func (o *UserIdentities) MarshalJSON() ([]byte, error) {
// 	if o == nil || o.UserIdentities == nil {
// 		return []byte("null"), nil
// 	}
// 	return o.UserIdentities.MarshalJSON()
// }

func MakeUserIdentitiesByProto(p *oidcapiv1.UserIdentities) *UserIdentities {
	return (*UserIdentities)(p)
}
func (m *UserIdentities) UserIdentitiesToProto() *oidcapiv1.UserIdentities {
	return (*oidcapiv1.UserIdentities)(m)
}

// GetId implements [daobase.Modal].
// Subtle: this method shadows the method (*BanPlayers).GetId of BanPlayersM.BanPlayers.
func (m *UserIdentities) GetId() daobase.IdType {
	return daobase.IdType(m.UserId)
}

// SetId implements [daobase.Modal].
func (m *UserIdentities) SetId(id daobase.IdType) {
	m.UserId = string(id)
}

func NewUserIdentities() UserIdentities {
	return UserIdentities{
		UserId: kits.Ids.NewXId(),
	}
}

func NewUserIdentitiesById(id daobase.IdType) UserIdentities {
	return UserIdentities{
		UserId: string(id),
	}
}

// Key implements [daobase.Modal].
func (m *UserIdentities) Key() []byte {
	return append(m.Prefix(), kits.StringToBytes(string(m.UserId))...)
}

// Prefix implements [daobase.Modal].
func (m *UserIdentities) Prefix() []byte {
	return daokeys.PrefixUserIdentities
}

// Value implements [daobase.Modal].
func (m *UserIdentities) Value() ([]byte, error) {
	return daobase.Value(m)
}

// FromValue implements [daobase.Modal].
func (m *UserIdentities) FromValue(bs []byte) error {
	return daobase.FromValue(m, bs)
}
