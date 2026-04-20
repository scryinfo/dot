// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/scryg/sutils/ssignal"

	"github.com/scryinfo/dot/dots/gindot"
	"github.com/scryinfo/dot/dots/sconfig"
)

type App struct {
	SConfig        *sconfig.SConfig
	Logger         *dot.LoggerType
	SampleCtroller *SampleCtroller
}

type AppConfig struct {
	Log    dot.LogConfig
	Router gindot.ConfigRouter
	Engine gindot.ConfigEngine
}

func NewAppConfig(config *sconfig.SConfig) (*AppConfig, error) {
	return sconfig.NewAppConfig[AppConfig](config)
}

var AppSet = wire.NewSet(
	NewAppConfig,
	wire.Struct(new(App), "*"),
	sconfig.NewConfig,
	dot.InitLogger,
	wire.FieldsOf(new(*AppConfig), "Log", "Router", "Engine"),
	gindot.NewRouter,
	gindot.NewGinDot,
	NewSampleCtroller,
)

func main() {
	app, close, err := InitializeService()
	if err != nil {
		dot.Logger.Error().Err(err).Msg("initialize service failed")
		return
	}
	if close != nil {
		defer close()
	}
	_ = app
	ssignal.WaitCtrlC(func(s os.Signal) bool { //third wait for exit
		return false
	})

	dot.Logger.Info().Msg("dot will stop")
}

// func add(l dot.Line) error {
// 	err := l.PreAdd(gindot.RouterTypeLives()...)
// 	//ReSetLiveEvents AddLiveEvents , they are different
// 	l.ToDotEventer().ReSetLiveEvents(dot.LiveID("6be39d0b-3f5b-47b4-818c-642c049f3166"), &dot.LiveEvents{AfterStart: func(live *dot.Live, l dot.Line) {
// 		//do any init for the router
// 		// router.Router().Use()
// 		// ....
// 	}})

// 	//add the SampleCtroller
// 	err = l.PreAdd(gindot.PreAddControlDot(reflect.TypeOf((*SampleCtroller)(nil)).Elem(), dot.LiveID("6be39d0b-3f5b-47b4-818c-642c049f3166")))

// 	l.ToDotEventer().AddLiveEvents(dot.LiveID(gindot.EngineLiveID), &dot.LiveEvents{
// 		AfterCreate: func(live *dot.Live, l dot.Line) {
// 			if g, ok := live.Dot.(*gindot.Engine); ok {
// 				//d.GinEngine().St
// 				_ = g
// 				dot.Logger().Infoln("sdf")
// 			}
// 			dot.Logger().Infoln("BeforeStop")

// 		},

// 		BeforeStop: func(live *dot.Live, l dot.Line) {
// 			dot.Logger().Infoln("BeforeStop")
// 		},

// 		BeforeStart: func(live *dot.Live, l dot.Line) {
// 			dot.Logger().Infoln("BeforeStart")
// 		},
// 	})

// 	return err
// }

func NewSampleCtroller(router *gindot.Router) (*SampleCtroller, error) {
	d := &SampleCtroller{
		Router: router,
	}
	err := d.Start()
	return d, err
}

type SampleCtroller struct {
	Router *gindot.Router `dot:"6be39d0b-3f5b-47b4-818c-642c049f3166"`
}

func (c *SampleCtroller) Start() error {

	//post := reflect.ValueOf(c.GinRouter_.Router()).MethodByName("POST")
	//gindot.RouterSelf(c,"sample", func(url string, gmethod reflect.Value) {
	//	vs := []reflect.Value{reflect.ValueOf(url), gmethod}
	//	post.Call(vs)
	//})
	c.Router.RouterGet(c, "sample")
	c.Router.Router().GET("/rpctest/*rpc", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "ok")
	})
	return nil
}

func (c *SampleCtroller) Hello(cxt *gin.Context) {
	cxt.JSON(http.StatusOK, "ok")
}
