package main

import "github.com/go-redis/redis/v8"

func conn(addr string) *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr:     addr,
        Username: "",
        Password: "",
        DB:       0,
    })
}
