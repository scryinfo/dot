package component

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/scryinfo/dot/dot"
	dot_redis "github.com/scryinfo/dot/dots/db/redis"
	"go.uber.org/zap"
	"os"
	"time"
)

const RedisDemoTypeId = "c5f966e2-147a-4b09-a5dc-8c74ff603d38"

type RedisDemo struct {
	Redis *dot_redis.Redis `dot:""`
}

func (c *RedisDemo) Start(_ bool) (err error) {
	c.Redis.FlushAll(c.Redis.Context())

	// simulate query in cache first (no result)
	v1, err := c.Redis.Get(c.Redis.Context(), "demo").Result()
	if err != redis.Nil {
		fmt.Println("Example: get value not run as suppose, error:", err)
		os.Exit(-1)
	}

	// skip query in db, only simulate update cache
	checkError("Example: set value failed", c.Redis.Set(c.Redis.Context(), "demo", "basic process demo", 0).Err())

	// suppose a request comes now, query in cache (has result)
	v2, err := c.Redis.Get(c.Redis.Context(), "demo").Result()
	checkError("Example: get value failed", err)

	go func() {
		time.Sleep(time.Second)
		fmt.Printf("v1: %s, v2: %s\n", v1, v2)
		fmt.Println("v1: , v2: basic process demo| <- suppose |")
	}()

	return nil
}

func checkError(msg string, err error) {
	if err != nil {
		dot.Logger().Errorln(msg, zap.NamedError("error", err))
		os.Exit(-1)
	}

	return
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
