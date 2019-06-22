// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"encoding/json"
	"os"
	"reflect"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/line"
	"github.com/scryinfo/scryg/sutils/ssignal"
)

func main() {
	l, err := line.BuildAndStart(add) //first step create line and dots
	if err != nil {
		dot.Logger().Errorln(err.Error())
		return
	}
	defer line.StopAndDestroy(l, true) //fourth step stop and destroy dots

	dot.Logger().Infoln("dot ok")
	t := &SomeUse{}

	l.ToInjecter().Inject(t)                    //second step use the injecter or others
	dot.GetDefaultLine().ToInjecter().Inject(t) //or second step, use the default line(in the sample, the default line  == l)

	ssignal.WaitCtrlC(func(s os.Signal) bool { //third wait for exit
		return false
	})

}

//how to new the dots of config
func add(l dot.Line) error {
	logger := dot.Logger()
	var err error
	{
		t := reflect.TypeOf(((*Dot1)(nil)))
		t = t.Elem()
		err = l.PreAdd(&dot.TypeLives{
			Meta: dot.Metadata{TypeId: "1", RefType: t}, Lives: []dot.Live{
				dot.Live{LiveId: "12"},
			},
		})

		// Point newer for typeid
		err = l.PreAdd(&dot.TypeLives{
			Meta: dot.Metadata{TypeId: "1", NewDoter: func(conf interface{}) (dot dot.Dot, err error) {
				return &Dot1{Name: "Create by type 1"}, nil
			}},
		})
	}

	{
		t := reflect.TypeOf(((*Dot2)(nil)))
		t = t.Elem()
		//If no newer assignment, then use reflect.newLine to create it
		err = l.PreAdd(&dot.TypeLives{
			Meta: dot.Metadata{TypeId: "2", RefType: t}, Lives: []dot.Live{
				dot.Live{LiveId: "21"}, dot.Live{LiveId: "22"},
			},
		})
	}

	{ // The following is Newer using LiveId，
		err = l.AddNewerByLiveId(dot.LiveId("31"), func(conf interface{}) (d dot.Dot, err error) {
			d = &Dot3{}
			err = nil
			t := reflect.ValueOf(conf)
			if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
				if t.Len() > 0 && t.Index(0).Kind() == reflect.Uint8 {
					v := t.Slice(0, t.Len())
					json.Unmarshal(v.Bytes(), d)
				}
			} else {
				err = dot.SError.Parameter
			}

			return
		})

		err = l.AddNewerByLiveId(dot.LiveId("32"), func(conf interface{}) (d dot.Dot, err error) {
			d = &Dot3{}
			err = nil
			t := reflect.ValueOf(conf)
			if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
				if t.Len() > 0 && t.Index(0).Kind() == reflect.Uint8 {
					v := t.Slice(0, t.Len())
					json.Unmarshal(v.Bytes(), d)
				}
			} else {
				err = dot.SError.Parameter
			}

			return
		})
	}

	{ // The following is Newer using typeid and LiveId，if both are provided, then use Liveid priorly
		err = l.AddNewerByLiveId(dot.LiveId("41"), func(conf interface{}) (d dot.Dot, err error) {
			d = &Dot4{}
			err = nil
			t := reflect.ValueOf(conf)
			if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
				if t.Len() > 0 && t.Index(0).Kind() == reflect.Uint8 {
					v := t.Slice(0, t.Len())
					json.Unmarshal(v.Bytes(), d)
				}
			} else {
				err = dot.SError.Parameter
			}

			return
		})

		err = l.AddNewerByTypeId(dot.TypeId("type_live"), func(conf interface{}) (d dot.Dot, err error) {
			d = &Dot4{}
			err = nil
			t := reflect.ValueOf(conf)
			if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
				if t.Len() > 0 && t.Index(0).Kind() == reflect.Uint8 {
					v := t.Slice(0, t.Len())
					json.Unmarshal(v.Bytes(), d)
				}
			} else {
				err = dot.SError.Parameter
			}
			return
		})
	}
	return err
}

//Add assigned type to container directly
func addDot(l dot.Line) {
	l.ToInjecter().ReplaceOrAddByType(&Dot1{Name: "null"})
	l.ToInjecter().ReplaceOrAddByLiveId(&Dot1{Name: "6666"}, dot.LiveId("6666"))
}

type Dot1 struct {
	Name string
}

type SomeUse struct {
	DotLive  *Dot1 `dot:""`
	DotLive2 *Dot1 `dot:"12"`

	DotLive3 *Dot2 `dot:"21"`
	DotLive4 *Dot2 `dot:"22"`

	DotLive5 *Dot3 `dot:"31"`
	DotLive6 *Dot3 `dot:"32"`

	DotLive10 *Dot4 `dot:"41"`
	DotLive11 *Dot4 `dot:"42"`

	Logger dot.SLogger `dot:""`
}

type Dot2 struct {
	T string
}

type Dot3 struct {
	T string
}

type Dot4 struct {
	T string
}

//Create use this method to initialize, run or monitor same content, better realize it in start method
func (c *Dot3) Create(conf dot.SConfig) error {

	return nil
}

//Start
//ignore When calliing other Lifer, if true erred then continue, if false erred return directly
func (c *Dot3) Start(ignore bool) error {
	return nil
}

//Stop
//ignore When calliing other Lifer, if true erred then continue, if false erred return directly
func (c *Dot3) Stop(ignore bool) error {
	return nil
}

//Destroy Destroy Dot
//ignore When calliing other Lifer, if true erred then continue, if false erred return directly
func (c *Dot3) Destroy(ignore bool) error {
	return nil
}
