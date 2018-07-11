package line

import (
	"sync"

	"github.com/scryinfo/dot"
	"github.com/scryinfo/dot/dots/slog"
)

var (
	_ dot.Lifer = (*line)(nil)
)

//Line line
type Line interface {
	//Line的接口
	Config() *Config
	//SConfig 通用配置接口
	SConfig() dot.SConfig
	//PreAdd 增加dot的liveid及meta信息, 这里还没有真正创建dot，计算依赖后才生成
	//如果是单例的话，可以不指定实例信息，实例id直接为typeid
	//如果配置文件在有配置实例，那么会自动增加来，如果实例id已经存在，则配置更优先

	PreAdd(ac *TypeLives) error
	//Rely  检查依赖关系是否都存在
	Rely() error
	//CreateDots create dots
	CreateDots() error
	//ToLifer to lifer
	ToLifer() dot.Lifer

	//GetDotLive get
	GetDotLive(liveid dot.LiveId) *dot.Live
	//GetDotConfig get
	GetDotConfig(liveid dot.LiveId) *DotConfig
}

//TypeLives living
type TypeLives struct {
	Meta  dot.MetaData
	Lives []dot.Live
}

//NewTypeLives new living
func NewTypeLives() *TypeLives {
	live := &TypeLives{}
	live.Lives = make([]dot.Live, 0)
	return live
}

//Clone the TypeLives, do not clone dot
func (c *TypeLives) Clone() *TypeLives {
	cc := *c
	cc.Lives = make([]dot.Live, len(c.Lives))
	copy(cc.Lives, c.Lives)
	cc.Meta.RelyTypeIds = make([]dot.TypeId, len(c.Meta.RelyTypeIds))
	copy(cc.Meta.RelyTypeIds, c.Meta.RelyTypeIds)
	return &cc
}

type line struct {
	dot.Lifer
	Line
	logger  slog.SLogger
	sConfig dot.SConfig
	config  Config
	metas   Metas
	lives   Lives

	mutex sync.Mutex
}

//NewError new
func New() Line {
	a := &line{}
	return a
}

//PreAdd the dot is nil, do not create it
func (c *line) PreAdd(livings *TypeLives) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	clone := livings.Clone()
	err := c.metas.Add(&clone.Meta)
	if err == nil {
		var dconf []*DotConfig = make([]*DotConfig, 0, 10)
		//get info from config
		{
			for _, it := range c.config.Dots {
				if it.TypeId == clone.Meta.TypeId {
					dconf = append(dconf, &it)
				}
			}
		}
		//get info from config end
		if len(dconf) > 0 {
			ls := make([]dot.Live, 0, len(dconf))
			lsm := make(map[dot.LiveId]*DotConfig, 0)
			for _, it := range dconf {
				lsm[it.LiveId] = it
				ls = append(ls, dot.Live{TypeId: it.TypeId, LiveId: it.LiveId, RelyLives: it.RelyLives})
			}
			if len(clone.Lives) > 0 { //加入配置中不存在的
				for _, it := range clone.Lives {
					if _, ok := lsm[it.LiveId]; ok {
						continue
					}
					ls = append(ls, it)
				}
			}
			clone.Lives = ls
		} else {
			if len(clone.Lives) > 0 {

			} else {
				clone.Lives = []dot.Live{dot.Live{TypeId: clone.Meta.TypeId, LiveId: dot.LiveId(clone.Meta.TypeId)}}
			}
		}

		for _, it := range clone.Lives {
			c.lives.Add(&it)
		}
	}

	return err
}

func (c *line) Rely() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var err error
	var live *dot.Live
LIVES:
	for _, it := range c.lives.liveIdMap {
		for _, dit := range it.RelyLives {
			live, err = c.lives.Get(dit)
			if err != nil {
				break LIVES
			}

			if live.TypeId != it.TypeId {
				err = dot.NewError(dot.SError.ErrRelyTypeNotMatch.Code(), dot.SError.ErrRelyTypeNotMatch.Error()+live.LiveId.String())
				break LIVES
			}
		}
	}

	return err
}

func (c *line) CreateDots() error {
	var err error

	return err
}

func (c *line) ToLifer() dot.Lifer {
	return c
}

func (c *line) GetDotLive(liveid dot.LiveId) *dot.Live {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	l, err := c.lives.Get(liveid)
	if err != nil {
		l = nil
	}
	return l
}

func (c *line) GetDotConfig(liveid dot.LiveId) *DotConfig {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var co *DotConfig
	for _, it := range c.config.Dots {
		if it.LiveId == liveid {
			co = &it
			break
		}
	}
	return co
}

//Create create
//如果 liveid为空， 直接赋值为 typeid
//如果 liveid重复，直接返回 dot.SError.ErrExistedLiveId
func (c *line) Create(conf dot.SConfig) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	var err error
	for {
		//first create config
		c.sConfig = dot.NewConfiger()
		if err = c.sConfig.Create(nil); err != nil {
			break
		}
		if err = c.sConfig.Unmarshal(c.config); err != nil {
			break
		}

		//check config
		{
			vt := make(map[dot.LiveId]*DotConfig, 0)
			for _, it := range c.config.Dots {
				if len(it.LiveId.String()) < 1 {
					it.LiveId = dot.LiveId(it.TypeId)
				}
				if _, ok := vt[it.LiveId]; ok {
					err = dot.NewError(dot.SError.ErrExisted.Code(), dot.SError.ErrExisted.Error())
					break
				}
				vt[it.LiveId] = &it
			}

			if err != nil {
				break
			}
		}

		//create log

		//create others
		break
	}

	return err
}

//Start
func (c *line) Start(ignore bool) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var err error
	for {
		//start config
		if err = c.sConfig.Start(ignore); err != nil {
			break
		}

		//start log

		//start other
		break
	}

	return nil
}

//Stop
func (c *line) Stop(ignore bool) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var err error
	//stop others

	//stop log

	//stop config
	err = c.sConfig.Stop(ignore)

	return err
}

//Destroy 销毁 Dot
func (c *line) Destroy(ignore bool) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	//Destroy others

	//Destroy log

	//Destroy config
	c.sConfig.Destroy(ignore)
	return nil
}

///////////////
