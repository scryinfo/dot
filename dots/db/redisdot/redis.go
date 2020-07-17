package redisdot

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/scryinfo/dot/dot"
	"go.uber.org/zap"
	"regexp"
	"strings"
	"sync"
)

const RedisTypeId = "0ae35550-7e37-4afe-866e-b129099759b7"

type configRedis struct {
	Addr           string `json:"addr"`
	KeepVersionNum int    `json:"keepVersionNum"`
}
type RedisClient struct {
	conf configRedis

	subscribe *redis.PubSub
	currentVersion sync.Map
	*redis.Client
}

// todo：缓存安全
// todo：思考服务器端redis配置，包括但不限于分配多大的内存、内存占用已满时的策略（报错或删除一部分内存，一共6种策略的那个）

// GetVersion get current version by key, return version and error.
func (c *RedisClient) GetVersion(key string) (string, error) {
	// get current version by key,
	versionI, ok := c.currentVersion.Load(key)
	if !ok {
		return "", errors.New(fmt.Sprintf("key: %s is not exist", key))
	}

	version, ok := versionI.(string)
	if !ok {
		return "", errors.New(fmt.Sprintf("key: %s is not string type", key))
	}

	return version, nil
}

func (c *RedisClient) SetVersion(key, version string) error {
	// tx: current version / all versions
	pipe := c.TxPipeline()
	pipe.Set(c.Context(), addCVPrefix(key), version, 0)
	pipe.LPush(c.Context(), addAVsLPrefix(key), version)
	pipe.Publish(c.Context(), VersionControlChannelName, MarshalKeyAndVersion(key, version))

	_, err := pipe.Exec(c.Context())
	if err != nil {
		dot.Logger().Errorln("redis set version failed", zap.NamedError("error", err))
		return err
	}

	c.currentVersion.Store(key, version)

	// limit keep versions' length
	// set version 函数本身重点不在于维护历史版本list，某一次维护失败对业务没有影响，
	// 所以维护操作没有加入事务中，且出错也只是打印日志，而不认为函数执行错误
	length, err := c.LLen(c.Context(), addAVsLPrefix(key)).Result()
	if err != nil {
		dot.Logger().Errorln("redis get list.length failed",
			zap.NamedError("error", err),
			zap.String("list name", key))
		return nil
	}
	if length > int64(c.conf.KeepVersionNum) {
		for i := 0; i < int(length)-c.conf.KeepVersionNum; i++ {
			if err = c.RPop(c.Context(), addAVsLPrefix(key)).Err(); err != nil {
				dot.Logger().Errorln("redis.rPop failed", zap.NamedError("error", err))
			}
		}
	}

	return nil
}

func addCVPrefix(key string) string {
	return CurrentVersionPrefix +key
}

func addAVsLPrefix(key string) string {
	return AllVersionsListPrefix +key
}

func MarshalKeyAndVersion(key, version string) string {
	return fmt.Sprintf(KeyWithVersionTemplate, version,  key)
}

func UnmarshalKeyWithVersion(keyWithVersion string) (string, string, error) {
	pattern, err := regexp.Compile(KeyWithVersionTemplateRE)
	if err != nil {
		dot.Logger().Errorln("make pattern failed", zap.NamedError("error", err))
		return "", "", err
	}

	strs := pattern.FindStringSubmatch(keyWithVersion)
	if len(strs) != 3 {
		dot.Logger().Errorln("regular experience match failed",
			zap.String("str", keyWithVersion),
			zap.String("pattern", pattern.String()),
			zap.Strings("res", strs))
		return "", "", errors.New("regular experience match failed")
	}

	version := strs[1]
	key := strs[2]

	return key, version, nil
}

// DeleteVersion del versions in key
func (c *RedisClient) DeleteVersion(key string, versions ...string) error {
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

func (c *RedisClient) GetAllVersions(key string) (string, error) {
	versions, err := c.LRange(c.Context(), addAVsLPrefix(key), 0, -1).Result()
	if err != nil {
		dot.Logger().Errorln("redis get all versions failed", zap.NamedError("error", err))
		return "", err
	}

	return strings.Join(versions, " <- "), nil
}

func (c *RedisClient) Create(_ dot.Line) error {
	c.Client = redis.NewClient(&redis.Options{
		Addr: c.conf.Addr,
	})

	c.subscribe = c.Subscribe(c.Context(), VersionControlChannelName)

	_, err := c.subscribe.Receive(c.Context())
	if err != nil {
		dot.Logger().Errorln("subscription redis failed", zap.NamedError("error", err))
	}
	go func() {
		ch := c.subscribe.Channel()
		for msg := range ch {
			key, version, err := UnmarshalKeyWithVersion(msg.Payload)
			if err != nil {
				dot.Logger().Errorln("unmarshal key with version failed",
					zap.NamedError("error", err),
					zap.String("key with version", msg.Payload))
				continue
			}

			c.currentVersion.Store(key, version)
		}
	}()

	return nil
}

func (c *RedisClient) AfterAllIDestroy(_ dot.Line) {
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

	d := &RedisClient{conf: *dconf}

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

// GenerateRedis func is for unit test and example
func GenerateRedis(conf string) *RedisClient {
	res := &RedisClient{conf: configRedis{}}
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
