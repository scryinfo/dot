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

func NewSubDao(conn *buns.DaoBase) *SubDao {
	return &SubDao{
		DaoBase: conn,
	}
}

type SubDao struct {
	*buns.DaoBase
}

// if find|update|insert nothing, sql.ErrNoRows error may returned

func (c *SubDao) GetByIDWithLock(conn bun.IDB, id string) (m *model.Sub, err error) {
	m = &model.Sub{}
	m.ID = id
	err = conn.NewSelect().Model(m).WherePK().For("UPDATE").Scan(context.TODO())
	if err != nil {
		m = nil
	}
	return
}
func (c *SubDao) GetByID(conn bun.IDB, id string) (m *model.Sub, err error) {
	m = &model.Sub{}
	m.ID = id
	err = conn.NewSelect().Model(m).WherePK().Scan(context.TODO())
	if err != nil {
		m = nil
	}
	return
}

// update before
// you must get OptimisticLockVersion value
func (c *SubDao) GetLockByID(conn bun.IDB, ids ...string) (ms []*model.Sub, err error) {
	for i, _ := range ids {
		m := &model.Sub{}
		m.ID = ids[i]
		ms = append(ms, m)
	}
	err = conn.NewSelect().Model(&ms).WherePK().Column(model.Sub_OptimisticLockVersion, model.Sub_Struct+"."+model.Sub_ID).For("UPDATE").Scan(context.TODO())
	if err != nil {
		ms = nil
	}
	return
}
func (c *SubDao) GetLockByModelID(conn bun.IDB, ms ...*model.Sub) error {
	return conn.NewSelect().Model(&ms).WherePK().Column(model.Sub_OptimisticLockVersion, model.Sub_Struct+"."+model.Sub_ID).For("UPDATE").Scan(context.TODO())
}

