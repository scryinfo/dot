package redisdot

import (
    "fmt"
    "github.com/go-redis/redis/v8"
)

var confStr = `{"addr":"192.168.1.65:6379"}`

func ExampleBasicProcess() {
    rdb := GenerateRedis(confStr)

    rdb.FlushAll(rdb.Context())

    // simulate query in cache first (no result)
    v1, err := rdb.Get(rdb.Context(), "demo").Result()
    if err != redis.Nil {
        fmt.Println("Example: get value not run as suppose, error:", err)
        return
    }

    // skip query in db, only simulate update cache
    if err = rdb.Set(rdb.Context(), "demo", "basic process demo", 0).Err(); err != nil {
        fmt.Println("Example: set value failed, error:", err)
        return
    }

    // suppose a request comes now, query in cache (has result)
    v2, err := rdb.Get(rdb.Context(), "demo").Result()
    if err != nil {
        fmt.Println("Example: get value failed, error:", err)
        return
    }

    fmt.Printf("v1: %s, v2: %s\n", v1, v2)

    // Output: v1: , v2: basic process demo

    return
}
