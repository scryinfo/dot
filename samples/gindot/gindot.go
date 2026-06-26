// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line"
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
	Log    dot.LogConfig       `json:"log" toml:"log" yaml:"log" mapstructure:"log"`
	Router gindot.RouterConfig `json:"router" toml:"router" yaml:"router" mapstructure:"router"`
	Engine gindot.EngineConfig `json:"engine" toml:"engine" yaml:"engine" mapstructure:"engine"`
}

func NewLineConfig(config *sconfig.SConfig) (*LineConfig, error) {
	lineConfig, err := sconfig.NewLineConfig[LineConfig](config)
	if err != nil {
		return nil, err
	}
	return sconfig.GenerateConfigWithArgs(config, lineConfig)
}

var LineSet = wire.NewSet(
	wire.Struct(new(Line), "*"),
	wire.FieldsOf(new(*LineConfig), "Log", "Router", "Engine"),
	NewLineConfig,
	line.SconfigNewConfig,
	dot.NewLogger,
	line.GindotNewRouter,
	line.GindotNewGinDot,
	NewSampleCtroller,
)

func main() {
	line, clean, err := InitializeService()
	if err != nil {
		if line != nil {
			dot.Logger.Error().Err(err).Msg("initialize service failed")
		} else {
			dot.Logger.Info().Msg(err.Error())
			fmt.Printf("\n\n")
		}
		return
	}
	if clean != nil {
		defer clean()
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
