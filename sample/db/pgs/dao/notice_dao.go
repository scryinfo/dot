package dao

import (
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/db/pgs"
	"github.com/scryinfo/dot/sample/db/pgs/model"
	"github.com/scryinfo/scryg/sutils/uuid"
)

const NoticeDaoTypeID = "aee8a20e-4f3e-464f-9fea-a13537de0b5f"

type NoticeDao struct {
	*pgs.DaoBase `dot:""`
}

//NoticeDaoTypeLives
func NoticeDaoTypeLives() []*dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{
			Name:   "NoticeDao",
			TypeID: NoticeDaoTypeID,
			NewDoter: func(conf []byte) (dot.Dot, error) {
				return &NoticeDao{}, nil
			},
		},
		Lives: []dot.Live{
			{
				LiveID: NoticeDaoTypeID,
				RelyLives: map[string]dot.LiveID{
					"DaoBase": pgs.DaoBaseTypeID,
				},
			},
		},
	}

	lives := pgs.DaoBaseTypeLives()
	lives = append(lives, tl)
	return lives
}

// if find nothing, return pg.ErrNoRows
func (c *NoticeDao) GetByIDWithLock(conn orm.DB, id string) (m *model.Notice, err error) {
	m = &model.Notice{}
	m.ID = id
	err = conn.Model(m).WherePK().For("UPDATE").Select()
	if err != nil {
		m = nil
	}
	return
}
func (c *NoticeDao) GetByID(conn orm.DB, id string) (m *model.Notice, err error) {
	m = &model.Notice{}
	m.ID = id
	err = conn.Model(m).WherePK().Select()
	if err != nil {
		m = nil
	}
	return
}

//update before
//you must get OptimisticLockVersion value
func (c *NoticeDao) GetLockByID(conn orm.DB, ids ...string) (ms []*model.Notice, err error) {
	for i, _ := range ids {
		m := &model.Notice{}
		m.ID = ids[i]
		ms = append(ms, m)
	}
	err = conn.Model(&ms).WherePK().Column(model.Notice_OptimisticLockVersion, model.Notice_ID).For("UPDATE").Select()
	if err != nil {
		ms = nil
	}
	return
}
func (c *NoticeDao) GetLockByModelID(conn orm.DB, ms ...*model.Notice) error {
	return conn.Model(&ms).WherePK().Column(model.Notice_OptimisticLockVersion, model.Notice_ID).For("UPDATE").Select()
}

// if find nothing, return pg.ErrNoRows
func (c *NoticeDao) QueryWithLock(conn orm.DB, condition string, params ...interface{}) (ms []*model.Notice, err error) {
	if len(condition) < 1 {
		err = conn.Model(&ms).For("UPDATE").Select()
	} else {
		err = conn.Model(&ms).Where(condition, params...).For("UPDATE").Select()
	}
	if err != nil { //be sure
		ms = nil
	}
	return
}
func (c *NoticeDao) Query(conn orm.DB, condition string, params ...interface{}) (ms []*model.Notice, err error) {
	if len(condition) < 1 {
		err = conn.Model(&ms).Select()
	} else {
		err = conn.Model(&ms).Where(condition, params...).Select()
	}
	if err != nil { //be sure
		ms = nil
	}
	return
}

// if find nothing, return pg.ErrNoRows
func (c *NoticeDao) ListWithLock(conn orm.DB) (ms []*model.Notice, err error) {
	err = conn.Model(&ms).For("UPDATE").Select()
	if err != nil { //be sure
		ms = nil
	}
	return
}
func (c *NoticeDao) List(conn orm.DB) (ms []*model.Notice, err error) {
	err = conn.Model(&ms).Select()
	if err != nil { //be sure
		ms = nil
	}
	return
}

func (c *NoticeDao) Count(conn orm.DB, condition string, params ...interface{}) (count int, err error) {
	if len(condition) < 1 {
		count, err = conn.Model(&model.Notice{}).Count()
	} else {
		count, err = conn.Model(&model.Notice{}).Where(condition, params...).Count()
	}
	return
}

