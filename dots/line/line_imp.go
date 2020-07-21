// Scry Info.  All rights reserved.
// license that can be found in the license file.

package line

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"go.uber.org/zap"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/sconfig"
	"github.com/scryinfo/dot/dots/slog"
	"github.com/scryinfo/scryg/sutils/skit"
)

var (
	_ dot.Lifer    = (*lineImp)(nil)
	_ dot.Line     = (*lineImp)(nil)
	_ dot.Injecter = (*lineImp)(nil)
)

type lineImp struct {
	logger     dot.SLogger
	sConfig    dot.SConfig //External general config
	config     dot.Config  //Config object of component
	metas      *Metas      //
	lives      *Lives      //
	types      map[reflect.Type]dot.Dot
	newerLives map[dot.LiveId]dot.Newer
	newerTypes map[dot.TypeId]dot.Newer

	parent dot.Injecter
	mutex  sync.Mutex

	lineBuilder *dot.Builder

	//dot event
	dotEvent dotEventerImp
}

//newLine new
func newLine(builder *dot.Builder) *lineImp {
	a := &lineImp{metas: NewMetas(),
		lives: NewLives(), types: make(map[reflect.Type]dot.Dot),
		newerLives:  make(map[dot.LiveId]dot.Newer),
		newerTypes:  make(map[dot.TypeId]dot.Newer),
		lineBuilder: builder,
	}
	a.dotEvent.Init()

	if dot.GetDefaultLine() == nil {
		dot.SetDefaultLine(a)
	}

	return a
}

func (c *lineImp) Id() string {
	return c.lineBuilder.LineLiveId
}

//AddNewerByLiveId add new for live id
func (c *lineImp) AddNewerByLiveId(liveId dot.LiveId, newDot dot.Newer) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.newerLives[liveId]; ok {
		return dot.SError.Existed.AddNewError(liveId.String())
	}

	c.newerLives[liveId] = newDot

	return nil
}

//AddNewerByTypeId add new for type
func (c *lineImp) AddNewerByTypeId(typeId dot.TypeId, newDot dot.Newer) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, ok := c.newerTypes[typeId]; ok {
		return dot.SError.Existed.AddNewError(typeId.String())
	}

	c.newerTypes[typeId] = newDot

	return nil
}

//RemoveNewerByLiveId remove
func (c *lineImp) RemoveNewerByLiveId(liveId dot.LiveId) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.newerLives, liveId)
}

//RemoveNewerByTypeId remove
func (c *lineImp) RemoveNewerByTypeId(typeId dot.TypeId) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.newerTypes, typeId)
}

//PreAdd the dot is nil, do not create it
func (c *lineImp) PreAdd(typeLives ...*dot.TypeLives) error {
	logger := dot.Logger()
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var err error

	for _, typeLive := range typeLives {

		err2 := c.metas.UpdateOrAdd(&typeLive.Meta)
		if err2 == nil {
			if len(typeLive.Lives) > 0 {
				for i := range typeLive.Lives {
					it := &typeLive.Lives[i]
					if len(it.TypeId.String()) < 1 {
						it.TypeId = typeLive.Meta.TypeId
					}

					//live := dot.Live{TypeId: it.TypeId, LiveId: it.LiveId, Dot: nil}
					//live.RelyLives = CloneRelyLiveId(it.RelyLives)
					err2 = c.lives.UpdateOrAdd(it)
					if err2 != nil {
						if err != nil {
							logger.Debug(func() string {
								return fmt.Sprintf("lineImp - meta: %v", typeLive.Meta)
							})
							logger.Errorln("lineImp", zap.Error(err)) //write it into logfile, otherwise it will gone
						}
						err = err2
					}
				}
			} else {
				//do nothing
			}
		} else {
			if err != nil {
				logger.Debug(func() string {
					return fmt.Sprintf("lineImp - meta: %v", typeLive.Meta)
				})
				logger.Errorln("lineImp", zap.Error(err)) //write it into logfile, otherwise it will gone
			}
			err = err2
		}
	}

	return err
}

