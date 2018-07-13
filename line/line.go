package line

import (
	"reflect"

	"github.com/scryinfo/dot/dot"
)

var gline Line

func GetDefaultLine() Line {
	return gline
}

func SetDefaultLine(line Line) {
	gline = line
}

type Injecter interface {
	//Inject inject
	Inject(obj interface{}) error
	//GetByType get by type
	GetByType(t reflect.Type) (d dot.Dot, err error)
	//GetByLiveId get by liveid
	GetByLiveId(liveId dot.LiveId) (d dot.Dot, err error)
	//GetByTypeId get by typeid
	GetByTypeId(typeId dot.TypeId) (d dot.Dot, err error)
}

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
