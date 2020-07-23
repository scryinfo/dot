package gorms

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/scryinfo/dot/dot"
	"go.uber.org/zap"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const (
	//GormTypeID type id
	GormTypeID = "d2b575cd-e38f-4002-b4bd-9dc85fe13fe6"
)

type config struct {
	//sample:  "user:password@/dbname?charset=utf8&parseTime=True&loc=Local"
	//see https://github.com/go-sql-driver/mysql#parameters
	DBParameters string `json:"dbParameters"`
	ShowSQL      bool   `json:"showSql"` //是否显示sql
	Dialect      string `json:"dialect"`
}

type Gorms struct {
	conf config
	Db   *gorm.DB
}

func (c *Gorms) Destroy(ignore bool) error {
	if c.Db != nil {
		c.Db.Close()
		c.Db = nil
	}
	return nil
}

func (c *Gorms) Create(l dot.Line) (err error) {
	logger := dot.Logger()
	if len(c.conf.Dialect) < 1 {
		err = errors.New("not in (sqlite3 mysql postgres)")
		logger.Errorln("", zap.Error(err))
		return err
	}
	c.Db, err = gorm.Open(c.conf.Dialect, c.conf.DBParameters)
	if err != nil {
		logger.Errorln("Gorms", zap.Error(err))
		c.Db = nil
	} else {
		if c.conf.ShowSQL {
			c.Db.LogMode(c.conf.ShowSQL)
		}
		c.Db.SingularTable(true) //不使用表名复数

		if l != nil {
			l.ToInjecter().ReplaceOrAddByType(c.Db)
		}
	}
	return err
}

func newGorms(conf []byte) (d dot.Dot, err error) {
	dconf := &config{}
	err = dot.UnMarshalConfig(conf, dconf)
	if err != nil {
		return nil, err
	}

	d = &Gorms{conf: *dconf}
	return d, err
}

//GormsTypeLives return type lives
func GormsTypeLives() []*dot.TypeLives {
	lives := []*dot.TypeLives{
		{
			Meta: dot.Metadata{TypeID: GormTypeID, NewDoter: func(conf []byte) (dot dot.Dot, err error) {
				return newGorms(conf)
			}},
		},
	}
	return lives
}

//GormsConfigTypeLive return config
func GormsConfigTypeLive() *dot.ConfigTypeLive {
	return &dot.ConfigTypeLive{
		TypeIDConfig: GormTypeID,
		ConfigInfo:   &config{},
	}
}

//NewGormsTest just for test
func NewGormsTest(dbParameters string, showSQL bool, dialect string) *Gorms {
	conf := config{
		DBParameters: dbParameters,
		ShowSQL:      showSQL,
		Dialect:      dialect,
	}
	bs, _ := json.Marshal(conf)
	d, err := newGorms(bs)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	db := d.(*Gorms)
	err = db.Create(nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return db
}