//order the dots by relay relation
//first return, ordered dots
//second return, circle relay dots
func (c *lineImp) relyOrder() ([]*dot.Live, []*dot.Live) {

	var cloneLives map[dot.LiveId]*dot.Live
	var cloneMetas map[dot.TypeId]*dot.Metadata
	logger := dot.Logger()
	{ //clone live and type
		c.mutex.Lock()
		cloneLives = make(map[dot.LiveId]*dot.Live, len(c.lives.LiveIdMap))
		for k, v := range c.lives.LiveIdMap {
			cloneLives[k] = v
		}
		cloneMetas = make(map[dot.TypeId]*dot.Metadata, len(c.metas.metas))
		for k, v := range c.metas.metas {
			cloneMetas[k] = v
		}
		c.mutex.Unlock()
	}

	order := make([]*dot.Live, 0, len(cloneLives))
	var circle []*dot.Live
	{
		relyed := make(map[dot.LiveId][]dot.LiveId, len(cloneLives))
		for _, it := range cloneLives {
			if _, ok := relyed[it.LiveId]; !ok {
				relyed[it.LiveId] = []dot.LiveId{}
			}

			for _, liveId := range it.RelyLives {
				r, ok := relyed[liveId]
				if !ok {
					r = []dot.LiveId{it.LiveId}
				} else {
					r = append(r, it.LiveId)
				}
				relyed[liveId] = r
			}
		}
		done := make(map[dot.LiveId]bool, len(relyed))          //know orderId
		remain := make(map[dot.LiveId]bool, len(relyed))        //do not know orderId
		levels := make([]map[dot.LiveId]bool, 0, len(relyed)/3) //
		{
			for liveId := range relyed {
				remain[liveId] = true
			}
			level0 := make(map[dot.LiveId]bool)
			for liveId, live := range cloneLives { //level 0
				if len(live.RelyLives) < 1 {
					level0[liveId] = true
					delete(remain, liveId)
					done[liveId] = true
				}
			}
			levels = append(levels, level0)
			levelCurrent := level0
			//todo if level0 is zero
			for curFor := 0; curFor <= len(relyed); curFor++ { //
				levelNext := make(map[dot.LiveId]bool, len(remain))
				for liveId := range levelCurrent {
					des := relyed[liveId]
					if len(des) < 1 {
						delete(remain, liveId)
						done[liveId] = true
					} else {
						for _, liveId2 := range des {
							allDone := true
							for _, liveId3 := range cloneLives[liveId2].RelyLives { //check all RelyLives
								if _, ok := done[liveId3]; !ok {
									allDone = false
									break
								}
							}
							if allDone {
								if _, ok := done[liveId2]; !ok { //just do it once
									levelNext[liveId2] = true
									done[liveId2] = true
								}

								delete(remain, liveId2)

							} else {
								//levelNext[liveId2] = true //put next level
							}
						}
					}
				}
				levels = append(levels, levelNext)
				if len(remain) < 1 || len(levelNext) < 1 {
					break
				}
				levelCurrent = levelNext
			}
		}

		//todo type relay
		{
			for i, lev := range levels {
				logger.Debugln(fmt.Sprintf("level : %d", i))
				for liveId := range lev {
					logger.Debugln(cloneLives[liveId].LiveId.String())
					live, ok := cloneLives[liveId]
					if ok {
						order = append(order, live)
					} else {
						logger.Warnln("", zap.String("", "dot not find the dot live id: "+liveId.String()))
					}
				}
			}
			if len(remain) > 0 {
				circle = make([]*dot.Live, 0, len(remain))
				for liveId := range remain { //append to tail
					live, ok := cloneLives[liveId]
					if ok {
						order = append(order, live)
						circle = append(circle, live)
					} else {
						logger.Warnln("", zap.String("", "dot not find the dot live id: "+liveId.String()))
					}
				}
			}
		}
	}

	return order, circle
}