// if find nothing, return pg.ErrNoRows
func (c *NoticeDao) QueryPageWithLock(conn orm.DB, pageSize int, page int, condition string, params ...interface{}) (ms []*model.Notice, err error) {
	if len(condition) < 1 {
		err = conn.Model(&ms).Limit(pageSize).Offset((page - 1) * pageSize).For("UPDATE").Select()
	} else {
		err = conn.Model(&ms).Where(condition, params...).Limit(pageSize).Offset((page - 1) * pageSize).For("UPDATE").Select()
	}
	if err != nil { //be sure
		ms = nil
	}
	return
}
func (c *NoticeDao) QueryPage(conn orm.DB, pageSize int, page int, condition string, params ...interface{}) (ms []*model.Notice, err error) {
	if len(condition) < 1 {
		err = conn.Model(&ms).Limit(pageSize).Offset((page - 1) * pageSize).Select()
	} else {
		err = conn.Model(&ms).Where(condition, params...).Limit(pageSize).Offset((page - 1) * pageSize).Select()
	}
	if err != nil { //be sure
		ms = nil
	}
	return
}

// count counts valid records which after conditions filter, rather than whole table's count
func (c *NoticeDao) QueryPageWithCount(
	conn orm.DB,
	pageSize,
	pageNum int,
	condition string,
	params ...interface{},
) (ms []*model.Notice, count int, err error) {
	if len(condition) < 1 {
		count, err = conn.Model(&ms).Limit(pageSize).Offset((pageNum - 1) * pageSize).SelectAndCount()
	} else {
		count, err = conn.Model(&ms).Where(condition, params...).Limit(pageSize).Offset((pageNum - 1) * pageSize).SelectAndCount()
	}

	if err != nil { //be sure
		ms = nil
	}
	return
}

// if find nothing, return pg.ErrNoRows
func (c *NoticeDao) QueryOneWithLock(conn orm.DB, condition string, params ...interface{}) (m *model.Notice, err error) {
	m = &model.Notice{}
	if len(condition) < 1 {
		err = conn.Model(m).For("UPDATE").First()
	} else {
		err = conn.Model(m).Where(condition, params...).For("UPDATE").First()
	}
	if err != nil { //be sure
		m = nil
	}
	return
}
func (c *NoticeDao) QueryOne(conn orm.DB, condition string, params ...interface{}) (m *model.Notice, err error) {
	m = &model.Notice{}
	if len(condition) < 1 {
		err = conn.Model(m).First()
	} else {
		err = conn.Model(m).Where(condition, params...).First()
	}
	if err != nil { //be sure
		m = nil
	}
	return
}

//if insert nothing, then return pg.ErrNoRows
func (c *NoticeDao) Insert(conn orm.DB, m *model.Notice) (err error) {
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuid()
	}
	m.CreateTime = time.Now().Unix()
	m.UpdateTime = m.CreateTime
	_, err = conn.Model(m).Insert()
	return
}

//if insert nothing, then return pg.ErrNoRows
func (c *NoticeDao) InsertReturn(conn orm.DB, m *model.Notice) (mnew *model.Notice, err error) {
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuid()
	}
	m.CreateTime = time.Now().Unix()
	m.UpdateTime = m.CreateTime

	mnew = &model.Notice{}
	_, err = conn.Model(m).Returning("*").Insert(mnew)
	if err != nil {
		mnew = nil
	}
	return
}

//if update nothing, then return pg.ErrNoRows
func (c *NoticeDao) Upsert(conn orm.DB, m *model.Notice) (err error) {
	m.UpdateTime = time.Now().Unix()
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuid()
		m.CreateTime = m.UpdateTime
	} else if m.CreateTime == 0 {
		m.CreateTime = m.UpdateTime
	}
	m.OptimisticLockVersion++
	om := conn.Model(m).OnConflict("(id) DO UPDATE").Where(model.Notice_Struct+"."+model.Notice_OptimisticLockVersion+" = ?", m.OptimisticLockVersion-1)
	for _, it := range m.ToUpsertSet() {
		om.Set(it)
	}
	res, err := om.Insert()
	if res.RowsAffected() == 0 {
		//err = pg.ErrNoRows
		newm, err := c.GetLockByID(conn, m.ID)
		if err != nil {
			return err
		}
		m.OptimisticLockVersion = newm[0].OptimisticLockVersion
		err = c.Update(conn, m)
		return err
	}
	return err
}

