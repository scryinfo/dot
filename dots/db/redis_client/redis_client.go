package redis_client

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/scryinfo/dot/dot"
	"go.uber.org/zap"
	"regexp"
	"strings"
	"sync"
)

const RedisClientTypeId = "0ae35550-7e37-4afe-866e-b129099759b7"

type configRedis struct {
	Addr                string `json:"addr"`
	KeepVersionNum      int    `json:"keepVersionNum"`
	GetVersionFromRedis bool   `json:"getVersionFromRedis"`
}
type RedisClient struct {
	conf configRedis

	subscribe      *redis.PubSub
	currentVersion sync.Map
	clientV8       *redis.Client
}

func (c *RedisClient) ClientV8() *redis.Client {
	return c.clientV8
}

// todo：缓存安全
// todo：思考服务器端redis配置，包括但不限于分配多大的内存、内存占用已满时的策略（报错或删除一部分内存，一共6种策略的那个）

// GetVersion get current version by key, return version and error.
func (c *RedisClient) GetVersion(key string) (string, error) {
	if c.conf.GetVersionFromRedis {
		return c.clientV8.Get(c.clientV8.Context(), addCVPrefix(key)).Result()
	}

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
	if c.conf.GetVersionFromRedis {
		return c.setVersionDirectlyRedis(key, version)
	}

	return c.setVersion(key, version)
}

func (c *RedisClient) setVersionDirectlyRedis(key, version string) error {
	// tx: current version / all versions
	pipe := c.clientV8.TxPipeline()
	pipe.Set(c.clientV8.Context(), addCVPrefix(key), version, 0)
	pipe.LPush(c.clientV8.Context(), addAVsLPrefix(key), version)

	_, err := pipe.Exec(c.clientV8.Context())
	if err != nil {
		dot.Logger().Errorln("redis set version failed", zap.NamedError("error", err))
		return err
	}

	// limit keep versions' length
	// set version 函数本身重点不在于维护历史版本list，某一次维护失败对业务没有影响，
	// 所以维护操作没有加入事务中，且出错也只是打印日志，而不认为函数执行错误
	length, err := c.clientV8.LLen(c.clientV8.Context(), addAVsLPrefix(key)).Result()
	if err != nil {
		dot.Logger().Errorln("redis get list.length failed",
			zap.NamedError("error", err),
			zap.String("list name", key))
		return nil
	}
	if length > int64(c.conf.KeepVersionNum) {
		for i := 0; i < int(length)-c.conf.KeepVersionNum; i++ {
			if err = c.clientV8.RPop(c.clientV8.Context(), addAVsLPrefix(key)).Err(); err != nil {
				dot.Logger().Errorln("redis.rPop failed", zap.NamedError("error", err))
			}
		}
	}

	return nil
}

func (c *RedisClient) setVersion(key, version string) error { //todo review 没有订阅
	// tx: current version / all versions / (publish)
	pipe := c.clientV8.TxPipeline()
	pipe.Set(c.clientV8.Context(), addCVPrefix(key), version, 0)
	pipe.LPush(c.clientV8.Context(), addAVsLPrefix(key), version)
	pipe.Publish(c.clientV8.Context(), VersionControlChannelName, MarshalKeyAndVersion(key, version))

	_, err := pipe.Exec(c.clientV8.Context())
	if err != nil {
		dot.Logger().Errorln("redis set version failed", zap.NamedError("error", err))
		return err
	}

	c.currentVersion.Store(key, version)

	// limit keep versions' length
	// set version 函数本身重点不在于维护历史版本list，某一次维护失败对业务没有影响，
	// 所以维护操作没有加入事务中，且出错也只是打印日志，而不认为函数执行错误
	length, err := c.clientV8.LLen(c.clientV8.Context(), addAVsLPrefix(key)).Result()
	if err != nil {
		dot.Logger().Errorln("redis get list.length failed",
			zap.NamedError("error", err),
			zap.String("list name", key))
		return nil
	}
	if length > int64(c.conf.KeepVersionNum) {
		for i := 0; i < int(length)-c.conf.KeepVersionNum; i++ { //todo review 使用事务， 且设定一个事务中最多执行多少条命名; 版本对应的数据怎么处理？
			if err = c.clientV8.RPop(c.clientV8.Context(), addAVsLPrefix(key)).Err(); err != nil {
				dot.Logger().Errorln("redis.rPop failed", zap.NamedError("error", err))
			}
		}
	}

	return nil
}

