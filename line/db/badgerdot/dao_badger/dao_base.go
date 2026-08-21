package dao_badger

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/db/badgerdot"
	daobase "github.com/scryinfo/dot/line/db/dao/dao_base"
)

type Daobase[T any, PT daobase.ModalPtr[T]] struct {
	Db     *badger.DB
	MakePT func(id daobase.IdType) T
	Logger *dot.LoggerType
}

type DaoBodybase[T any, PT daobase.ModalBodyPtr[T]] struct {
	Db     *badger.DB
	MakePT func(id daobase.IdType) T
	Logger *dot.LoggerType
}

func NewDaobase[T any, PT daobase.ModalPtr[T]](db *badgerdot.BadgerDbDot, logger *dot.LoggerType, makePT func(id daobase.IdType) T) Daobase[T, PT] {
	return Daobase[T, PT]{
		Db:     db.Db(),
		MakePT: makePT,
		Logger: logger,
	}
}
func NewPointDaobase[T any, PT daobase.ModalPtr[T]](db *badgerdot.BadgerDbDot, logger *dot.LoggerType, makePT func(id daobase.IdType) T) *Daobase[T, PT] {
	return &Daobase[T, PT]{
		Db:     db.Db(),
		MakePT: makePT,
		Logger: logger,
	}
}

func NewDaoBodybase[T any, PT daobase.ModalBodyPtr[T]](db *badgerdot.BadgerDbDot, logger *dot.LoggerType, makePT func(id daobase.IdType) T) DaoBodybase[T, PT] {
	return DaoBodybase[T, PT]{
		Db:     db.Db(),
		MakePT: makePT,
		Logger: logger,
	}
}

func NewPointDaoBodybase[T any, PT daobase.ModalBodyPtr[T]](db *badgerdot.BadgerDbDot, logger *dot.LoggerType, makePT func(id daobase.IdType) T) *DaoBodybase[T, PT] {
	return &DaoBodybase[T, PT]{
		Db:     db.Db(),
		MakePT: makePT,
		Logger: logger,
	}
}
