package main

import (
    "context"
    "fmt"
    "github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func ExampleNewClient() {
    rdb := redis.NewClient(&redis.Options{
        Addr:     "192.168.1.65:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    pong, err := rdb.Ping(ctx).Result()
    fmt.Println(pong, err)
    // Output: PONG <nil>
}

func ExampleClient() {
    rdb := redis.NewClient(&redis.Options{
        Addr:     "192.168.1.65:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })
    err := rdb.Set(ctx, "key", "value", 0).Err()
    if err != nil {
        panic(err)
    }

    val, err := rdb.Get(ctx, "key").Result()
    if err != nil {
        panic(err)
    }
    fmt.Println("key", val)

    val2, err := rdb.Get(ctx, "key2").Result()
    if err == redis.Nil {
        fmt.Println("key2 does not exist")
    } else if err != nil {
        panic(err)
    } else {
        fmt.Println("key2", val2)
    }
    // Output: key value
    // key2 does not exist
}
