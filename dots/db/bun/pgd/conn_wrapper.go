package pgd

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/scryinfo/dot/dot"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"go.uber.org/zap"
)

const (
	//ConnWrapperTypeID type id
	ConnWrapperTypeID = "ffc08507-dd5f-456c-84ea-cdbf00b220b0"
)

type config struct {
	Host     dot.StringFromEnv `json:"host" yaml:"host" mapstructure:"host"`
	Port     dot.StringFromEnv `json:"port" yaml:"port" mapstructure:"port"`
	User     dot.StringFromEnv `json:"user" yaml:"user" mapstructure:"user"`
	Password dot.StringFromEnv `json:"password" yaml:"password" mapstructure:"password"`
	Database dot.StringFromEnv `json:"database" yaml:"database" mapstructure:"database"`
	ShowSQL  bool              `json:"showSql" yaml:"showSql" mapstructure:"showSql"`
}

// ConnWrapper connect wrapper
type ConnWrapper struct {
	db   *bun.DB
	conf config
}

func (c *ConnWrapper) Create(dot.Line) error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?&timeout=5&sslmode=disable",
		string(c.conf.User),
		string(c.conf.Password),
		string(c.conf.Host),
		string(c.conf.Port),
		string(c.conf.Database),
	)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New(), bun.WithDiscardUnknownColumns())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(c.conf.ShowSQL)))
	c.db = db

	return nil
}

func (c *ConnWrapper) AfterAllDestroy(dot.Line) {
	if c.db != nil {
		_ = c.db.Close()
		c.db = nil
	}
}

// GetDb get db
func (c *ConnWrapper) GetDb() *bun.DB {
	return c.db
}

func (c *ConnWrapper) RunInTx(task func(db bun.IDB) error) error {
	var err error
	if task != nil {
		err = c.db.RunInTx(context.TODO(), nil, func(ctx context.Context, tx bun.Tx) error {
			err = task(tx)
			return err
		})
	}
	return err
}

func (c *ConnWrapper) RunInNoTx(task func(db bun.IDB) error) error {
	var err error
	if task != nil {
		err = task(c.db)
	}

	return err
}

// TestConn test the connect
func (c *ConnWrapper) TestConn() bool {
	n := -1
	c.db.NewSelect().ColumnExpr("1 AS n").Scan(context.TODO(), &n)
	return n == 1
}

// construct dot
func newConnWrapper(conf []byte) (dot.Dot, error) {
	dconf := &config{}

	dot.Logger().Infoln(fmt.Sprintf("database connect conf before UnMarshall %+v", string(conf)))

	err := dot.UnMarshalConfig(conf, dconf)
	if err != nil {
		dot.Logger().Errorln("database connect conf ", zap.Error(err))
		return nil, err
	}
	d := &ConnWrapper{conf: *dconf}
	return d, err
}

// ConnWrapperTypeLives make all type lives
func ConnWrapperTypeLives() []*dot.TypeLives {
	return []*dot.TypeLives{{
		Meta: dot.Metadata{TypeID: ConnWrapperTypeID, NewDoter: newConnWrapper},
	}}
}

// GenerateConnWrapper this func is for test
func GenerateConnWrapper(conf string) *ConnWrapper {
	conn := &ConnWrapper{}
	err := json.Unmarshal([]byte(conf), &conn.conf)
	if err != nil {
		dot.Logger().Errorln("database connect conf ", zap.Error(err))
		return nil
	}
	_ = conn.Create(nil)
	return conn
}

// GenerateConnWrapperByDb this func is for test
func GenerateConnWrapperByDb(db *bun.DB) *ConnWrapper {
	conn := &ConnWrapper{db, config{}}
	return conn
}

type pgLogger struct{}

func (d pgLogger) BeforeQuery(c context.Context, _ *pgdriver.Listener) (context.Context, error) {
	return c, nil
}

func (d pgLogger) AfterQuery(_ context.Context, q *bun.DB) error {
	dot.Logger().Debug(func() string {
		// q.Formatter().FormatQuery()
		return ""
	})
	return nil
}
