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

type CodeAuthRequestDao struct {
	dao_pebble2.Daobase[CodeAuthRequest, *CodeAuthRequest]
}

func NewCodeAuthRequestDao(db *pebble2dot.Pebble2, logger *dot.LoggerType) *CodeAuthRequestDao {
	return &CodeAuthRequestDao{
		Daobase: dao_pebble2.NewDaobase(db, logger, NewCodeAuthRequestById),
	}
}

type CodeAuthRequest oidcapiv1.CodeAuthRequest

var _ daobase.Modal = (*CodeAuthRequest)(nil)

// func (o *CodeAuthRequest) UnmarshalJSON(data []byte) error {
// 	if o.CodeAuthRequest == nil {
// 		o.CodeAuthRequest = &oidcapiv1.CodeAuthRequest{}
// 	}
// 	return o.CodeAuthRequest.UnmarshalJSON(data)
// }

// func (o *CodeAuthRequest) MarshalJSON() ([]byte, error) {
// 	if o == nil || o.CodeAuthRequest == nil {
// 		return []byte("null"), nil
// 	}
// 	return o.CodeAuthRequest.MarshalJSON()
// }

func MakeCodeAuthRequestByProto(p *oidcapiv1.CodeAuthRequest) *CodeAuthRequest {
	return (*CodeAuthRequest)(p)
}
func (m *CodeAuthRequest) CodeAuthRequestToProto() *oidcapiv1.CodeAuthRequest {
	return (*oidcapiv1.CodeAuthRequest)(m)
}

// GetId implements [daobase.Modal].
// Subtle: this method shadows the method (*BanPlayers).GetId of BanPlayersM.BanPlayers.
func (m *CodeAuthRequest) GetId() daobase.IdType {
	return daobase.IdType(m.Code)
}

// SetId implements [daobase.Modal].
func (m *CodeAuthRequest) SetId(id daobase.IdType) {
	m.Code = string(id)
}

func NewCodeAuthRequest() CodeAuthRequest {
	return CodeAuthRequest{
		Code: kits.Ids.NewAuthCode(),
	}
}

func NewCodeAuthRequestById(id daobase.IdType) CodeAuthRequest {
	return CodeAuthRequest{
		Code: string(id),
	}
}

// Key implements [daobase.Modal].
func (m *CodeAuthRequest) Key() []byte {
	return append(m.Prefix(), kits.StringToBytes(string(m.Code))...)
}

// Prefix implements [daobase.Modal].
func (m *CodeAuthRequest) Prefix() []byte {
	return daokeys.PrefixAuthCode
}

// Value implements [daobase.Modal].
func (m *CodeAuthRequest) Value() ([]byte, error) {
	return daobase.Value(m)
}

// FromValue implements [daobase.Modal].
func (m *CodeAuthRequest) FromValue(bs []byte) error {
	return daobase.FromValue(m, bs)
}