// DeleteVersion del versions in key
func (c *RedisClient) DeleteVersion(key string, versions ...string) error {
	// re-curse del target versions
	for i := range versions { //todo 使用事务， 且设定一个事务中最多执行多少条命名
		if err := c.clientV8.LRem(c.clientV8.Context(), key, 1, versions[i]).Err(); err != nil {
			dot.Logger().Errorln("redis del versions failed",
				zap.NamedError("error", err),
				zap.Int("index of versions slice", i))
			return err
		}
	}
	//todo review , scan命令，参见 http://doc.redisfans.com/key/scan.html
	//清理版本开头的所有key

	return nil
}

func (c *RedisClient) GetAllVersions(key string) (string, error) {
	versions, err := c.clientV8.LRange(c.clientV8.Context(), addAVsLPrefix(key), 0, -1).Result()
	if err != nil {
		dot.Logger().Errorln("redis get all versions failed", zap.NamedError("error", err))
		return "", err
	}

	return strings.Join(versions, " <- "), nil
}

func (c *RedisClient) Create(_ dot.Line) error {
	c.clientV8 = redis.NewClient(&redis.Options{
		Addr: c.conf.Addr,
	})

	// directly return if not use component to manage version
	if c.conf.GetVersionFromRedis {
		return nil
	}

	c.subscribe = c.clientV8.Subscribe(c.clientV8.Context(), VersionControlChannelName) //todo review 确认只会有第一层的 "channel"吗？ 如： v:key:id,这种key的数据会进来吗

	// wait until subscribe success
	_, err := c.subscribe.Receive(c.clientV8.Context())
	if err != nil {
		dot.Logger().Errorln("subscription redis failed", zap.NamedError("error", err))
	}

	go func() {
		ch := c.subscribe.Channel() //todo review 没有退出机制
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

func (c *RedisClient) AfterAllDestroy(_ dot.Line) {
	if c.clientV8 != nil {
		_ = c.clientV8.Close()
		c.clientV8 = nil
	}

	if c.subscribe != nil { //todo review, sub 与 client的先后关系
		_ = c.subscribe.Close()
	}

	return
}

//construct dot
func newRedisClient(conf []byte) (dot.Dot, error) {
	dconf := &configRedis{}

	err := dot.UnMarshalConfig(conf, dconf)
	if err != nil {
		return nil, err
	}

	d := &RedisClient{conf: *dconf}

	return d, nil
}

//RedisClientTypeLives
func RedisClientTypeLives() []*dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{TypeId: RedisClientTypeId, NewDoter: func(conf []byte) (dot.Dot, error) {
			return newRedisClient(conf)
		}},
	}

	lives := []*dot.TypeLives{tl}

	return lives
}

func GenerateKey(keys ...string) string {
	return strings.Join(keys, KeySplitChar)
}

func addCVPrefix(key string) string {
	return CurrentVersionPrefix + key
}

func addAVsLPrefix(key string) string {
	return AllVersionsListPrefix + key
}

func MarshalKeyAndVersion(key, version string) string {
	return fmt.Sprintf(KeyWithVersionTemplate, version, key) //todo review, fmt的性能不好，如果不想自己拼接字符串，可以使用strings.Join
}

func UnmarshalKeyWithVersion(keyWithVersion string) (string, string, error) {
	pattern, err := regexp.Compile(KeyWithVersionTemplateRE) //todo review, 不建议使用regex，如果要使用也要先编译好 pattern，不用每次使用都编译一次
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
