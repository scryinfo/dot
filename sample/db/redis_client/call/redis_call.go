package call

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/db/redis_client"
)

const RedisDemoTypeId = "c5f966e2-147a-4b09-a5dc-8c74ff603d38"

type RedisDemo struct {
	Redis *redis_client.RedisClient `dot:""`
}

func (c *RedisDemo) Start(_ bool) (err error) {
	c.Redis.ClientV8().FlushAll(c.Redis.ClientV8().Context())

	go func() {
		c.basicDemo()
		c.versionControlDemo()
	}()

	return nil
}

func (c *RedisDemo) basicDemo() {
	// simulate query in cache first (no result)
	v1, err := c.Redis.ClientV8().Get(c.Redis.ClientV8().Context(), basicDemoKey).Result()
	if err != redis.Nil {
		fmt.Println("Example: get value not run as suppose, error:", err)
		os.Exit(-1)
	}

	// skip query in db, only simulate update cache
	checkError("Example: set value failed", c.Redis.ClientV8().Set(c.Redis.ClientV8().Context(), basicDemoKey, basicDemoValue, 0).Err())

	// suppose a request comes now, query in cache (has result)
	v2, err := c.Redis.ClientV8().Get(c.Redis.ClientV8().Context(), basicDemoKey).Result()
	checkError("Example: get value failed", err)

	time.Sleep(time.Second)
	fmt.Printf("v1: , v2: %s| <- suppose |\n", basicDemoValue)
	fmt.Printf("v1: %s, v2: %s\n", v1, v2)

	fmt.Println("Node: Basic demo passed.")
	fmt.Println("-------")

	return
}

func (c *RedisDemo) versionControlDemo() {
	// register
	checkError("Example: set version failed", c.Redis.SetVersion(versionControlDemoKey, versionControlDemoVersionOne))

	// wait for subscribe goroutine
	time.Sleep(time.Second)

	// get version
	v, err := c.Redis.GetVersion(versionControlDemoKey)
	checkError("Example: get version failed", err)
	fmt.Printf("Version get suppose: %s, actually: %s\n", versionControlDemoVersionOne, v)

	// set(update) version
	checkError("Example: set version failed", c.Redis.SetVersion(versionControlDemoKey, versionControlDemoVersionTwo))

	// get version
	v, err = c.Redis.GetVersion(versionControlDemoKey)
	checkError("Example: get version failed", err)
	fmt.Printf("Version get suppose: %s, actually: %s\n", versionControlDemoVersionTwo, v)

	// set version 3 times
	for i := 3; i < 6; i++ {
		checkError("Example: set version failed", c.Redis.SetVersion(versionControlDemoKey, strconv.Itoa(i)))
	}

	// get all versions
	versions, err := c.Redis.GetAllVersions(versionControlDemoKey)
	checkError("Example: get all versions failed", err)

	fmt.Println("All versions: 5 <- 4 <- 3 <- 2| <- suppose |")
	fmt.Println("All versions:", versions)

	fmt.Println("Node: version control demo passed.")
	fmt.Println("-------")

	return
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
				RelyLives: map[string]dot.LiveId{"RedisClient": redis_client.RedisClientTypeId},
			},
		},
	}

	lives := []*dot.TypeLives{tl}
	lives = append(lives, redis_client.RedisClientTypeLives()...)

	return lives
}
