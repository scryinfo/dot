// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/line"
	"github.com/scryinfo/scryg/sutils/ssignal"
)

func main() {
	l, err := line.BuildAndStart(add) //first step create line and dots
	if err != nil {
		fmt.Println(err)
		return
	}
	defer line.StopAndDestroy(l, true) //fourth step stop and destroy dots

	dot.Logger().Infoln("dot ok")
	t := &SomeUse{}

	l.ToInjecter().Inject(t)                    //second step use the injecter or others
	dot.GetDefaultLine().ToInjecter().Inject(t) //or second step, use the default line(in the sample, the default line  == l)

	ssignal.WatiCtrlC(func(s os.Signal) bool { //third wait for exit
		return false
	})

}

//how to new the dots of config
func add(l dot.Line) error {
	var err error
	{
		t := reflect.TypeOf(((*Dot1)(nil)))
		t = t.Elem()
		fmt.Println("  ", t)
		err = l.PreAdd(&dot.TypeLives{
			Meta: dot.Metadata{TypeId: "1", RefType: t}, Lives: []dot.Live{
				dot.Live{LiveId: "12"},
			},
		})

		// 给typeid指定newer
		err = l.PreAdd(&dot.TypeLives{
			Meta: dot.Metadata{TypeId: "1", NewDoter: func(conf interface{}) (dot dot.Dot, err error) {
				return &Dot1{Name: "Create by type 1"}, nil
			}},
		})
	}

	{
		t := reflect.TypeOf(((*Dot2)(nil)))
		t = t.Elem()
		fmt.Println("  ", t)
		//这里没有指定 newer, 那么会直接使用反射 reflect.NewLine 来创建
		err = l.PreAdd(&dot.TypeLives{
			Meta: dot.Metadata{TypeId: "2", RefType: t}, Lives: []dot.Live{
				dot.Live{LiveId: "21"}, dot.Live{LiveId: "22"},
			},
		})
	}

	{ // 以下为使用 LiveId对就应的Newer，
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

	{ // 以下为使用 typeid 与 LiveId对就应的Newer，如果两个都提供，那么优先使用liveid对应的
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

//直接向容器中加入指定的类型
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

//Create 在这个方法在进行初始，也运行或监听相同内容，最好放在Start方法中实现
func (c *Dot3) Create(conf dot.SConfig) error {

	return nil
}

//Start
//ignore 在调用其它Lifer时，true 出错出后继续，false 出现一个错误直接返回
func (c *Dot3) Start(ignore bool) error {
	return nil
}

//Stop
//ignore 在调用其它Lifer时，true 出错出后继续，false 出现一个错误直接返回
func (c *Dot3) Stop(ignore bool) error {
	return nil
}

//Destroy 销毁 Dot
//ignore 在调用其它Lifer时，true 出错出后继续，false 出现一个错误直接返回
func (c *Dot3) Destroy(ignore bool) error {
	return nil
}
