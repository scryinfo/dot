package lineimp

import (
	"reflect"
	"sync"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/slog"
	"github.com/scryinfo/dot/line"
	"github.com/scryinfo/scryg/sutils/skit"
)

var (
	_ dot.Lifer     = (*lineimp)(nil)
	_ line.Line     = (*lineimp)(nil)
	_ line.Injecter = (*lineimp)(nil)
)

type lineimp struct {
	dot.Lifer
	line.Line
	line.Injecter
	logger  slog.SLogger
	sConfig dot.SConfig
	config  line.Config
	metas   *line.Metas
	lives   *line.Lives

	mutex sync.Mutex
}

//New new
func New() line.Line {
	a := &lineimp{}
	if line.GetDefaultLine() == nil {
		line.SetDefaultLine(a)
	}
	a.metas = line.NewMetas()
	a.lives = line.NewLives()
	return a
}

//PreAdd the dot is nil, do not create it
func (c *lineimp) PreAdd(livings *line.TypeLives) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	clone := livings
	err := c.metas.UpdateOrAdd(&clone.Meta)
	if err == nil {
		for _, it := range clone.Lives {

			live := dot.Live{TypeId: it.TypeId, LiveId: it.LiveId}
			live.RelyLives = make([]dot.LiveId, len(it.RelyLives))
			copy(live.RelyLives, it.RelyLives)
			c.lives.UpdateOrAdd(&it)
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
			if err != nil {
				break LIVES
			}

			if live.TypeId != it.TypeId {
				err = dot.SError.RelyTypeNotMatch.AddNewError(live.LiveId.String())
				break LIVES
			}
		}
	}

	return err
}

func (c *lineimp) CreateDots() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	var err error
LIVES:
	for _, it := range c.lives.LiveIdMap {

		if skit.IsNil(it.Dot) {
			var m *dot.MetaData
			m, err = c.metas.Get(it.TypeId)
			if m != nil {
				break LIVES
			}

			if m.NewDoter == nil {
				err = dot.SError.NoDotNewer.AddNewError(m.TypeId.String())
			}

			var bconfig []byte
			{
				config := c.config.FindConfig(it.TypeId, it.LiveId)
				if config != nil {
					bconfig, err = config.Json.MarshalJSON()
					if err != nil {
						break LIVES
					}
				}
			}

			it.Dot, err = m.NewDot(bconfig)
		}
	}

	return err
}

func (c *lineimp) ToLifer() dot.Lifer {
	return c
}

func (c *lineimp) GetDotConfig(liveid dot.LiveId) *line.DotConfig {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var co *line.DotConfig
	c.config.FindConfig("", liveid)
	return co
}

/////injecter

//Inject see https://github.com/facebookgo/inject
func (c *lineimp) Inject(obj interface{}) error {
	var err error
	if skit.IsNil(obj) {
		return dot.SError.NilParameter
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()

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
		f := v.Field(i)
		structField := t.Field(i)
		if f.CanSet() && (structField.Tag == dot.TagDot || structField.Tag.Get(dot.TagDot) != "") {
			ft := f.Type()
			var d dot.Dot
			d, errt = c.GetByType(ft)
			if errt != nil && err == nil {
				err = errt
			}
			if err == nil {
				vv := reflect.ValueOf(d)
				if vv.IsValid() {
					f.Set(vv)
				} else if err == nil {
					err = dot.SError.DotInvalid.AddNewError(ft.Name())
				}
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
	defer c.mutex.Unlock()

	return
}

//GetByLiveId get by liveid
func (c *lineimp) GetByLiveId(liveId dot.LiveId) (d dot.Dot, err error) {
	d = nil
	err = nil
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var l *dot.Live
	l, err = c.lives.Get(liveId)
	if err == nil {
		d = l.Dot
	}

	return
}

//GetByTypeId get by typeid
func (c *lineimp) GetByTypeId(typeId dot.TypeId) (d dot.Dot, err error) {
	d = nil
	err = nil
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if m, err2 := c.metas.Get(typeId); err2 == nil {
		if l, err2 := c.lives.Get(dot.LiveId(m.TypeId)); err2 == nil {
			d = l.Dot
		} else {
			err = err2
		}
	} else {
		err = err2
	}
	return
}

////injecter end

//Create create
//如果 liveid为空， 直接赋值为 typeid
//如果 liveid重复，直接返回 dot.SError.ErrExistedLiveId
func (c *lineimp) Create(conf dot.SConfig) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	var err error
FOR_FUN:
	for {
		//first create config
		c.sConfig = dot.NewConfiger()
		if err = c.sConfig.Create(nil); err != nil {
			break FOR_FUN
		}

		if err = c.sConfig.Unmarshal(c.config); err != nil {
			break FOR_FUN
		}

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
					live := dot.Live{TypeId: it.MetaData.TypeId, LiveId: dot.LiveId(it.MetaData.TypeId)}
					if len(it.MetaData.RelyTypeIds) > 0 {
						live.RelyLives = make([]dot.LiveId, len(it.MetaData.RelyTypeIds))
						for i, li := range it.MetaData.RelyTypeIds {
							live.RelyLives[i] = dot.LiveId(li)
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
						live := dot.Live{TypeId: it.MetaData.TypeId, LiveId: li.LiveId, RelyLives: li.RelyLives}
						if err = c.lives.Add(&live); err != nil {
							break FOR_FUN
						}
					}
				}
			}
		}

		//create log

		//create others
		break
	}

	return err
}

//Start
func (c *lineimp) Start(ignore bool) error {
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
func (c *lineimp) Stop(ignore bool) error {
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
func (c *lineimp) Destroy(ignore bool) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	//Destroy others

	//Destroy log

	//Destroy config
	c.sConfig.Destroy(ignore)
	return nil
}

///////////////
