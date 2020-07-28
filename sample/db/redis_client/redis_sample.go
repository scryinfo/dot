package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/db/redis_client"
)

const (
	basicDemoKey   = "basic_demo"
	basicDemoValue = "basic value"

	versionControlDemoKey        = "vcd"
	versionControlDemoVersionOne = "1"
	versionControlDemoVersionTwo = "2"
)

type RedisSample struct {
	RedisClient *redis_client.RedisClient `dot:""`
}

func (c *RedisSample) basicDemo() {
	// simulate query in cache first (no result)
	v1, err := c.RedisClient.ClientV8().Get(c.RedisClient.ClientV8().Context(), basicDemoKey).Result()
	if err != redis.Nil {
		fmt.Println("Example: get value not run as suppose, error:", err)
		os.Exit(-1)
	}

	// skip query in db, only simulate update cache
	checkError("Example: set value failed", c.RedisClient.ClientV8().Set(c.RedisClient.ClientV8().Context(), basicDemoKey, basicDemoValue, 0).Err())

	// suppose a request comes now, query in cache (has result)
	v2, err := c.RedisClient.ClientV8().Get(c.RedisClient.ClientV8().Context(), basicDemoKey).Result()
	checkError("Example: get value failed", err)

	time.Sleep(time.Second)
	fmt.Printf("v1: , v2: %s| <- suppose |\n", basicDemoValue)
	fmt.Printf("v1: %s, v2: %s\n", v1, v2)

	fmt.Println("Node: Basic demo passed.")
	fmt.Println("-------")

	return
}

func (c *RedisSample) versionControlDemo() {
	// register
	checkError("Example: set version failed", c.RedisClient.SetVersion(versionControlDemoKey, versionControlDemoVersionOne))

	// wait for subscribe goroutine
	time.Sleep(time.Second)

	// get version
	v, err := c.RedisClient.GetVersion(versionControlDemoKey)
	checkError("Example: get version failed", err)
	fmt.Printf("Version get suppose: %s, actually: %s\n", versionControlDemoVersionOne, v)

	// set(update) version
	checkError("Example: set version failed", c.RedisClient.SetVersion(versionControlDemoKey, versionControlDemoVersionTwo))

	// get version
	v, err = c.RedisClient.GetVersion(versionControlDemoKey)
	checkError("Example: get version failed", err)
	fmt.Printf("Version get suppose: %s, actually: %s\n", versionControlDemoVersionTwo, v)

	// set version 3 times
	for i := 3; i < 6; i++ {
		checkError("Example: set version failed", c.RedisClient.SetVersion(versionControlDemoKey, strconv.Itoa(i)))
	}

	// get all versions
	versions, err := c.RedisClient.GetVersions(versionControlDemoKey)
	checkError("Example: get all versions failed", err)

	fmt.Println("All versions: 5 <- 4 <- 3 <- 2| <- suppose |")
	fmt.Println("All versions:", versions)

	fmt.Println("Node: version control demo passed.")
	fmt.Println("-------")

	return
}

func checkError(msg string, err error) {
	if err != nil {
		dot.Logger().Errorln(msg, zap.NamedError("error", err))
		os.Exit(-1)
	}

	return
}
