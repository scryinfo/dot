package oidc_storage

import (
	"github.com/scryinfo/dot/lib/kits"
	daobase "github.com/scryinfo/dot/line/db/dao/dao_base"
	daokeys "github.com/scryinfo/dot/line/db/dao/dao_keys"
)

func (p *UserDaoPebble2) Add(m *User) error {
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
	if m.Email != "" {
		indexKey := append(daokeys.PrefixIndexUser, kits.StringToBytes(m.Email)...)
		if err = txn.Set(indexKey, userIdBs, p.Db.DefaultWriteOpt()); err != nil {
			return err
		}
	}
	if m.Phone != "" {
		indexKey := append(daokeys.PrefixIndexUser, kits.StringToBytes(m.Phone)...)
		if err = txn.Set(indexKey, userIdBs, p.Db.DefaultWriteOpt()); err != nil {
			return err
		}
	}
	if m.Username != "" {
		indexKey := append(daokeys.PrefixIndexUser, kits.StringToBytes(m.Username)...)
		if err = txn.Set(indexKey, userIdBs, p.Db.DefaultWriteOpt()); err != nil {
			return err
		}
	}
	return txn.Commit(p.Db.DefaultWriteOpt())
}
func (p *UserDaoPebble2) Remove(m *User) error {
	txn := p.Db.Db().NewBatch()
	defer txn.Close()
	if err := txn.Delete(m.Key(), p.Db.DefaultWriteOpt()); err != nil {
		return err
	}
	if m.Email != "" {
		indexKey := append(daokeys.PrefixIndexUser, kits.StringToBytes(m.Email)...)
		if err := txn.Delete(indexKey, p.Db.DefaultWriteOpt()); err != nil {
			return err
		}
	}
	if m.Phone != "" {
		indexKey := append(daokeys.PrefixIndexUser, kits.StringToBytes(m.Phone)...)
		if err := txn.Delete(indexKey, p.Db.DefaultWriteOpt()); err != nil {
			return err
		}
	}
	if m.Username != "" {
		indexKey := append(daokeys.PrefixIndexUser, kits.StringToBytes(m.Username)...)
		if err := txn.Delete(indexKey, p.Db.DefaultWriteOpt()); err != nil {
			return err
		}
	}
	return txn.Commit(p.Db.DefaultWriteOpt())
}

func (p *UserDaoPebble2) RemoveBy(id daobase.IdType) error {
	txn := p.Db.Db().NewBatch()
	defer txn.Close()

	m := &User{Id: string(id)}
	v, closer, err := txn.Get(m.Key())
	if err != nil {
		return err
	}
	err = m.FromValue(v)
	closer.Close()
	if err := txn.Delete(m.Key(), p.Db.DefaultWriteOpt()); err != nil {
		return err
	}
	if m.Email != "" {
		indexKey := append(daokeys.PrefixIndexUser, kits.StringToBytes(m.Email)...)
		if err := txn.Delete(indexKey, p.Db.DefaultWriteOpt()); err != nil {
			return err
		}
	}
	if m.Phone != "" {
		indexKey := append(daokeys.PrefixIndexUser, kits.StringToBytes(m.Phone)...)
		if err := txn.Delete(indexKey, p.Db.DefaultWriteOpt()); err != nil {
			return err
		}
	}
	if m.Username != "" {
		indexKey := append(daokeys.PrefixIndexUser, kits.StringToBytes(m.Username)...)
		if err := txn.Delete(indexKey, p.Db.DefaultWriteOpt()); err != nil {
			return err
		}
	}
	return txn.Commit(p.Db.DefaultWriteOpt())
}

// FindByKeyname finds a user by their keyname (username, email, or phone).
func (p *UserDaoPebble2) FindByKeyname(key string) (*User, error) {
	indexKey := append(daokeys.PrefixIndexUser, kits.StringToBytes(key)...)
	bs, closer, err := p.Db.Db().Get(indexKey)
	if err != nil {
		return nil, err
	}
	userId := kits.BytesToString(bs)
	closer.Close()
	m, err := p.Find(daobase.IdType(userId))
	return &m, err
}
func (p *UserDaoPebble2) FindByUsername(username string) (*User, error) {
	return p.FindByKeyname(username)
}
func (p *UserDaoPebble2) FindByEmail(email string) (*User, error) {
	return p.FindByKeyname(email)
}
func (p *UserDaoPebble2) FindByPhone(phone string) (*User, error) {
	return p.FindByKeyname(phone)
}
