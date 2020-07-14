package component

import (
	"fmt"
	"github.com/scryinfo/dot/dot"
	dot_redis "github.com/scryinfo/dot/dots/db/redis"
	"go.uber.org/zap"
	"os"
	"strconv"
)

const RedisDemoTypeId = "c5f966e2-147a-4b09-a5dc-8c74ff603d38"

type RedisDemo struct {
	Collections []*dot_redis.Collection

	Redis       *dot_redis.Redis `dot:""`
}

type Demo struct {
	Str                string `redis:"key_str"`
	I                  int    `redis:"key_int"`
	dot_redis.ModelImp `redis:"-"`
}

func (c *RedisDemo) AfterAllInject(_ dot.Line) {
	demoInsWithName := &Demo{}
	demoInsWithName.SetModelName("demo")

	c.Collections = make([]*dot_redis.Collection, 0)
	c.Collections = append(c.Collections, c.Redis.RegisterCollections([]dot_redis.Model{demoInsWithName})...)
}

func (c *RedisDemo) Start(_ bool) (err error) {
	dot.Logger().Infoln("-------")

	// transplant orm demo here, use component format
	if len(c.Collections) != 1 {
		dot.Logger().Errorln("redis collection register failed", zap.Int("collection.length should be 1 but actually:", len(c.Collections)))
		return nil
	}

	collection := c.Collections[0]

	// clean redis
	n, err := collection.DeleteAll()
	checkError("delete all failed", err)
	dot.Logger().Infoln("del history", zap.Int("del items num", 3))

	// set some value for test
	for i := 0; i < 3; i++ {
		demoIns := &Demo{
			Str: "index: " + strconv.Itoa(i),
			I:   i + 1,
		}
		checkError("preset values for test failed, index:"+strconv.Itoa(i), collection.Save(demoIns))
	}

	// save value for following test (hereinafter called 'payload')
	demoInsSave := &Demo{
		Str: "redis call simulate",
		I:   -1,
	}
	checkError("save value failed", collection.Save(demoInsSave))

	// find payload
	demoInsFind := Demo{}
	checkError("find value failed", collection.Find(demoInsSave.Id, &demoInsFind))
	dot.Logger().Infoln("find res", zap.String("payload", fmt.Sprintf("%#v", demoInsFind)))

	// update certain fields
	demoInsSave.Str = "update"
	checkError("update certain fields failed", collection.SaveFields([]string{"Str"}, demoInsSave))

	// find certain fields, attention on 'demo.I' field
	demoInsFindFields := Demo{}
	checkError("find certain field failed", collection.FindFields(demoInsSave.Id, []string{"Str"}, &demoInsFindFields))
	dot.Logger().Infoln("find res", zap.String("certain fields", fmt.Sprintf("%#v", demoInsFindFields)))

	// del
	ok, err := collection.Delete(demoInsSave.Id)
	if !ok || err != nil {
		dot.Logger().Errorln("del id failed", zap.Bool("operation success?", ok), zap.NamedError("error", err))
		os.Exit(-1)
	}

	// find all with count
	demoInsFindAll := make([]*Demo, 0)
	checkError("find all failed", collection.FindAll(&demoInsFindAll))

	n, err = collection.Count()
	if err != nil {
		dot.Logger().Errorln("count summary failed", zap.NamedError("error", err))
		os.Exit(-1)
	}
	dot.Logger().Infoln("find all:", zap.Int("total", n))

	for i, payload := range demoInsFindAll {
		dot.Logger().Infoln("find all res, index: "+strconv.Itoa(i), zap.String("value", fmt.Sprintf("%#v", payload)))
	}

	dot.Logger().Infoln("-------")

	return nil
}

func checkError(msg string, err error) {
	if err != nil {
		dot.Logger().Errorln(msg, zap.NamedError("error", err))
		os.Exit(-1)
	}

	return
}

func (c *RedisDemo) Stop(_ bool) error {
	c.Collections = make([]*dot_redis.Collection, 0)

	return nil
}

//RedisDemoTypeLives
func RedisDemoTypeLives() []*dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{TypeId: RedisDemoTypeId, NewDoter: func(_ []byte) (dot.Dot, error) {
			return &RedisDemo{}, nil
		}},
		Lives: []dot.Live{
			{
				LiveId:    RedisDemoTypeId,
				RelyLives: map[string]dot.LiveId{"Redis": dot_redis.RedisTypeId},
			},
		},
	}

	lives := []*dot.TypeLives{tl}
	lives = append(lives, dot_redis.RedisTypeLives()...)

	return lives
}
