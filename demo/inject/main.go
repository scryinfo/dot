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
	{ //指针
		t := &T{fPtr: nil, fInterface: nil}
		field := reflect.ValueOf(t).Elem().FieldByName("fPtr")
		if field.CanAddr() { // have to true
			fp := field.Addr().Pointer()             //得到字段的地址
			fpp := ((**uintptr)(unsafe.Pointer(fp))) //转换为指指针
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

			{ // 方式一，需要在编译时确定interface的类型
				var t2 Inter2 = &Inter2Imp{F: 20}
				fp2 := (*interface{})(unsafe.Pointer(&t2))
				*fpp = *fp2
			}

			{ // 方式二，通用使用反射实现
				var t2 interface{} = &Inter2Imp{F: 15}
				v := reflect.ValueOf(t2).Convert(field.Type()) //一定使用Convert 转换为字段的类型， 因为t2的中的类型 为interface{},  不是 Inter2类型
				v2 := reflect.ValueOf(v)                       //通过反射取到 ptr的值， 这里不能使用Pointer，会panic，因为 t2不是指针类型
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
