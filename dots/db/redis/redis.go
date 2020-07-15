package redis

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/scryinfo/dot/dot"
)

const RedisTypeId = "0ae35550-7e37-4afe-866e-b129099759b7"

type configRedis struct {
	Addr string `json:"addr"`
}
type Redis struct {
	conf configRedis

	*redis.Client
}

// todo?：考虑修改包名和结构体名，但是不知道改成啥
// todo：+版本管理，思考：所有数据都要管理吗？还是只有复杂数据需要管理？

// todo：缓存安全
// todo：思考服务器端redis配置，包括但不限于分配多大的内存、内存占用已满时的策略（报错或删除一部分内存，一共6种策略的那个）

func (c *Redis) Create(_ dot.Line) error {
	c.Client = redis.NewClient(&redis.Options{
		Addr:     c.conf.Addr,
	})

	return nil
}

func (c *Redis) AfterAllIDestroy(_ dot.Line) {
	if c.Client != nil {
		_ = c.Client.Close()
		c.Client = nil
	}

	return
}

//construct dot
func newRedis(conf []byte) (dot.Dot, error) {
	dconf := &configRedis{}

	err := dot.UnMarshalConfig(conf, dconf)
	if err != nil {
		return nil, err
	}

	d := &Redis{conf: *dconf}

	return d, nil
}

//RedisTypeLives
func RedisTypeLives() []*dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{TypeId: RedisTypeId, NewDoter: func(conf []byte) (dot.Dot, error) {
			return newRedis(conf)
		}},
	}

	lives := []*dot.TypeLives{tl}

	return lives
}

//RedisConfigTypeLive
func RedisConfigTypeLive() *dot.ConfigTypeLives {
	paths := make([]string, 0)
	paths = append(paths, "")
	return &dot.ConfigTypeLives{
		TypeIdConfig: RedisTypeId,
		ConfigInfo:   &configRedis{
			//todo
		},
	}
}

// GenerateRedis func is for unit test and example
func GenerateRedis(conf string) *Redis {
	res := &Redis{conf: configRedis{}}
	if err := json.Unmarshal([]byte(conf), &res.conf); err != nil {
		fmt.Println("Test: json unmarshal failed, error:", err)
		return nil
	}
	if err := res.Create(nil); err != nil {
		fmt.Println("Test: res.create failed, error:", err)
		return nil
	}

	return res
}
