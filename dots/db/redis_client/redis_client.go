package redis_client //nolint:golint

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/scryinfo/dot/dot"
)

//RedisClientTypeID type id
const RedisClientTypeID = "0ae35550-7e37-4afe-866e-b129099759b7"

type configRedis struct {
	Addr                string `json:"addr"`
	KeepMaxVersionCount int64  `json:"keepMaxVersionCount"` //分类中版本的最大数据量
	VersionFromRedis    bool   `json:"versionFromRedis"`    //get version from redis or not
}

//RedisClient redis client dot
type RedisClient struct {
	conf configRedis

	subscribe *redis.PubSub
	versions  sync.Map
	clientV8  *redis.Client
	ctx       context.Context    //用于退出程序
	cancelFun context.CancelFunc //ctx的取消
}

// Operation defines the corresponding operations
type Operation uint8

const (
	//Set set version
	Set Operation = iota
	//Clean Clean version
	Clean
)

// VersionControl defines a version edit.
type VersionControl struct {
	// Op indicates the operation of the update.
	Op      Operation `json:"op"`
	Key     string    `json:"key"`
	Version string    `json:"version"`
}

//Marshal to json string
func (c *VersionControl) Marshal() string {
	bs, err := json.Marshal(c)
	if err == nil {
		return string(bs)
	}
	return ""
}

//Unmarshal json string to VersionControl
func (c *VersionControl) Unmarshal(jsonStr string) (*VersionControl, error) {
	vc := c
	err := json.Unmarshal([]byte(jsonStr), vc)
	if err != nil {
		vc = nil
	}
	return vc, err
}

//ClientV8 return redis.Client
func (c *RedisClient) ClientV8() *redis.Client {
	return c.clientV8
}

// GetVersion get current version by category, return version and error.
func (c *RedisClient) GetVersion(category string) (string, error) {
	categoryKey := CategoryKey(category)
	if c.conf.VersionFromRedis {
		return c.clientV8.Get(c.ctx, categoryKey).Result()
	}

	versionCache, ok := c.versions.Load(categoryKey)
	if !ok {
		return "", errors.New(fmt.Sprintf("category: %s is not exist", category))
	}

	version, ok := versionCache.(string)
	if !ok {
		return "", errors.New(fmt.Sprintf("category: %s is not string type", category))
	}

	return version, nil
}

//SetVersion set version for category
func (c *RedisClient) SetVersion(category, version string) error {
	// tx: current version / all versions
	categoryKey := CategoryKey(category)
	categoriesKey := CategoriesKey(category)

	pipe := c.clientV8.TxPipeline()
	pipe.Set(c.ctx, categoryKey, version, 0)
	pipe.ZAdd(c.ctx, categoriesKey, &redis.Z{
		Score:  0,
		Member: version,
	}) //zadd will ignore the same version

	vcJSON := (&VersionControl{
		Op: Set, Key: categoryKey, Version: version,
	}).Marshal()
	pipe.Publish(c.ctx, VersionControlChannelName, vcJSON)
	_, err := pipe.Exec(c.ctx)
	if err != nil {
		dot.Logger().Errorln("redis set version failed", zap.NamedError("error", err))
		return err
	}
	if !c.conf.VersionFromRedis {
		c.versions.Store(categoryKey, version)
	}

	// limit keep versions' length
	// set version 函数本身重点不在于维护历史版本list，某一次维护失败对业务没有影响，
	// 所以维护操作没有加入事务中，且出错也只是打印日志，而不认为函数执行错误
	length, err := c.clientV8.ZCard(c.ctx, categoriesKey).Result()
	if err != nil {
		dot.Logger().Errorln("redis get list.length failed",
			zap.NamedError("error", err),
			zap.String("list name", category))
		return nil
	}
	if length > c.conf.KeepMaxVersionCount {
		if err = c.clientV8.ZRemRangeByRank(c.ctx, categoriesKey, length-c.conf.KeepMaxVersionCount, -1).Err(); err != nil {
			dot.Logger().Errorln("redis.rPop failed", zap.NamedError("error", err))
		}
	}

	return nil
}

//Get get value
func (c *RedisClient) Get(category, version, itemKey string) (string, error) {
	key := VersionItemKey(category, version, itemKey)
	return c.clientV8.Get(c.ctx, key).Result()
}

//Get get value
func (c *RedisClient) GetJsonSerialize(category, version, itemKey string, value interface{}) error {
	key := VersionItemKey(category, version, itemKey)
	str, err := c.clientV8.Get(c.ctx, key).Result()
	if err == nil {
		err = json.Unmarshal([]byte(str), value)
	}
	return err
}

//Set set value
func (c *RedisClient) Set(category, version, itemKey, value string, expiration time.Duration) (string, error) {
	key := VersionItemKey(category, version, itemKey)
	return c.clientV8.Set(c.ctx, key, value, expiration).Result()
}

//Set set value
func (c *RedisClient) SetJsonSerialize(category, version, itemKey string, value interface{}, expiration time.Duration) (string, error) {
	bs, err := json.Marshal(value)
	if err == nil {
		key := VersionItemKey(category, version, itemKey)
		jsonValue := string(bs)
		return c.clientV8.Set(c.ctx, key, jsonValue, expiration).Result()
	} else {
		return "", err
	}
}

//Del del key
func (c *RedisClient) Del(category, version, itemKey string) (int64, error) {
	key := VersionItemKey(category, version, itemKey)
	return c.clientV8.Del(c.ctx, key).Result()
}

