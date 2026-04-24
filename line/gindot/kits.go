// Scry Info.  All rights reserved.
// license that can be found in the license file.

package gindot

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

// for each funcs that like “gin.HandlerFunc”
// sampe pre = "scry", and  "func (c * SampleCtroller) Hello(cxt *gin.Context) {}", the url is "/scry/hello"
func RouterSelf(h any, pre string, call func(url string, gmethod reflect.Value)) {
	hf := reflect.TypeFor[gin.HandlerFunc]()
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
func RouterPost(g *gin.RouterGroup, h any, pre string) {
	post := reflect.ValueOf(g).MethodByName("POST")
	RouterSelf(h, pre, func(url string, gmethod reflect.Value) {
		vs := []reflect.Value{reflect.ValueOf(url), gmethod}
		post.Call(vs)
	})
}

// all get
func RouterGet(g *gin.RouterGroup, h any, pre string) {
	get := reflect.ValueOf(g).MethodByName("GET")
	RouterSelf(h, pre, func(url string, gmethod reflect.Value) {
		vs := []reflect.Value{reflect.ValueOf(url), gmethod}
		get.Call(vs)
	})
}
