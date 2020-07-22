package redis_client

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

const RedisClientTypeId = "0ae35550-7e37-4afe-866e-b129099759b7"

type configRedis struct {
	Addr                string `json:"addr"`
	KeepMaxVersionCount int64  `json:"keepMaxVersionCount"` //分类中版本的最大数据量
	VersionFromRedis    bool   `json:"versionFromRedis"`    //get version from redis or not
	TrySeconds          int64  `json:"trySeconds"`          //订阅失败后，再次尝试的时间，单位秒， 默认为10秒
}
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
	Set   Operation = iota // set version
	Clean                  // Clean version
)

// VersionControl defines a version edit.
type VersionControl struct {
	// Op indicates the operation of the update.
	Op      Operation `json:"op"`
	Key     string    `json:"key"`
	Version string    `json:"version"`
}

func (c *RedisClient) ClientV8() *redis.Client {
	return c.clientV8
}

// 缓存版本缓存安全
// todo：思考服务器端redis配置，包括但不限于分配多大的内存、内存占用已满时的策略（报错或删除一部分内存，一共6种策略的那个）
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

//SetVersion
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

	pipe.Publish(c.ctx, VersionControlChannelName, MarshalVersionControl(&VersionControl{
		Op: Set, Key: categoryKey, Version: version,
	}))
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
		if err = c.clientV8.ZRemRangeByRank(c.ctx, categoriesKey, length-c.conf.KeepMaxVersionCount, -1).Err(); err != nil { //todo the stop -1 work
			dot.Logger().Errorln("redis.rPop failed", zap.NamedError("error", err))
		}
	}

	return nil
}

//Get
func (c *RedisClient) Get(category string, version string, itemKey string) (string, error) {
	key := VersionItemKey(category, version, itemKey)
	return c.clientV8.Get(c.ctx, key).Result()
}

//Set
func (c *RedisClient) Set(category string, version string, itemKey string, value string, expiration time.Duration) (string, error) {
	key := VersionItemKey(category, version, itemKey)
	return c.clientV8.Set(c.ctx, key, value, expiration).Result()
}

//Del
func (c *RedisClient) Del(category string, version string, itemKey string) (int64, error) {
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
		//scan命令，参见 http://doc.redisfans.com/key/scan.html
		keys, cursor, err = c.clientV8.Scan(c.ctx, cursor, versionKey+KeySplitChar+"*", 100).Result()
		pipe := c.clientV8.TxPipeline()
		pipe.Del(c.ctx, keys...)
		_, err = pipe.Exec(c.ctx)
		if err != nil {
			dot.Logger().Errorln("RedisClient CleanVersion", zap.NamedError("error", err))
			return err
		}
		if cursor < 1 { //scan is over
			break
		}
	}
	pipe := c.clientV8.TxPipeline()
	pipe.Del(c.ctx, categoryKey)
	pipe.ZRem(c.ctx, categoriesKey, version)
	pipe.Publish(c.ctx, VersionControlChannelName, MarshalVersionControl(&VersionControl{
		Op: Clean, Key: categoryKey, Version: "",
	}))
	_, err = pipe.Exec(c.ctx)
	if err != nil {
		dot.Logger().Errorln("RedisClient CleanVersion", zap.NamedError("error", err))
	}
	return err
}

//GetVersions,  返回版本号按按照字典顺序由小到大排列
func (c *RedisClient) GetVersions(category string) ([]string, error) {
	categoryKey := CategoriesKey(category)
	versions, err := c.clientV8.ZRange(c.ctx, categoryKey, 0, -1).Result()
	if err != nil {
		dot.Logger().Errorln("redis get all versions failed", zap.NamedError("error", err))
	}

	return versions, err
}

//Create
func (c *RedisClient) Create(dot.Line) error {
	c.clientV8 = redis.NewClient(&redis.Options{
		Addr: c.conf.Addr,
	})

	// directly return if not use component to manage version
	if c.conf.VersionFromRedis {
		return nil
	}

	go func() {
		c.subscribe = c.clientV8.Subscribe(c.ctx, VersionControlChannelName)

		// wait until subscribe success
		_, err := c.subscribe.Receive(c.ctx)
		if err != nil {
			dot.Logger().Errorln("subscription redis failed", zap.NamedError("error", err))
		}
		for {
			select {
			case <-c.ctx.Done():
				return
			case msg := <-c.subscribe.Channel():
				vc, err := UnmarshalVersionControl(msg.Payload)
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

//AfterAllDestroy
func (c *RedisClient) AfterAllDestroy(_ dot.Line) {
	if c.cancelFun != nil {
		c.cancelFun() //此方法可以调用多次
	}

	if c.subscribe != nil {
		c.subscribe.Close()
		c.subscribe = nil
	}

	if c.clientV8 != nil {
		c.clientV8.Close()
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
	if d.conf.TrySeconds < 1 {
		d.conf.TrySeconds = 10
	}

	d.ctx, d.cancelFun = context.WithCancel(context.Background())

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

func RedisClientTest(jsonConfig string) *RedisClient {
	d, err := newRedisClient([]byte(jsonConfig))
	if err != nil || d == nil {
		return nil
	}

	redisClient := d.(*RedisClient)
	redisClient.Create(nil)
	return redisClient
}

func GenerateKey(keys ...string) string {
	return strings.Join(keys, KeySplitChar)
}

func CategoryKey(category string) string {
	return GenerateKey(CategoryPrefix, category)
}
func CategoriesKey(category string) string {
	return GenerateKey(CategoriesPrefix, category)
}

//VersionKey
func VersionKey(category string, version string) string {
	return GenerateKey(CategoryPrefix, category, version)
}

//VersionItemKey 使用category,version,item key生成redis的key
func VersionItemKey(category string, version string, itemKey string) string {
	return GenerateKey(CategoryPrefix, category, version, itemKey)
}

func MarshalVersionControl(vc *VersionControl) string {
	bs, err := json.Marshal(vc)
	if err == nil {
		return string(bs)
	}
	return ""
}

func UnmarshalVersionControl(jsonStr string) (*VersionControl, error) {
	vc := &VersionControl{}
	err := json.Unmarshal([]byte(jsonStr), vc)
	if err != nil {
		vc = nil
	}
	return vc, err
}
