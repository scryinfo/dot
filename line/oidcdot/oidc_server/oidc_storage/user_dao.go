package oidc_storage

import (
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/lib/kits"
	"github.com/scryinfo/dot/line/db/badgerdot"
	"github.com/scryinfo/dot/line/db/badgerdot/dao_badger"
	daobase "github.com/scryinfo/dot/line/db/dao/dao_base"
	daokeys "github.com/scryinfo/dot/line/db/dao/dao_keys"
	"github.com/scryinfo/dot/line/db/pebble2dot"
	"github.com/scryinfo/dot/line/db/pebble2dot/dao_pebble2"
	oidcapiv1 "github.com/scryinfo/dot/line/oidcdot/oidc_gen/oidcapi/v1"
)

type UserDaoPebble2 struct {
	dao_pebble2.Daobase[User, *User]
}
type UserDaoBadger struct {
	dao_badger.Daobase[User, *User]
}

func NewUserDaoPebble2(db *pebble2dot.Pebble2, logger *dot.LoggerType) *UserDaoPebble2 {
	return &UserDaoPebble2{
		Daobase: dao_pebble2.NewDaobase(db, logger, NewUserById),
	}
}

func NewUserDaoBadger(db *badgerdot.BadgerDbDot, logger *dot.LoggerType) *UserDaoBadger {
	return &UserDaoBadger{
		Daobase: dao_badger.NewDaobase(db, logger, NewUserById),
	}
}

type User oidcapiv1.User

// Expire implements [daobase.Modal].
func (m *User) Expire() bool {
	return daobase.ModalExpire(m.ExpireTs)
}

// GetExpireTs implements [daobase.Modal].
func (m *User) GetExpireTs() uint64 {
	return m.ExpireTs
}

// SetExpireTs implements [daobase.Modal].
func (m *User) SetExpireTs(ts uint64) {
	m.ExpireTs = ts
}

var _ daobase.Modal = (*User)(nil)

// func (o *User) UnmarshalJSON(data []byte) error {
// 	if o.User == nil {
// 		o.User = &oidcapiv1.User{}
// 	}
// 	return o.User.UnmarshalJSON(data)
// }

// func (o *User) MarshalJSON() ([]byte, error) {
// 	if o == nil || o.User == nil {
// 		return []byte("null"), nil
// 	}
// 	return o.User.MarshalJSON()
// }

func MakeUserByProto(p *oidcapiv1.User) *User {
	return (*User)(p)
}
func (m *User) UserToProto() *oidcapiv1.User {
	return (*oidcapiv1.User)(m)
}

// GetId implements [daobase.Modal].
// Subtle: this method shadows the method (*BanPlayers).GetId of BanPlayersM.BanPlayers.
func (m *User) GetId() daobase.IdType {
	return daobase.IdType(m.Id)
}

// SetId implements [daobase.Modal].
func (m *User) SetId(id daobase.IdType) {
	m.Id = string(id)
}

func NewUser() User {
	return User{
		Id: kits.Ids.NewXId(),
	}
}

func NewUserById(id daobase.IdType) User {
	return User{
		Id: string(id),
	}
}

// Key implements [daobase.Modal].
func (m *User) Key() []byte {
	return append(m.Prefix(), kits.StringToBytes(string(m.Id))...)
}

// Prefix implements [daobase.Modal].
func (m *User) Prefix() []byte {
	return daokeys.PrefixUser
}

// Value implements [daobase.Modal].
func (m *User) Value() ([]byte, error) {
	return daobase.Value(m)
}

// FromValue implements [daobase.Modal].
func (m *User) FromValue(bs []byte) error {
	return daobase.FromValue(m, bs)
}
