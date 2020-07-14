package main

import (
	"fmt"
	"github.com/albrow/zoom"
	"strconv"
)

func setSomeValues(num int) {
	for i := 0; i < num; i++ {
		save(sc, &S{Str: "index " + strconv.Itoa(i), I: i+1})
	}

	return
}

func NewCollection() {
	// new collection, option param is for 'find all' func
	sc, err = pool.NewCollectionWithOptions(&S{}, zoom.CollectionOptions{Index: true})
	if err != nil {
		fmt.Println("create collection failed, error:", err)
	}
}

func save(collection *zoom.Collection, model *S) {
	if err = collection.Save(model); err != nil {
		fmt.Println("set value failed, error:", err)
	}
}

func update(model *S) {
	if err = sc.SaveFields([]string{"Str"}, model); err != nil {
		fmt.Println("update value failed, error:", err)
	}
}

func find(id string) *S {
	res := &S{}
	if err = sc.Find(id, res); err != nil {
		fmt.Println("find value failed, error:", err)
	}

	return res
}

func findCertain(id string, fieldNames []string) *S {
	res := &S{}
	if err = sc.FindFields(id, fieldNames, res); err != nil {
		fmt.Println("find certain values failed, error:", err)
	}

	return res
}

func findAll() {
	res := make([]*S, 0)
	if err = sc.FindAll(&res); err != nil {
		fmt.Println("find all failed, error:", err)
	}
	var i int
	if i, err = sc.Count(); err != nil {
		fmt.Println("count summary failed, error:", err)
	}
	fmt.Printf("Node: show find all res: (total: %d)\n", i)
	for i := range res {
		fmt.Printf("index: %d, value:%#v\n", i, res[i])
	}

	return
}

func del(id string) {
	var ok bool
	if ok, err = sc.Delete(id); !ok || err != nil {
		fmt.Println("del id failed, error:", err)
	}

	return
}
