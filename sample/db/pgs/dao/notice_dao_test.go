package dao

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/scryinfo/dot/dots/db/pgs"
	"github.com/scryinfo/dot/sample/db/pgs/model"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

const dbconf_test = `{
			"addr": "localhost:5432",
			"user": "postgres",
			"password": "scry",
			"database": "postgres",
			"showSql": true
		}`

func TestWithTx(t *testing.T) {
	noticeDao := &NoticeDao{}
	noticeDao.DaoBase = pgs.GenerateDaoBase(dbconf_test)
	err := Create(noticeDao.DaoBase.Wrapper.GetDb())
	assert.Equal(t, nil, err)
	//create table
	//Create(noticeDao.DaoBase.Wrapper.GetDb())
	noticeDao.DaoBase.WithTx(func(conn orm.DB) error {
		notice := &model.Notice{
			Status: 1,
		}
		mInsertReturn, err := noticeDao.InsertReturn(conn, notice)
		assert.Equal(t, nil, err)
		assert.NotEqual(t, nil, mInsertReturn.ID)
		assert.Equal(t, 1, mInsertReturn.Status)

		mGet, err := noticeDao.GetByID(conn, mInsertReturn.ID)
		assert.Equal(t, nil, err)
		assert.Equal(t, mInsertReturn.ID, mGet.ID)
		assert.Equal(t, mInsertReturn.UpdateTime, mGet.UpdateTime)
		assert.Equal(t, mInsertReturn.CreateTime, mGet.CreateTime)
		assert.Equal(t, mInsertReturn.OptimisticLockVersion, mGet.OptimisticLockVersion)
		assert.Equal(t, mInsertReturn.Status, mGet.Status)

		ms, err := noticeDao.Query(conn, "")
		assert.Equal(t, nil, err)
		assert.Equal(t, 1, len(ms))
		mGet = ms[0]
		assert.Equal(t, mInsertReturn.ID, mGet.ID)
		assert.Equal(t, mInsertReturn.UpdateTime, mGet.UpdateTime)
		assert.Equal(t, mInsertReturn.CreateTime, mGet.CreateTime)
		assert.Equal(t, mInsertReturn.OptimisticLockVersion, mGet.OptimisticLockVersion)
		assert.Equal(t, mInsertReturn.Status, mGet.Status)

		mGet, err = noticeDao.QueryOne(conn, "")
		assert.Equal(t, nil, err)
		assert.Equal(t, mInsertReturn.ID, mGet.ID)
		assert.Equal(t, mInsertReturn.UpdateTime, mGet.UpdateTime)
		assert.Equal(t, mInsertReturn.CreateTime, mGet.CreateTime)
		assert.Equal(t, mInsertReturn.OptimisticLockVersion, mGet.OptimisticLockVersion)
		assert.Equal(t, mInsertReturn.Status, mGet.Status)

		ms, err = noticeDao.List(conn)
		assert.Equal(t, nil, err)
		assert.Equal(t, 1, len(ms))
		mGet = ms[0]
		assert.Equal(t, mInsertReturn.ID, mGet.ID)
		assert.Equal(t, mInsertReturn.UpdateTime, mGet.UpdateTime)
		assert.Equal(t, mInsertReturn.CreateTime, mGet.CreateTime)
		assert.Equal(t, mInsertReturn.OptimisticLockVersion, mGet.OptimisticLockVersion)
		assert.Equal(t, mInsertReturn.Status, mGet.Status)

		ms, err = noticeDao.QueryPage(conn, 10, 1, "")
		assert.Equal(t, nil, err)
		assert.Equal(t, 1, len(ms))
		mGet = ms[0]
		assert.Equal(t, mInsertReturn.ID, mGet.ID)
		assert.Equal(t, mInsertReturn.UpdateTime, mGet.UpdateTime)
		assert.Equal(t, mInsertReturn.CreateTime, mGet.CreateTime)
		assert.Equal(t, mInsertReturn.OptimisticLockVersion, mGet.OptimisticLockVersion)
		assert.Equal(t, mInsertReturn.Status, mGet.Status)

		ms, count, err := noticeDao.QueryPageWithCount(conn, 10, 1, "")
		assert.Equal(t, nil, err)
		assert.Equal(t, 1, len(ms))
		assert.Equal(t, 1, count)
		mGet = ms[0]
		assert.Equal(t, mInsertReturn.ID, mGet.ID)
		assert.Equal(t, mInsertReturn.UpdateTime, mGet.UpdateTime)
		assert.Equal(t, mInsertReturn.CreateTime, mGet.CreateTime)
		assert.Equal(t, mInsertReturn.OptimisticLockVersion, mGet.OptimisticLockVersion)
		assert.Equal(t, mInsertReturn.Status, mGet.Status)

		mGet.Status = 10
		time.Sleep(time.Second * 1) //
		err = noticeDao.Update(conn, mGet)
		assert.Equal(t, nil, err)
		mGet, err = noticeDao.QueryOne(conn, fmt.Sprintf(" %s = ?", model.Notice_ID), mGet.ID)
		assert.Equal(t, nil, err)
		assert.Equal(t, 10, mGet.Status)
		assert.Equal(t, mInsertReturn.CreateTime, mGet.CreateTime)
		assert.NotEqual(t, mInsertReturn.UpdateTime, mGet.UpdateTime)

		mGet.Status = 11
		time.Sleep(time.Second * 1) //
		mGet, err = noticeDao.UpdateReturn(conn, mGet)
		assert.Equal(t, nil, err)
		mGet, err = noticeDao.GetByID(conn, mGet.ID)
		assert.Equal(t, nil, err)
		assert.Equal(t, 11, mGet.Status)
		assert.Equal(t, mInsertReturn.CreateTime, mGet.CreateTime)
		assert.NotEqual(t, mInsertReturn.UpdateTime, mGet.UpdateTime)

		//err =  noticeDao.Insert(conn, mInsertReturn)
		//assert.NotEqual(t, nil, err)

		return err
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

	err := db.Model((*model.Notice)(nil)).DropTable(&orm.DropTableOptions{IfExists: true, Cascade: true})
	if err != nil {
		log.Print(err)
	}
	err = db.Model((*model.Data)(nil)).DropComposite(&orm.DropCompositeOptions{
		IfExists: true,
	})
	if err != nil {
		log.Print(err)
	}
	err = db.Model((*model.Data)(nil)).CreateComposite(nil)
	if err != nil {
		log.Print(err)
	}
	err = db.Model((*model.Notice)(nil)).CreateTable(nil)
	if err != nil {
		log.Print(err)
	}
	return err
}
