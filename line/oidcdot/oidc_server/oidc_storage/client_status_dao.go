package oidc_storage

import (
	"github.com/go-jose/go-jose/v4"
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

type ClientStatusDaoPebble2 struct {
	dao_pebble2.Daobase[ClientStatus, *ClientStatus]
}

type ClientStatusDaoBadger struct {
	dao_badger.Daobase[ClientStatus, *ClientStatus]
}

func NewClientStatusDaoPebble2(db *pebble2dot.Pebble2, logger *dot.LoggerType) *ClientStatusDaoPebble2 {
	return &ClientStatusDaoPebble2{
		Daobase: dao_pebble2.NewDaobase(db, logger, NewClientStatusById),
	}
}

func NewClientStatusDaoBadger(db *badgerdot.BadgerDbDot, logger *dot.LoggerType) *ClientStatusDaoBadger {
	return &ClientStatusDaoBadger{
		Daobase: dao_badger.NewDaobase(db, logger, NewClientStatusById),
	}
}

type ClientStatus oidcapiv1.ClientStatus

// Expire implements [daobase.Modal].
func (m *ClientStatus) Expire() bool {
	return daobase.ModalExpire(m.ExpireTs)
}

// GetExpireTs implements [daobase.Modal].
func (m *ClientStatus) GetExpireTs() uint64 {
	return m.ExpireTs
}

// SetExpireTs implements [daobase.Modal].
func (m *ClientStatus) SetExpireTs(ts uint64) {
	m.ExpireTs = ts
}

var _ daobase.Modal = (*ClientStatus)(nil)

// func (o *ClientStatus) UnmarshalJSON(data []byte) error {
// 	if o.ClientStatus == nil {
// 		o.ClientStatus = &oidcapiv1.ClientStatus{}
// 	}
// 	return o.ClientStatus.UnmarshalJSON(data)
// }

// func (o *ClientStatus) MarshalJSON() ([]byte, error) {
// 	if o == nil || o.ClientStatus == nil {
// 		return []byte("null"), nil
// 	}
// 	return o.ClientStatus.MarshalJSON()
// }

func MakeClientStatusByProto(p *oidcapiv1.ClientStatus) *ClientStatus {
	return (*ClientStatus)(p)
}
func (m *ClientStatus) ClientStatusToProto() *oidcapiv1.ClientStatus {
	return (*oidcapiv1.ClientStatus)(m)
}

// GetId implements [daobase.Modal].
func (m *ClientStatus) GetId() daobase.IdType {
	return daobase.IdType(m.Id)
}

// SetId implements [daobase.Modal].
func (m *ClientStatus) SetId(id daobase.IdType) {
	m.Id = string(id)
}

func NewClientStatus() ClientStatus {
	return ClientStatus{
		Id: kits.Ids.NewXId(),
	}
}

func NewClientStatusById(id daobase.IdType) ClientStatus {
	return ClientStatus{
		Id: string(id),
	}
}

// Key implements [daobase.Modal].
func (m *ClientStatus) Key() []byte {
	return append(m.Prefix(), kits.StringToBytes(string(m.Id))...)
}

// Prefix implements [daobase.Modal].
func (m *ClientStatus) Prefix() []byte {
	return daokeys.PrefixClientStatus
}

// Value implements [daobase.Modal].
func (m *ClientStatus) Value() ([]byte, error) {
	return daobase.Value(m)
}

// FromValue implements [daobase.Modal].
func (m *ClientStatus) FromValue(bs []byte) error {
	return daobase.FromValue(m, bs)
}

type _SignAlgorithmEx int

var SignAlgorithmEx = _SignAlgorithmEx(0)

func (m *_SignAlgorithmEx) ToSignatureAlgorithm(alg oidcapiv1.SignAlgorithm) jose.SignatureAlgorithm {
	switch alg {
	case oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_UNSPECIFIED:
		return ""
	case oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_EDDSA:
		return jose.EdDSA
	case oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_HS256:
		return jose.HS256
	case oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_HS384:
		return jose.HS384
	case oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_HS512:
		return jose.HS512
	case oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_RS256:
		return jose.RS256
	case oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_RS384:
		return jose.RS384
	case oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_RS512:
		return jose.RS512
	case oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_ES256:
		return jose.ES256
	case oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_ES384:
		return jose.ES384
	case oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_ES512:
		return jose.ES512
	case oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_PS256:
		return jose.PS256
	case oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_PS384:
		return jose.PS384
	case oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_PS512:
		return jose.PS512
	default:
		dot.Logger.Error().Msgf("unknown sign algorithm, alg=%d", alg)
		return ""
	}
}

func (m *_SignAlgorithmEx) ToProto(alg jose.SignatureAlgorithm) oidcapiv1.SignAlgorithm {
	switch alg {
	case jose.SignatureAlgorithm(""):
		return oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_UNSPECIFIED
	case jose.EdDSA:
		return oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_EDDSA
	case jose.HS256:
		return oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_HS256
	case jose.HS384:
		return oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_HS384
	case jose.HS512:
		return oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_HS512
	case jose.RS256:
		return oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_RS256
	case jose.RS384:
		return oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_RS384
	case jose.RS512:
		return oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_RS512
	case jose.ES256:
		return oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_ES256
	case jose.ES384:
		return oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_ES384
	case jose.ES512:
		return oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_ES512
	case jose.PS256:
		return oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_PS256
	case jose.PS384:
		return oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_PS384
	case jose.PS512:
		return oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_PS512
	default:
		dot.Logger.Error().Msgf("unknown sign algorithm, alg=%s", alg)
		return oidcapiv1.SignAlgorithm_SIGN_ALGORITHM_UNSPECIFIED
	}
}
