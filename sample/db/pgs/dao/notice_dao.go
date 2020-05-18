package dao

import (
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/db/pgs"
	"github.com/scryinfo/dot/sample/db/pgs/model"
	"github.com/scryinfo/scryg/sutils/uuid"
)

const NoticeDaoTypeId = "72c713a6-7bd1-4594-a17e-8e4d55794a40"

type NoticeDao struct {
	*pgs.DaoBase `dot:""`
}

//NoticeDaoTypeLives
func NoticeDaoTypeLives() []*dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{
			Name:   "NoticeDao",
			TypeId: NoticeDaoTypeId,
			NewDoter: func(conf []byte) (dot.Dot, error) {
				return &NoticeDao{}, nil
			},
		},
		Lives: []dot.Live{
			{
				LiveId: NoticeDaoTypeId,
				RelyLives: map[string]dot.LiveId{
					"DaoBase": pgs.DaoBaseTypeId,
				},
			},
		},
	}

	lives := pgs.DaoBaseTypeLives()
	lives = append(lives, tl)
	return lives
}

// if find nothing, return pg.ErrNoRows
func (c *NoticeDao) GetById(conn *pg.Conn, id string) (m *model.Notice, err error) {
	m = &model.Notice{Id: id}
	err = conn.Model(m).WherePK().For("UPDATE").Select()
	if err != nil {
		m = nil
	}
	return
}