//CreateDots create dots
func (c *lineImp) CreateDots(orderedDots []*dot.Live) error {
	logger := dot.Logger()
	createDotFun := func(it *dot.Live) error {
		{ // Check whether special info needed before Create
			if nl, ok := it.Dot.(dot.SetterLine); ok {
				nl.SetLine(c)
			}

			if nl, ok := it.Dot.(dot.SetterTypeAndLiveId); ok {
				nl.SetTypeId(it.TypeId, it.LiveId)
			}
		}

		if b := c.dotEvent.TypeEvents(it.TypeId); len(b) > 0 { // do before create for type
			for i := range b {
				e := &b[i]
				if e.BeforeCreate != nil {
					e.BeforeCreate(it, c)
				}
			}
		}

		if b := c.dotEvent.LiveEvents(it.LiveId); len(b) > 0 { // do before create for live
			for i := range b {
				e := &b[i]
				if e.BeforeCreate != nil {
					e.BeforeCreate(it, c)
				}
			}
		}

		if creator, ok := it.Dot.(dot.Creator); ok { //do createDotFun
			if err := creator.Create(c); err != nil {
				return err
			}
		}

		if a := c.dotEvent.LiveEvents(it.LiveId); len(a) > 0 { //do after create for live
			for i := range a {
				e := &a[i]
				if e.AfterCreate != nil {
					e.AfterCreate(it, c)
				}
			}
		}

		if a := c.dotEvent.TypeEvents(it.TypeId); len(a) > 0 { // dot after create for type
			for i := range a {
				e := &a[i]
				if e.AfterCreate != nil {
					e.AfterCreate(it, c)
				}
			}
		}

		return nil
	}
	var err error
	var dotLive *dot.Live //just for debug
CreateLives:
	for _, dotLive = range orderedDots {
		logger.Debug(func() string {
			m, _ := c.metas.Get(dotLive.TypeId)
			if m != nil {
				return fmt.Sprintf("Create dot, type id: %s, live id: %s, name: %s", dotLive.TypeId, dotLive.LiveId, m.Name)
			} else {
				return fmt.Sprintf("Create dot, type id: %s, live id: %s", dotLive.TypeId, dotLive.LiveId)
			}
		})

		if skit.IsNil(&dotLive.Dot) {
			var bytesConfig []byte
			var config *dot.LiveConfig
			{
				config = c.config.FindConfig(dotLive.TypeId, dotLive.LiveId)
				bytesConfig, err = dot.MarshalConfig(config)
				if err != nil {
					break CreateLives
				}
			}

			//new by liveId
			{
				if newer, ok := c.newerLives[dotLive.LiveId]; ok {
					dotLive.Dot, err = newer(bytesConfig)
					if err != nil {
						break CreateLives
					} else {
						if err = createDotFun(dotLive); err != nil {
							break CreateLives
						}
						continue CreateLives
					}
				}
			}
			//new by typeId
			{
				if newer, ok := c.newerTypes[dotLive.TypeId]; ok {
					dotLive.Dot, err = newer(bytesConfig)
					if err != nil {
						break CreateLives
					} else {
						if err = createDotFun(dotLive); err != nil {
							break CreateLives
						}
						continue CreateLives
					}
				}
			}

			//new by metadata
			{
				var m *dot.Metadata
				m, err = c.metas.Get(dotLive.TypeId)
				if err != nil {
					break CreateLives
				}

				if m.NewDoter == nil && m.RefType == nil {
					err = dot.SError.NoDotNewer.AddNewError(m.TypeId.String())
					break CreateLives
				}

				dotLive.Dot, err = m.NewDot(bytesConfig)
				if err == nil {
					if err = createDotFun(dotLive); err != nil {
						break CreateLives
					}
					continue CreateLives
				} else {
					break CreateLives
				}
			}
		}
	}

	if err != nil {
		logger.Debug(func() string {
			m, _ := c.metas.Get(dotLive.TypeId)
			if m != nil {
				return fmt.Sprintf("Create dot, meta: %v\n live: %v", m, dotLive)
			}
			return fmt.Sprintf("Create dot, live: %v", dotLive)
		})
		logger.Errorln("lineImp", zap.Error(err))
		return err
	}
	//Add logger and config
	{
		c.mutex.Lock()
		c.types[reflect.TypeOf(c.logger)] = c.logger
		{
			t := reflect.TypeOf((*dot.SLogger)(nil)).Elem()
			c.types[t] = c.logger
		}
		c.types[reflect.TypeOf(c.config)] = c.config
		{
			t := reflect.TypeOf((*dot.SConfig)(nil)).Elem()
			c.types[t] = c.sConfig
		}
		c.mutex.Unlock()
	}

	//Add type and relationships with dot, only record whose type id == live Id
	for _, it := range orderedDots {
		if !skit.IsNil(&it.Dot) && ((string)(it.TypeId) == (string)(it.LiveId)) {
			t := reflect.TypeOf(it.Dot)
			c.mutex.Lock()
			c.types[t] = it.Dot
			//if the dot implements the dot.GetInterfaceType, put the interface into types too
			if getter, ok := it.Dot.(dot.GetInterfaceType); ok {
				c.types[getter.GetInterfaceType()] = it.Dot
			}
			c.mutex.Unlock()
		}
	}

	{
		afterInjects := make([]dot.AfterAllInjecter, 0, 20)
		err = nil
		for _, it := range orderedDots {
			if it.Dot != nil { //todo not success
				_ = c.injectInLine(it.Dot, it)
				if ed, ok := it.Dot.(dot.Injected); ok {
					err = ed.Injected(c)
					if err != nil {
						logger.Debug(func() string {
							m, _ := c.metas.Get(it.TypeId)
							if m != nil {
								return fmt.Sprintf("Create dot, meta: %v\n live: %v", m, it)
							} else {
								return fmt.Sprintf("Create dot, live: %v", it)
							}
						})
						logger.Errorln("lineImp", zap.Error(err))
						break
					}
				}
				if s, ok := it.Dot.(dot.AfterAllInjecter); ok {
					afterInjects = append(afterInjects, s)
				}
			}
		}
		if err == nil {
			for _, v := range afterInjects {
				v.AfterAllInject(c)
			}
		}
	}

	return err
}

func (c *lineImp) Config() *dot.Config {
	return &c.config
}

func (c *lineImp) SLogger() dot.SLogger {
	return c.logger
}

func (c *lineImp) SConfig() dot.SConfig {
	return c.sConfig
}

func (c *lineImp) ToLifer() dot.Lifer {
	return c
}

//ToInjecter to injecter
func (c *lineImp) ToInjecter() dot.Injecter {
	return c
}

func (c *lineImp) ToDotEventer() dot.Eventer {
	return &c.dotEvent
}

func (c *lineImp) GetDotConfig(liveId dot.LiveId) *dot.LiveConfig {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var co *dot.LiveConfig
	co = c.config.FindConfig("", liveId)
	return co
}

