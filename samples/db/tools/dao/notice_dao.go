package dao

import (
	"context"
	"database/sql"
	"time"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/db/buns"
	"github.com/scryinfo/dot/samples/db/tools/model"
	"github.com/scryinfo/scryg/sutils/uuid"
	"github.com/uptrace/bun"
)

func NewNoticeDao(conn *buns.DaoBase) *NoticeDao {
	return &NoticeDao{
		DaoBase: conn,
	}
}

type NoticeDao struct {
	*buns.DaoBase
}

// if find|update|insert nothing, sql.ErrNoRows error may returned

func (c *NoticeDao) GetByIDWithLock(conn bun.IDB, id string) (m *model.Notice, err error) {
	m = &model.Notice{}
	m.ID = id
	err = conn.NewSelect().Model(m).WherePK().For("UPDATE").Scan(context.TODO())
	if err != nil {
		m = nil
	}
	return
}
func (c *NoticeDao) GetByID(conn bun.IDB, id string) (m *model.Notice, err error) {
	m = &model.Notice{}
	m.ID = id
	err = conn.NewSelect().Model(m).WherePK().Scan(context.TODO())
	if err != nil {
		m = nil
	}
	return
}

// update before
// you must get OptimisticLockVersion value
func (c *NoticeDao) GetLockByID(conn bun.IDB, ids ...string) (ms []*model.Notice, err error) {
	for i, _ := range ids {
		m := &model.Notice{}
		m.ID = ids[i]
		ms = append(ms, m)
	}
	err = conn.NewSelect().Model(&ms).WherePK().Column(model.Notice_OptimisticLockVersion, model.Notice_Struct+"."+model.Notice_ID).For("UPDATE").Scan(context.TODO())
	if err != nil {
		ms = nil
	}
	return
}
func (c *NoticeDao) GetLockByModelID(conn bun.IDB, ms ...*model.Notice) error {
	return conn.NewSelect().Model(&ms).WherePK().Column(model.Notice_OptimisticLockVersion, model.Notice_Struct+"."+model.Notice_ID).For("UPDATE").Scan(context.TODO())
}

func (c *NoticeDao) QueryWithLock(conn bun.IDB, condition string, params ...interface{}) (ms []*model.Notice, err error) {
	if len(condition) < 1 {
		err = conn.NewSelect().Model(&ms).For("UPDATE").Scan(context.TODO())
	} else {
		err = conn.NewSelect().Model(&ms).Where(condition, params...).For("UPDATE").Scan(context.TODO())
	}
	if err != nil { //be sure
		ms = nil
	}
	return
}
func (c *NoticeDao) Query(conn bun.IDB, condition string, params ...interface{}) (ms []*model.Notice, err error) {
	if len(condition) < 1 {
		err = conn.NewSelect().Model(&ms).Scan(context.TODO())
	} else {
		err = conn.NewSelect().Model(&ms).Where(condition, params...).Scan(context.TODO())
	}
	if err != nil { //be sure
		ms = nil
	}
	return
}

func (c *NoticeDao) ListWithLock(conn bun.IDB) (ms []*model.Notice, err error) {
	err = conn.NewSelect().Model(&ms).For("UPDATE").Scan(context.TODO())
	if err != nil { //be sure
		ms = nil
	}
	return
}
func (c *NoticeDao) List(conn bun.IDB) (ms []*model.Notice, err error) {
	err = conn.NewSelect().Model(&ms).Scan(context.TODO())
	if err != nil { //be sure
		ms = nil
	}
	return
}

func (c *NoticeDao) Count(conn bun.IDB, condition string, params ...interface{}) (count int, err error) {
	if len(condition) < 1 {
		count, err = conn.NewSelect().Model(&model.Notice{}).Count(context.TODO())
	} else {
		count, err = conn.NewSelect().Model(&model.Notice{}).Where(condition, params...).Count(context.TODO())
	}
	return
}

