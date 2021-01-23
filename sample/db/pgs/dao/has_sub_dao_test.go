package dao

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/scryinfo/dot/dots/db/pgs"
	"github.com/scryinfo/dot/sample/db/pgs/model"
	"github.com/scryinfo/scryg/sutils/uuid"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

//const dbconf_test = `{
//			"addr": "localhost:5432",
//			"user": "postgres",
//			"password": "scry",
//			"database": "postgres",
//			"showSql": true
//		}`

func TestHasSubDao(t *testing.T) {
	hasSubDao := &HasSubDao{}
	hasSubDao.DaoBase = pgs.GenerateDaoBase(dbconf_test)
	subDao := &SubDao{DaoBase: hasSubDao.DaoBase}
	err := createHasSub(hasSubDao.DaoBase.Wrapper.GetDb())
	assert.Equal(t, nil, err)

	hasSubDao.DaoBase.WithNoTx(func(conn orm.DB) error {
		hasSub := &model.HasSub{}
		hasSub.SubData = &model.Sub{}
		mReturn, err := hasSubDao.InsertReturn(conn, hasSub)
		assert.Equal(t, nil, err)
		assert.NotEqual(t, nil, mReturn.ID)
		assert.Equal(t, "", mReturn.SubId)
		assert.Equal(t, true, mReturn.SubData == nil)
		err = hasSubDao.DeleteByID(conn, mReturn.ID)
		assert.Equal(t, nil, err)

		hasSub = &model.HasSub{}
		hasSub.SubData = &model.Sub{}
		hasSub.SubData.ID = uuid.GetUuid()
		hasSub.SubId = hasSub.SubData.ID
		mReturn, err = hasSubDao.InsertReturn(conn, hasSub)
		assert.Equal(t, nil, err)
		assert.NotEqual(t, nil, mReturn.ID)
		assert.Equal(t, hasSub.SubId, mReturn.SubId)
		assert.Equal(t, true, mReturn.SubData == nil)
		mInsertSub, err := subDao.InsertReturn(conn, hasSub.SubData)
		assert.Equal(t, nil, err)
		assert.Equal(t, hasSub.SubData.ID, mInsertSub.ID)

		mGet, err := hasSubDao.GetByID(conn, mReturn.ID)
		assert.Equal(t, nil, err)
		assert.Equal(t, mReturn.ID, mGet.ID)
		assert.Equal(t, mReturn.UpdateTime, mGet.UpdateTime)
		assert.Equal(t, mReturn.CreateTime, mGet.CreateTime)
		assert.Equal(t, mReturn.OptimisticLockVersion, mGet.OptimisticLockVersion)

		ms, err := hasSubDao.Query(conn, "")
		assert.Equal(t, nil, err)
		assert.Equal(t, 1, len(ms))
		mGet = ms[0]
		assert.Equal(t, mReturn.ID, mGet.ID)
		assert.Equal(t, mReturn.UpdateTime, mGet.UpdateTime)
		assert.Equal(t, mReturn.CreateTime, mGet.CreateTime)
		assert.Equal(t, mReturn.OptimisticLockVersion, mGet.OptimisticLockVersion)

		mGet, err = hasSubDao.QueryOne(conn, "")
		assert.Equal(t, nil, err)
		assert.Equal(t, mReturn.ID, mGet.ID)
		assert.Equal(t, mReturn.UpdateTime, mGet.UpdateTime)
		assert.Equal(t, mReturn.CreateTime, mGet.CreateTime)
		assert.Equal(t, mReturn.OptimisticLockVersion, mGet.OptimisticLockVersion)

		ms, err = hasSubDao.List(conn)
		assert.Equal(t, nil, err)
		assert.Equal(t, 1, len(ms))
		mGet = ms[0]
		assert.Equal(t, mReturn.ID, mGet.ID)
		assert.Equal(t, mReturn.UpdateTime, mGet.UpdateTime)
		assert.Equal(t, mReturn.CreateTime, mGet.CreateTime)
		assert.Equal(t, mReturn.OptimisticLockVersion, mGet.OptimisticLockVersion)

		ms, err = hasSubDao.QueryPage(conn, 10, 1, "")
		assert.Equal(t, nil, err)
		assert.Equal(t, 1, len(ms))
		mGet = ms[0]
		assert.Equal(t, mReturn.ID, mGet.ID)
		assert.Equal(t, mReturn.UpdateTime, mGet.UpdateTime)
		assert.Equal(t, mReturn.CreateTime, mGet.CreateTime)
		assert.Equal(t, mReturn.OptimisticLockVersion, mGet.OptimisticLockVersion)

		ms, count, err := hasSubDao.QueryPageWithCount(conn, 10, 1, "")
		assert.Equal(t, nil, err)
		assert.Equal(t, 1, len(ms))
		assert.Equal(t, 1, count)
		mGet = ms[0]
		assert.Equal(t, mReturn.ID, mGet.ID)
		assert.Equal(t, mReturn.UpdateTime, mGet.UpdateTime)
		assert.Equal(t, mReturn.CreateTime, mGet.CreateTime)
		assert.Equal(t, mReturn.OptimisticLockVersion, mGet.OptimisticLockVersion)

		time.Sleep(time.Second * 1) //
		err = hasSubDao.Update(conn, mGet)
		assert.Equal(t, nil, err)
		mGet, err = hasSubDao.QueryOne(conn, fmt.Sprintf(" %s = ?", model.Notice_ID), mGet.ID)
		assert.Equal(t, nil, err)
		//assert.Equal(t, 10, mGet.Status)
		assert.Equal(t, mReturn.CreateTime, mGet.CreateTime)
		assert.NotEqual(t, mReturn.UpdateTime, mGet.UpdateTime)

		mGet.Count = 11
		time.Sleep(time.Second * 1) //
		mReturn, err = hasSubDao.UpdateReturn(conn, mGet)
		assert.Equal(t, nil, err)
		assert.Equal(t, mGet.Count, mReturn.Count)
		assert.Equal(t, 11, mReturn.Count)
		assert.Equal(t, mReturn.CreateTime, mGet.CreateTime)
		assert.Equal(t, mReturn.UpdateTime, mGet.UpdateTime)
		mReturn, err = hasSubDao.GetByID(conn, mGet.ID)
		assert.Equal(t, nil, err)
		assert.Equal(t, 11, mReturn.Count)
		assert.Equal(t, mReturn.CreateTime, mGet.CreateTime)
		assert.Equal(t, mReturn.UpdateTime, mGet.UpdateTime)

		mGet.Count = 12
		time.Sleep(time.Second * 1) //
		mReturn, err = hasSubDao.UpsertReturn(conn, mGet)
		assert.Equal(t, nil, err)
		assert.Equal(t, 12, mReturn.Count)
		assert.Equal(t, mReturn.CreateTime, mGet.CreateTime)
		assert.Equal(t, mReturn.UpdateTime, mGet.UpdateTime)
		return err
	})
}

func createHasSub(db *pg.DB) error {

	err := db.Model((*model.Sub)(nil)).DropTable(&orm.DropTableOptions{IfExists: true, Cascade: true})
	if err != nil {
		log.Print(err)
	}
	err = db.Model((*model.HasSub)(nil)).DropTable(&orm.DropTableOptions{
		IfExists: true,
	})
	if err != nil {
		log.Print(err)
	}
	err = db.Model((*model.Sub)(nil)).CreateTable(nil)
	if err != nil {
		log.Print(err)
	}
	err = db.Model((*model.HasSub)(nil)).CreateTable(nil)
	if err != nil {
		log.Print(err)
	}
	return err
}
