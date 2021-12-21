package pgs

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/scryinfo/dot/dots/db/pgs"
	"github.com/scryinfo/dot/sample/db/tools/model"
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

func TestNoticeDao(t *testing.T) {
	noticeDao := &NoticeDao{}
	noticeDao.DaoBase = pgs.GenerateDaoBase(dbconf_test)
	err := createNotice(noticeDao.DaoBase.Wrapper.GetDb())
	assert.Equal(t, nil, err)
	//create table
	//createNotice(noticeDao.DaoBase.Wrapper.GetDb())
	noticeDao.DaoBase.WithNoTx(func(conn orm.DB) error {
		notice := &model.Notice{
			Status: 1,
		}
		mReturn, err := noticeDao.InsertReturn(conn, notice)
		assert.Equal(t, nil, err)
		assert.NotEqual(t, nil, mReturn.ID)
		assert.Equal(t, 1, mReturn.Status)

		mGet, err := noticeDao.GetByID(conn, mReturn.ID)
		assert.Equal(t, nil, err)
		assert.Equal(t, mReturn.ID, mGet.ID)
		assert.Equal(t, mReturn.UpdateTime, mGet.UpdateTime)
		assert.Equal(t, mReturn.CreateTime, mGet.CreateTime)
		assert.Equal(t, mReturn.OptimisticLockVersion, mGet.OptimisticLockVersion)
		assert.Equal(t, mReturn.Status, mGet.Status)

		ms, err := noticeDao.Query(conn, "")
		assert.Equal(t, nil, err)
		assert.Equal(t, 1, len(ms))
		mGet = ms[0]
		assert.Equal(t, mReturn.ID, mGet.ID)
		assert.Equal(t, mReturn.UpdateTime, mGet.UpdateTime)
		assert.Equal(t, mReturn.CreateTime, mGet.CreateTime)
		assert.Equal(t, mReturn.OptimisticLockVersion, mGet.OptimisticLockVersion)
		assert.Equal(t, mReturn.Status, mGet.Status)

		mGet, err = noticeDao.QueryOne(conn, "")
		assert.Equal(t, nil, err)
		assert.Equal(t, mReturn.ID, mGet.ID)
		assert.Equal(t, mReturn.UpdateTime, mGet.UpdateTime)
		assert.Equal(t, mReturn.CreateTime, mGet.CreateTime)
		assert.Equal(t, mReturn.OptimisticLockVersion, mGet.OptimisticLockVersion)
		assert.Equal(t, mReturn.Status, mGet.Status)

		ms, err = noticeDao.List(conn)
		assert.Equal(t, nil, err)
		assert.Equal(t, 1, len(ms))
		mGet = ms[0]
		assert.Equal(t, mReturn.ID, mGet.ID)
		assert.Equal(t, mReturn.UpdateTime, mGet.UpdateTime)
		assert.Equal(t, mReturn.CreateTime, mGet.CreateTime)
		assert.Equal(t, mReturn.OptimisticLockVersion, mGet.OptimisticLockVersion)
		assert.Equal(t, mReturn.Status, mGet.Status)

		ms, err = noticeDao.QueryPage(conn, 10, 1, "")
		assert.Equal(t, nil, err)
		assert.Equal(t, 1, len(ms))
		mGet = ms[0]
		assert.Equal(t, mReturn.ID, mGet.ID)
		assert.Equal(t, mReturn.UpdateTime, mGet.UpdateTime)
		assert.Equal(t, mReturn.CreateTime, mGet.CreateTime)
		assert.Equal(t, mReturn.OptimisticLockVersion, mGet.OptimisticLockVersion)
		assert.Equal(t, mReturn.Status, mGet.Status)

		ms, count, err := noticeDao.QueryPageWithCount(conn, 10, 1, "")
		assert.Equal(t, nil, err)
		assert.Equal(t, 1, len(ms))
		assert.Equal(t, 1, count)
		mGet = ms[0]
		assert.Equal(t, mReturn.ID, mGet.ID)
		assert.Equal(t, mReturn.UpdateTime, mGet.UpdateTime)
		assert.Equal(t, mReturn.CreateTime, mGet.CreateTime)
		assert.Equal(t, mReturn.OptimisticLockVersion, mGet.OptimisticLockVersion)
		assert.Equal(t, mReturn.Status, mGet.Status)

		mGet.Status = 10
		time.Sleep(time.Second * 1) //
		err = noticeDao.Update(conn, mGet)
		assert.Equal(t, nil, err)
		mGet, err = noticeDao.QueryOne(conn, fmt.Sprintf(" %s = ?", model.Notice_ID), mGet.ID)
		assert.Equal(t, nil, err)
		assert.Equal(t, 10, mGet.Status)
		assert.Equal(t, mReturn.CreateTime, mGet.CreateTime)
		assert.NotEqual(t, mReturn.UpdateTime, mGet.UpdateTime)

		mGet.Status = 11
		time.Sleep(time.Second * 1) //
		mGet, err = noticeDao.UpdateReturn(conn, mGet)
		assert.Equal(t, nil, err)
		mGet, err = noticeDao.GetByID(conn, mGet.ID)
		assert.Equal(t, nil, err)
		assert.Equal(t, 11, mGet.Status)
		assert.Equal(t, mReturn.CreateTime, mGet.CreateTime)
		assert.NotEqual(t, mReturn.UpdateTime, mGet.UpdateTime)

		mGet.Data.Name = "upsert"
		mGet.Status = 12
		time.Sleep(time.Second * 1) //
		updateTime := mGet.UpdateTime
		mUpsert, err := noticeDao.UpsertReturn(conn, mGet)
		assert.Equal(t, nil, err)
		assert.Equal(t, 12, mUpsert.Status)
		assert.Equal(t, mUpsert.CreateTime, mGet.CreateTime)
		assert.Equal(t, mUpsert.UpdateTime, mGet.UpdateTime)
		assert.NotEqual(t, updateTime, mUpsert.UpdateTime)
		assert.Equal(t, mGet.Data.Name, mUpsert.Data.Name)

		return err
	})
}

func createNotice(db *pg.DB) error {

	err := db.Model((*model.Notice)(nil)).DropTable(&orm.DropTableOptions{IfExists: true, Cascade: true})
	if err != nil {
		log.Print(err)
	}
	err = db.Model((*model.DataType)(nil)).DropComposite(&orm.DropCompositeOptions{
		IfExists: true,
	})
	if err != nil {
		log.Print(err)
	}
	err = db.Model((*model.DataType)(nil)).CreateComposite(nil)
	if err != nil {
		log.Print(err)
	}
	err = db.Model((*model.Notice)(nil)).CreateTable(nil)
	if err != nil {
		log.Print(err)
	}
	return err
}