/////injecter
//Inject see https://github.com/facebookgo/inject
func (c *lineImp) Inject(obj interface{}) error {
	logger := dot.Logger()
	var err error
	if skit.IsNil(obj) {
		return dot.SError.NilParameter
	}
	multiErr := func(er error) {
		if err != nil {
			logger.Errorln("lineImp", zap.Error(err))
		}
		err = er
	}

	value := reflect.ValueOf(obj)

	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return dot.SError.NotStruct
	}

	rType := value.Type()

	for i := 0; i < value.NumField(); i++ {
		var err2 error
		valueField := value.Field(i)
		if !valueField.CanSet() {
			continue
		}

		typeField := rType.Field(i)
		tagName, ok := typeField.Tag.Lookup(dot.TagDot)
		if !ok {
			continue
		}

		var d dot.Dot
		{
			if len(tagName) < 1 || tagName == dot.CanNull { //by type
				d, err2 = c.GetByType(valueField.Type())
			} else { //by live id
				liveId := strings.TrimPrefix(tagName, dot.CanNull) // like ?live id
				d, err2 = c.GetByLiveId(dot.LiveId(liveId))
			}

			if err2 != nil {
				if tagName == dot.CanNull { //组件为可选
					if dotError, ok := err2.(dot.Errorer); !ok || dotError.Code() != dot.SError.NotExisted.Code() { //如果error code不为不存在
						logger.Debugln(fmt.Sprintf("lineImp, find dot error, field: %s, err: %value", typeField.Name, err2))
						multiErr(err2)
					}
				} else {
					logger.Debugln(fmt.Sprintf("lineImp, find dot error, field: %s, err: %value", typeField.Name, err2))
					multiErr(err2)
				}
			}

			if d == nil {
				logger.Debugln(fmt.Sprintf("lineImp, can not find dot error, field: %s", typeField.Name))
				continue
			}
		}

		if err2 == nil {
			valueDot := reflect.ValueOf(d)
			if valueDot.IsValid() && valueDot.Type().AssignableTo(valueField.Type()) {
				valueField.Set(valueDot)
			} else {
				multiErr(dot.SError.DotInvalid.AddNewError(typeField.Type.String() + "  " + tagName))
			}
		}
	}

	return err
}

// like Inject func
func (c *lineImp) injectInLine(obj interface{}, live *dot.Live) error {
	logger := dot.Logger()
	var err error
	if skit.IsNil(obj) {
		return dot.SError.NilParameter
	}
	multiErr := func(er error) {
		if err != nil {
			logger.Errorln("lineImp", zap.Error(err))
		}
		err = er
	}

	value := reflect.ValueOf(obj)

	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return dot.SError.NotStruct
	}

	rType := value.Type()

	for i := 0; i < value.NumField(); i++ {
		var err2 error = nil
		valueField := value.Field(i)
		if !valueField.CanSet() {
			continue
		}

		tField := rType.Field(i)
		tagName, ok := tField.Tag.Lookup(dot.TagDot)
		if !ok {
			continue
		}

		var d dot.Dot
		{
			if len(live.RelyLives) > 0 { //Config prior
				if liveId, ok := live.RelyLives[tField.Name]; ok {
					d, err2 = c.GetByLiveId(liveId)
				}
			}
			if d == nil {
				if len(tagName) < 1 || tagName == dot.CanNull { //by type
					d, err2 = c.GetByType(valueField.Type())
				} else { //by live id
					liveId := strings.TrimPrefix(tagName, dot.CanNull) // like ?live id
					d, err2 = c.GetByLiveId(dot.LiveId(liveId))
				}
			}

			if err2 != nil {
				if tagName == dot.CanNull { //组件为可选
					if dotError, ok := err2.(dot.Errorer); !ok || dotError.Code() != dot.SError.NotExisted.Code() { //如果error code不为不存在
						logger.Debugln(fmt.Sprintf("lineImp, find dot error, field: %s, err: %value", tField.Name, err2))
						multiErr(err2)
					}
				} else {
					logger.Debugln(fmt.Sprintf("lineImp, find dot error, field: %s, err: %value", tField.Name, err2))
					multiErr(err2)
				}
			}

			if d == nil {
				logger.Debugln(fmt.Sprintf("lineImp, can not find dot error, field: %s, live: %value", tField.Name, live))
				continue
			}
		}

		if err2 == nil {
			valueDot := reflect.ValueOf(d)
			//fmt.Println("valueDot: ", valueDot.Type(), "valueField: ", valueField.Type(), "dd: ", reflect.TypeOf(d))
			if valueDot.IsValid() && valueDot.Type().AssignableTo(valueField.Type()) {
				valueField.Set(valueDot)
			} else if err == nil {
				multiErr(dot.SError.DotInvalid.AddNewError(tField.Type.String() + "  " + tagName))
			}
		}
	}

	return err
}

