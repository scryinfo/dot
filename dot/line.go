package dot

import (
	"reflect"
)

var gline Line

func GetDefaultLine() Line {
	return gline
}

func SetDefaultLine(line Line) {
	gline = line
}

//
type Injecter interface {
	//Inject inject
	//obj只支持结构体
	//字段中带有 dot.TagDot (dot) 的 tag
	//如果tag为空，那么以字段的类型来注入，不为空以tag的值（dot.LiveId）进行注入
	//在整个过程中如果出错不会退出， 返回的错误是发生的第一个错误
	Inject(obj interface{}) error
	//GetByType get by type
	//如果在当前的容器中没有找到对应的，会调用 parent 查找
	//type 的容器与 liveid的容器是单独分开的
	GetByType(t reflect.Type) (d Dot, err error)
	//GetByLiveId get by liveid
	//如果在当前的容器中没有找到对应的，会调用 parent 查找
	//type 的容器与liveid的容器是单独分开的
	GetByLiveId(id LiveId) (d Dot, err error)

	//ReplaceOrAddByType update
	//不会操作prarent
	ReplaceOrAddByType(d Dot) error
	//不会操作prarent
	ReplaceOrAddByParamType(d Dot, t reflect.Type) error
	//ReplaceOrAddByLiveId update
	//不会操作prarent
	ReplaceOrAddByLiveId(d Dot, id LiveId) error
	//RemoveByType remove
	RemoveByType(t reflect.Type) error
	//RemoveByLiveId remove
	RemoveByLiveId(id LiveId) error

	//SetParent set parent injecter
	SetParent(p Injecter)
	//GetParent get parent injecter
	GetParent() Injecter
}

//Line line
type Line interface {

	//Line的接口
	Config() *Config
	//SConfig 通用配置接口
	SConfig() SConfig

	SLogger() SLogger

	//order
	//1,查找liveid对应的 newer
	//2,查找typid对应的 newer
	//3,查找meta上对的 newer
	//AddNewerByLiveId add new for liveid
	AddNewerByLiveId(liveid LiveId, newDot Newer) error
	//AddNewerByTypeId add new for type
	AddNewerByTypeId(typeid TypeId, newDot Newer) error
	//RemoveNewerByLiveId remove
	RemoveNewerByLiveId(liveid LiveId)
	//RemoveNewerByTypeId remove
	RemoveNewerByTypeId(typeid TypeId)

	//PreAdd 增加dot的liveid及meta信息, 这里还没有真正创建dot，计算依赖后才生成
	//如果是单例的话，可以不指定实例信息，实例id直接为typeid
	//如果配置文件在有配置实例，那么会自动增加来，如果实例id已经存在，则配置更优先
	PreAdd(ac *TypeLives) error
	//Rely  检查依赖关系是否都存在
	Rely() error
	//CreateDots create dots
	CreateDots() error
	//ToLifer to lifer
	ToLifer() Lifer
	//ToInjecter to injecter
	ToInjecter() Injecter

	//GetDotConfig get
	GetDotConfig(liveid LiveId) *LiveConfig

	GetLineBuilder() *Builder
}

// 如果组件需要知道当前的line那么，就实现这个接口，此接口会在组件的Create方法之前调用
type NeedLine interface {
	SetLine(l Line)
}

//TypeLives living
type TypeLives struct {
	Meta  Metadata
	Lives []Live
}

//为配置文件中的dot增加typeid， newer
//这个函数是在 line的create之后运行的，也可以增加其它的初始内容
type BuildNewer func(l Line) error
type LifeEvent func(l Line)

type Builder struct {
	Add         BuildNewer

	BeforeCeate LifeEvent  //line的create之前
	AfterCreate LifeEvent  //line的create之后
	BeforeStart LifeEvent  //line的start之前
	AfterStart  LifeEvent  //line的start之后

	BeforeStop    LifeEvent //line的stop 之前
	AfterStop     LifeEvent //line的stop 之后
	BeforeDestroy LifeEvent //line的destroy 之前
	AfterDestroy  LifeEvent //line的destroy 之后

	BeforeCreateDots LifeEvent // 在create dots之前
	AfterCreateDots LifeEvent  // 在create dots之后
}


//NewTypeLives new living
func NewTypeLives() *TypeLives {
	live := &TypeLives{}
	live.Lives = make([]Live, 0)
	return live
}

//Clone the TypeLives, do not clone dot
func (c *TypeLives) Clone() *TypeLives {
	cc := *c
	cc.Lives = make([]Live, len(c.Lives))
	copy(cc.Lives, c.Lives)
	cc.Meta.RelyTypeIds = make([]TypeId, len(c.Meta.RelyTypeIds))
	copy(cc.Meta.RelyTypeIds, c.Meta.RelyTypeIds)
	return &cc
}
