package oidc_storage

import (
	"github.com/scryinfo/dot/lib/kits"
	daobase "github.com/scryinfo/dot/line/db/dao/dao_base"
	daokeys "github.com/scryinfo/dot/line/db/dao/dao_keys"
)

func (p *CodeAuthRequestDaoPebble2) Add(m *CodeAuthRequest) error {
	if m.GetId() == "" {
		m.SetId(daobase.IdType(kits.Ids.NewAuthCode()))
	}

	txn := p.Db.Db().NewBatch()
	defer txn.Close()
	bs, err := m.Value()
	if err != nil {
		return err
	}
	if err = txn.Set(m.Key(), bs, p.Db.DefaultWriteOpt()); err != nil {
		return err
	}
	userIdBs := kits.StringToBytes(m.Code)
	if m.AuthRequestId != "" {
		indexKey := append(daokeys.PrefixIndexAuthCode, kits.StringToBytes(m.AuthRequestId)...)
		if err = txn.Set(indexKey, userIdBs, p.Db.DefaultWriteOpt()); err != nil {
			return err
		}
	}
	return txn.Commit(p.Db.DefaultWriteOpt())
}
func (p *CodeAuthRequestDaoPebble2) Remove(m *CodeAuthRequest) error {
	txn := p.Db.Db().NewBatch()
	defer txn.Close()
	if err := txn.Delete(m.Key(), p.Db.DefaultWriteOpt()); err != nil {
		return err
	}
	if m.AuthRequestId != "" {
		indexKey := append(daokeys.PrefixIndexAuthCode, kits.StringToBytes(m.AuthRequestId)...)
		if err := txn.Delete(indexKey, p.Db.DefaultWriteOpt()); err != nil {
			return err
		}
	}
	return txn.Commit(p.Db.DefaultWriteOpt())
}

func (p *CodeAuthRequestDaoPebble2) RemoveBy(id daobase.IdType) error {
	txn := p.Db.Db().NewBatch()
	defer txn.Close()

	m := &CodeAuthRequest{Code: string(id)}
	v, closer, err := txn.Get(m.Key())
	if err != nil {
		return err
	}
	err = m.FromValue(v)
	closer.Close()
	if err := txn.Delete(m.Key(), p.Db.DefaultWriteOpt()); err != nil {
		return err
	}
	if m.AuthRequestId != "" {
		indexKey := append(daokeys.PrefixIndexAuthCode, kits.StringToBytes(m.AuthRequestId)...)
		if err := txn.Delete(indexKey, p.Db.DefaultWriteOpt()); err != nil {
			return err
		}
	}
	return txn.Commit(p.Db.DefaultWriteOpt())
}

func (p *CodeAuthRequestDaoPebble2) FindByAuthRequestId(authRequestId string) (*CodeAuthRequest, error) {
	indexKey := append(daokeys.PrefixIndexAuthCode, kits.StringToBytes(authRequestId)...)
	bs, closer, err := p.Db.Db().Get(indexKey)
	if err != nil {
		return nil, err
	}
	userId := kits.BytesToString(bs)
	closer.Close()
	m, err := p.Find(daobase.IdType(userId))
	return &m, err
}
