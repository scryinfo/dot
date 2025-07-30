// Scry Info.  All rights reserved.
// license that can be found in the license file.

package gindot

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/scryinfo/dot/dot"
)

// for each funcs that like “gin.HandlerFunc”
// sampe pre = "scry", and  "func (c * SampleCtroller) Hello(cxt *gin.Context) {}", the url is "/scry/hello"
func RouterSelf(h interface{}, pre string, call func(url string, gmethod reflect.Value)) {
	hf := reflect.TypeOf(gin.HandlerFunc(nil))
	vr := reflect.ValueOf(h)
	tr := reflect.TypeOf(h)
	pre = strings.TrimSpace(pre)
	if pre == "" || pre == "/" {
		pre = "/"
	} else {
		if pre[0] != '/' {
			pre = "/" + pre
		}
		if pre[len(pre)-1] != '/' {
			pre = pre + "/"
		}
	}
	for i := 0; i < vr.NumMethod(); i++ {
		vm := vr.Method(i)
		if vm.Type().AssignableTo(hf) {
			lname := strings.ToLower(tr.Method(i).Name)
			if lname == "index" {
				lname = ""
			}
			p := pre + lname
			call(p, vm)
		}
	}
}

// all post
func RouterPost(g *gin.RouterGroup, h interface{}, pre string) {
	post := reflect.ValueOf(g).MethodByName("POST")
	RouterSelf(h, pre, func(url string, gmethod reflect.Value) {
		vs := []reflect.Value{reflect.ValueOf(url), gmethod}
		post.Call(vs)
	})
}

// all get
func RouterGet(g *gin.RouterGroup, h interface{}, pre string) {
	get := reflect.ValueOf(g).MethodByName("GET")
	RouterSelf(h, pre, func(url string, gmethod reflect.Value) {
		vs := []reflect.Value{reflect.ValueOf(url), gmethod}
		get.Call(vs)
	})
}

// GinDotTypeLives generate data for structural  dot
// routerID: is the liveid of  gindot/router component
func PreAddControlDot(ctype reflect.Type, routerID dot.LiveID) *dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{TypeID: dot.TypeID(ctype.Name()), RefType: ctype, NewDoter: nil},
		Lives: []dot.Live{dot.Live{
			LiveID:    dot.LiveID(ctype.Name()),
			RelyLives: map[string]dot.LiveID{"GinRouter_": routerID},
		}},
	}
	return tl
}
