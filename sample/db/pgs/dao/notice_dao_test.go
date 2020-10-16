package dao

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/scryinfo/dot/dots/db/pgs"
	"github.com/scryinfo/dot/sample/db/pgs/model"
	"testing"
)

const dbconf_test = `{
			"addr": "localhost:5432",
			"user": "postgres",
			"password": "postgres",
			"database": "postgres",
			"showSql": true
		}`

func TestWithTx(t *testing.T) {
	m := &NoticeDao{}
	m.DaoBase = pgs.GenerateDaoBase(dbconf_test)
	//create table
	//Create(m.DaoBase.Wrapper.GetDb())
	m.DaoBase.WithTx(func(conn orm.DB) error {
		notice := &model.Notice{
			Status: 1,
		}
		mnew, err := m.InsertReturn(conn, notice)
		if err != nil {
			return err
		}
		return m.Insert(conn, mnew)
	})
}

func TestWithNoTx(t *testing.T) {
	m := &NoticeDao{}
	m.DaoBase = pgs.GenerateDaoBase(dbconf_test)
	m.DaoBase.WithNoTx(func(conn orm.DB) error {
		notice := &model.Notice{
			Status: 1,
		}
		mnew, err := m.InsertReturn(conn, notice)
		if err != nil {
			return err
		}
		return m.Insert(conn, mnew)
	})
}

func Create(db *pg.DB) error {

	return db.CreateTable(&model.Notice{}, &orm.CreateTableOptions{
		IfNotExists:   true,
		FKConstraints: true,
	})
}
