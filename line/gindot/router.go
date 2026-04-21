// Scry Info.  All rights reserved.
// license that can be found in the license file.

package gindot

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	RelativePath string `json:"relativePath" yaml:"relativePath"`
}

// Router  gin router
type Router struct {
	Engine_ *Engine
	router  *gin.RouterGroup
	config  RouterConfig
}

// construct dot
func NewRouter(conf *RouterConfig, engine *Engine) (*Router, error) {

	d := &Router{config: *conf, router: engine.GinEngine().Group(conf.RelativePath), Engine_: engine}

	return d, nil
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
