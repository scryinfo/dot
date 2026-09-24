package oidc_storage

import (
	"github.com/scryinfo/dot/lib/kits"
	daobase "github.com/scryinfo/dot/line/db/dao/dao_base"
	daokeys "github.com/scryinfo/dot/line/db/dao/dao_keys"
)

func (p *AppSessionDaoPebble2) Add(m *AppSession) error {
	if m.GetId() == "" {
		m.SetId(daobase.IdType(kits.Ids.NewXId()))
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
	userIdBs := kits.StringToBytes(m.Id)
	if m.WebId == "" {
		m.WebId = kits.Ids.WebSessionID()
		indexKey := append(daokeys.PrefixIndexAppSession, kits.StringToBytes(m.WebId)...)
		if err = txn.Set(indexKey, userIdBs, p.Db.DefaultWriteOpt()); err != nil {
			return err
		}
	}
	return txn.Commit(p.Db.DefaultWriteOpt())
}
func (p *AppSessionDaoPebble2) Remove(m *AppSession) error {
	txn := p.Db.Db().NewBatch()
	defer txn.Close()
	if err := txn.Delete(m.Key(), p.Db.DefaultWriteOpt()); err != nil {
		return err
	}
	if m.WebId != "" {
		indexKey := append(daokeys.PrefixIndexAppSession, kits.StringToBytes(m.WebId)...)
		if err := txn.Delete(indexKey, p.Db.DefaultWriteOpt()); err != nil {
			return err
		}
	}
	return txn.Commit(p.Db.DefaultWriteOpt())
}

func (p *AppSessionDaoPebble2) RemoveBy(id daobase.IdType) error {
	txn := p.Db.Db().NewBatch()
	defer txn.Close()

	m := &AppSession{Id: string(id)}
	v, closer, err := txn.Get(m.Key())
	if err != nil {
		return err
	}
	err = m.FromValue(v)
	closer.Close()
	if err := txn.Delete(m.Key(), p.Db.DefaultWriteOpt()); err != nil {
		return err
	}
	if m.WebId != "" {
		indexKey := append(daokeys.PrefixIndexAppSession, kits.StringToBytes(m.WebId)...)
		if err := txn.Delete(indexKey, p.Db.DefaultWriteOpt()); err != nil {
			return err
		}
	}
	return txn.Commit(p.Db.DefaultWriteOpt())
}

// FindByKeyname finds a user by their keyname (username, email, or phone).
func (p *AppSessionDaoPebble2) FindByKeyname(key string) (*AppSession, error) {
	indexKey := append(daokeys.PrefixIndexAppSession, kits.StringToBytes(key)...)
	bs, closer, err := p.Db.Db().Get(indexKey)
	if err != nil {
		return nil, err
	}
	id := kits.BytesToString(bs)
	closer.Close()
	m, err := p.Find(daobase.IdType(id))
	return &m, err
}
func (p *AppSessionDaoPebble2) FindByWebId(webId string) (*AppSession, error) {
	return p.FindByKeyname(webId)
}
