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

type AppSessionDaoPebble2 struct {
	dao_pebble2.Daobase[AppSession, *AppSession]
}
type AppSessionDaoBadger struct {
	dao_badger.Daobase[AppSession, *AppSession]
}

func NewAppSessionDaoPebble2(db *pebble2dot.Pebble2, logger *dot.LoggerType) *AppSessionDaoPebble2 {
	return &AppSessionDaoPebble2{
		Daobase: dao_pebble2.NewDaobase(db, logger, NewAppSessionById),
	}
}

func NewAppSessionDaoBadger(db *badgerdot.BadgerDbDot, logger *dot.LoggerType) *AppSessionDaoBadger {
	return &AppSessionDaoBadger{
		Daobase: dao_badger.NewDaobase(db, logger, NewAppSessionById),
	}
}

type AppSession oidcapiv1.AppSession

// Expire implements [daobase.Modal].
func (m *AppSession) Expire() bool {
	return daobase.ModalExpire(m.ExpireTs)
}

// GetExpireTs implements [daobase.Modal].
func (m *AppSession) GetExpireTs() uint64 {
	return m.ExpireTs
}

// SetExpireTs implements [daobase.Modal].
func (m *AppSession) SetExpireTs(ts uint64) {
	m.ExpireTs = ts
}

var _ daobase.Modal = (*AppSession)(nil)

func MakeAppSessionByProto(p *oidcapiv1.AppSession) *AppSession {
	return (*AppSession)(p)
}
func (m *AppSession) AppSessionToProto() *oidcapiv1.AppSession {
	return (*oidcapiv1.AppSession)(m)
}

// GetId implements [daobase.Modal].
// Subtle: this method shadows the method (*BanPlayers).GetId of BanPlayersM.BanPlayers.
func (m *AppSession) GetId() daobase.IdType {
	return daobase.IdType(m.Id)
}

// SetId implements [daobase.Modal].
func (m *AppSession) SetId(id daobase.IdType) {
	m.Id = string(id)
}

func NewAppSession() AppSession {
	return AppSession{
		Id: kits.Ids.NewXId(),
	}
}

func NewAppSessionById(id daobase.IdType) AppSession {
	return AppSession{
		Id: string(id),
	}
}

// Key implements [daobase.Modal].
func (m *AppSession) Key() []byte {
	return append(m.Prefix(), kits.StringToBytes(string(m.Id))...)
}

// Prefix implements [daobase.Modal].
func (m *AppSession) Prefix() []byte {
	return daokeys.PrefixAppSession
}

// Value implements [daobase.Modal].
func (m *AppSession) Value() ([]byte, error) {
	return daobase.Value(m)
}

// FromValue implements [daobase.Modal].
func (m *AppSession) FromValue(bs []byte) error {
	return daobase.FromValue(m, bs)
}
