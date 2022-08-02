package pgs

import (
	"context"
	"encoding/json"
	"github.com/go-pg/pg/v10"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/utils"
)

const (
	//ConnWrapperTypeID type id
	ConnWrapperTypeID = "ffc08507-dd5f-456c-84ea-cdae00b220bf"
)

type config struct {
	Addr     string `json:"addr"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	ShowSQL  bool   `json:"showSql"`
}

//ConnWrapper connect wrapper
type ConnWrapper struct {
	db   *pg.DB
	conf config
}

func (c *ConnWrapper) Create(dot.Line) error {
	c.db = pg.Connect(&pg.Options{
		Addr:     c.conf.Addr,
		User:     c.conf.User,
		Password: c.conf.Password,
		Database: c.conf.Database,
	})
	if c.conf.ShowSQL {
		c.db.AddQueryHook(pgLogger{})
	}
	return nil
}

func (c *ConnWrapper) AfterAllDestroy(dot.Line) {
	if c.db != nil {
		_ = c.db.Close()
		c.db = nil
	}
}

//GetDb get db
func (c *ConnWrapper) GetDb() *pg.DB {
	return c.db
}

//TestConn test the connect
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
	dconf.Database = utils.UnmarshalENV(dconf.Database)
	dconf.Addr = utils.UnmarshalENV(dconf.Addr)
	dconf.User = utils.UnmarshalENV(dconf.User)
	dconf.Password = utils.UnmarshalENV(dconf.Password)
	d := &ConnWrapper{conf: *dconf}
	return d, err
}

//ConnWrapperTypeLives make all type lives
func ConnWrapperTypeLives() []*dot.TypeLives {
	return []*dot.TypeLives{{
		Meta: dot.Metadata{TypeID: ConnWrapperTypeID, NewDoter: newConnWrapper},
	}}
}

//GenerateConnWrapper this func is for test
func GenerateConnWrapper(conf string) *ConnWrapper {
	conn := &ConnWrapper{}
	_ = json.Unmarshal([]byte(conf), &conn.conf)
	_ = conn.Create(nil)
	return conn
}

//GenerateConnWrapperByDb this func is for test
func GenerateConnWrapperByDb(db *pg.DB) *ConnWrapper {
	conn := &ConnWrapper{db, config{}}
	return conn
}

type pgLogger struct{}

func (d pgLogger) BeforeQuery(c context.Context, _ *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d pgLogger) AfterQuery(_ context.Context, q *pg.QueryEvent) error {
	dot.Logger().Debug(func() string {
		temp, _ := q.FormattedQuery()
		return string(temp)
	})
	return nil
}
