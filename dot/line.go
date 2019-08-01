// Scry Info.  All rights reserved.
// license that can be found in the license file.

package dot

import (
	"reflect"
)

//
type Injecter interface {
	//Inject inject
	//obj only support structure
	//dot.TagDot (dot) tag is in the field
	//If tag is empty, then input with field type, otherwise input with tag value（dot.LiveId）
	//In the process if error occurred, it will not quit, returned error is the first one occurred
	Inject(obj interface{}) error
	//GetByType get by type
	//If no related object in current container, then will call parent to search
	//type container is seperate with liveid container
	GetByType(t reflect.Type) (d Dot, err error)
	//GetByLiveId get by liveid
	//If no related object in current container, then will call parent to search
	//type container is seperate with liveid container
	GetByLiveId(id LiveId) (d Dot, err error)

	//ReplaceOrAddByType update
	//Cannot operate parent
	ReplaceOrAddByType(d Dot) error
	//Cannot operate parent
	ReplaceOrAddByParamType(d Dot, t reflect.Type) error
	//ReplaceOrAddByLiveId update
	//Cannot operate parent
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

	//Return unique Line name
	Id() string
	//Line API
	Config() *Config
	//SConfig general config API
	SConfig() SConfig

	SLogger() SLogger

	//order
	//1,Search liveid corresponding newer
	//2,Search typid corresponding newer
	//3,Search right newer in meta
	//AddNewerByLiveId add new for liveid
	AddNewerByLiveId(liveid LiveId, newDot Newer) error
	//AddNewerByTypeId add new for type
	AddNewerByTypeId(typeid TypeId, newDot Newer) error
	//RemoveNewerByLiveId remove
	RemoveNewerByLiveId(liveid LiveId)
	//RemoveNewerByTypeId remove
	RemoveNewerByTypeId(typeid TypeId)

	//PreAdd Add dot liveid and meta info, here no dot is created, it will be generated after Computing dependence
	//If it is the single sample, don't need to point sample info, sample id is typeid
	//If config file has config sample, then it will be added automatically, if sample id already existing, then config is prior
	PreAdd(typeLives ...*TypeLives) error
	//RelyOrder  Check whether dependency existing
	//RelyOrder() error
	////CreateDots create dots
	//CreateDots() error
	//ToLifer to lifer
	ToLifer() Lifer
	//ToInjecter to injecter
	ToInjecter() Injecter

	//ToDotEventer to Eventer
	ToDotEventer() Eventer

	//GetDotConfig get
	GetDotConfig(liveid LiveId) *LiveConfig

	GetLineBuilder() *Builder
	//InfoAllTypeAdnLives just for debug, log info all types and lives
	InfoAllTypeAdnLives()
}

// If component need to know current line, then realize this API, and this API Will be called before component Create
type SetterLine interface {
	SetLine(l Line)
}

// If component need to know current TypeId or LiveId, then realize this API, and this API Will be called before component Create
type SetterTypeAndLiveId interface {
	SetTypeId(tid TypeId, lid LiveId)
}

// After all start, before builder AfterStart
type AfterAllStarter interface {
	AfterAllStart(l Line)
}

// After all inject, before builder AfterStart
type AfterAllInjecter interface {
	AfterAllInject(l Line)
}

// After all destroy, before builder AfterDestroy
type AfterAllIDestroyer interface {
	AfterAllIDestroy(l Line)
}

// Call before all stop, after Builder BeforeStop
type BeforeAllStopper interface {
	BeforeAllStop(l Line)
}

//TypeLives living
type TypeLives struct {
	Meta  Metadata
	Lives []Live
}

//jayce edit
//config json
type ConfigTypeLives struct {
	TypeIdConfig TypeId      `json:"typeId"`
	ConfigInfo   interface{} `json:"json"`
}

//Add typeid, newer for dot in config file
//This function is run after line create, also you can add other initialized content
type BuildNewer func(l Line) error
type AllEvent func(l Line)

type Builder struct {
	Add BuildNewer

	BeforeCreate  AllEvent //Before all dot create
	AfterCreate   AllEvent //after  all dot create
	BeforeStart   AllEvent //Before  all dot start
	AfterStart    AllEvent //After  all dot start
	BeforeStop    AllEvent //Before  all dot stop
	AfterStop     AllEvent //After  all dot stop
	BeforeDestroy AllEvent //Before  all dot destroy
	AfterDestroy  AllEvent //After  all dot destroy

	LineLiveId string //line unique id， default value is “default”
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