//GetByType get by type
func (c *lineImp) GetByType(t reflect.Type) (d dot.Dot, err error) {
	d = nil
	err = nil
	c.mutex.Lock()
	d, ok := c.types[t]
	c.mutex.Unlock()
	if !ok {
		if c.parent != nil {
			d, err = c.parent.GetByType(t)
		} else {
			err = dot.SError.NotExisted.AddNewError(t.String())
		}
	}

	return
}

//GetByLiveId get by live id
func (c *lineImp) GetByLiveId(liveId dot.LiveId) (d dot.Dot, err error) {
	d = nil
	err = nil
	c.mutex.Lock()
	var l *dot.Live
	l, err = c.lives.Get(liveId)
	c.mutex.Unlock()
	if err != nil {
		if c.parent != nil {
			d, err = c.parent.GetByLiveId(liveId)
		}
	} else {
		d = l.Dot
	}

	return
}

//ReplaceOrAddByType update
func (c *lineImp) ReplaceOrAddByType(d dot.Dot) error {
	var err error
	typeDot := reflect.TypeOf(d)
	//for typeDot.Kind() == reflect.Ptr || typeDot.Kind() == reflect.Interface {
	//	typeDot = typeDot.Elem()
	//}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.types[typeDot] = d
	return err
}

//ReplaceOrAddByParamType update
func (c *lineImp) ReplaceOrAddByParamType(d dot.Dot, rType reflect.Type) error {
	var err error
	//for rType.Kind() == reflect.Ptr || rType.Kind() == reflect.Interface {
	//	rType = rType.Elem()
	//}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.types[rType] = d
	return err
}

//ReplaceOrAddByLiveId update
func (c *lineImp) ReplaceOrAddByLiveId(d dot.Dot, liveId dot.LiveId) error {
	var err error
	c.mutex.Lock()
	defer c.mutex.Unlock()

	l := dot.Live{LiveId: liveId, TypeId: "", Dot: d, RelyLives: nil}
	//c.lives.Remove(&l)
	err = c.lives.UpdateOrAdd(&l)

	return err
}

//RemoveByType remove
func (c *lineImp) RemoveByType(rType reflect.Type) error {
	var err error
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.types, rType)
	return err
}

//RemoveByLiveId remove
func (c *lineImp) RemoveByLiveId(liveId dot.LiveId) error {
	var err error
	c.mutex.Lock()
	defer c.mutex.Unlock()
	err = c.lives.RemoveById(liveId)
	return err
}

//SetParent set parent injecter
func (c *lineImp) SetParent(parent dot.Injecter) {
	c.parent = parent
}

//GetParent get parent injecter
func (c *lineImp) GetParent() dot.Injecter {
	return c.parent
}

////injecter end

//Create create
//If live id is empty， directly assign type id
//If live id repeated，directly return dot.SError.ErrExistedLiveId
func (c *lineImp) Create(l dot.Line) error {
	var err error

	//first create config
	c.sConfig = sconfig.NewConfiger()
	c.sConfig.RootPath()
	if s, ok := c.sConfig.(dot.Creator); ok {
		if err = s.Create(l); err != nil {
			createLog(c)
			return err
		} else if len(c.sConfig.ConfigFile()) < 1 { //no config file return
			createLog(c)
			return err
		}
	}

	if err = c.sConfig.Unmarshal(&c.config); err != nil {
		createLog(c)
		return err
	}
	if len(c.config.Dots) < 1 { //no config
		createLog(c)
		return err
	}

	//create log
	createLog(c)
	return err
}

func (c *lineImp) makeDotMetaFromConfig() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var err error
	var dotConfig *dot.DotConfig
ForConfigDots: //handle config
	for i := range c.config.Dots {
		dotConfig = &c.config.Dots[i]
		if len(dotConfig.MetaData.TypeId.String()) < 1 {
			err = dot.SError.Config.AddNewError("type id is null")
			break ForConfigDots
		}

		if err = c.metas.UpdateOrAdd(&dotConfig.MetaData); err != nil {
			break ForConfigDots
		}

		if len(dotConfig.Lives) < 1 { //create the single live
			live := dot.Live{TypeId: dotConfig.MetaData.TypeId, LiveId: dot.LiveId(dotConfig.MetaData.TypeId), Dot: nil}
			if len(dotConfig.MetaData.RelyTypeIds) > 0 {
				live.RelyLives = make(map[string]dot.LiveId, len(dotConfig.MetaData.RelyTypeIds))
				for i := range dotConfig.MetaData.RelyTypeIds {
					li := &dotConfig.MetaData.RelyTypeIds[i]
					live.RelyLives[li.String()] = dot.LiveId(*li)
				}
			}
			if err = c.lives.UpdateOrAdd(&live); err != nil {
				break ForConfigDots
			}
		} else {
			for _, li := range dotConfig.Lives {
				if len(li.LiveId.String()) < 1 {
					li.LiveId = dot.LiveId(dotConfig.MetaData.TypeId)
				}
				live := dot.Live{TypeId: dotConfig.MetaData.TypeId, LiveId: li.LiveId, RelyLives: li.RelyLives, Dot: nil}
				if err = c.lives.UpdateOrAdd(&live); err != nil {
					break ForConfigDots
				}
			}
		}
	}
	if err != nil {
		dot.Logger().Debugln(fmt.Sprintf("lineImp, %v", dotConfig))
	}

	return err
}

