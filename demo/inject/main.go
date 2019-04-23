package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main()  {
	SetPrivatef()
}


func SetPrivatef()  {
	t := &T{fPtr: nil,fInterface:nil}
	f := reflect.ValueOf(t).Elem().FieldByName("fPtr")
	fmt.Println(f.CanAddr()) // have to true
	if f.CanAddr() {
		fp := f.Addr().Pointer()
		fpp := ((**int)(unsafe.Pointer(fp)))

		t2 := &T2{f2: 10}
		fp2 := (*int)(unsafe.Pointer(t2))

		*fpp = fp2
		fmt.Println(t.fPtr.f2)
	}

	f = reflect.ValueOf(t).Elem().FieldByName("fInterface")
	if f.CanAddr() {
		fp := f.Addr().Pointer()
		fpp := ((*interface{})(unsafe.Pointer(fp)))

		//var t2 Int2 = &Int2Imp{F: 20}
		//fp2 := (*interface{})(unsafe.Pointer(&t2))
		//*fpp = *fp2

		var t2 interface{} = &Int2Imp{F: 20}
		fp2 := (*interface{})(unsafe.Pointer(reflect.ValueOf(t2).Pointer()))

		*fpp = fp2
		fmt.Println(t.fInterface.Get())
	}

}


type T struct {
	fPtr       *T2
	fInterface Int2
}

type T2 struct {
	f2 int
}

type Int2 interface {
	Get() int
}

type Int2Imp struct {
	F int
}

func (c *Int2Imp) Get() int {
	return c.F
}