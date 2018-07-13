package main

import (
	"fmt"

	"github.com/scryinfo/dot/line"
	"github.com/scryinfo/dot/line/lineimp"
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

	defer func() {
		l.ToLifer().Stop(true)
		l.ToLifer().Destroy(true)
	}()

}

func add(l line.Line) {
	//l.PreAdd()
}

type Dot1 struct {
}

type SomeUse struct {
	DotLive  Dot1 `dot:""`
	DotLive2 Dot1 `dot:""`
}
