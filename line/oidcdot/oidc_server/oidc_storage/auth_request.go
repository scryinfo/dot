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

type AuthRequestDao struct {
	dao_pebble2.Daobase[AuthRequest, *AuthRequest]
}

func NewBanPlayersDao(db *pebble2dot.Pebble2, logger *dot.LoggerType) *AuthRequestDao {
	return &AuthRequestDao{
		Daobase: dao_pebble2.NewDaobase(db, logger, func(id daobase.IdType) AuthRequest {
			return AuthRequest{
				AuthRequest: &oidcapiv1.AuthRequest{
					Id: string(id),
					// Ts: kits.Tss.Ts(),
				},
			}
		}),
	}
}

// dao

var _ daobase.Modal = (*AuthRequest)(nil)

func (o *AuthRequest) UnmarshalJSON(data []byte) error {
	if o.AuthRequest == nil {
		o.AuthRequest = &oidcapiv1.AuthRequest{}
	}
	return o.AuthRequest.UnmarshalJSON(data)
}

func (o *AuthRequest) MarshalJSON() ([]byte, error) {
	if o == nil || o.AuthRequest == nil {
		return []byte("null"), nil
	}
	return o.AuthRequest.MarshalJSON()
}
func MakeByProto(p *oidcapiv1.AuthRequest) *AuthRequest {
	return &AuthRequest{AuthRequest: p}
}
func (m *AuthRequest) ToProto() *oidcapiv1.AuthRequest {
	return m.AuthRequest
}

// GetId implements [daobase.Modal].
// Subtle: this method shadows the method (*BanPlayers).GetId of BanPlayersM.BanPlayers.
func (m *AuthRequest) GetId() daobase.IdType {
	return daobase.IdType(m.AuthRequest.Id)
}

// SetId implements [daobase.Modal].
func (m *AuthRequest) SetId(id daobase.IdType) {
	m.AuthRequest.Id = string(id)
}

func NewAuthRequest() AuthRequest {
	return AuthRequest{
		AuthRequest: &oidcapiv1.AuthRequest{
			Id: kits.Ids.NewXId(),
			// Ts: kits.Tss.Ts(),
		},
	}
}

func NewAuthRequestById(id daobase.IdType) AuthRequest {
	return AuthRequest{
		AuthRequest: &oidcapiv1.AuthRequest{
			Id: string(id),
			// Ts: kits.Tss.Ts(),
		},
	}
}

// Key implements [daobase.Modal].
func (m *AuthRequest) Key() []byte {
	return append(m.Prefix(), kits.UnsafeToBytes(string(m.AuthRequest.Id))...)
}

// Prefix implements [daobase.Modal].
func (m *AuthRequest) Prefix() []byte {
	return daokeys.PrefixAuthRequest
}

// Value implements [daobase.Modal].
func (m *AuthRequest) Value() ([]byte, error) {
	return daobase.Value(m)
}

// FromValue implements [daobase.Modal].
func (m *AuthRequest) FromValue(bs []byte) error {
	return daobase.FromValue(m, bs)
}
