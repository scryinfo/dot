// Scry Info.  All rights reserved.
// license that can be found in the license file.

package dot

import (
	"reflect"
)

// TypeID dot type guid
type TypeID string

// LiveID dot live guid
type LiveID string

// String convert typeid to string
func (c *TypeID) String() string {
	return string(*c)
}

// String convert liveid to string
func (c *LiveID) String() string {
	return string(*c)
}

// Event dot event
type Event = func(live *Live, l Line)

// Events dot events
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

// TypeEvents dot events for type
type TypeEvents = Events

// LiveEvents dot events for live
type LiveEvents = Events

// Eventer event interface
type Eventer interface {
	//
	ReSetLiveEvents(lid LiveID, liveEvents *LiveEvents)
	//
	AddLiveEvents(lid LiveID, liveEvents *LiveEvents)
	//
	LiveEvents(lid LiveID) []LiveEvents
	//
	ReSetTypeEvents(tid TypeID, typeEvents *TypeEvents)
	//
	AddTypeEvents(tid TypeID, typeEvents *TypeEvents)
	//
	TypeEvents(tid TypeID) []TypeEvents
}

// Metadata dot metadata
type Metadata struct {
	TypeID      TypeID       `json:"typeId"`
	Version     string       `json:"version"`
	Name        string       `json:"name"`
	ShowName    string       `json:"showName"`
	Single      bool         `json:"single"`
	RelyTypeIDs []TypeID     `json:"relyTypeIds"`
	NewDoter    Newer        `json:"-"`
	RefType     reflect.Type `json:"-"` //if the dot is interface, set the type of interface
}

// Live live/instance
type Live struct {
	TypeID    TypeID            `json:"typeId"`
	LiveID    LiveID            `json:"liveId"`
	RelyLives map[string]LiveID `json:"relyLives"`
	Dot       Dot               `json:"-"`
}

// Clone clone Metadata
func (m *Metadata) Clone() *Metadata {
	c := *m
	c.RelyTypeIDs = append(m.RelyTypeIDs[:0:0], m.RelyTypeIDs...)
	return &c
}

// Merge merge metadata
func (m *Metadata) Merge(m2 *Metadata) {
	if len(m2.TypeID) > 0 {
		m.TypeID = m2.TypeID
	}
	if len(m2.Version) > 0 {
		m.Version = m2.Version
	}
	if len(m2.Name) > 0 {
		m.Name = m2.Name
	}
	if len(m2.ShowName) > 0 {
		m.ShowName = m2.ShowName
	}
	m.Single = m2.Single
	if len(m2.RelyTypeIDs) > 0 {
		mergeIDs := make([]TypeID, 0, len(m.RelyTypeIDs)+len(m2.RelyTypeIDs))
		mergeIDs = append(mergeIDs, m.RelyTypeIDs...)
		hasID := make(map[TypeID]bool, cap(mergeIDs))
		for i := range m.RelyTypeIDs {
			hasID[m.RelyTypeIDs[i]] = true
		}

		for _, relyTypeID := range m2.RelyTypeIDs {
			if _, ok := hasID[relyTypeID]; !ok {
				hasID[relyTypeID] = true
				mergeIDs = append(mergeIDs, relyTypeID)
			}
		}
		m.RelyTypeIDs = mergeIDs[:]
	}
	if m2.NewDoter != nil {
		m.NewDoter = m2.NewDoter
	}
	if m2.RefType != nil {
		m.RefType = m2.RefType
	}
}

// NewDot new a dot
func (m *Metadata) NewDot(args []byte) (dot Dot, err error) {
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

// Newer instance for new dot
type Newer = func(args []byte) (dot Dot, err error)

// Dot component
type Dot interface{}

// Lifer life cycle
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

// Creator create interface
type Creator interface {
	//Create When this method is initializing, running or monitoring same content, better to realize it in Start method
	Create(l Line) error
}

// Injected inject interface
type Injected interface {
	//Injected call the function after inject
	Injected(l Line) error
}

// Starter start interface
type Starter interface {
	//ignore When calling other Lifer, if true erred will continue, if false erred will return directly
	Start(ignore bool) error
}

// Stopper stop interface
type Stopper interface {
	//ignore When calling other Lifer, if true erred will continue, if false erred will return directly
	Stop(ignore bool) error
}

// Destroyer destroy interface
type Destroyer interface {
	//Destroy Dot
	//ignore When calling other Lifer, if true erred will continue, if false erred will return directly
	Destroy(ignore bool) error
}

// Tager dot signature data, used by dot
type Tager interface {
	//SetTag set tag
	SetTag(tag interface{})
	//GetTag get tag
	GetTag() (tag interface{})
}

// GetInterfaceType return the interface of dot
type GetInterfaceType interface {
	//get interface type
	GetInterfaceType() reflect.Type
}

// StatusType status type
type StatusType int

// Statuser Status
type Statuser interface {
	Status() StatusType
}

// HotConfig hot change config
type HotConfig interface {
	//Update Update config info, return true means successful
	HotConfig(newConf SConfig) bool
}

// Checker Check dot，run some verification or test data, return the result
type Checker interface {
	Check(args interface{}) interface{}
}

const (
	//TagDot tag dot
	TagDot = "dot"
	//CanNull allow nil, if do not find the dot , set it nil
	CanNull = "?"
)
