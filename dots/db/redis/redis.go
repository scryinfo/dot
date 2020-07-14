package redis

import (
	"encoding/json"
	"github.com/albrow/zoom"
	"github.com/scryinfo/dot/dot"
	"go.uber.org/zap"
)

const RedisTypeId = "0ae35550-7e37-4afe-866e-b129099759b7"

type configRedis struct {
	Addr string `json:"addr"`
}
type Redis struct {
	conf configRedis

	pool *zoom.Pool
}

// todo?：考虑修改包名和结构体名，但是不知道改成啥
// todo：+版本管理，思考：所有数据都要管理吗？还是只有复杂数据需要管理？

// todo：缓存安全
// todo：思考服务器端redis配置，包括但不限于分配多大的内存、内存占用已满时的策略（报错或删除一部分内存，一共6种策略的那个）

// all models can use 'find all' support from zoom, because of the param: 'Index: true'
func (c *Redis) RegisterCollections(models []Model) []*Collection {
	res := make([]*Collection, 0)
	for i, m := range models {
		c, err := c.pool.NewCollectionWithOptions(m, zoom.CollectionOptions{
			FallbackMarshalerUnmarshaler: m,
			Index:                        true,
			Name:                         m.GetName(),
		})
		if err != nil {
			dot.Logger().Errorln("register redis collection failed",
				zap.Int("index", i),
				zap.NamedError("error", err))
			continue
		}
		res = append(res, &Collection{Collection: c})
	}

	if len(models) != len(res) {
		dot.Logger().Errorln("something works not as suppose",
			zap.Int("models.length", len(models)),
			zap.Int("result.length", len(res)))
	}

	return res
}

func (c *Redis) Create(_ dot.Line) error {
	c.pool = zoom.NewPool(c.conf.Addr)

	return nil
}

func (c *Redis) AfterAllIDestroy(_ dot.Line) {
	if err := c.pool.Close(); err != nil {
		dot.Logger().Errorln("redis close failed", zap.NamedError("error", err))
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

type Model interface {
	zoom.Model // zoom.RandomID implements this interface
	zoom.MarshalerUnmarshaler
	GetName() string
}

var _ Model = (*ModelImp)(nil)

type ModelImp struct {
	Id   string
	Name string
}

func (m *ModelImp) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (m *ModelImp) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (m *ModelImp) ModelID() string {
	if m.Id == "" {
		m.Id = (&zoom.RandomID{}).ModelID()
	}

	return m.Id
}

func (m *ModelImp) SetModelID(id string) {
	m.Id = id
}

func (m *ModelImp) GetName() string {
	if m.Name != "" {
	} else if m.Id != "" {
		m.Name = m.Id
	} else {
		m.Name = m.ModelID()
	}

	return m.Name
}

type Collection struct {
	*zoom.Collection
}
