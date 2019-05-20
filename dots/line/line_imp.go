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
	_ dot.Lifer    = (*lineimp)(nil)
	_ dot.Line     = (*lineimp)(nil)
	_ dot.Injecter = (*lineimp)(nil)
)

type lineimp struct {
	dot.Lifer
	dot.Line
	dot.Injecter
	logger      dot.SLogger
	sConfig     dot.SConfig //对外的通用配置
	config      dot.Config  //组件的配置对象
	metas       *Metas
	lives       *Lives
	types       map[reflect.Type]dot.Dot
	newerLiveid map[dot.LiveId]dot.Newer
	newerTypeid map[dot.TypeId]dot.Newer

	parent dot.Injecter
	mutex  sync.Mutex

	lineBuilder *dot.Builder
}

//New new
func New(builer *dot.Builder) dot.Line {
	a := &lineimp{metas: NewMetas(),
		lives: NewLives(), types: make(map[reflect.Type]dot.Dot),
		newerLiveid: make(map[dot.LiveId]dot.Newer),
		newerTypeid: make(map[dot.TypeId]dot.Newer),
		lineBuilder: builer,
	}

	if dot.GetDefaultLine() == nil {
		dot.SetDefaultLine(a)
	}

	return a
}

//AddNewerByLiveId add new for liveid
func (c *lineimp) AddNewerByLiveId(liveid dot.LiveId, newDot dot.Newer) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.newerLiveid[liveid]; ok {
		return dot.SError.Existed.AddNewError(liveid.String())
	}

	c.newerLiveid[liveid] = newDot

	return nil
}

//AddNewerByTypeId add new for type
func (c *lineimp) AddNewerByTypeId(typeid dot.TypeId, newDot dot.Newer) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, ok := c.newerTypeid[typeid]; ok {
		return dot.SError.Existed.AddNewError(typeid.String())
	}

	c.newerTypeid[typeid] = newDot

	return nil
}

//RemoveNewerByLiveId remove
func (c *lineimp) RemoveNewerByLiveId(liveid dot.LiveId) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.newerLiveid, liveid)
}

//RemoveNewerByTypeId remove
func (c *lineimp) RemoveNewerByTypeId(typeid dot.TypeId) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.newerTypeid, typeid)
}

//PreAdd the dot is nil, do not create it
func (c *lineimp) PreAdd(livings *dot.TypeLives) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	clone := livings
	err := c.metas.UpdateOrAdd(&clone.Meta)
	if err == nil {
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
	}

	return err
}

func (c *lineimp) Rely() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var err error
	var live *dot.Live
LIVES:
	for _, it := range c.lives.LiveIdMap {
		for _, dit := range it.RelyLives {
			live, err = c.lives.Get(dit)
			if err != nil || live == nil {
				break LIVES
			}
			//todo check more about for live
		}
	}

	return err
}

//CreateDots create dots
func (c *lineimp) CreateDots() error {
	var tdots []*dot.Live
	{
		c.mutex.Lock()
		tdots = make([]*dot.Live, 0, len(c.lives.LiveIdMap))
		for _, it := range c.lives.LiveIdMap {
			if it != nil {
				tdots = append(tdots, it)
			}
		}
		c.mutex.Unlock()
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
						{ // 在Create 之前检测是否需要特定的信息
							if nl, ok := it.Dot.(dot.SetterLine); ok {
								nl.SetLine(c)
							}

							if nl, ok := it.Dot.(dot.SetterTypeAndLiveId); ok {
								nl.SetTypeId(it.TypeId, it.LiveId)
							}
						}
						if l, ok := it.Dot.(dot.Creator); ok {
							//c.mutex.Unlock()
							l.Create(c)
							//c.mutex.Lock()
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
						{ // 在Create 之前检测是否需要特定的信息
							if nl, ok := it.Dot.(dot.SetterLine); ok {
								nl.SetLine(c)
							}

							if nl, ok := it.Dot.(dot.SetterTypeAndLiveId); ok {
								nl.SetTypeId(it.TypeId, it.LiveId)
							}
						}
						if l, ok := it.Dot.(dot.Creator); ok {
							//c.mutex.Unlock()
							l.Create(c)
							//c.mutex.Lock()
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
					{ // 在Create 之前检测是否需要特定的信息
						if nl, ok := it.Dot.(dot.SetterLine); ok {
							nl.SetLine(c)
						}

						if nl, ok := it.Dot.(dot.SetterTypeAndLiveId); ok {
							nl.SetTypeId(it.TypeId, it.LiveId)
						}
					}
					if l, ok := it.Dot.(dot.Creator); ok {
						//c.mutex.Unlock()
						l.Create(c)
						//c.mutex.Lock()
					}
				}
			}
		}
	}

	//增加logger 与 config
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

	//增加类型与 dot的对应关系, 只记录typeid == liveId的
	for _, it := range tdots {
		if !skit.IsNil(&it.Dot) && ((string)(it.TypeId) == (string)(it.LiveId)) {
			t := reflect.TypeOf(it.Dot)
			c.mutex.Lock()
			c.types[t] = it.Dot
			c.mutex.Unlock()
		}
	}

	for _, it := range tdots {
		if it.Dot != nil {
			c.injectInLine(it.Dot, it)
		}
	}

	return err
}

