// Scry Info.  All rights reserved.
// license that can be found in the license file.

package line

import (
	"fmt"
	"reflect"
	"sync"

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
	logger      dot.SLogger
	sConfig     dot.SConfig //External general config
	config      dot.Config  //Config object of component
	metas       *Metas
	lives       *Lives
	types       map[reflect.Type]dot.Dot
	newerLiveid map[dot.LiveId]dot.Newer
	newerTypeid map[dot.TypeId]dot.Newer

	parent dot.Injecter
	mutex  sync.Mutex

	lineBuilder *dot.Builder

	//dot event
	dotEventer dotEventerImp
}

//newLine new
func newLine(builer *dot.Builder) *lineImp {
	a := &lineImp{metas: NewMetas(),
		lives: NewLives(), types: make(map[reflect.Type]dot.Dot),
		newerLiveid: make(map[dot.LiveId]dot.Newer),
		newerTypeid: make(map[dot.TypeId]dot.Newer),
		lineBuilder: builer,
	}
	a.dotEventer.Init()

	if dot.GetDefaultLine() == nil {
		dot.SetDefaultLine(a)
	}

	return a
}

func (c *lineImp) Id() string {
	return c.lineBuilder.LineLiveId
}

//AddNewerByLiveId add new for liveid
func (c *lineImp) AddNewerByLiveId(liveid dot.LiveId, newDot dot.Newer) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.newerLiveid[liveid]; ok {
		return dot.SError.Existed.AddNewError(liveid.String())
	}

	c.newerLiveid[liveid] = newDot

	return nil
}

//AddNewerByTypeId add new for type
func (c *lineImp) AddNewerByTypeId(typeid dot.TypeId, newDot dot.Newer) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, ok := c.newerTypeid[typeid]; ok {
		return dot.SError.Existed.AddNewError(typeid.String())
	}

	c.newerTypeid[typeid] = newDot

	return nil
}

//RemoveNewerByLiveId remove
func (c *lineImp) RemoveNewerByLiveId(liveid dot.LiveId) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.newerLiveid, liveid)
}

//RemoveNewerByTypeId remove
func (c *lineImp) RemoveNewerByTypeId(typeid dot.TypeId) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.newerTypeid, typeid)
}

//PreAdd the dot is nil, do not create it
func (c *lineImp) PreAdd(typeLives ...*dot.TypeLives) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var err error

	for _, clone := range typeLives {

		err2 := c.metas.UpdateOrAdd(&clone.Meta)
		if err2 == nil {
			if clone.Lives != nil {
				for i := range clone.Lives {
					it := &clone.Lives[i]
					if len(it.TypeId.String()) < 1 {
						it.TypeId = clone.Meta.TypeId
					}

					//live := dot.Live{TypeId: it.TypeId, LiveId: it.LiveId, Dot: nil}
					//live.RelyLives = CloneRelyLiveId(it.RelyLives)
					c.lives.UpdateOrAdd(it)
				}
			} else {
				lid := (dot.LiveId)(clone.Meta.TypeId)
				live := dot.Live{TypeId: clone.Meta.TypeId, LiveId: lid, Dot: nil, RelyLives: nil}
				c.lives.UpdateOrAdd(&live)
			}
		} else {
			if err != nil {
				dot.Logger().Errorln(err.Error()) //write it into logfile, otherwise it will gone
				err = err2
			}
		}
	}

	return err
}

