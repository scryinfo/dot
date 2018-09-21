package main

import (
	"encoding/json"
	"fmt"

	"reflect"

	"github.com/scryinfo/dot-0/dot"
	"github.com/scryinfo/dot-0/line"
	"github.com/scryinfo/dot-0/line/lineimp"
)

func main() {
	l := lineimp.New()
	l.ToLifer().Create(nil)

	add(l)

	err := l.Rely()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = l.CreateDots()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = l.ToLifer().Start(false)
	if err != nil {
		fmt.Println(err)
		return
	}

	addDot(l)

	t := &SomeUse{}

	l.ToInjecter().Inject(t)

	t = nil

	defer func() {
		l.ToLifer().Stop(true)
		l.ToLifer().Destroy(true)
	}()

}

func add(l line.Line) {
	//t := reflect.TypeOf(((*Dot2)(nil)))
	//t = t.Elem()
	//fmt.Println("  ", t)

	l.AddNewerByLiveId(dot.LiveId("668"), func(conf interface{}) (d dot.Dot, err error) {
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

	//l.PreAdd(&line.TypeLives{
	//	Meta: dot.Metadata{TypeId: "789", RefType: t}, Lives: []dot.Live{
	//		dot.Live{LiveId: "1234"},
	//	},
	//})
	//
	//l.PreAdd(&line.TypeLives{
	//	Meta: dot.Metadata{TypeId: "668"},
	//})
}

//直接向容器中加入指定的类型
func addDot(l line.Line) {
	l.ToInjecter().ReplaceOrAddByType(&Dot1{Name: "null"})
	l.ToInjecter().ReplaceOrAddByLiveId(&Dot1{Name: "6666"}, dot.LiveId("6666"))
}

type Dot1 struct {
	Name string
}

type SomeUse struct {
	DotLive  *Dot1 `dot:""`
	DotLive2 *Dot1 `dot:"6666"`
	DotLive3 *Dot2 `dot:"789"`
	DotLive4 *Dot2 `dot:"1234"`

	DotLive5 *Dot3 `dot:"668"`
}

type Dot2 struct {
	T string
}

type Dot3 struct {
	T string
}
