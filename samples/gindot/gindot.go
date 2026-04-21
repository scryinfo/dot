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

	"github.com/scryinfo/dot/line/gindot"
	"github.com/scryinfo/dot/line/sconfig"
)

type Line struct {
	SConfig        *sconfig.SConfig
	Logger         *dot.LoggerType
	SampleCtroller *SampleCtroller
}

type LineConfig struct {
	Log    dot.LogConfig
	Router gindot.RouterConfig
	Engine gindot.EngineConfig
}

func NewLineConfig(config *sconfig.SConfig) (*LineConfig, error) {
	return sconfig.NewLiceConfig[LineConfig](config)
}

var LineSet = wire.NewSet(
	NewLineConfig,
	wire.Struct(new(Line), "*"),
	sconfig.NewConfig,
	dot.NewLogger,
	wire.FieldsOf(new(*LineConfig), "Log", "Router", "Engine"),
	gindot.NewRouter,
	gindot.NewGinDot,
	NewSampleCtroller,
)

func main() {
	line, clear, err := InitializeService()
	if err != nil {
		dot.Logger.Error().Err(err).Msg("initialize service failed")
		return
	}
	if clear != nil {
		defer clear()
	}
	_ = line
	ssignal.WaitCtrlC(func(s os.Signal) bool { //third wait for exit
		return false
	})

	dot.Logger.Info().Msg("dot will stop")
}

func NewSampleCtroller(router *gindot.Router) (*SampleCtroller, error) {
	d := &SampleCtroller{
		Router: router,
	}
	err := d.Start()
	return d, err
}

type SampleCtroller struct {
	Router *gindot.Router
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
