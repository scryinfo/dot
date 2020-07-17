package main

import (
    "fmt"
    "time"
)

// authenticate basic set and get funcs, and other demo will ignore error returned if not necessary
func basicDemo() {
    fmt.Println("-------")

    // query in cache
    v, err := getValue("demo")
    if err != nil {
        fmt.Printf("get value failed, error: %v, value is empty str?: %t\n", err, v == "")
    }

    // simulate query not find, skip save to db and update cache
    err = setValue("demo", "basic demo")
    if err != nil {
        fmt.Println("set value failed, error:", err)
    }

    v, err = getValue("demo")
    if err != nil {
        fmt.Println("get value failed, error:", err)
    }

    fmt.Println("show value after set:", v)
}

func expireDemo() {
    fmt.Println("-------")

    err := setValue("demo", "expire demo", time.Second*2)
    if err != nil {
        fmt.Println("set value failed, error:", err)
    }
    fmt.Println("expire time: 2s")

    v, err := getValue("demo")
    fmt.Printf("get value immediately, value: |%s|, error: %v\n", v, err)

    time.Sleep(time.Second)

    v, err = getValue("demo")
    fmt.Printf("get value after 2s   , value: |%s|, error: %v\n", v, err)
}

func updateDemo() {
    fmt.Println("-------")

    _ = setValue("demo", "update demo", 0)
    v, _ := getValue("demo")
    fmt.Println("value:", v)

    err := setValue("demo", "UPDATE VALUE", 0)
    if err != nil {
        fmt.Println("update value failed, error:", err)
    }

    v, _ = getValue("demo")
    fmt.Println("value(updated):", v)

    err = setValue("demo", "UPDATE VALUE ONCE AGAIN", time.Second*2)
    if err != nil {
        fmt.Println("update value or set expire time failed, error:", err)
    }

    v, _ = getValue("demo")
    fmt.Println("value(updated):", v)
}
