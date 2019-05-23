package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	SetPrivateField()
}

func SetPrivateField() {
	{ //Pointer
		t := &T{fPtr: nil, fInterface: nil}
		field := reflect.ValueOf(t).Elem().FieldByName("fPtr")
		if field.CanAddr() { // have to true
			fp := field.Addr().Pointer()             //Get the field address
			fpp := ((**uintptr)(unsafe.Pointer(fp))) //Convert to pointer
			t2 := &T2{f2: 10}
			fp2 := (*uintptr)(unsafe.Pointer(t2))
			*fpp = fp2
			fmt.Println(t.fPtr.f2)
		}
	}
	{ //interface
		t := &T{fPtr: nil, fInterface: nil}
		field := reflect.ValueOf(t).Elem().FieldByName("fInterface")
		if field.CanAddr() {
			fp := field.Addr().Pointer()
			fpp := ((*interface{})(unsafe.Pointer(fp)))

			{ // Method 1, need to confirm interface type when compiling
				var t2 Inter2 = &Inter2Imp{F: 20}
				fp2 := (*interface{})(unsafe.Pointer(&t2))
				*fpp = *fp2
			}

			{ // Method 2, Generally use reflection to realize it
				var t2 interface{} = &Inter2Imp{F: 15}
				v := reflect.ValueOf(t2).Convert(field.Type()) // Must convert to  field type,since type in t2 is interface{}, not Inter2 type
				v2 := reflect.ValueOf(v)                       // Get ptr value through reflection, cannot use Pointer here，it will panic，Since t2 is not pointer type
				ptr := v2.FieldByName("ptr")
				fp2 := (*interface{})(unsafe.Pointer(ptr.Pointer()))
				*fpp = *fp2
			}
		}
		if t.fInterface != nil {
			fmt.Println(t.fInterface.Get())
		}
	}

}

type T struct {
	fPtr       *T2
	fInterface Inter2
}

type T2 struct {
	f2 int
}

type Inter2 interface {
	Get() int
}

type Inter2Imp struct {
	F int
}

func (c *Inter2Imp) Get() int {
	return c.F
}