func (c *lineimp) SLogger() dot.SLogger {
	return c.logger
}

func (c *lineimp) SConfig() dot.SConfig {
	return c.sConfig
}

func (c *lineimp) ToLifer() dot.Lifer {
	return c
}

//ToInjecter to injecter
func (c *lineimp) ToInjecter() dot.Injecter {
	return c
}

func (c *lineimp) GetDotConfig(liveid dot.LiveId) *dot.LiveConfig {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var co *dot.LiveConfig
	co = c.config.FindConfig("", liveid)
	return co
}

/////injecter

//Inject see https://github.com/facebookgo/inject
func (c *lineimp) Inject(obj interface{}) error {
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

func (c *lineimp) injectInLine(obj interface{}, live *dot.Live) error {
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
			if len(live.RelyLives) > 0 { //配置优先
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
func (c *lineimp) GetByType(t reflect.Type) (d dot.Dot, err error) {
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
func (c *lineimp) GetByLiveId(liveId dot.LiveId) (d dot.Dot, err error) {
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
func (c *lineimp) ReplaceOrAddByType(d dot.Dot) error {
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
func (c *lineimp) ReplaceOrAddByParamType(d dot.Dot, t reflect.Type) error {
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
func (c *lineimp) ReplaceOrAddByLiveId(d dot.Dot, id dot.LiveId) error {
	var err error
	c.mutex.Lock()
	defer c.mutex.Unlock()

	l := dot.Live{LiveId: id, TypeId: "", Dot: d, RelyLives: nil}
	//c.lives.Remove(&l)
	err = c.lives.UpdateOrAdd(&l)

	return err
}

//RemoveByType remove
func (c *lineimp) RemoveByType(t reflect.Type) error {
	var err error
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.types, t)
	return err
}

//RemoveByLiveId remove
func (c *lineimp) RemoveByLiveId(id dot.LiveId) error {
	var err error
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.lives.RemoveById(id)
	return err
}

//SetParent set parent injecter
func (c *lineimp) SetParent(p dot.Injecter) {
	c.parent = p
}

//GetParent get parent injecter
func (c *lineimp) GetParent() dot.Injecter {
	return c.parent
}

////injecter end

//Create create
//如果 liveid为空， 直接赋值为 typeid
//如果 liveid重复，直接返回 dot.SError.ErrExistedLiveId
func (c *lineimp) Create(l dot.Line) error {
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

//todo 这个方法就为私有，且按照组件的方式来实现
func createLog(c *lineimp) {
	c.logger = slog.NewSLogger(&(c.config.Log), c)
	dot.SetLogger(c.logger)
}

//Start
func (c *lineimp) Start(ignore bool) error {
	var err error
	for {
		//start config
		if s, ok := c.sConfig.(dot.Srater); ok {
			if err = s.Start(ignore); err != nil {
				break
			}
		}
		//start log
		if s, ok := c.logger.(dot.Srater); ok {
			if err = s.Start(ignore); err != nil {
				break
			}
		}

		//start other
		{

			var tdots []*dot.Live
			c.mutex.Lock()
			tdots = make([]*dot.Live, 0, len(c.lives.LiveIdMap))
			for _, it := range c.lives.LiveIdMap {
				if it != nil && it.Dot != nil {
					tdots = append(tdots, it)
				}
			}
			c.mutex.Unlock()

			afterStarts := make([]dot.AfterStarter, 0, 20)
			for _, it := range tdots {
				if d, ok := it.Dot.(dot.Srater); ok {
					d.Start(ignore)
				}

				if s, ok := it.Dot.(dot.AfterStarter); ok {
					afterStarts = append(afterStarts, s)
				}
			}

			for _, s := range afterStarts {
				s.AfterStart(c)
			}
		}

		break
	}

	return nil
}

//Stop
func (c *lineimp) Stop(ignore bool) error {
	var err error
	//stop others
	{
		var tdots []*dot.Live
		beforeStops := make([]dot.BeforeStopper, 0, 20)

		c.mutex.Lock()
		tdots = make([]*dot.Live, 0, len(c.lives.LiveIdMap))
		for _, it := range c.lives.LiveIdMap {
			if it != nil && it.Dot != nil {
				tdots = append(tdots, it)

				if s, ok := it.Dot.(dot.BeforeStopper); ok {
					beforeStops = append(beforeStops, s)
				}
			}
		}
		c.mutex.Unlock()

		for _, it := range beforeStops {
			it.BeforeStop(c)
		}

		for _, it := range tdots {
			if d, ok := it.Dot.(dot.Stopper); ok {
				d.Stop(ignore)
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

//Destroy 销毁 Dot
func (c *lineimp) Destroy(ignore bool) error {
	//Destroy others
	{
		var tdots []*dot.Live
		c.mutex.Lock()
		tdots = make([]*dot.Live, 0, len(c.lives.LiveIdMap))
		for _, it := range c.lives.LiveIdMap {
			if it != nil && it.Dot != nil {
				tdots = append(tdots, it)
			}
		}
		c.mutex.Unlock()
		for _, it := range tdots {
			if d, ok := it.Dot.(dot.Destroyer); ok {
				d.Destroy(ignore)
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

	return nil
}

func (c *lineimp) GetLineBuilder() *dot.Builder {
	return c.lineBuilder
}

///////////////