func (c *NoticeDao) QueryPageWithLock(conn bun.IDB, pageSize int, page int, condition string, params ...interface{}) (ms []*model.Notice, err error) {
	if len(condition) < 1 {
		err = conn.NewSelect().Model(&ms).Limit(pageSize).Offset((page - 1) * pageSize).For("UPDATE").Scan(context.TODO())
	} else {
		err = conn.NewSelect().Model(&ms).Where(condition, params...).Limit(pageSize).Offset((page - 1) * pageSize).For("UPDATE").Scan(context.TODO())
	}
	if err != nil { //be sure
		ms = nil
	}
	return
}
func (c *NoticeDao) QueryPage(conn bun.IDB, pageSize int, page int, condition string, params ...interface{}) (ms []*model.Notice, err error) {
	if len(condition) < 1 {
		err = conn.NewSelect().Model(&ms).Limit(pageSize).Offset((page - 1) * pageSize).Scan(context.TODO())
	} else {
		err = conn.NewSelect().Model(&ms).Where(condition, params...).Limit(pageSize).Offset((page - 1) * pageSize).Scan(context.TODO())
	}
	if err != nil { //be sure
		ms = nil
	}
	return
}

// count counts valid records which after conditions filter, rather than whole table's count
func (c *NoticeDao) QueryPageWithCount(
	conn bun.IDB,
	pageSize,
	pageNum int,
	condition string,
	params ...interface{},
) (ms []*model.Notice, count int, err error) {
	if len(condition) < 1 {
		count, err = conn.NewSelect().Model(&ms).Limit(pageSize).Offset((pageNum - 1) * pageSize).ScanAndCount(context.TODO())
	} else {
		count, err = conn.NewSelect().Model(&ms).Where(condition, params...).Limit(pageSize).Offset((pageNum - 1) * pageSize).ScanAndCount(context.TODO())
	}

	if err != nil { //be sure
		ms = nil
	}
	return
}

func (c *NoticeDao) QueryOneWithLock(conn bun.IDB, condition string, params ...interface{}) (m *model.Notice, err error) {
	m = &model.Notice{}
	if len(condition) < 1 {
		err = conn.NewSelect().Model(m).For("UPDATE").Limit(1).Scan(context.TODO())
	} else {
		err = conn.NewSelect().Model(m).Where(condition, params...).For("UPDATE").Limit(1).Scan(context.TODO())
	}
	if err != nil { //be sure
		m = nil
	}
	return
}

func (c *NoticeDao) QueryOne(conn bun.IDB, condition string, params ...interface{}) (m *model.Notice, err error) {
	m = &model.Notice{}
	if len(condition) < 1 {
		err = conn.NewSelect().Model(m).Scan(context.TODO())
	} else {
		err = conn.NewSelect().Model(m).Where(condition, params...).Scan(context.TODO())
	}
	if err != nil { //be sure
		m = nil
	}
	return
}

func (c *NoticeDao) Insert(conn bun.IDB, m *model.Notice) (err error) {
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuidForDB()
	}
	m.CreateTime = time.Now().Unix()
	m.UpdateTime = m.CreateTime
	_, err = conn.NewInsert().Model(m).Exec(context.TODO())
	return
}

func (c *NoticeDao) InsertReturn(conn bun.IDB, m *model.Notice) (mnew *model.Notice, err error) {
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuidForDB()
	}
	m.CreateTime = time.Now().Unix()
	m.UpdateTime = m.CreateTime

	mnew = &model.Notice{}
	_, err = conn.NewInsert().Model(m).Returning("*").Exec(context.TODO(), mnew)
	if err != nil {
		mnew = nil
	}
	return
}

func (c *NoticeDao) Upsert(conn bun.IDB, m *model.Notice) (err error) {
	m.UpdateTime = time.Now().Unix()
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuidForDB()
		m.CreateTime = m.UpdateTime
	} else if m.CreateTime == 0 {
		m.CreateTime = m.UpdateTime
	}
	m.OptimisticLockVersion++
	om := conn.NewInsert().Model(m).On("CONFLICT (id) DO UPDATE").Where(model.Notice_Struct+"."+model.Notice_OptimisticLockVersion+" = ?", m.OptimisticLockVersion-1)
	for _, it := range m.ToUpsertSet() {
		om.Set(it)
	}
	res, err := om.Exec(context.TODO())
	if n, _ := res.RowsAffected(); n == 0 {
		newm, err := c.GetLockByID(conn, m.ID)
		if err != nil {
			return err
		}
		m.OptimisticLockVersion = newm[0].OptimisticLockVersion
		err = c.Update(conn, m)
	}
	return err
}

