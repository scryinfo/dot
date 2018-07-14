package main

import (
	"fmt"

	"github.com/scryinfo/dot/line"
	"github.com/scryinfo/dot/line/lineimp"
	"github.com/scryinfo/dot/dot"
	"reflect"
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
	t := reflect.TypeOf(((*Dot2)(nil)))
	t = t.Elem()
	fmt.Println("  ", t)
	l.PreAdd(&line.TypeLives{
		Meta:dot.MetaData{TypeId:"789", RefType: t,},
	})
}

//直接向容器中加入指定的类型
func addDot(l line.Line) {
	l.ToInjecter().ReplaceOrAddByType(&Dot1{Name: "null"})
	l.ToInjecter().ReplaceOrAddByLiveId(&Dot1{Name:"6666"}, dot.LiveId("6666"))
}

type Dot1 struct {
	Name string
}

type SomeUse struct {
	DotLive  *Dot1 `dot:""`
	DotLive2 *Dot1 `dot:"6666"`
	DotLive3 *Dot2 `dot:"789"`
}

type Dot2 struct {
	T string
}


