package pgs

import (
	"context"
	"encoding/json"
	"github.com/go-pg/pg/v9"
	"github.com/scryinfo/dot/dot"
)

const (
	ConnWrapperTypeId = "ffc08507-dd5f-456c-84ea-cdae00b220bf"
)

type config struct { //todo 连接池等配置
	Addr     string `json:"addr"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	ShowSql  bool   `json:"showSql"`
}

type ConnWrapper struct {
	db   *pg.DB
	conf config
}

func (c *ConnWrapper) Create(l dot.Line) error {
	c.db = pg.Connect(&pg.Options{
		Addr:     c.conf.Addr,
		User:     c.conf.User,
		Password: c.conf.Password,
		Database: c.conf.Database,
	})
	if c.conf.ShowSql {
		c.db.AddQueryHook(pgLogger{})
	}
	return nil
}

func (c *ConnWrapper) AfterAllIDestroy(l dot.Line) {
	if c.db != nil {
		_ = c.db.Close()
		c.db = nil
	}
}

func (c *ConnWrapper) GetDb() *pg.DB {
	return c.db
}

func (c *ConnWrapper) TestConn() bool {
	n := -1
	_, _ = c.db.QueryOne(pg.Scan(&n), "select 1")
	return n == 1
}

//construct dot
func newConnWrapper(conf []byte) (dot.Dot, error) {
	dconf := &config{}
	err := dot.UnMarshalConfig(conf, dconf)
	if err != nil {
		return nil, err
	}

	d := &ConnWrapper{conf: *dconf}

	return d, err
}

//ConnWrapperTypeLives make all type lives
func ConnWrapperTypeLives() []*dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{TypeId: ConnWrapperTypeId, NewDoter: func(conf []byte) (dot.Dot, error) {
			return newConnWrapper(conf)
		}},
	}

	lives := []*dot.TypeLives{tl}

	return lives
}

//GenarateConnWrapper this func is for test
func GenarateConnWrapper(conf string) *ConnWrapper {
	conn := &ConnWrapper{}
	_ = json.Unmarshal([]byte(conf), &conn.conf)
	_ = conn.Create(nil)
	return conn
}

//GenarateConnWrapper this func is for test
func GenarateConnWrapperByDb(db *pg.DB) *ConnWrapper {
	conn := &ConnWrapper{db, config{}}
	return conn
}

type pgLogger struct{}

func (d pgLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d pgLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	dot.Logger().Debug(func() string {
		temp, _ := q.FormattedQuery()
		return temp
	})
	return nil
}
