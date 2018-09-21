package main

import (
	"github.com/scryinfo/dot-0/line/lineimp"
	"github.com/scryinfo/dot-0/line"
	"fmt"
	"github.com/scryinfo/dot-0/dot"
)

// 创建容器
//添加fmt组件

func main()  {
	l := lineimp.New()
	l.ToLifer().Create(nil)

	addFmt(l)

	f := &Fmt2{}

	l.ToInjecter().Inject(f)

	f.F.Println("12312312313")

}

//注册fmt组件
func addFmt(l line.Line)  {
	l.ToInjecter().ReplaceOrAddByLiveId(&Fmt1{},dot.LiveId("1"))
}

//实际对象
type Fmt1 struct {}

//创建使用对象
type Fmt2 struct {
	F *Fmt1 `dot:"1"`
}

//组件方法
func (f Fmt1) Println(s string) (n int, err error) {
	return  fmt.Println(s)
}