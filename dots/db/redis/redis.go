package redis

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/scryinfo/dot/dot"
)

const RedisTypeId = "0ae35550-7e37-4afe-866e-b129099759b7"

type configRedis struct {
	Addr     string `json:"addr"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database int    `json:"database"`
}
type Redis struct {
	conf configRedis

	DB *redis.Client
}

// todo：支持带描述的key（多级key），如:(k-v)  product:weddingRing:[index] - [product details]
// todo?：考虑封装常用方法，例如让调用者不用输入ctx、将key过期设置为可选项并附加随机值等
// todo：封装常用方法组合，例如 缓存查询不到 与 查询数据库后更新缓存 可以放在一起
// todo：为 缓存优先 的业务提供统一的管理方法
// todo?：考虑修改包名和结构体名

// todo：缓存安全

func (c *Redis) Create(_ dot.Line) error {
	c.DB = redis.NewClient(&redis.Options{
		Addr:     c.conf.Addr,
		Username: c.conf.User,
		Password: c.conf.Password,
		DB:       c.conf.Database,
	})

	return nil
}

func (c *Redis) AfterAllIDestroy(_ dot.Line) {
	if c.DB != nil {
		_ = c.DB.Close()
		c.DB = nil
	}
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
