package dao_pebble2

import (
	"github.com/scryinfo/dot/dot"
	daobase "github.com/scryinfo/dot/line/db/dao/dao_base"
	"github.com/scryinfo/dot/line/db/pebble2dot"
)

type Daobase[T any, PT daobase.ModalPtr[T]] struct {
	Db     *pebble2dot.Pebble2
	MakePT func(id daobase.IdType) T
	Logger *dot.LoggerType
}

type DaoBodybase[T any, PT daobase.ModalBodyPtr[T]] struct {
	Db     *pebble2dot.Pebble2
	MakePT func(id daobase.IdType) T
	Logger *dot.LoggerType
}

func NewDaobase[T any, PT daobase.ModalPtr[T]](db *pebble2dot.Pebble2, logger *dot.LoggerType, makePT func(id daobase.IdType) T) Daobase[T, PT] {
	return Daobase[T, PT]{
		Db:     db,
		MakePT: makePT,
		Logger: logger,
	}
}

func NewDaoBodybase[T any, PT daobase.ModalBodyPtr[T]](db *pebble2dot.Pebble2, logger *dot.LoggerType, makePT func(id daobase.IdType) T) DaoBodybase[T, PT] {
	return DaoBodybase[T, PT]{
		Db:     db,
		MakePT: makePT,
		Logger: logger,
	}
}

// prefix "user:" 的上界是 "user;"）
func prefixUpperBound(prefix []byte) []byte {
	if len(prefix) == 0 {
		return nil
	}
	upper := make([]byte, len(prefix))
	copy(upper, prefix)
	for i := len(upper) - 1; i >= 0; i-- {
		upper[i]++
		if upper[i] != 0 {
			return upper[:i+1]
		}
	}
	return nil
}
