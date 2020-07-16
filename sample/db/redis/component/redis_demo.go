package component

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/db/redisdot"
	"go.uber.org/zap"
	"os"
	"time"
)

const RedisDemoTypeId = "c5f966e2-147a-4b09-a5dc-8c74ff603d38"

type RedisDemo struct {
	Redis *redisdot.Redis `dot:""`
}

func (c *RedisDemo) Start(_ bool) (err error) {
	c.Redis.FlushAll(c.Redis.Context())

	go func() {
		c.basicDemo()
		c.versionControlDemo()
		c.versionControlExceptionProcessDemo()
	}()

	return nil
}

func (c *RedisDemo) basicDemo() {
	// simulate query in cache first (no result)
	v1, err := c.Redis.Get(c.Redis.Context(), basicDemoKey).Result()
	if err != redis.Nil {
		fmt.Println("Example: get value not run as suppose, error:", err)
		os.Exit(-1)
	}

	// skip query in db, only simulate update cache
	checkError("Example: set value failed", c.Redis.Set(c.Redis.Context(), basicDemoKey, basicDemoValue, 0).Err())

	// suppose a request comes now, query in cache (has result)
	v2, err := c.Redis.Get(c.Redis.Context(), basicDemoKey).Result()
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

	// get version
	v, err := c.Redis.GetVersion(versionControlDemoKey)
	checkError("Example: get version failed", err)
	fmt.Printf("Version get suppose: %d, actually: %d\n", versionControlDemoVersionOne, v)

	// set(update) version
	checkError("Example: set version failed", c.Redis.SetVersion(versionControlDemoKey, versionControlDemoVersionTwo))

	// get version
	v, err = c.Redis.GetVersion(versionControlDemoKey)
	checkError("Example: get version failed", err)
	fmt.Printf("Version get suppose: %d, actually: %d\n", versionControlDemoVersionTwo, v)

	// set version 3 times
	for i := 3; i < 6; i++ {
		checkError("Example: set version failed", c.Redis.SetVersion(versionControlDemoKey, i))
	}

	// get all versions
	versions, err := c.Redis.GetAllVersions(versionControlDemoKey)
	checkError("Example: get all versions failed", err)
	fmt.Println("All versions:", versions)

	fmt.Println("Node: version control demo passed.")
	fmt.Println("-------")

	return
}

// 异常流程：没有注册时，获取版本号；获取不存在的版本号；获取不存在的版本号并注册
func (c *RedisDemo) versionControlExceptionProcessDemo() {
	// get version not exist and not register
	flag, err := c.Redis.GetVersion(versionControlExceptionProcessDemoKey)
	if flag != redisdot.GetVersionNotExistFlag {
		fmt.Println("Example: get version not exist and not register failed, error:", err)
		os.Exit(-1)
	}

	// get version not exist and register
	flag, err = c.Redis.GetVersion(versionControlExceptionProcessDemoKey, versionControlExceptionProcessDemoVersion)
	if flag != redisdot.GetVersionNotExistAndRegisterFlag || err != nil {
		fmt.Println("Example: get version not exist and register failed, error:", err)
		os.Exit(-1)
	}

	// get version authenticate
	version, err := c.Redis.GetVersion(versionControlExceptionProcessDemoKey)
	if version != versionControlExceptionProcessDemoVersion || err != nil {
		fmt.Println("Example: register by calling 'GetValue' failed, error:", err)
		os.Exit(-1)
	}

	fmt.Println("Node: version control exception process demo passed.")
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
				RelyLives: map[string]dot.LiveId{"Redis": redisdot.RedisTypeId},
			},
		},
	}

	lives := []*dot.TypeLives{tl}
	lives = append(lives, redisdot.RedisTypeLives()...)

	return lives
}
