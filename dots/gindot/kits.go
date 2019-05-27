package gindot

import (
	"github.com/gin-gonic/gin"
	"github.com/scryinfo/dot/dot"
	"reflect"
	"strings"
)

//for each funcs that like “gin.HandlerFunc”
//sampe pre = "scry", and  "func (c * SampleCtroller) Hello(cxt *gin.Context) {}", the url is "/scry/hello"
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

//all post
func RouterPost(g *gin.RouterGroup, h interface{}, pre string) {
	post := reflect.ValueOf(g).MethodByName("POST")
	RouterSelf(h, pre, func(url string, gmethod reflect.Value) {
		vs := []reflect.Value{reflect.ValueOf(url), gmethod}
		post.Call(vs)
	})
}

//all get
func RouterGet(g *gin.RouterGroup, h interface{}, pre string) {
	get := reflect.ValueOf(g).MethodByName("GET")
	RouterSelf(h, pre, func(url string, gmethod reflect.Value) {
		vs := []reflect.Value{reflect.ValueOf(url), gmethod}
		get.Call(vs)
	})
}

//TypeLiveGinDot generate data for structural  dot
//routerId: is the liveid of  gindot/router component
func PreAddControlDot(ctype reflect.Type, routerId dot.LiveId) *dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{TypeId: dot.TypeId(ctype.Name()), RefType: ctype, NewDoter: nil},
		Lives: []dot.Live{dot.Live{
			LiveId:    dot.LiveId(ctype.Name()),
			RelyLives: map[string]dot.LiveId{"GinRouter_": routerId},
		}},
	}
	return tl
}
