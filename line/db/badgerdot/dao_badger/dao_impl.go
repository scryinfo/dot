package dao_badger

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"
	"github.com/scryinfo/dot/lib/kits"
	daobase "github.com/scryinfo/dot/line/db/dao/dao_base"
)

var _ daobase.Dao[daobase.ModalSample, *daobase.ModalSample] = (*Daobase[daobase.ModalSample, *daobase.ModalSample])(nil)

func (p *Daobase[T, PT]) Add(m PT) error {
	return p.Db.Update(func(txn *badger.Txn) error {
		if m.GetId() == "" {
			m.SetId(daobase.IdType(kits.Ids.NewXId()))
		}

		bs, err := m.Value()
		if err != nil {
			return err
		}
		return txn.Set(m.Key(), bs)
	})
}

func (p *Daobase[T, PT]) Remove(m PT) error {
	return p.Db.Update(func(txn *badger.Txn) error {
		return txn.Delete(m.Key())
	})
}

func (p *Daobase[T, PT]) RemoveBy(id daobase.IdType) error {
	return p.Db.Update(func(txn *badger.Txn) error {
		return txn.Delete(daobase.ModalKey[T, PT](id))
	})
}

func (p *Daobase[T, PT]) Find(id daobase.IdType) (T, error) {
	m := p.MakePT(id)
	err := p.Db.View(func(txn *badger.Txn) error {
		var pt PT = &m
		v, err2 := txn.Get(pt.Key())
		if err2 == nil {
			err2 = v.Value(func(val []byte) error {
				return pt.FromValue(val)
			})
		}
		return err2
	})
	return m, err
}

func (p *Daobase[T, PT]) NextPage(page *daobase.PageMeta) ([]T, error) {
	if page == nil {
		return nil, fmt.Errorf("the page parameter is nil")
	}
	if page.PageSize < 1 {
		page.PageSize = daobase.DefaultPageSize
	}
	ms := make([]T, 0, page.PageSize)
	err := p.Db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		opts.Prefix = daobase.ModalPrefix[T, PT]()
		iterator := txn.NewIterator(opts)
		defer iterator.Close()
		if len(page.LastId) > 0 {
			iterator.Seek(daobase.ModalKey[T, PT](page.LastId))
			if iterator.ValidForPrefix(opts.Prefix) {
				iterator.Next()
			}
		} else {
			iterator.Rewind()
		}

		var count int32
		for ; iterator.ValidForPrefix(opts.Prefix); iterator.Next() {
			item := iterator.Item()
			m := p.MakePT("")
			var pt PT = &m
			err2 := item.Value(func(val []byte) error {
				return pt.FromValue(val)
			})
			if err2 != nil {
				if len(ms) > 0 {
					p.Logger.Info().Err(err2).Send()
					return nil
				}
				return err2
			}
			if len(pt.GetId()) < 1 {
				continue
			}
			ms = append(ms, m)
			count++
			if count >= page.PageSize {
				break
			}
		}
		return nil
	})
	return ms, err
}

func (p *Daobase[T, PT]) PrevPage(page *daobase.PageMeta) ([]T, error) {
	if page == nil {
		return nil, fmt.Errorf("the page parameter is nil")
	}
	if page.PageSize < 1 {
		page.PageSize = daobase.DefaultPageSize
	}
	ms := make([]T, 0, page.PageSize)
	err := p.Db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		opts.Prefix = daobase.ModalPrefix[T, PT]()
		opts.Reverse = true
		iterator := txn.NewIterator(opts)
		defer iterator.Close()
		if len(page.FirstId) > 0 {
			iterator.Seek(daobase.ModalKey[T, PT](page.FirstId))
			if iterator.ValidForPrefix(opts.Prefix) {
				iterator.Next()
			}
		} else {
			max := make([]byte, 0, len(opts.Prefix)+1)
			max = append(max, opts.Prefix...)
			max = append(max, 0xff)
			iterator.Seek(max)
		}

		var count int32
		for ; iterator.ValidForPrefix(opts.Prefix); iterator.Next() {
			item := iterator.Item()
			m := p.MakePT("")
			var pt PT = &m
			err2 := item.Value(func(val []byte) error {
				return pt.FromValue(val)
			})
			if err2 != nil {
				if len(ms) > 0 {
					p.Logger.Info().Err(err2).Send()
					return nil
				}
				return err2
			}
			if len(pt.GetId()) < 1 {
				continue
			}
			ms = append(ms, m)
			count++
			if count >= page.PageSize {
				break
			}
		}
		return nil
	})
	return ms, err
}
