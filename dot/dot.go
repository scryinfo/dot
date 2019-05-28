// Scry Info.  All rights reserved.
// license that can be found in the license file.

package dot

import (
	"reflect"
)

//TypeId dot type guid
type TypeId string

//LiveId dot live guid
type LiveId string

//String convert typeid to string
func (c *TypeId) String() string {
	return string(*c)
}

//String convert liveid to string
func (c *LiveId) String() string {
	return string(*c)
}

//Event
type Event func(live *Live, l Line)

//Events
type Events struct {
	//Before the dot create
	BeforeCreate Event
	//after the dot create
	AfterCreate Event
	//Before the dot start
	BeforeStart Event
	//After the dot start
	AfterStart Event
	//Before the dot stop
	BeforeStop Event
	//After the dot stop
	AfterStop Event
	//Before the dot destroy
	BeforeDestroy Event
	//After the dot destroy
	AfterDestroy Event
}

type TypeEvents = Events
type LiveEvents = Events

//Eventer
type Eventer interface {
	//
	ReSetLiveEvents(lid LiveId, liveEvents *LiveEvents)
	//
	AddLiveEvents(lid LiveId, liveEvents *LiveEvents)
	//
	LiveEvents(lid LiveId) []LiveEvents
	//
	ReSetTypeEvents(tid TypeId, typeEvents *TypeEvents)
	//
	AddTypeEvents(tid TypeId, typeEvents *TypeEvents)
	//
	TypeEvents(tid TypeId) []TypeEvents
}

//Metadata dot metadata
type Metadata struct {
	TypeId      TypeId       `json:"typeId"`
	Version     string       `json:"version"`
	Name        string       `json:"name"`
	ShowName    string       `json:"showName"`
	Single      bool         `json:"single"`
	RelyTypeIds []TypeId     `json:"relyTypeIds"`
	NewDoter    Newer        `json:"-"`
	RefType     reflect.Type `json:"-"`
}

//Live live/instance
type Live struct {
	TypeId    TypeId
	LiveId    LiveId
	RelyLives map[string]LiveId
	Dot       Dot
}

//NewMetadata new Metadata
func NewMetadata() *Metadata {
	return &Metadata{}
}

//Clone clone Metadata
func (m *Metadata) Clone() *Metadata {
	c := *m
	c.RelyTypeIds = make([]TypeId, len(m.RelyTypeIds))
	copy(c.RelyTypeIds, m.RelyTypeIds)
	return &c
}

//NewDot new a dot
func (m *Metadata) NewDot(args interface{}) (dot Dot, err error) {
	dot = nil
	err = nil
	if m.NewDoter != nil {
		dot, err = m.NewDoter(args)
	} else if m.RefType != nil {
		v := reflect.New(m.RefType)
		dot = v.Interface()
	}
	return
}

//Newer instace for new dot
type Newer = func(args interface{}) (dot Dot, err error)

//Dot componet
type Dot interface{}

//Lifer life cycle
// Create, Start,Stop,Destroy
// Create and Start are separate, in order to resolve the dependencies between different dot instances,
// if there is no problem with the dependencies, then you can directly null in Start
// All methods of Lifer cannot be stucked while running, now the realization of line is sync call
type Lifer interface {
	Creator
	Starter
	Stopper
	Destroyer
}

type Creator interface {
	//Create When this method is initializing, running or monitoring same content, better to realize it in Start method
	Create(l Line) error
}

type Starter interface {
	//ignore When calling other Lifer, if true erred will continue, if false erred will return directly
	Start(ignore bool) error
}

type Stopper interface {
	//ignore When calling other Lifer, if true erred will continue, if false erred will return directly
	Stop(ignore bool) error
}

type Destroyer interface {
	//Destroy Dot
	//ignore When calling other Lifer, if true erred will continue, if false erred will return directly
	Destroy(ignore bool) error
}

//Tager dot signature data, used by dot
type Tager interface {
	//SetTag set tag
	SetTag(tag interface{})
	//GetTag get tag
	GetTag() (tag interface{})
}

//StatusType status type
type StatusType int

//Statuser Status
type Statuser interface {
	Status() StatusType
}

//HotConfig hot change config
type HotConfig interface {
	//Update Update config info, return true means successful
	HotConfig(newConf SConfig) bool
}

//Checker Check dot，run some verification or test data, return the result
type Checker interface {
	Check(args interface{}) interface{}
}

const (
	//TagDot tag dot
	TagDot = "dot"
)
