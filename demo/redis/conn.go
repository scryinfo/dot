package main

import "github.com/redis/go-redis/v9"

func conn(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: "",
		Password: "",
		DB:       0,
	})
}
