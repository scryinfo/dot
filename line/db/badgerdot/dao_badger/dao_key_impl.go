package dao_badger

import (
	"errors"

	"github.com/dgraph-io/badger/v4"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/db/badgerdot"
	daobase "github.com/scryinfo/dot/line/db/dao/dao_base"
)

var _ daobase.DaoKey[daobase.ModalSample] = (*DaoKeyBase[daobase.ModalSample])(nil)

type DaoKeyBase[T any] struct {
	Db           *badger.DB
	key          []byte
	DefaultValue func() T
	Logger       *dot.LoggerType
}

func NewDaoKeyBase[T any](db *badgerdot.BadgerDbDot, key []byte, logger *dot.LoggerType, defaultValue func() T) *DaoKeyBase[T] {
	return &DaoKeyBase[T]{
		Db:           db.Db(),
		key:          key,
		DefaultValue: defaultValue,
		Logger:       logger,
	}
}

func (d *DaoKeyBase[T]) Set(m *T) error {
	return d.Db.Update(func(txn *badger.Txn) error {
		bs, err := daobase.Value(m)
		if err != nil {
			return err
		}
		return txn.Set(d.key, bs)
	})
}

// if not found, return default T and error is nil
func (d *DaoKeyBase[T]) Get() (*T, error) {
	var m T
	err := d.Db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(d.key)
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return daobase.FromValue(&m, val)
		})
	})
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			err = nil
			if d.DefaultValue != nil {
				m = d.DefaultValue()
			} else {
				d.Logger.Error().Msg("DefaultValue is not set")
			}
		} else {
			d.Logger.Error().AnErr("get error", err).Send()
		}
	}
	return &m, err

}
func (d *DaoKeyBase[T]) Remove() error {
	return d.Db.Update(func(txn *badger.Txn) error {
		return txn.Delete(d.key)
	})
}

func (d *DaoKeyBase[T]) Key() []byte {
	return d.key
}