//if update nothing, then return pg.ErrNoRows
func (c *NoticeDao) UpsertReturn(conn orm.DB, m *model.Notice) (mnew *model.Notice, err error) {
	m.UpdateTime = time.Now().Unix()
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuid()
		m.CreateTime = m.UpdateTime
	} else if m.CreateTime == 0 {
		m.CreateTime = m.UpdateTime
	}
	m.OptimisticLockVersion++
	om := conn.Model(m).OnConflict("(id) DO UPDATE").Where(model.Notice_Struct+"."+model.Notice_OptimisticLockVersion+" = ?", m.OptimisticLockVersion-1)
	for _, it := range m.ToUpsertSet() {
		om.Set(it)
	}
	mnew = &model.Notice{}
	_, err = om.Returning("*").Insert(mnew)
	if err == pg.ErrNoRows {
		new_m, err := c.GetLockByID(conn, m.ID)
		if err != nil {
			return nil, err
		}
		m.OptimisticLockVersion = new_m[0].OptimisticLockVersion
		mnew, err = c.UpdateReturn(conn, m)
		return mnew, err
	}
	return
}

//if update nothing, then return pg.ErrNoRows
func (c *NoticeDao) Update(conn orm.DB, m *model.Notice) (err error) {
	m.UpdateTime = time.Now().Unix()
	m.OptimisticLockVersion++
	//err = conn.Update(m)
	res, err := conn.Model(m).Where(model.Notice_ID+" = ? and "+model.Notice_OptimisticLockVersion+" = ?", m.ID, m.OptimisticLockVersion-1).Update()
	if res.RowsAffected() == 0 {
		err = pg.ErrNoRows
	}
	return
}

//if update nothing, then return pg.ErrNoRows
func (c *NoticeDao) UpdateReturn(conn orm.DB, m *model.Notice) (mnew *model.Notice, err error) {
	m.UpdateTime = time.Now().Unix()
	m.OptimisticLockVersion++
	mnew = &model.Notice{}
	res, err := conn.Model(m).Where(model.Notice_ID+" = ? and "+model.Notice_OptimisticLockVersion+" = ?", m.ID, m.OptimisticLockVersion-1).Returning("*").Update(mnew)
	if err != nil {
		mnew = nil
	}
	if res.RowsAffected() == 0 {
		err = pg.ErrNoRows
	}
	return
}

//if delete nothing, then return pg.ErrNoRows
func (c *NoticeDao) Delete(conn orm.DB, m *model.Notice) (err error) {
	_, err = conn.Model(m).Delete()
	return
}

//if delete nothing, then return pg.ErrNoRows
func (c *NoticeDao) DeleteByID(conn orm.DB, id string) (err error) {
	_, err = conn.Model((*model.Notice)(nil)).Where(model.Notice_ID+" = ?", id).Delete()
	return
}

//if delete nothing, then return pg.ErrNoRows
func (c *NoticeDao) DeleteByIDs(conn orm.DB, ids []string, oneMax int) (err error) {
	m := (*model.Notice)(nil)
	max := oneMax
	times := len(ids) / max
	for i := 1; i < times; i++ {
		oneIDs := ids[(i-1)*max : i*max-1]
		_, err = conn.Model(m).Where(model.Notice_ID+" in (?)", pg.In(oneIDs)).Delete()
		if err != nil {
			return
		}
	}

	if max*times < len(ids) {
		oneIDs := ids[max*times:]
		_, err = conn.Model(m).Where(model.Notice_ID+" in (?)", pg.In(oneIDs)).Delete()
	}
	return
}

//if delete nothing, then return pg.ErrNoRows
func (c *NoticeDao) DeleteReturn(conn orm.DB, m *model.Notice) (mnew *model.Notice, err error) {
	mnew = &model.Notice{}
	_, err = conn.Model(m).WherePK().Returning("*").Delete(mnew)
	if err != nil {
		mnew = nil
	}
	return
}

//example,please edit it
//update designated column with Optimistic Lock
func (c *NoticeDao) UpdateNoticeSomeColumn(conn orm.DB, ids []string /*todo: update parameters*/) (err error) {

	ms, err := c.GetLockByID(conn, ids...)
	if err != nil {
		return
	}

	for i, _ := range ms {
		ms[i].UpdateTime = time.Now().Unix()
		ms[i].OptimisticLockVersion++
		//todo ms[i].xx=parameter
		_, err = conn.Model(ms[i]).Where(model.Notice_ID+" = ? and "+model.Notice_OptimisticLockVersion+" = ?", ms[i].ID, ms[i].OptimisticLockVersion-1).Column( /*model.Notice_xx,*/ model.Notice_OptimisticLockVersion, model.Notice_UpdateTime).Update()
		if err != nil {
			dot.Logger().Debugln(err.Error())
			return
		}
	}
	return
}
