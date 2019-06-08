// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"fmt"
	"net/http"
	"os"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/gindot"
	"github.com/scryinfo/dot/dots/line"
	"github.com/scryinfo/scryg/sutils/ssignal"
)

func main() {
	//l, err := line.BuildAndStart(add) //first step create line and dots
	l, err := line.BuildAndStartBy(&dot.Builder{
		AfterCreate: func(l dot.Line) {
			fmt.Println("AfterCreate")
		},
		Add: add,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer line.StopAndDestroy(l, true) //fourth step stop and destroy dots
	dot.Logger().Infoln("dot ok")
	//third step
	//....

	ssignal.WaitCtrlC(func(s os.Signal) bool { //third wait for exit
		return false
	})
	dot.Logger().Infoln("dot will stop")
}

func add(l dot.Line) error {
	err := l.PreAdd(gindot.TypeLiveRouter()...)
	//ReSetLiveEvents AddLiveEvents , they are different
	l.ToDotEventer().ReSetLiveEvents(dot.LiveId("6be39d0b-3f5b-47b4-818c-642c049f3166"), &dot.LiveEvents{AfterStart: func(live *dot.Live, l dot.Line) {
		//do any init for the router
		// router.Router().Use()
		// ....
	}})

	//add the SampleCtroller
	err = l.PreAdd(gindot.PreAddControlDot(reflect.TypeOf((*SampleCtroller)(nil)).Elem(), dot.LiveId("6be39d0b-3f5b-47b4-818c-642c049f3166")))

	l.ToDotEventer().AddLiveEvents(dot.LiveId(gindot.EngineLiveId), &dot.LiveEvents{
		AfterCreate: func(live *dot.Live, l dot.Line) {
			if _, ok := live.Dot.(*gindot.Engine); ok {
				//d.GinEngine().St
				dot.Logger().Infoln("sdf")
			}
			fmt.Println("BeforeStop")
		},

		BeforeStop: func(live *dot.Live, l dot.Line) {
			fmt.Println("BeforeStop")
		},

		BeforeStart: func(live *dot.Live, l dot.Line) {
			fmt.Println("BeforeStart")
		},
	})

	return err
}

type SampleCtroller struct {
	GinRouter_ *gindot.Router `dot:"6be39d0b-3f5b-47b4-818c-642c049f3166"`
}

func (c *SampleCtroller) Start(ignore bool) error {

	//post := reflect.ValueOf(c.GinRouter_.Router()).MethodByName("POST")
	//gindot.RouterSelf(c,"sample", func(url string, gmethod reflect.Value) {
	//	vs := []reflect.Value{reflect.ValueOf(url), gmethod}
	//	post.Call(vs)
	//})
	c.GinRouter_.RouterGet(c, "sample")
	return nil
}

func (c *SampleCtroller) Hello(cxt *gin.Context) {
	cxt.JSON(http.StatusOK, "ok")
}