//case #17
func (c *lineImp) autoMakeLiveId() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	hasType := make(map[dot.TypeId]bool, len(c.lives.LiveIdMap))
	for _, v := range c.lives.LiveIdMap {
		hasType[v.TypeId] = true
	}

	for typeId := range c.metas.metas {
		if _, ok := hasType[typeId]; !ok {
			liveId := (dot.LiveId)(typeId)
			live := dot.Live{TypeId: typeId, LiveId: liveId, Dot: nil, RelyLives: nil}
			err2 := c.lives.UpdateOrAdd(&live)
			if err2 != nil && c.logger != nil {
				c.logger.Debugln(fmt.Sprintf("err: %v, live: %v", err2, live))
			}
		}
	}
}

//get relay from the tag
//merge type relay
//merge live relay
//verify the key of relay lives(the name have to eq the field name)
func (c *lineImp) makeRelays() {
	var setExists = struct{}{} //only for set value
	c.mutex.Lock()
	defer c.mutex.Unlock()
	logger := dot.Logger()
	liveIdToType := make(map[dot.LiveId]dot.TypeId)
	typeIdToLiveIds := make(map[dot.TypeId]map[dot.LiveId]struct{}, len(c.lives.LiveIdMap))
	for _, live := range c.lives.LiveIdMap {
		liveIdToType[live.LiveId] = live.TypeId
		if _, ok := typeIdToLiveIds[live.TypeId]; !ok {
			typeIdToLiveIds[live.TypeId] = make(map[dot.LiveId]struct{})
		}
		typeIdToLiveIds[live.TypeId][live.LiveId] = setExists
	}
	typeToTypeId := make(map[reflect.Type]dot.TypeId)
	for _, meta := range c.metas.metas {
		if meta.RefType != nil {
			typeToTypeId[meta.RefType] = meta.TypeId
		}
	}
	for _, meta := range c.metas.metas { //relay type
		dotFields := make(map[string]dot.LiveId)
		if meta.RefType != nil {
			relayTypeIds := make(map[dot.TypeId]struct{})
			for i := 0; i < meta.RefType.NumField(); i++ { // relay by tag
				tField := meta.RefType.Field(i)
				liveId, ok := tField.Tag.Lookup(dot.TagDot)
				if !ok || liveId == dot.CanNull { //not find or eq ?
					continue
				}
				dotFields[tField.Name] = dot.LiveId(liveId)
				if len(liveId) < 1 {
					if typeId, ok := typeToTypeId[tField.Type]; ok {
						relayTypeIds[typeId] = setExists
					} else {
						//do not find typeId
					}
				} else {
					if typeId, ok := liveIdToType[dot.LiveId(liveId)]; ok {
						relayTypeIds[typeId] = setExists
					}
				}
			}

			if lives, ok := typeIdToLiveIds[meta.TypeId]; ok {
				//for liveId := range lives { relay type and relay live, sometimes they are not same
				//	if live, ok := c.lives.LiveIdMap[liveId]; ok {
				//		relayTypeIds[live.TypeId] = setExists
				//	}
				//}

				//verify the field name for relay live
				for liveId := range lives {
					if live, ok := c.lives.LiveIdMap[liveId]; ok {
						for name := range live.RelyLives {
							if _, ok := dotFields[name]; !ok {
								if logger != nil {
									logger.Warn(func() string {
										return fmt.Sprintf("make relay of dot, meta: %v\n live: %v", meta, live)
									})
								}
							}
						}

						for name, nameLiveId := range dotFields {
							if _, ok := live.RelyLives[name]; !ok {
								if len(nameLiveId.String()) > 0 {
									live.RelyLives[name] = nameLiveId
								} else {
									//do nothing, do not know live id
								}
							}
						}
					}
				}
			}

			if len(relayTypeIds) > 0 {
				for _, typeId := range meta.RelyTypeIds { //distinct
					relayTypeIds[typeId] = setExists
				}
				meta.RelyTypeIds = make([]dot.TypeId, 0, len(relayTypeIds))
				for typeId := range relayTypeIds {
					meta.RelyTypeIds = append(meta.RelyTypeIds, typeId)
				}
			}
		}
	}
}

func createLog(c *lineImp) {
	c.logger = slog.NewSLogger(&(c.config.Log), c)
	dot.SetLogger(c.logger)
}

