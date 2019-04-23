package main

import (
	"fmt"
	"reflect"
)

func main()  {
	t := &T{f:10}

	v := reflect.ValueOf(t).Elem()

	f := v.FieldByName("f")
	fmt.Println(f.CanSet())

	f.Set(reflect.ValueOf(10))

	fmt.Println(t.f)

}

type T struct {
	f int
}