func (c *SubDao) QueryWithLock(conn bun.IDB, condition string, params ...interface{}) (ms []*model.Sub, err error) {
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
func (c *SubDao) Query(conn bun.IDB, condition string, params ...interface{}) (ms []*model.Sub, err error) {
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

func (c *SubDao) ListWithLock(conn bun.IDB) (ms []*model.Sub, err error) {
	err = conn.NewSelect().Model(&ms).For("UPDATE").Scan(context.TODO())
	if err != nil { //be sure
		ms = nil
	}
	return
}
func (c *SubDao) List(conn bun.IDB) (ms []*model.Sub, err error) {
	err = conn.NewSelect().Model(&ms).Scan(context.TODO())
	if err != nil { //be sure
		ms = nil
	}
	return
}

func (c *SubDao) Count(conn bun.IDB, condition string, params ...interface{}) (count int, err error) {
	if len(condition) < 1 {
		count, err = conn.NewSelect().Model(&model.Sub{}).Count(context.TODO())
	} else {
		count, err = conn.NewSelect().Model(&model.Sub{}).Where(condition, params...).Count(context.TODO())
	}
	return
}

func (c *SubDao) QueryPageWithLock(conn bun.IDB, pageSize int, page int, condition string, params ...interface{}) (ms []*model.Sub, err error) {
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
func (c *SubDao) QueryPage(conn bun.IDB, pageSize int, page int, condition string, params ...interface{}) (ms []*model.Sub, err error) {
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
func (c *SubDao) QueryPageWithCount(
	conn bun.IDB,
	pageSize,
	pageNum int,
	condition string,
	params ...interface{},
) (ms []*model.Sub, count int, err error) {
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

func (c *SubDao) QueryOneWithLock(conn bun.IDB, condition string, params ...interface{}) (m *model.Sub, err error) {
	m = &model.Sub{}
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

func (c *SubDao) QueryOne(conn bun.IDB, condition string, params ...interface{}) (m *model.Sub, err error) {
	m = &model.Sub{}
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

func (c *SubDao) Insert(conn bun.IDB, m *model.Sub) (err error) {
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuidForDB()
	}
	m.CreateTime = time.Now().Unix()
	m.UpdateTime = m.CreateTime
	_, err = conn.NewInsert().Model(m).Exec(context.TODO())
	return
}

func (c *SubDao) InsertReturn(conn bun.IDB, m *model.Sub) (mnew *model.Sub, err error) {
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuidForDB()
	}
	m.CreateTime = time.Now().Unix()
	m.UpdateTime = m.CreateTime

	mnew = &model.Sub{}
	_, err = conn.NewInsert().Model(m).Returning("*").Exec(context.TODO(), mnew)
	if err != nil {
		mnew = nil
	}
	return
}

func (c *SubDao) Upsert(conn bun.IDB, m *model.Sub) (err error) {
	m.UpdateTime = time.Now().Unix()
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuidForDB()
		m.CreateTime = m.UpdateTime
	} else if m.CreateTime == 0 {
		m.CreateTime = m.UpdateTime
	}
	m.OptimisticLockVersion++
	om := conn.NewInsert().Model(m).On("CONFLICT (id) DO UPDATE").Where(model.Sub_Struct+"."+model.Sub_OptimisticLockVersion+" = ?", m.OptimisticLockVersion-1)
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

func (c *SubDao) UpsertReturn(conn bun.IDB, m *model.Sub) (mnew *model.Sub, err error) {
	m.UpdateTime = time.Now().Unix()
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuidForDB()
		m.CreateTime = m.UpdateTime
	} else if m.CreateTime == 0 {
		m.CreateTime = m.UpdateTime
	}
	m.OptimisticLockVersion++
	om := conn.NewInsert().Model(m).On("CONFLICT (id) DO UPDATE").Where(model.Sub_Struct+"."+model.Sub_OptimisticLockVersion+" = ?", m.OptimisticLockVersion-1)
	for _, it := range m.ToUpsertSet() {
		om.Set(it)
	}
	mnew = &model.Sub{}
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

func (c *SubDao) Update(conn bun.IDB, m *model.Sub) (err error) {
	m.UpdateTime = time.Now().Unix()
	m.OptimisticLockVersion++
	res, err := conn.NewUpdate().Model(m).Where(model.Sub_ID+" = ? and "+model.Sub_OptimisticLockVersion+" = ?", m.ID, m.OptimisticLockVersion-1).Exec(context.TODO())
	if n, _ := res.RowsAffected(); n == 0 {
		err = sql.ErrNoRows
	}
	return
}

func (c *SubDao) UpdateReturn(conn bun.IDB, m *model.Sub) (mnew *model.Sub, err error) {
	m.UpdateTime = time.Now().Unix()
	m.OptimisticLockVersion++
	mnew = &model.Sub{}
	res, err := conn.NewUpdate().Model(m).Where(model.Sub_ID+" = ? and "+model.Sub_OptimisticLockVersion+" = ?", m.ID, m.OptimisticLockVersion-1).Returning("*").Exec(context.TODO(), mnew)
	if err != nil {
		mnew = nil
	}
	if n, _ := res.RowsAffected(); n == 0 {
		err = sql.ErrNoRows
	}
	return
}

func (c *SubDao) Delete(conn bun.IDB, m *model.Sub) (err error) {
	return c.DeleteByID(conn, m.ID)
}

func (c *SubDao) DeleteByID(conn bun.IDB, id string) (err error) {
	_, err = conn.NewDelete().Model((*model.Sub)(nil)).Where(model.Sub_ID+" = ?", id).Exec(context.TODO())
	return
}

func (c *SubDao) DeleteByIDs(conn bun.IDB, ids []string, oneMax int) (err error) {
	m := (*model.Sub)(nil)
	max := oneMax
	times := len(ids) / max
	for i := 1; i < times; i++ {
		oneIDs := ids[(i-1)*max : i*max-1]
		_, err = conn.NewDelete().Model(m).Where(model.Sub_ID+" in (?)", bun.In(oneIDs)).Exec(context.TODO())
		if err != nil {
			return
		}
	}

	if max*times < len(ids) {
		oneIDs := ids[max*times:]
		_, err = conn.NewDelete().Model(m).Where(model.Sub_ID+" in (?)", bun.In(oneIDs)).Exec(context.TODO())
	}
	return
}

func (c *SubDao) DeleteReturn(conn bun.IDB, m *model.Sub) (mnew *model.Sub, err error) {
	mnew = &model.Sub{}
	_, err = conn.NewDelete().Model(m).WherePK().Returning("*").Exec(context.TODO(), mnew)
	if err != nil {
		mnew = nil
	}
	return
}

// example,please edit it
// update designated column with Optimistic Lock
func (c *SubDao) UpdateSubSomeColumn(conn bun.IDB, ids []string /*todo: update parameters*/) (err error) {

	ms, err := c.GetLockByID(conn, ids...)
	if err != nil {
		return
	}
	ctx := context.TODO()
	condition := model.Sub_ID + " = ? and " + model.Sub_OptimisticLockVersion + " = ?"
	for i, _ := range ms {
		ms[i].UpdateTime = time.Now().Unix()
		ms[i].OptimisticLockVersion++
		_, err = conn.NewUpdate().Model(ms[i]).Where(condition, ms[i].ID, ms[i].OptimisticLockVersion-1).Column( /*model.Sub_xx,*/ model.Sub_OptimisticLockVersion, model.Sub_UpdateTime).Exec(ctx)
		if err != nil {
			dot.Logger.Debug().Err(err).Send()
			return
		}
	}
	return
}