func (c *NoticeDao) UpsertReturn(conn bun.IDB, m *model.Notice) (mnew *model.Notice, err error) {
	m.UpdateTime = time.Now().Unix()
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuidForDB()
		m.CreateTime = m.UpdateTime
	} else if m.CreateTime == 0 {
		m.CreateTime = m.UpdateTime
	}
	m.OptimisticLockVersion++
	om := conn.NewInsert().Model(m).On("CONFLICT (id) DO UPDATE").Where(model.Notice_Struct+"."+model.Notice_OptimisticLockVersion+" = ?", m.OptimisticLockVersion-1)
	for _, it := range m.ToUpsertSet() {
		om.Set(it)
	}
	mnew = &model.Notice{}
	res, err := om.Returning("*").Exec(context.TODO(), mnew)
	if n, _ := res.RowsAffected(); n == 0 {
		ms, err := c.GetLockByID(conn, m.ID)
		if err != nil {
			return nil, err
		}
		m.OptimisticLockVersion = ms[0].OptimisticLockVersion
		mnew, err = c.UpdateReturn(conn, m)
	}
	return
}

func (c *NoticeDao) Update(conn bun.IDB, m *model.Notice) (err error) {
	m.UpdateTime = time.Now().Unix()
	m.OptimisticLockVersion++
	res, err := conn.NewUpdate().Model(m).Where(model.Notice_ID+" = ? and "+model.Notice_OptimisticLockVersion+" = ?", m.ID, m.OptimisticLockVersion-1).Exec(context.TODO())
	if n, _ := res.RowsAffected(); n == 0 {
		err = sql.ErrNoRows
	}
	return
}

func (c *NoticeDao) UpdateReturn(conn bun.IDB, m *model.Notice) (mnew *model.Notice, err error) {
	m.UpdateTime = time.Now().Unix()
	m.OptimisticLockVersion++
	mnew = &model.Notice{}
	res, err := conn.NewUpdate().Model(m).Where(model.Notice_ID+" = ? and "+model.Notice_OptimisticLockVersion+" = ?", m.ID, m.OptimisticLockVersion-1).Returning("*").Exec(context.TODO(), mnew)
	if err != nil {
		mnew = nil
	}
	if n, _ := res.RowsAffected(); n == 0 {
		err = sql.ErrNoRows
	}
	return
}

func (c *NoticeDao) Delete(conn bun.IDB, m *model.Notice) (err error) {
	return c.DeleteByID(conn, m.ID)
}

func (c *NoticeDao) DeleteByID(conn bun.IDB, id string) (err error) {
	_, err = conn.NewDelete().Model((*model.Notice)(nil)).Where(model.Notice_ID+" = ?", id).Exec(context.TODO())
	return
}

func (c *NoticeDao) DeleteByIDs(conn bun.IDB, ids []string, oneMax int) (err error) {
	m := (*model.Notice)(nil)
	max := oneMax
	times := len(ids) / max
	for i := 1; i < times; i++ {
		oneIDs := ids[(i-1)*max : i*max-1]
		_, err = conn.NewDelete().Model(m).Where(model.Notice_ID+" in (?)", bun.In(oneIDs)).Exec(context.TODO())
		if err != nil {
			return
		}
	}

	if max*times < len(ids) {
		oneIDs := ids[max*times:]
		_, err = conn.NewDelete().Model(m).Where(model.Notice_ID+" in (?)", bun.In(oneIDs)).Exec(context.TODO())
	}
	return
}

func (c *NoticeDao) DeleteReturn(conn bun.IDB, m *model.Notice) (mnew *model.Notice, err error) {
	mnew = &model.Notice{}
	_, err = conn.NewDelete().Model(m).WherePK().Returning("*").Exec(context.TODO(), mnew)
	if err != nil {
		mnew = nil
	}
	return
}

// example,please edit it
// update designated column with Optimistic Lock
func (c *NoticeDao) UpdateNoticeSomeColumn(conn bun.IDB, ids []string /*todo: update parameters*/) (err error) {

	ms, err := c.GetLockByID(conn, ids...)
	if err != nil {
		return
	}
	ctx := context.TODO()
	condition := model.Notice_ID + " = ? and " + model.Notice_OptimisticLockVersion + " = ?"
	for i, _ := range ms {
		ms[i].UpdateTime = time.Now().Unix()
		ms[i].OptimisticLockVersion++
		_, err = conn.NewUpdate().Model(ms[i]).Where(condition, ms[i].ID, ms[i].OptimisticLockVersion-1).Column( /*model.Notice_xx,*/ model.Notice_OptimisticLockVersion, model.Notice_UpdateTime).Exec(ctx)
		if err != nil {
			dot.Logger.Debug().Err(err).Send()
			return
		}
	}
	return
}