// if find nothing, return pg.ErrNoRows
func (c *NoticeDao) Query(conn *pg.Conn, condition string, params ...interface{}) (ms []*model.Notice, err error) {
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

// if find nothing, return pg.ErrNoRows
func (c *NoticeDao) List(conn *pg.Conn) (ms []*model.Notice, err error) {
	err = conn.Model(&ms).For("UPDATE").Select()
	if err != nil { //be sure
		ms = nil
	}
	return
}

func (c *NoticeDao) Count(conn *pg.Conn, condition string, params ...interface{}) (count int, err error) {
	if len(condition) < 1 {
		count, err = conn.Model(&model.Notice{}).Count()
	} else {
		count, err = conn.Model(&model.Notice{}).Where(condition, params...).Count()
	}
	return
}

// if find nothing, return pg.ErrNoRows
func (c *NoticeDao) QueryPage(conn *pg.Conn, pageSize int, page int, condition string, params ...interface{}) (ms []*model.Notice, err error) {
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

// if find nothing, return pg.ErrNoRows
func (c *NoticeDao) QueryOne(conn *pg.Conn, condition string, params ...interface{}) (m *model.Notice, err error) {
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

//if insert nothing, then return pg.ErrNoRows
func (c *NoticeDao) Insert(conn *pg.Conn, m *model.Notice) (err error) {
	if len(m.Id) < 1 {
		m.Id = uuid.GetUuid()
	}
	m.CreateTime = time.Now().Unix()
	m.UpdateTime = m.CreateTime
	m.Version = 1
	err = conn.Insert(m)
	return
}

//if insert nothing, then return pg.ErrNoRows
func (c *NoticeDao) InsertReturn(conn *pg.Conn, m *model.Notice) (mnew *model.Notice, err error) {
	if len(m.Id) < 1 {
		m.Id = uuid.GetUuid()
	}
	m.CreateTime = time.Now().Unix()
	m.UpdateTime = m.CreateTime
	m.Version = 1

	mnew = &model.Notice{}
	_, err = conn.Model(m).Returning("*").Insert(mnew)
	if err != nil {
		mnew = nil
	}
	return
}

//if update nothing, then return pg.ErrNoRows
func (c *NoticeDao) Upsert(conn *pg.Conn, m *model.Notice) (err error) {
	m.UpdateTime = time.Now().Unix()
	if len(m.Id) < 1 {
		m.Id = uuid.GetUuid()
		m.CreateTime = m.UpdateTime
	} else if m.CreateTime == 0 {
		m.CreateTime = m.UpdateTime
	}
	if m.Version == 0 {
		m.Version = 1
	}
	om := conn.Model(m).OnConflict("(id) DO UPDATE").Where("Notice."+model.Notice_Version+" = ?", m.Version)
	for _, it := range m.ToUpsertSet() {
		om.Set(it)
	}
	res, err := om.Insert()
	if res.RowsAffected() == 0 {
		err = pg.ErrNoRows
	}
	return err
}

//if update nothing, then return pg.ErrNoRows
func (c *NoticeDao) UpsertReturn(conn *pg.Conn, m *model.Notice) (mnew *model.Notice, err error) {
	m.UpdateTime = time.Now().Unix()
	if len(m.Id) < 1 {
		m.Id = uuid.GetUuid()
		m.CreateTime = m.UpdateTime
	} else if m.CreateTime == 0 {
		m.CreateTime = m.UpdateTime
	}
	if m.Version == 0 {
		m.Version = 1
	}

	om := conn.Model(m).OnConflict("(id) DO UPDATE").Where("Notice."+model.Notice_Version+" = ?", m.Version)
	for _, it := range m.ToUpsertSet() {
		om.Set(it)
	}
	mnew = &model.Notice{}
	_, err = om.Returning("*").Insert(mnew)
	if err != nil {
		mnew = nil
	}
	return
}

//if update nothing, then return pg.ErrNoRows
func (c *NoticeDao) Update(conn *pg.Conn, m *model.Notice) (err error) {
	m.UpdateTime = time.Now().Unix()
	m.Version++
	//err = conn.Update(m)
	res, err := conn.Model(m).Where(model.Notice_Id+" = ? and "+model.Notice_Version+" = ?", m.Id, m.Version-1).Update()
	if res.RowsAffected() == 0 {
		err = pg.ErrNoRows
	}
	return
}

//if update nothing, then return pg.ErrNoRows
func (c *NoticeDao) UpdateReturn(conn *pg.Conn, m *model.Notice) (mnew *model.Notice, err error) {
	m.UpdateTime = time.Now().Unix()
	m.Version++
	mnew = &model.Notice{}
	res, err := conn.Model(m).Where(model.Notice_Id+" = ? and "+model.Notice_Version+" = ?", m.Id, m.Version-1).Returning("*").Update(mnew)
	if err != nil {
		mnew = nil
	}
	if res.RowsAffected() == 0 {
		err = pg.ErrNoRows
	}
	return
}

//if delete nothing, then return pg.ErrNoRows
func (c *NoticeDao) Delete(conn *pg.Conn, m *model.Notice) (err error) {
	err = conn.Delete(m)
	return
}

//if delete nothing, then return pg.ErrNoRows
func (c *NoticeDao) DeleteById(conn *pg.Conn, id string) (err error) {
	_, err = conn.Model((*model.Notice)(nil)).Where(model.Notice_Id+" = ?", id).Delete()
	return
}

//if delete nothing, then return pg.ErrNoRows
func (c *NoticeDao) DeleteByIds(conn *pg.Conn, ids []string, oneMax int) (err error) {
	m := (*model.Notice)(nil)
	max := oneMax
	times := len(ids) / max
	for i := 1; i < times; i++ {
		oneIds := ids[(i-1)*max : i*max-1]
		_, err = conn.Model(m).Where(model.Notice_Id+" in (?)", pg.In(oneIds)).Delete()
		if err != nil {
			return
		}
	}

	if max*times < len(ids) {
		oneIds := ids[max*times:]
		_, err = conn.Model(m).Where(model.Notice_Id+" in (?)", pg.In(oneIds)).Delete()
	}
	return
}

//if delete nothing, then return pg.ErrNoRows
func (c *NoticeDao) DeleteReturn(conn *pg.Conn, m *model.Notice) (mnew *model.Notice, err error) {
	mnew = &model.Notice{}
	_, err = conn.Model(m).WherePK().Returning("*").Delete(mnew)
	if err != nil {
		mnew = nil
	}
	return
}
