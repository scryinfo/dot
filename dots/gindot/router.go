// Scry Info.  All rights reserved.
// license that can be found in the license file.

package gindot

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/scryinfo/dot/dot"
)

const (
	//RouterTypeID for gin dot
	RouterTypeID = "6be39d0b-3f5b-47b4-818c-642c049f3166"
)

type configRouter struct {
	RelativePath string `json:"relativePath" yaml:"relativePath"`
}

// Router  gin router
type Router struct {
	Engine_ *Engine `dot:""`
	router  *gin.RouterGroup
	config  configRouter
	liveID  dot.LiveID
}

// construct dot
func newRouter(conf []byte) (*Router, error) {
	dconf := &configRouter{}
	err := dot.UnMarshalConfig(conf, dconf)
	if err != nil {
		return nil, err
	}

	d := &Router{config: *dconf}

	return d, err
}

// RouterTypeLives generate data for structural  dot,  include gindot.Engine
func RouterTypeLives() []*dot.TypeLives {
	lives := []*dot.TypeLives{{
		Meta: dot.Metadata{TypeID: RouterTypeID, NewDoter: func(conf []byte) (dot.Dot, error) {
			return newRouter(conf)
		}},
		Lives: []dot.Live{
			{
				LiveID:    RouterTypeID,
				RelyLives: map[string]dot.LiveID{"Engine_": EngineLiveID},
			},
		},
	}}
	lives = append(lives, GinDotTypeLives()...)
	return lives
}

// jayce edit
// return config of Router
func RouterConfigTypeLive() *dot.ConfigTypeLive {
	return &dot.ConfigTypeLive{
		TypeIDConfig: RouterTypeID,
		ConfigInfo:   &configRouter{},
	}
}

func (c *Router) SetTypeID(_ dot.TypeID, liveID dot.LiveID) {
	c.liveID = liveID
}

func (c *Router) AfterAllInject(l dot.Line) {
	c.router = c.Engine_.GinEngine().Group(c.config.RelativePath)
}

// Start start the gin
func (c *Router) Start(ignore bool) error {

	return nil
}

func (c *Router) Router() *gin.RouterGroup {
	return c.router
}

func (c *Router) RelativePath() string {
	return c.config.RelativePath
}

// all post
func (c *Router) RouterPost(h interface{}, pre string) {
	post := reflect.ValueOf(c.router).MethodByName("POST")
	RouterSelf(h, pre, func(url string, gmethod reflect.Value) {
		vs := []reflect.Value{reflect.ValueOf(url), gmethod}
		post.Call(vs)
	})
}

// all get
func (c *Router) RouterGet(h interface{}, pre string) {
	get := reflect.ValueOf(c.router).MethodByName("GET")
	RouterSelf(h, pre, func(url string, gmethod reflect.Value) {
		vs := []reflect.Value{reflect.ValueOf(url), gmethod}
		get.Call(vs)
	})
}
