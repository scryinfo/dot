package dao

import (
	"github.com/go-pg/pg"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/db/pgs"
	"github.com/scryinfo/dot/sample/db/pgs/model"
)

const (
	NoticeTypeId = "9fb3c050-f116-4fb4-b41a-e5a8dbb77bf9"
)

type Notice struct {
	*pgs.DaoBase `dot:""`
}

//NoticeTypeLives make all type lives
func NoticeTypeLives() []*dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{TypeId: NoticeTypeId, NewDoter: func(conf []byte) (dot.Dot, error) {
			return &Notice{}, nil
		}},
		Lives: []dot.Live{
			dot.Live{
				LiveId:    NoticeTypeId,
				RelyLives: map[string]dot.LiveId{"DaoBase": pgs.DaoBaseTypeId},
			},
		},
	}

	lives := pgs.DaoBaseTypeLives()
	lives = append(lives, tl)

	return lives
}
func (c *Notice) Query(conn *pg.Conn, condition string, params ...interface{}) (ms []model.Notice, err error) {
	err = conn.Model(&ms).Where(condition, params...).Select()
	return
}

func (c *Notice) QueryOne(conn *pg.Conn, condition string, params ...interface{}) (m *model.Notice, err error) {
	m = &model.Notice{}
	err = conn.Model(m).Where(condition, params...).First()
	return
}

func (c *Notice) Insert(conn *pg.Conn, m *model.Notice) (err error) {
	err = conn.Insert(m)
	return
}

func (c *Notice) Upsert(conn *pg.Conn, m *model.Notice) (err error) {
	om := conn.Model(m).OnConflict("(id) DO UPDATE")
	for _, it := range m.ToUpsertSet() {
		om.Set(it)
	}
	_, err = om.Insert()
	return err
}

func (c *Notice) Delete(conn *pg.Conn, m *model.Notice) (err error) {
	err = conn.Delete(m)
	return
}
