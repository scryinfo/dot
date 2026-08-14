package dao_pebble2

import (
	"errors"

	"github.com/dgraph-io/badger/v4"
	"github.com/scryinfo/dot/dot"
	daobase "github.com/scryinfo/dot/line/db/dao/dao_base"
	"github.com/scryinfo/dot/line/db/pebble2dot"
)

var _ daobase.DaoKey[daobase.ModalSample] = (*DaoKeyBase[daobase.ModalSample])(nil)

type DaoKeyBase[T any] struct {
	Db           *pebble2dot.Pebble2
	key          []byte
	DefaultValue func() T
	Logger       *dot.LoggerType
}

func NewDaoKeyBase[T any](db *pebble2dot.Pebble2, key []byte, logger *dot.LoggerType, defaultValue func() T) *DaoKeyBase[T] {
	return &DaoKeyBase[T]{
		Db:           db,
		key:          key,
		DefaultValue: defaultValue,
		Logger:       logger,
	}
}

func (d *DaoKeyBase[T]) Set(m *T) error {
	bs, err := daobase.Value(m)
	if err != nil {
		return err
	}
	return d.Db.Db().Set(d.key, bs, d.Db.DefaultWriteOpt())
}

// if not found, return default T and error is nil
func (d *DaoKeyBase[T]) Get() (*T, error) {
	var m T
	item, closer, err := d.Db.Db().Get(d.key)
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
	} else {
		err = daobase.FromValue(&m, item)
		closer.Close()
	}
	return &m, err
}
func (d *DaoKeyBase[T]) Remove() error {
	return d.Db.Db().Delete(d.key, d.Db.DefaultWriteOpt())
}

func (d *DaoKeyBase[T]) Key() []byte {
	return d.key
}
