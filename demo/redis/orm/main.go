package main

import (
	"fmt"
	"github.com/albrow/zoom"
)

type S struct {
	Str string `redis:"key_string"`
	I   int    `redis:"key_int"`
	zoom.RandomID
}

var (
	url = "192.168.1.65:6379"

	pool = zoom.NewPool(url)
	sc   *zoom.Collection

	err error
)

// zoom has nothing about key expire
func main() {
	defer pool.Close()

	NewCollection()

	// Attention: all demo are under collection 'sc'.

	var i int
	i, err = sc.DeleteAll()
	if err != nil {
		fmt.Println("delete all failed, error:", err)
	}
	fmt.Println("del items:", i)

	setSomeValues(3)

	// save
	sInsSave := &S{Str: "zoom demo", I: -1}
	save(sc, sInsSave)

	// update certain fields
	sInsSave.Str = "update"
	update(sInsSave)

	// find
	sInsFind := find(sInsSave.ID)
	fmt.Printf("Node: show find res, %#v\n", sInsFind)

	// find certain fields, attention on 'S.I' item, is default value
	sInsFind2 := findCertain(sInsSave.ID, []string{"Str"})
	fmt.Printf("Node: show find certain fields res, %#v\n", sInsFind2)

	// find all and count
	findAll()

	// del
	del(sInsSave.ID)

	findAll()

	// 关于zoom条件查询与事务的写法：
	// https://github.com/albrow/zoom#using-query-modifiers
	// https://github.com/albrow/zoom#transactions
}
