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

type RefreshTokenDaoPebble2 struct {
	dao_pebble2.Daobase[RefreshToken, *RefreshToken]
}

func NewRefreshTokenDaoPebble2(db *pebble2dot.Pebble2, logger *dot.LoggerType) *RefreshTokenDaoPebble2 {
	return &RefreshTokenDaoPebble2{
		Daobase: dao_pebble2.NewDaobase(db, logger, NewRefreshTokenById),
	}
}

type RefreshToken oidcapiv1.RefreshToken

// Expire implements [daobase.Modal].
func (m *RefreshToken) Expire() bool {
	return daobase.ModalExpire(m.ExpireTs)
}

// GetExpireTs implements [daobase.Modal].
func (m *RefreshToken) GetExpireTs() uint64 {
	return m.ExpireTs
}

// SetExpireTs implements [daobase.Modal].
func (m *RefreshToken) SetExpireTs(ts uint64) {
	m.ExpireTs = ts
}

var _ daobase.Modal = (*RefreshToken)(nil)

// func (o *RefreshToken) UnmarshalJSON(data []byte) error {
// 	if o.RefreshToken == nil {
// 		o.RefreshToken = &oidcapiv1.RefreshToken{}
// 	}
// 	return o.RefreshToken.UnmarshalJSON(data)
// }

// func (o *RefreshToken) MarshalJSON() ([]byte, error) {
// 	if o == nil || o.RefreshToken == nil {
// 		return []byte("null"), nil
// 	}
// 	return o.RefreshToken.MarshalJSON()
// }

func MakeRefreshTokenByProto(p *oidcapiv1.RefreshToken) *RefreshToken {
	return (*RefreshToken)(p)
}
func (m *RefreshToken) RefreshTokenToProto() *oidcapiv1.RefreshToken {
	return (*oidcapiv1.RefreshToken)(m)
}

// GetId implements [daobase.Modal].
// Subtle: this method shadows the method (*BanPlayers).GetId of BanPlayersM.BanPlayers.
func (m *RefreshToken) GetId() daobase.IdType {
	return daobase.IdType(m.Id)
}

// SetId implements [daobase.Modal].
func (m *RefreshToken) SetId(id daobase.IdType) {
	m.Id = string(id)
}

func NewRefreshToken() RefreshToken {
	return RefreshToken{
		Id: kits.Ids.NewXId(),
	}
}

func NewRefreshTokenById(id daobase.IdType) RefreshToken {
	return RefreshToken{
		Id: string(id),
	}
}

// Key implements [daobase.Modal].
func (m *RefreshToken) Key() []byte {
	return append(m.Prefix(), kits.StringToBytes(string(m.Id))...)
}

// Prefix implements [daobase.Modal].
func (m *RefreshToken) Prefix() []byte {
	return daokeys.PrefixRefreshToken
}

// Value implements [daobase.Modal].
func (m *RefreshToken) Value() ([]byte, error) {
	return daobase.Value(m)
}

// FromValue implements [daobase.Modal].
func (m *RefreshToken) FromValue(bs []byte) error {
	return daobase.FromValue(m, bs)
}