//Start
func (c *lineImp) Start(ignore bool) error {
	var err error
	logger := dot.Logger()

	for {
		//start config
		if s, ok := c.sConfig.(dot.Starter); ok {
			if err = s.Start(ignore); err != nil {
				break
			}
		}
		//start log
		if s, ok := c.logger.(dot.Starter); ok {
			if err = s.Start(ignore); err != nil {
				break
			}
		}

		//start other
		{
			//recount the order, maybe the "Ceate" change it
			tdots, _ := c.relyOrder() //do not care the circle
			afterStarts := make([]dot.AfterAllStarter, 0, 20)
			for _, it := range tdots {
				logger.Debug(func() string {
					m, _ := c.metas.Get(it.TypeId)
					if m != nil {
						return fmt.Sprintf("Start dot, type id: %s, live id: %s, name: %s", it.TypeId, it.LiveId, m.Name)
					} else {
						return fmt.Sprintf("Start dot, type id: %s, live id: %s", it.TypeId, it.LiveId)
					}
				})
				if b := c.dotEvent.TypeEvents(it.TypeId); len(b) > 0 {
					for i := range b {
						e := &b[i]
						if e.BeforeStart != nil {
							e.BeforeStart(it, c)
						}
					}

				}

				if b := c.dotEvent.LiveEvents(it.LiveId); len(b) > 0 {
					for i := range b {
						e := &b[i]
						if e.BeforeStart != nil {
							e.BeforeStart(it, c)
						}
					}
				}

				if d, ok := it.Dot.(dot.Starter); ok {
					err2 := d.Start(ignore)
					if err2 != nil {
						logger.Debug(func() string {
							m, _ := c.metas.Get(it.TypeId)
							if m != nil {
								return fmt.Sprintf("Start dot, meta: %v\n live: %v\n %v", m, it, d)
							} else {
								return fmt.Sprintf("Start dot, live: %v\n %v", it, d)
							}
						})
						if err != nil {
							logger.Errorln("lineImp", zap.Error(err))
						}
						err = err2
						if !ignore {
							return err
						}
					}
				}

				if a := c.dotEvent.LiveEvents(it.LiveId); len(a) > 0 {
					for i := range a {
						e := &a[i]
						if e.AfterStart != nil {
							e.AfterStart(it, c)
						}
					}
				}

				if a := c.dotEvent.TypeEvents(it.TypeId); len(a) > 0 {
					for i := range a {
						e := &a[i]
						if e.AfterStart != nil {
							e.AfterStart(it, c)
						}
					}
				}

				if s, ok := it.Dot.(dot.AfterAllStarter); ok {
					afterStarts = append(afterStarts, s)
				}
			}

			for _, s := range afterStarts {
				s.AfterAllStart(c)
			}
		}

		break
	}

	return err
}

//Stop
func (c *lineImp) Stop(ignore bool) error {
	var err error
	logger := dot.Logger()
	//stop others
	{
		//recount the order, maybe the "Ceate" change it
		tdots, _ := c.relyOrder() //do not care the circle
		{
			beforeStops := make([]dot.BeforeAllStopper, 0, 20)

			for i := len(tdots) - 1; i >= 0; i-- {
				if b, ok := tdots[i].Dot.(dot.BeforeAllStopper); ok {
					beforeStops = append(beforeStops, b)
				}
			}

			for _, it := range beforeStops {
				it.BeforeAllStop(c)
			}
		}

		for i := len(tdots) - 1; i >= 0; i-- {
			it := tdots[i]
			logger.Debug(func() string {
				m, _ := c.metas.Get(it.TypeId)
				if m != nil {
					return fmt.Sprintf("Stop dot, type id: %s, live id: %s, name: %s", it.TypeId, it.LiveId, m.Name)
				} else {
					return fmt.Sprintf("Stop dot, type id: %s, live id: %s", it.TypeId, it.LiveId)
				}
			})
			if b := c.dotEvent.TypeEvents(it.TypeId); len(b) > 0 {
				for i := range b {
					e := &b[i]
					if e.BeforeStop != nil {
						e.BeforeStop(it, c)
					}
				}
			}

			if b := c.dotEvent.LiveEvents(it.LiveId); len(b) > 0 {
				for i := range b {
					e := &b[i]
					if e.BeforeStop != nil {
						e.BeforeStop(it, c)
					}
				}
			}

			if d, ok := it.Dot.(dot.Stopper); ok {
				err2 := d.Stop(ignore)
				if err2 != nil {
					logger.Debugln(fmt.Sprintf("lineImp, Stop dot: %v", d))
					if err != nil {
						logger.Errorln("", zap.Error(err))
					}
					err = err2
				}

				if !ignore {
					return err
				}
			}

			if a := c.dotEvent.LiveEvents(it.LiveId); len(a) > 0 {
				for i := range a {
					e := &a[i]
					if e.AfterStop != nil {
						e.AfterStop(it, c)
					}
				}
			}

			if a := c.dotEvent.TypeEvents(it.TypeId); len(a) > 0 {
				for i := range a {
					e := &a[i]
					if e.AfterStop != nil {
						e.AfterStop(it, c)
					}
				}
			}
		}
	}
	//stop log
	if d, ok := c.logger.(dot.Stopper); ok {
		err2 := d.Stop(ignore)
		if err2 != nil {
			logger.Debugln(fmt.Sprintf("lineImp, Stop dot: %v", d))
			if err != nil {
				logger.Errorln("", zap.Error(err))
			}
			err = err2
		}
	}

	//stop config
	if d, ok := c.sConfig.(dot.Stopper); ok {
		err2 := d.Stop(ignore)
		if err2 != nil {
			logger.Debugln(fmt.Sprintf("lineImp, Stop dot: %v", d))
			if err != nil {
				logger.Errorln("", zap.Error(err))
			}
			err = err2
		}
	}

	return err
}

