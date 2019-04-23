package main

import (
	"fmt"
	"github.com/scryInfo/dot/dot"
	"github.com/scryInfo/dot/dots/line"
)

var (
	_ dot.SLogger = (*dot.ULog)(nil)
)

func main() {

	//if sv,ok := slog.sLogger.(slog.SLogger); ok {
	//	fmt.Printf("v implements String(): %s\n", sv)
	//}

	l := line.New()
	l.ToLifer().Create(nil)

	//fmt.Println(ll)

	//dot.Add(l)

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

	f := &dot.DotLog{}
	l.ToInjecter().Inject(f)

	//lineimp.CreateLog()

	f.Log.Debug(func() string {
		return "6666666666"
	})

	f.Log.Debugln("ssssssssssssssssssssss")

	fmt.Println(f.Log.GetLevel())

	//var m runtime.MemStats
	//
	//runtime.ReadMemStats(&m)
	//
	//fmt.Printf("%d KB\n",m.Alloc/1024)

}

//func addUlog(l line.Line)  {
//	l.AddNewerByLiveId(dot.LiveId("3"), func(conf interface{}) (d dot.Dot, err error) {
//		d = &slog.sLogger{
//			Level:slog.DebugLevel,
//			OutputPath:"out1.log",
//		}
//		err = nil
//		//t := reflect.ValueOf(conf)
//		return
//	})
//	l.PreAdd(&line.TypeLives{
//		Meta: dot.Metadata{TypeId: "3"},Lives: []dot.Live{
//			dot.Live{LiveId: "3"}},
//	})
//}
