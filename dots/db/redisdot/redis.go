package redisdot

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/scryinfo/dot/dot"
	"go.uber.org/zap"
	"strings"
)

const RedisTypeId = "0ae35550-7e37-4afe-866e-b129099759b7"

type configRedis struct {
	Addr         string `json:"addr"`
	KeepVersions int    `json:"keepVersions"`
}
type Redis struct {
	conf configRedis

	currentVersion map[string]int
	*redis.Client
}

// todo：缓存安全
// todo：思考服务器端redis配置，包括但不限于分配多大的内存、内存占用已满时的策略（报错或删除一部分内存，一共6种策略的那个）

// GetVersion get current version by key, return version and error.
// If the key is not exist and 'registerNX' is given (registerNX: set not exist), we will register the key for you.
func (c *Redis) GetVersion(key string, registerNX ...int) (int, error) {
	// get current version by key,
	version, ok := c.currentVersion[key]
	if ok {
		return version, nil
	}

	// if registerNX is given, call set version
	if len(registerNX) != 1 {
		return GetVersionNotExistFlag, errors.New(fmt.Sprintf("key: %s is not exist", key))
	}

	version = registerNX[0]
	return GetVersionNotExistAndRegisterFlag, c.SetVersion(key, version)
}

func (c *Redis) SetVersion(key string, version int) error {
	// tx: current version / all versions
	pipe := c.TxPipeline()
	pipe.Set(c.Context(), key+RegisteredKeySuffix, version, 0)
	pipe.LPush(c.Context(), key+RegisteredKeysListSuffix, version)

	_, err := pipe.Exec(c.Context())
	if err != nil {
		dot.Logger().Errorln("redis set version failed", zap.NamedError("error", err))
		return err
	}

	c.currentVersion[key] = version

	// limit keep versions' length
	// set version 函数本身重点不在于维护历史版本list，某一次维护失败对业务没有影响，
	// 所以维护操作没有加入事务中，且出错也只是打印日志，而不认为函数执行错误
	length, err := c.LLen(c.Context(), key+RegisteredKeysListSuffix).Result()
	if err != nil {
		dot.Logger().Errorln("redis get list.length failed",
			zap.NamedError("error", err),
			zap.String("list name", key))
		return nil
	}
	if length > int64(c.conf.KeepVersions) {
		for i := 0; i < int(length)-c.conf.KeepVersions; i++ {
			if err = c.RPop(c.Context(), key+RegisteredKeysListSuffix).Err(); err != nil {
				dot.Logger().Errorln("redis.rPop failed", zap.NamedError("error", err))
			}
		}
	}

	return nil
}

func MarshalKeyAndVersion(key string, version int) string {
	return fmt.Sprintf(KeyWithVersionTemplate, key, version)
}

func UnmarshalKeyWithVersion(keyWithVersion string) (key string, version int, err error) {
	n, err := fmt.Sscanf(keyWithVersion, KeyWithVersionTemplate, key, version)
	if n != 2 || err != nil {
		dot.Logger().Errorln("redis unmarshal key with version failed",
			zap.NamedError("error", err),
			zap.Int("unmarshal variables, expect 2, actually:", n))
		return "", 0, err
	}

	return key, version, nil
}

// DeleteVersion del versions in key
func (c *Redis) DeleteVersion(key string, versions ...int) error {
	// re-curse del target versions
	for i := range versions {
		if err := c.LRem(c.Context(), key, 1, versions[i]).Err(); err != nil {
			dot.Logger().Errorln("redis del versions failed",
				zap.NamedError("error", err),
				zap.Int("index of versions slice", i))
			return err
		}
	}

	return nil
}

func (c *Redis) GetAllVersions(key string) (string, error) {
	versions, err := c.LRange(c.Context(), key+RegisteredKeysListSuffix, 0, -1).Result()
	if err != nil {
		dot.Logger().Errorln("redis get all versions failed", zap.NamedError("error", err))
		return "", err
	}

	return strings.Join(versions, " <- "), nil
}

func (c *Redis) Create(_ dot.Line) error {
	c.Client = redis.NewClient(&redis.Options{
		Addr: c.conf.Addr,
	})

	c.currentVersion = make(map[string]int)

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