//Destroy Destroy Dot
func (c *lineImp) Destroy(ignore bool) error {
	var err error
	logger := dot.Logger()
	//Destroy others
	{
		afterAll := make([]dot.AfterAllDestroyer, 0, 20)
		//recount the order, maybe the "Ceate" change it
		tdots, _ := c.relyOrder() //do not care the circle
		for i := len(tdots) - 1; i >= 0; i-- {
			it := tdots[i]
			logger.Debug(func() string {
				m, _ := c.metas.Get(it.TypeId)
				if m != nil {
					return fmt.Sprintf("Destroy dot, type id: %s, live id: %s, name: %s", it.TypeId, it.LiveId, m.Name)
				} else {
					return fmt.Sprintf("Destroy dot, type id: %s, live id: %s", it.TypeId, it.LiveId)
				}
			})
			if b := c.dotEvent.TypeEvents(it.TypeId); len(b) > 0 {
				for i := range b {
					e := &b[i]
					if e.BeforeDestroy != nil {
						e.BeforeDestroy(it, c)
					}
				}
			}

			if b := c.dotEvent.LiveEvents(it.LiveId); len(b) > 0 {
				for i := range b {
					e := &b[i]
					if e.BeforeDestroy != nil {
						e.BeforeDestroy(it, c)
					}
				}
			}

			if d, ok := it.Dot.(dot.Destroyer); ok {
				err2 := d.Destroy(ignore)
				if err2 != nil {
					logger.Debugln(fmt.Sprintf("lineImp, Destroy dot: %v", d))
					if err != nil {
						logger.Errorln("lineImp", zap.Error(err))
					}
					err = err2
				}
				if !ignore {
					return err
				}
			}

			if a := c.dotEvent.LiveEvents(it.LiveId); len(a) > 0 {
				for i := range a {
					e := &a[i]
					if e.AfterDestroy != nil {
						e.AfterDestroy(it, c)
					}
				}
			}

			if a := c.dotEvent.TypeEvents(it.TypeId); len(a) > 0 {
				for i := range a {
					e := &a[i]
					if e.AfterDestroy != nil {
						e.AfterDestroy(it, c)
					}
				}
			}

			if all, ok := it.Dot.(dot.AfterAllDestroyer); ok {
				afterAll = append(afterAll, all)
			}

		}

		for _, it := range afterAll {
			it.AfterAllDestroy(c)
		}
	}

	//Destroy log
	if d, ok := c.logger.(dot.Destroyer); ok {
		err2 := d.Destroy(ignore)
		if err2 != nil {
			logger.Debugln(fmt.Sprintf("lineImp, Destroy dot: %v", d))
			if err != nil {
				logger.Errorln("lineImp", zap.Error(err))
			}
			err = err2
		}
	}

	//Destroy config
	if d, ok := c.sConfig.(dot.Destroyer); ok {
		err2 := d.Destroy(ignore)
		if err2 != nil {
			logger.Debugln(fmt.Sprintf("lineImp, Destroy dot: %v", d))
			if err != nil {
				logger.Errorln("lineImp", zap.Error(err))
			}
			err = err2
		}
	}

	return err
}

func (c *lineImp) GetLineBuilder() *dot.Builder {
	return c.lineBuilder
}
func (c *lineImp) InfoAllTypeAdnLives() {
	logger := c.logger
	logger.Info(func() string {
		return fmt.Sprintf("lives - %d: %v, types - %d: %v", len(c.lives.LiveIdMap), c.lives.LiveIdMap, len(c.types), c.types)
	})
}

func (c *lineImp) EachLives(call func(live *dot.Live, meta *dot.Metadata) bool) {
	if call != nil {
		c.mutex.Lock()
		liveIds := make([]dot.LiveId, 0, len(c.lives.LiveIdMap))
		for liveId := range c.lives.LiveIdMap {
			liveIds = append(liveIds, liveId)
		}
		typeIds := make([]dot.TypeId, 0, len(c.metas.metas))
		for typeId := range c.metas.metas {
			typeIds = append(typeIds, typeId)
		}
		c.mutex.Unlock()
		for _, liveId := range liveIds {
			var live *dot.Live
			var meta *dot.Metadata
			c.mutex.Lock()
			live = c.lives.LiveIdMap[liveId] //if the key do not exist, return nil
			if live != nil {
				meta = c.metas.metas[live.TypeId] //if the key do not exist, return nil
			}
			c.mutex.Unlock() //unlock mutex to avoid the dead lock
			if !call(live, meta) {
				break
			}
		}
	}
}

///////////////