//CleanVersion del versions in key
func (c *RedisClient) CleanVersion(category string, version string) error {
	categoryKey := CategoryKey(category)
	categoriesKey := CategoriesKey(category)
	versionKey := VersionKey(category, version)
	// re-curse del target versions
	var keys []string
	var cursor uint64
	var err error
	for {
		//scan see http://doc.redisfans.com/key/scan.html
		keys, cursor, err = c.clientV8.Scan(c.ctx, cursor, versionKey+KeySplitChar+"*", 100).Result()
		if err != nil {
			dot.Logger().Errorln("RedisClient CleanVersion", zap.NamedError("error", err))
			return err
		}
		if len(keys) > 0 {
			pipe := c.clientV8.TxPipeline()
			pipe.Del(c.ctx, keys...)
			_, err = pipe.Exec(c.ctx)
			if err != nil {
				dot.Logger().Errorln("RedisClient CleanVersion", zap.NamedError("error", err))
				return err
			}
		}
		if cursor < 1 { //scan is over
			break
		}
	}
	pipe := c.clientV8.TxPipeline()
	pipe.Del(c.ctx, categoryKey)
	pipe.ZRem(c.ctx, categoriesKey, version)
	vsJSON := (&VersionControl{
		Op: Clean, Key: categoryKey, Version: "",
	}).Marshal()
	pipe.Publish(c.ctx, VersionControlChannelName, vsJSON)
	_, err = pipe.Exec(c.ctx)
	if err != nil {
		dot.Logger().Errorln("RedisClient CleanVersion", zap.NamedError("error", err))
	}
	return err
}

//GetVersions return versions (dictionary order)
func (c *RedisClient) GetVersions(category string) ([]string, error) {
	categoryKey := CategoriesKey(category)
	versions, err := c.clientV8.ZRange(c.ctx, categoryKey, 0, -1).Result()
	if err != nil {
		dot.Logger().Errorln("redis get all versions failed", zap.NamedError("error", err))
	}

	return versions, err
}

//Create create dot
func (c *RedisClient) Create(dot.Line) error {
	c.clientV8 = redis.NewClient(&redis.Options{
		Addr: c.conf.Addr,
	})

	// directly return if not use component to manage version
	if c.conf.VersionFromRedis {
		return nil
	}

	go func() {
		//see https://redis.uptrace.dev/#pubsub the subscribe automatic reconnect
		c.subscribe = c.clientV8.Subscribe(c.ctx, VersionControlChannelName)

		for {
			select {
			case <-c.ctx.Done():
				return
			case msg, ok := <-c.subscribe.Channel(): //see
				if !ok || msg == nil {
					continue
				}
				vc := &VersionControl{}
				vc, err := vc.Unmarshal(msg.Payload)
				if err != nil {
					dot.Logger().Errorln("unmarshal key with version failed",
						zap.NamedError("error", err),
						zap.String("key with version", msg.Payload))
				} else {
					if vc.Op == Set {
						c.versions.Store(vc.Key, vc.Version)
					} else if vc.Op == Clean {
						c.versions.Delete(vc.Key)
					}
				}
			}
		}
	}()

	return nil
}

//AfterAllDestroy after all destroy
func (c *RedisClient) AfterAllDestroy(_ dot.Line) {
	if c.cancelFun != nil {
		c.cancelFun() //此方法可以调用多次
	}

	if c.subscribe != nil {
		_ = c.subscribe.Close()
		c.subscribe = nil
	}

	if c.clientV8 != nil {
		_ = c.clientV8.Close()
		c.clientV8 = nil
	}
}

//construct dot
func newRedisClient(conf []byte) (dot.Dot, error) {
	dconf := &configRedis{}

	err := dot.UnMarshalConfig(conf, dconf)
	if err != nil {
		return nil, err
	}

	d := &RedisClient{conf: *dconf}

	d.ctx, d.cancelFun = context.WithCancel(context.Background())

	return d, nil
}

//RedisClientTypeLives return type lives
func RedisClientTypeLives() []*dot.TypeLives {
	return []*dot.TypeLives{{
		Meta: dot.Metadata{TypeID: RedisClientTypeID, NewDoter: newRedisClient},
	}}
}

//RedisClientTest for test
func RedisClientTest(jsonConfig string) *RedisClient {
	d, err := newRedisClient([]byte(jsonConfig))
	if err != nil || d == nil {
		return nil
	}

	redisClient := d.(*RedisClient)
	_ = redisClient.Create(nil)
	time.Sleep(2 * time.Second)
	return redisClient
}

//RedisClientConfigTypeLive return config
func RedisClientConfigTypeLive() *dot.ConfigTypeLive {
	return &dot.ConfigTypeLive{
		TypeIDConfig: RedisClientTypeID,
		ConfigInfo:   &configRedis{},
	}
}

//GenerateKey generate key for redis
func GenerateKey(keys ...string) string {
	return strings.Join(keys, KeySplitChar)
}

//CategoryKey category key
func CategoryKey(category string) string {
	return GenerateKey(CategoryPrefix, category)
}

//CategoriesKey categories key
func CategoriesKey(category string) string {
	return GenerateKey(CategoriesPrefix, category)
}

//VersionKey version key
func VersionKey(category string, version string) string {
	return GenerateKey(CategoryPrefix, category, version)
}

//VersionItemKey use category,version,item key to generate key of redis
func VersionItemKey(category string, version string, itemKey string) string {
	return GenerateKey(CategoryPrefix, category, version, itemKey)
}
