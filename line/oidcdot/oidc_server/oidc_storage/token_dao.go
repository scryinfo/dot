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

// type TokenDao daobase.Dao[Token, *Token]
// func NewTokenDao(db *pebble2dot.Pebble2, logger *dot.LoggerType) TokenDao {
// 	return dao_pebble2.NewPointDaobase(db, logger, NewTokenById)
// }

type TokenDaoPebble2 struct {
	dao_pebble2.Daobase[Token, *Token]
}

func NewTokenDaoPebble2(db *pebble2dot.Pebble2, logger *dot.LoggerType) *TokenDaoPebble2 {
	return &TokenDaoPebble2{
		Daobase: dao_pebble2.NewDaobase(db, logger, NewTokenById),
	}
}

type TokenDaoBadger struct {
	dao_badger.Daobase[Token, *Token]
}

func NewTokenDaoBadger(db *badgerdot.BadgerDbDot, logger *dot.LoggerType) *TokenDaoBadger {
	return &TokenDaoBadger{
		Daobase: dao_badger.NewDaobase(db, logger, NewTokenById),
	}
}

type Token oidcapiv1.Token

var _ daobase.Modal = (*Token)(nil)

// func (o *Token) UnmarshalJSON(data []byte) error {
// 	if o.Token == nil {
// 		o.Token = &oidcapiv1.Token{}
// 	}
// 	return o.Token.UnmarshalJSON(data)
// }

// func (o *Token) MarshalJSON() ([]byte, error) {
// 	if o == nil || o.Token == nil {
// 		return []byte("null"), nil
// 	}
// 	return o.Token.MarshalJSON()
// }

func MakeTokenByProto(p *oidcapiv1.Token) *Token {
	return (*Token)(p)
}
func (m *Token) TokenToProto() *oidcapiv1.Token {
	return (*oidcapiv1.Token)(m)
}

// GetId implements [daobase.Modal].
// Subtle: this method shadows the method (*BanPlayers).GetId of BanPlayersM.BanPlayers.
func (m *Token) GetId() daobase.IdType {
	return daobase.IdType(m.Id)
}

// SetId implements [daobase.Modal].
func (m *Token) SetId(id daobase.IdType) {
	m.Id = string(id)
}

func NewToken() Token {
	return Token{
		Id: kits.Ids.NewXId(),
	}
}

func NewTokenById(id daobase.IdType) Token {
	return Token{
		Id: string(id),
	}
}

// Key implements [daobase.Modal].
func (m *Token) Key() []byte {
	return append(m.Prefix(), kits.StringToBytes(string(m.Id))...)
}

// Prefix implements [daobase.Modal].
func (m *Token) Prefix() []byte {
	return daokeys.PrefixToken
}

// Value implements [daobase.Modal].
func (m *Token) Value() ([]byte, error) {
	return daobase.Value(m)
}

// FromValue implements [daobase.Modal].
func (m *Token) FromValue(bs []byte) error {
	return daobase.FromValue(m, bs)
}