func (c *lineImp) RelyOrder() ([]*dot.Live, []*dot.Live) {

	var cloneLives map[dot.LiveId]*dot.Live
	var cloneMetas map[dot.TypeId]*dot.Metadata
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

			for _, lid := range it.RelyLives {
				r, ok := relyed[lid]
				if !ok {
					r = []dot.LiveId{it.LiveId}
				} else {
					r = append(r, it.LiveId)
				}
				relyed[lid] = r
			}
		}
		done := make(map[dot.LiveId]bool, len(relyed))          //know orderId
		remain := make(map[dot.LiveId]bool, len(relyed))        //do not know orderId
		levels := make([]map[dot.LiveId]bool, 0, len(relyed)/3) //
		{
			for lid, _ := range relyed {
				remain[lid] = true
			}
			level0 := make(map[dot.LiveId]bool)
			for lid, lids := range cloneLives { //level 0
				if len(lids.RelyLives) < 1 {
					level0[lid] = true
					delete(remain, lid)
					done[lid] = true
				}
			}
			levels = append(levels, level0)
			levelCurrent := level0
			//todo if level0 is zero
			for curFor := 0; curFor <= len(relyed); curFor++ { //
				levelNext := make(map[dot.LiveId]bool, len(remain))
				for lid, _ := range levelCurrent {
					des := relyed[lid]
					if len(des) < 1 {
						delete(remain, lid)
						done[lid] = true
					} else {
						for _, lid2 := range des {
							alldone := true
							for _, lid3 := range cloneLives[lid2].RelyLives { //check all RelyLives
								if _, ok := done[lid3]; !ok {
									alldone = false
									break
								}
							}
							if alldone {
								levelNext[lid2] = true
								delete(remain, lid2)
								done[lid2] = true
							} else {
								levelNext[lid] = true //put next level
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

		//todo type dependency

		{
			for _, lev := range levels {
				for lid, _ := range lev {
					order = append(order, cloneLives[lid])
				}
			}
			if len(remain) > 0 {
				circle = make([]*dot.Live, 0, len(remain))
				for lid, _ := range remain { //append to tail
					order = append(order, cloneLives[lid])
					circle = append(circle, cloneLives[lid])
				}
			}
		}
	}

	return order, circle
}

//CreateDots create dots
func (c *lineImp) CreateDots(order []*dot.Live) error {
	tdots := order
	creator := func(it *dot.Live) error {
		{ // Check whether special info needed before Create
			if nl, ok := it.Dot.(dot.SetterLine); ok {
				nl.SetLine(c)
			}

			if nl, ok := it.Dot.(dot.SetterTypeAndLiveId); ok {
				nl.SetTypeId(it.TypeId, it.LiveId)
			}
		}

		if b := c.dotEventer.TypeEvents(it.TypeId); len(b) > 0 { // dot not care the dot.Creator
			for i := range b {
				e := &b[i]
				if e.BeforeCreate != nil {
					e.BeforeCreate(it, c)
				}
			}
		}

		if b := c.dotEventer.LiveEvents(it.LiveId); len(b) > 0 { // dot not care the dot.Creator
			for i := range b {
				e := &b[i]
				if e.BeforeCreate != nil {
					e.BeforeCreate(it, c)
				}
			}
		}

		if creator, ok := it.Dot.(dot.Creator); ok {
			if err := creator.Create(c); err != nil {
				return err
			}
		}

		if a := c.dotEventer.LiveEvents(it.LiveId); len(a) > 0 { // dot not care the dot.Creator
			for i := range a {
				e := &a[i]
				if e.AfterCreate != nil {
					e.AfterCreate(it, c)
				}
			}
		}

		if a := c.dotEventer.TypeEvents(it.TypeId); len(a) > 0 { // dot not care the dot.Creator
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
LIVES:
	for _, it := range tdots {

		if skit.IsNil(&it.Dot) == true {
			var bconfig []byte
			var config *dot.LiveConfig
			if true {
				config = c.config.FindConfig(it.TypeId, it.LiveId)
				bconfig, err = dot.MarshalConfig(config)
				if err != nil {
					break LIVES
				}
			}
			//liveid
			{
				if newer, ok := c.newerLiveid[it.LiveId]; ok {
					it.Dot, err = newer(bconfig)
					if err != nil {
						break LIVES
					} else {
						if err = creator(it); err != nil {
							break LIVES
						}
						continue LIVES
					}
				}
			}
			//typeid
			{
				if newer, ok := c.newerTypeid[it.TypeId]; ok {
					it.Dot, err = newer(bconfig)
					if err != nil {
						break LIVES
					} else {
						if err = creator(it); err != nil {
							break LIVES
						}
						continue LIVES
					}
				}
			}

			//metadata
			{
				var m *dot.Metadata
				m, err = c.metas.Get(it.TypeId)
				if err != nil {
					break LIVES
				}

				if m.NewDoter == nil && m.RefType == nil {
					err = dot.SError.NoDotNewer.AddNewError(m.TypeId.String())
					break LIVES
				}

				it.Dot, err = m.NewDot(bconfig)
				if err == nil {
					if err = creator(it); err != nil {
						break LIVES
					}
					continue LIVES
				} else {
					break LIVES
				}
			}
		}
	}

	if err != nil {
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

	//Add type and relationships with dot, only record whose typeid == liveId
	for _, it := range tdots {
		if !skit.IsNil(&it.Dot) && ((string)(it.TypeId) == (string)(it.LiveId)) {
			t := reflect.TypeOf(it.Dot)
			c.mutex.Lock()
			c.types[t] = it.Dot
			c.mutex.Unlock()
		}
	}

	for _, it := range tdots {
		if it.Dot != nil { //todo not success
			c.injectInLine(it.Dot, it)
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
	return &c.dotEventer
}

func (c *lineImp) GetDotConfig(liveid dot.LiveId) *dot.LiveConfig {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var co *dot.LiveConfig
	co = c.config.FindConfig("", liveid)
	return co
}

/////injecter

//Inject see https://github.com/facebookgo/inject
func (c *lineImp) Inject(obj interface{}) error {
	var err error
	if skit.IsNil(obj) {
		return dot.SError.NilParameter
	}

	v := reflect.ValueOf(obj)

	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return dot.SError.NotStruct
	}

	t := v.Type()

	var errt error
	for i := 0; i < v.NumField(); i++ {
		errt = nil
		f := v.Field(i)
		if !f.CanSet() {
			continue
		}

		tField := t.Field(i)
		tname, ok := tField.Tag.Lookup(dot.TagDot)
		if !ok {
			continue
		}

		var d dot.Dot
		{
			if len(tname) < 1 { //by type
				d, errt = c.GetByType(f.Type())
			} else { //by liveid
				d, errt = c.GetByLiveId(dot.LiveId(tname))
			}

			if errt != nil && err == nil {
				err = errt
				fmt.Println("err:", err.Error())
			}

			if d == nil {
				fmt.Println("can not find the dot tname:{}", tname)
				continue
			}
		}

		if errt == nil {
			vv := reflect.ValueOf(d)
			//fmt.Println("vv: ", vv.Type(), "f: ", f.Type(), "dd: ", reflect.TypeOf(d))
			if vv.IsValid() && vv.Type().AssignableTo(f.Type()) {
				f.Set(vv)
			} else if err == nil {
				err = dot.SError.DotInvalid.AddNewError(tField.Type.String() + "  " + tname)
			}
		}
	}

	return err
}

func (c *lineImp) injectInLine(obj interface{}, live *dot.Live) error {
	var err error
	if skit.IsNil(obj) {
		return dot.SError.NilParameter
	}

	v := reflect.ValueOf(obj)

	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return dot.SError.NotStruct
	}

	t := v.Type()

	var errt error
	for i := 0; i < v.NumField(); i++ {
		errt = nil
		f := v.Field(i)
		if !f.CanSet() {
			continue
		}

		tField := t.Field(i)
		tname, ok := tField.Tag.Lookup(dot.TagDot)
		if !ok {
			continue
		}

		var d dot.Dot
		{
			if len(live.RelyLives) > 0 { //Config prior
				if lid, ok := live.RelyLives[tField.Name]; ok {
					d, errt = c.GetByLiveId(dot.LiveId(lid))
				}
			}
			if d == nil {
				if len(tname) < 1 { //by type
					d, errt = c.GetByType(f.Type())
				} else { //by liveid
					d, errt = c.GetByLiveId(dot.LiveId(tname))
				}
			}

			if errt != nil && err == nil {
				err = errt
				fmt.Println("err:", err.Error())
			}

			if d == nil {
				fmt.Println("can not find the dot tname:{}", tname)
				continue
			}
		}

		if errt == nil {
			vv := reflect.ValueOf(d)
			//fmt.Println("vv: ", vv.Type(), "f: ", f.Type(), "dd: ", reflect.TypeOf(d))
			if vv.IsValid() && vv.Type().AssignableTo(f.Type()) {
				f.Set(vv)
			} else if err == nil {
				err = dot.SError.DotInvalid.AddNewError(tField.Type.String() + "  " + tname)
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

//GetByLiveId get by liveid
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
	t := reflect.TypeOf(d)
	//for t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
	//	t = t.Elem()
	//}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.types[t] = d
	return err
}

//ReplaceOrAddByParamType update
func (c *lineImp) ReplaceOrAddByParamType(d dot.Dot, t reflect.Type) error {
	var err error
	//for t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
	//	t = t.Elem()
	//}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.types[t] = d
	return err
}

//ReplaceOrAddByLiveId update
func (c *lineImp) ReplaceOrAddByLiveId(d dot.Dot, id dot.LiveId) error {
	var err error
	c.mutex.Lock()
	defer c.mutex.Unlock()

	l := dot.Live{LiveId: id, TypeId: "", Dot: d, RelyLives: nil}
	//c.lives.Remove(&l)
	err = c.lives.UpdateOrAdd(&l)

	return err
}

//RemoveByType remove
func (c *lineImp) RemoveByType(t reflect.Type) error {
	var err error
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.types, t)
	return err
}

//RemoveByLiveId remove
func (c *lineImp) RemoveByLiveId(id dot.LiveId) error {
	var err error
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.lives.RemoveById(id)
	return err
}

//SetParent set parent injecter
func (c *lineImp) SetParent(p dot.Injecter) {
	c.parent = p
}

//GetParent get parent injecter
func (c *lineImp) GetParent() dot.Injecter {
	return c.parent
}

////injecter end

//Create create
//If liveid is empty， directly assign typeid
//If liveid repeated，directly return dot.SError.ErrExistedLiveId
func (c *lineImp) Create(l dot.Line) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	var err error

FOR_FUN:
	for {
		//first create config
		c.sConfig = sconfig.NewConfiger()
		c.sConfig.RootPath()
		if s, ok := c.sConfig.(dot.Creator); ok {
			if err = s.Create(l); err != nil {
				createLog(c)
				break FOR_FUN
			} else if len(c.sConfig.ConfigFile()) < 1 { //no config file return
				createLog(c)
				break FOR_FUN
			}
		}

		if err = c.sConfig.Unmarshal(&c.config); err != nil {
			createLog(c)
			break FOR_FUN
		}
		if len(c.config.Dots) < 1 { //no config
			createLog(c)
			break FOR_FUN
		}

		//create log
		createLog(c)

		{ //handle config
			for _, it := range c.config.Dots {

				if len(it.MetaData.TypeId.String()) < 1 {
					err = dot.SError.Config.AddNewError("typeid is null")
					break FOR_FUN
				}

				if err = c.metas.Add(&it.MetaData); err != nil {
					break FOR_FUN
				}

				if len(it.Lives) < 1 { //create the single live
					live := dot.Live{TypeId: it.MetaData.TypeId, LiveId: dot.LiveId(it.MetaData.TypeId), Dot: nil}
					if len(it.MetaData.RelyTypeIds) > 0 {
						live.RelyLives = make(map[string]dot.LiveId, len(it.MetaData.RelyTypeIds))
						for i := range it.MetaData.RelyTypeIds {
							li := &it.MetaData.RelyTypeIds[i]
							live.RelyLives[li.String()] = dot.LiveId(*li)
						}
					}
					if err = c.lives.Add(&live); err != nil {
						break FOR_FUN
					}
				} else {
					for _, li := range it.Lives {
						if len(li.LiveId.String()) < 1 {
							li.LiveId = dot.LiveId(it.MetaData.TypeId)
						}
						live := dot.Live{TypeId: it.MetaData.TypeId, LiveId: li.LiveId, RelyLives: li.RelyLives, Dot: nil}
						if err = c.lives.Add(&live); err != nil {
							break FOR_FUN
						}
					}
				}
			}
		}
		break
	}

	return err
}

//todo this method is private and will be realized with component method
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
			tdots, _ := c.RelyOrder() //do not care the circle
			afterStarts := make([]dot.AfterAllStarter, 0, 20)
			for _, it := range tdots {

				if b := c.dotEventer.TypeEvents(it.TypeId); len(b) > 0 {
					for i := range b {
						e := &b[i]
						if e.BeforeStart != nil {
							e.BeforeStart(it, c)
						}
					}

				}

				if b := c.dotEventer.LiveEvents(it.LiveId); len(b) > 0 {
					for i := range b {
						e := &b[i]
						if e.BeforeStart != nil {
							e.BeforeStart(it, c)
						}
					}
				}

				if d, ok := it.Dot.(dot.Starter); ok {
					terr := d.Start(ignore)
					if terr != nil {
						if err != nil {
							logger.Errorln(err.Error())
						}
						err = terr
						if !ignore {
							return err
						}
					}
				}

				if a := c.dotEventer.LiveEvents(it.LiveId); len(a) > 0 {
					for i := range a {
						e := &a[i]
						if e.AfterStart != nil {
							e.AfterStart(it, c)
						}
					}
				}

				if a := c.dotEventer.TypeEvents(it.TypeId); len(a) > 0 {
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
		tdots, _ := c.RelyOrder() //do not care the circle
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

		for idot := len(tdots) - 1; idot >= 0; idot-- {
			it := tdots[idot]
			if b := c.dotEventer.TypeEvents(it.TypeId); len(b) > 0 {
				for i := range b {
					e := &b[i]
					if e.BeforeStop != nil {
						e.BeforeStop(it, c)
					}
				}
			}

			if b := c.dotEventer.LiveEvents(it.LiveId); len(b) > 0 {
				for i := range b {
					e := &b[i]
					if e.BeforeStop != nil {
						e.BeforeStop(it, c)
					}
				}
			}

			if d, ok := it.Dot.(dot.Stopper); ok {
				terr := d.Stop(ignore)
				if terr != nil {
					if err != nil {
						logger.Error(err.Error)
					}
					err = terr
				}

				if !ignore {
					return err
				}
			}

			if a := c.dotEventer.LiveEvents(it.LiveId); len(a) > 0 {
				for i := range a {
					e := &a[i]
					if e.AfterStop != nil {
						e.AfterStop(it, c)
					}
				}
			}

			if a := c.dotEventer.TypeEvents(it.TypeId); len(a) > 0 {
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
		err = d.Stop(ignore)
	}

	//stop config
	if d, ok := c.sConfig.(dot.Stopper); ok {
		err = d.Stop(ignore)
	}

	return err
}

//Destroy Destroy Dot
func (c *lineImp) Destroy(ignore bool) error {
	var err error
	logger := dot.Logger()
	//Destroy others
	{
		//recount the order, maybe the "Ceate" change it
		tdots, _ := c.RelyOrder() //do not care the circle
		for idot := len(tdots) - 1; idot >= 0; idot-- {
			it := tdots[idot]

			if b := c.dotEventer.TypeEvents(it.TypeId); len(b) > 0 {
				for i := range b {
					e := &b[i]
					if e.BeforeDestroy != nil {
						e.BeforeDestroy(it, c)
					}
				}
			}

			if b := c.dotEventer.LiveEvents(it.LiveId); len(b) > 0 {
				for i := range b {
					e := &b[i]
					if e.BeforeDestroy != nil {
						e.BeforeDestroy(it, c)
					}
				}
			}

			if d, ok := it.Dot.(dot.Destroyer); ok {
				terr := d.Destroy(ignore)
				if terr != nil {
					if err != nil {
						logger.Errorln(err.Error())
					}
					err = terr
				}
				if !ignore {
					return err
				}
			}

			if a := c.dotEventer.LiveEvents(it.LiveId); len(a) > 0 {
				for i := range a {
					e := &a[i]
					if e.AfterDestroy != nil {
						e.AfterDestroy(it, c)
					}
				}
			}

			if a := c.dotEventer.TypeEvents(it.TypeId); len(a) > 0 {
				for i := range a {
					e := &a[i]
					if e.AfterDestroy != nil {
						e.AfterDestroy(it, c)
					}
				}
			}

		}
	}

	//Destroy log
	if d, ok := c.logger.(dot.Destroyer); ok {
		d.Destroy(ignore)
	}

	//Destroy config
	if d, ok := c.sConfig.(dot.Destroyer); ok {
		d.Destroy(ignore)
	}

	return err
}

func (c *lineImp) GetLineBuilder() *dot.Builder {
	return c.lineBuilder
}

///////////////
