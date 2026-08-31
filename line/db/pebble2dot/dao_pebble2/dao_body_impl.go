package dao_pebble2

import (
	"bytes"
	"fmt"

	"github.com/cockroachdb/pebble/v2"
	daobase "github.com/scryinfo/dot/line/db/dao/dao_base"
)

var _ daobase.DaoBody[daobase.ModalBodySample, *daobase.ModalBodySample] = (*DaoBodybase[daobase.ModalBodySample, *daobase.ModalBodySample])(nil)
var _ daobase.Dao[daobase.ModalBodySample, *daobase.ModalBodySample] = (*DaoBodybase[daobase.ModalBodySample, *daobase.ModalBodySample])(nil)

func (p *DaoBodybase[T, PT]) Add(m PT) error {
	txn := p.Db.Db().NewBatch()
	defer txn.Close()
	bs, err := m.Value()
	if err != nil {
		return err
	}
	if err = txn.Set(m.Key(), bs, p.Db.DefaultWriteOpt()); err != nil {
		return err
	}
	bs, err = m.ValueBody()
	if err != nil {
		return err
	}
	err = txn.Set(m.KeyBody(), bs, p.Db.DefaultWriteOpt())
	if err != nil {
		return err
	}
	return txn.Commit(p.Db.DefaultWriteOpt())
}

func (p *DaoBodybase[T, PT]) Remove(m PT) error {
	txn := p.Db.Db().NewBatch()
	defer txn.Close()
	if err := txn.Delete(m.Key(), p.Db.DefaultWriteOpt()); err != nil {
		return err
	}
	if err := txn.Delete(m.KeyBody(), p.Db.DefaultWriteOpt()); err != nil {
		return err
	}
	return txn.Commit(p.Db.DefaultWriteOpt())
}

func (p *DaoBodybase[T, PT]) RemoveBy(id daobase.IdType) error {
	txn := p.Db.Db().NewBatch()
	defer txn.Close()
	if err := txn.Delete(daobase.ModalKey[T, PT](id), p.Db.DefaultWriteOpt()); err != nil {
		return err
	}
	if err := txn.Delete(daobase.ModalKeyBody[T, PT](id), p.Db.DefaultWriteOpt()); err != nil {
		return err
	}
	return txn.Commit(p.Db.DefaultWriteOpt())
}

func (p *DaoBodybase[T, PT]) Find(id daobase.IdType) (T, error) {
	m := p.MakePT(id)
	txn := p.Db.Db().NewSnapshot()
	defer txn.Close()

	var pt PT = &m
	v, closer, err := txn.Get(pt.Key())
	if err != nil {
		return m, err
	}
	err = pt.FromValue(v)
	closer.Close()
	if err != nil {
		return m, err
	} else if pt.Expire() {
		err = p.Remove(pt)
		if err != nil {
			p.Logger.Error().AnErr("remove expired model", err).Send()
			return m, err
		} else {
			return m, pebble.ErrNotFound
		}
	}

	v, closer, err = txn.Get(pt.KeyBody())
	if err != nil {
		return m, err
	}
	err = pt.FromValueBody(v)
	closer.Close()

	return m, err
}

func (p *DaoBodybase[T, PT]) NextPage(page *daobase.PageMeta) ([]T, error) {
	if page == nil {
		return nil, fmt.Errorf("the page parameter is nil")
	}
	if page.PageSize < 1 {
		page.PageSize = daobase.DefaultPageSize
	}
	txn := p.Db.Db().NewSnapshot()
	defer txn.Close()

	prefix := daobase.ModalPrefix[T, PT]()
	upperBound := prefixUpperBound(prefix)
	iterator, err := txn.NewIter(&pebble.IterOptions{
		LowerBound: prefix,
		UpperBound: upperBound,
	})
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	ms := make([]T, 0, page.PageSize)
	{
		if len(page.LastId) > 0 {
			if iterator.SeekGE(daobase.ModalKey[T, PT](page.LastId)) {
				if bytes.Equal(iterator.Key(), []byte(page.LastId)) {
					iterator.Next()
				}
			}
		} else {
			iterator.First()
		}

		var count int32
		for ; iterator.Valid(); iterator.Next() {
			m := p.MakePT("")
			var pt PT = &m
			err := pt.FromValue(iterator.Value())
			if err != nil {
				if len(ms) > 0 {
					p.Logger.Info().Err(err).Send()
					return ms, nil
				} else {
					return nil, err
				}
			} else if pt.Expire() {
				err2 := p.Remove(pt)
				if err2 != nil {
					p.Logger.Error().AnErr("remove expired model", err2).Send()
				}
				continue
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
	}
	return ms, nil
}

func (p *DaoBodybase[T, PT]) PrevPage(page *daobase.PageMeta) ([]T, error) {
	if page == nil {
		return nil, fmt.Errorf("the page parameter is nil")
	} else if page.FirstId == "" {
		return nil, fmt.Errorf("the firstId parameter is empty")
	}

	if page.PageSize < 1 {
		page.PageSize = daobase.DefaultPageSize
	}
	txn := p.Db.Db().NewSnapshot()
	defer txn.Close()

	prefix := daobase.ModalPrefix[T, PT]()
	upperBound := prefixUpperBound(prefix)
	iterator, err := txn.NewIter(&pebble.IterOptions{
		LowerBound: prefix,
		UpperBound: upperBound,
	})
	if err != nil {
		return nil, err
	}
	defer iterator.Close()
	ms := make([]T, 0, page.PageSize)
	{
		if len(page.FirstId) > 0 {
			iterator.SeekLT(daobase.ModalKey[T, PT](page.FirstId))
			if bytes.Equal(iterator.Key(), []byte(page.FirstId)) {
				iterator.Prev()
			}
		} else {
			// no the case
		}

		var count int32
		for ; iterator.Valid(); iterator.Next() {
			m := p.MakePT("")
			var pt PT = &m
			err := pt.FromValue(iterator.Value())
			if err != nil {
				if len(ms) > 0 {
					p.Logger.Info().Err(err).Send()
					return ms, nil
				}
				return nil, err
			} else if pt.Expire() {
				err2 := p.Remove(pt)
				if err2 != nil {
					p.Logger.Error().AnErr("remove expired model", err2).Send()
				}
				continue
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
	}
	return ms, nil
}
