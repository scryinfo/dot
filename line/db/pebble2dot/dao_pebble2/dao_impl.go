package dao_pebble2

import (
	"bytes"
	"fmt"

	"github.com/cockroachdb/pebble/v2"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/lib/kits"
	daobase "github.com/scryinfo/dot/line/db/dao/dao_base"
)

var _ daobase.Dao[daobase.ModalSample, *daobase.ModalSample] = (*Daobase[daobase.ModalSample, *daobase.ModalSample])(nil)

func (p *Daobase[T, PT]) Add(m PT) error {
	if m.GetId() == "" {
		m.SetId(daobase.IdType(kits.Ids.NewXId()))
	}

	bs, err := m.Value()
	if err != nil {
		return err
	}
	return p.Db.Db().Set(m.Key(), bs, p.Db.DefaultWriteOpt())
}

func (p *Daobase[T, PT]) Remove(m PT) error {
	return p.Db.Db().Delete(m.Key(), p.Db.DefaultWriteOpt())
}

func (p *Daobase[T, PT]) RemoveBy(id daobase.IdType) error {
	return p.Db.Db().Delete(daobase.ModalKey[T, PT](id), p.Db.DefaultWriteOpt())
}

func (p *Daobase[T, PT]) Find(id daobase.IdType) (T, error) {
	m := p.MakePT(id)

	var pt PT = &m
	v, closer, err := p.Db.Db().Get(pt.Key())
	if err != nil {
		return m, err
	}
	err = pt.FromValue(v)
	closer.Close()
	if err == nil && pt.Expire() {
		err = p.Remove(pt)
		if err != nil {
			dot.Logger.Error().AnErr("remove expired modal", err).Send()
			return m, err
		} else {
			return m, pebble.ErrNotFound
		}
	}
	return m, err
}

func (p *Daobase[T, PT]) NextPage(page *daobase.PageMeta) ([]T, error) {
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
					p.Logger.Error().AnErr("remove expired modal", err2).Send()
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
	return ms, err
}

func (p *Daobase[T, PT]) PrevPage(page *daobase.PageMeta) ([]T, error) {
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
				iterator.Next()
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
				} else {
					return nil, err
				}
			} else if pt.Expire() {
				err2 := p.Remove(pt)
				if err2 != nil {
					p.Logger.Error().AnErr("remove expired modal", err2).Send()
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
