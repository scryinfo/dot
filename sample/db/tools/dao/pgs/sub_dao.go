package pgs

import (
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/db/pgs"
	"github.com/scryinfo/dot/sample/db/tools/model"
	"github.com/scryinfo/scryg/sutils/uuid"
)

const SubDaoTypeID = "6e3e5751-a618-410b-b06a-4ca960f2f4ac"

type SubDao struct {
	*pgs.DaoBase `dot:""`
}

//SubDaoTypeLives
func SubDaoTypeLives() []*dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{
			Name:   "SubDao",
			TypeID: SubDaoTypeID,
			NewDoter: func(conf []byte) (dot.Dot, error) {
				return &SubDao{}, nil
			},
		},
		Lives: []dot.Live{
			{
				LiveID: SubDaoTypeID,
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
func (c *SubDao) GetByIDWithLock(conn orm.DB, id string) (m *model.Sub, err error) {
	m = &model.Sub{}
	m.ID = id
	err = conn.Model(m).WherePK().For("UPDATE").Select()
	if err != nil {
		m = nil
	}
	return
}
func (c *SubDao) GetByID(conn orm.DB, id string) (m *model.Sub, err error) {
	m = &model.Sub{}
	m.ID = id
	err = conn.Model(m).WherePK().Select()
	if err != nil {
		m = nil
	}
	return
}

//update before
//you must get OptimisticLockVersion value
func (c *SubDao) GetLockByID(conn orm.DB, ids ...string) (ms []*model.Sub, err error) {
	for i, _ := range ids {
		m := &model.Sub{}
		m.ID = ids[i]
		ms = append(ms, m)
	}
	err = conn.Model(&ms).WherePK().Column(model.Sub_OptimisticLockVersion, model.Sub_ID).For("UPDATE").Select()
	if err != nil {
		ms = nil
	}
	return
}
func (c *SubDao) GetLockByModelID(conn orm.DB, ms ...*model.Sub) error {
	return conn.Model(&ms).WherePK().Column(model.Sub_OptimisticLockVersion, model.Sub_ID).For("UPDATE").Select()
}

// if find nothing, return pg.ErrNoRows
func (c *SubDao) QueryWithLock(conn orm.DB, condition string, params ...interface{}) (ms []*model.Sub, err error) {
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
func (c *SubDao) Query(conn orm.DB, condition string, params ...interface{}) (ms []*model.Sub, err error) {
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
func (c *SubDao) ListWithLock(conn orm.DB) (ms []*model.Sub, err error) {
	err = conn.Model(&ms).For("UPDATE").Select()
	if err != nil { //be sure
		ms = nil
	}
	return
}
func (c *SubDao) List(conn orm.DB) (ms []*model.Sub, err error) {
	err = conn.Model(&ms).Select()
	if err != nil { //be sure
		ms = nil
	}
	return
}

func (c *SubDao) Count(conn orm.DB, condition string, params ...interface{}) (count int, err error) {
	if len(condition) < 1 {
		count, err = conn.Model(&model.Sub{}).Count()
	} else {
		count, err = conn.Model(&model.Sub{}).Where(condition, params...).Count()
	}
	return
}

// if find nothing, return pg.ErrNoRows
func (c *SubDao) QueryPageWithLock(conn orm.DB, pageSize int, page int, condition string, params ...interface{}) (ms []*model.Sub, err error) {
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
func (c *SubDao) QueryPage(conn orm.DB, pageSize int, page int, condition string, params ...interface{}) (ms []*model.Sub, err error) {
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
func (c *SubDao) QueryPageWithCount(
	conn orm.DB,
	pageSize,
	pageNum int,
	condition string,
	params ...interface{},
) (ms []*model.Sub, count int, err error) {
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
func (c *SubDao) QueryOneWithLock(conn orm.DB, condition string, params ...interface{}) (m *model.Sub, err error) {
	m = &model.Sub{}
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
func (c *SubDao) QueryOne(conn orm.DB, condition string, params ...interface{}) (m *model.Sub, err error) {
	m = &model.Sub{}
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
func (c *SubDao) Insert(conn orm.DB, m *model.Sub) (err error) {
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuid()
	}
	m.CreateTime = time.Now().Unix()
	m.UpdateTime = m.CreateTime
	_, err = conn.Model(m).Insert()
	return
}

//if insert nothing, then return pg.ErrNoRows
func (c *SubDao) InsertReturn(conn orm.DB, m *model.Sub) (mnew *model.Sub, err error) {
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuid()
	}
	m.CreateTime = time.Now().Unix()
	m.UpdateTime = m.CreateTime

	mnew = &model.Sub{}
	_, err = conn.Model(m).Returning("*").Insert(mnew)
	if err != nil {
		mnew = nil
	}
	return
}

//if update nothing, then return pg.ErrNoRows
func (c *SubDao) Upsert(conn orm.DB, m *model.Sub) (err error) {
	m.UpdateTime = time.Now().Unix()
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuid()
		m.CreateTime = m.UpdateTime
	} else if m.CreateTime == 0 {
		m.CreateTime = m.UpdateTime
	}
	m.OptimisticLockVersion++
	om := conn.Model(m).OnConflict("(id) DO UPDATE").Where(model.Sub_Struct+"."+model.Sub_OptimisticLockVersion+" = ?", m.OptimisticLockVersion-1)
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
func (c *SubDao) UpsertReturn(conn orm.DB, m *model.Sub) (mnew *model.Sub, err error) {
	m.UpdateTime = time.Now().Unix()
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuid()
		m.CreateTime = m.UpdateTime
	} else if m.CreateTime == 0 {
		m.CreateTime = m.UpdateTime
	}
	m.OptimisticLockVersion++
	om := conn.Model(m).OnConflict("(id) DO UPDATE").Where(model.Sub_Struct+"."+model.Sub_OptimisticLockVersion+" = ?", m.OptimisticLockVersion-1)
	for _, it := range m.ToUpsertSet() {
		om.Set(it)
	}
	mnew = &model.Sub{}
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
func (c *SubDao) Update(conn orm.DB, m *model.Sub) (err error) {
	m.UpdateTime = time.Now().Unix()
	m.OptimisticLockVersion++
	//err = conn.Update(m)
	res, err := conn.Model(m).Where(model.Sub_ID+" = ? and "+model.Sub_OptimisticLockVersion+" = ?", m.ID, m.OptimisticLockVersion-1).Update()
	if res.RowsAffected() == 0 {
		err = pg.ErrNoRows
	}
	return
}

//if update nothing, then return pg.ErrNoRows
func (c *SubDao) UpdateReturn(conn orm.DB, m *model.Sub) (mnew *model.Sub, err error) {
	m.UpdateTime = time.Now().Unix()
	m.OptimisticLockVersion++
	mnew = &model.Sub{}
	res, err := conn.Model(m).Where(model.Sub_ID+" = ? and "+model.Sub_OptimisticLockVersion+" = ?", m.ID, m.OptimisticLockVersion-1).Returning("*").Update(mnew)
	if err != nil {
		mnew = nil
	}
	if res.RowsAffected() == 0 {
		err = pg.ErrNoRows
	}
	return
}

//if delete nothing, then return pg.ErrNoRows
func (c *SubDao) Delete(conn orm.DB, m *model.Sub) (err error) {
	_, err = conn.Model(m).Delete()
	return
}

//if delete nothing, then return pg.ErrNoRows
func (c *SubDao) DeleteByID(conn orm.DB, id string) (err error) {
	_, err = conn.Model((*model.Sub)(nil)).Where(model.Sub_ID+" = ?", id).Delete()
	return
}

//if delete nothing, then return pg.ErrNoRows
func (c *SubDao) DeleteByIDs(conn orm.DB, ids []string, oneMax int) (err error) {
	m := (*model.Sub)(nil)
	max := oneMax
	times := len(ids) / max
	for i := 1; i < times; i++ {
		oneIDs := ids[(i-1)*max : i*max-1]
		_, err = conn.Model(m).Where(model.Sub_ID+" in (?)", pg.In(oneIDs)).Delete()
		if err != nil {
			return
		}
	}

	if max*times < len(ids) {
		oneIDs := ids[max*times:]
		_, err = conn.Model(m).Where(model.Sub_ID+" in (?)", pg.In(oneIDs)).Delete()
	}
	return
}

//if delete nothing, then return pg.ErrNoRows
func (c *SubDao) DeleteReturn(conn orm.DB, m *model.Sub) (mnew *model.Sub, err error) {
	mnew = &model.Sub{}
	_, err = conn.Model(m).WherePK().Returning("*").Delete(mnew)
	if err != nil {
		mnew = nil
	}
	return
}

//example,please edit it
//update designated column with Optimistic Lock
func (c *SubDao) UpdateSubSomeColumn(conn orm.DB, ids []string /*todo: update parameters*/) (err error) {

	ms, err := c.GetLockByID(conn, ids...)
	if err != nil {
		return
	}

	for i, _ := range ms {
		ms[i].UpdateTime = time.Now().Unix()
		ms[i].OptimisticLockVersion++
		//todo ms[i].xx=parameter
		_, err = conn.Model(ms[i]).Where(model.Sub_ID+" = ? and "+model.Sub_OptimisticLockVersion+" = ?", ms[i].ID, ms[i].OptimisticLockVersion-1).Column( /*model.Sub_xx,*/ model.Sub_OptimisticLockVersion, model.Sub_UpdateTime).Update()
		if err != nil {
			dot.Logger().Debugln(err.Error())
			return
		}
	}
	return
}
