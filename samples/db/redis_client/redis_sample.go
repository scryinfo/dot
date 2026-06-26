package main

import (
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/db/redis_client"
)

const (
	basicDemoKey   = "basic_demo"
	basicDemoValue = "basic value"

	versionControlDemoKey        = "vcd"
	versionControlDemoVersionOne = "1"
	versionControlDemoVersionTwo = "2"
)

func NewRedisSample(client *redis_client.RedisClient) *RedisSample {
	return &RedisSample{RedisClient: client}
}

type RedisSample struct {
	RedisClient *redis_client.RedisClient
}

func (c *RedisSample) basicDemo() {
	// simulate query in cache first (no result)
	v1, err := c.RedisClient.ClientV9().Get(c.RedisClient.Context(), basicDemoKey).Result()
	if err != redis.Nil {
		dot.Logger.Error().AnErr("Example: get value not run as suppose, error:", err)
		os.Exit(-1)
	}

	// skip query in db, only simulate update cache
	checkError("Example: set value failed", c.RedisClient.ClientV9().Set(c.RedisClient.Context(), basicDemoKey, basicDemoValue, 0).Err())

	// suppose a request comes now, query in cache (has result)
	v2, err := c.RedisClient.ClientV9().Get(c.RedisClient.Context(), basicDemoKey).Result()
	checkError("Example: get value failed", err)

	time.Sleep(time.Second)
	dot.Logger.Info().Msgf("v1: %s, v2: %s| <- suppose |\n", v1, v2)
	dot.Logger.Info().Msgf("v1: %s, v2: %s\n", v1, v2)

	dot.Logger.Info().Msg("Node: Basic demo passed.")
	dot.Logger.Info().Msg("-------")

}

func (c *RedisSample) versionControlDemo() {
	// register
	checkError("Example: set version failed", c.RedisClient.SetVersion(versionControlDemoKey, versionControlDemoVersionOne))

	// wait for subscribe goroutine
	time.Sleep(time.Second)

	// get version
	v, err := c.RedisClient.GetVersion(versionControlDemoKey)
	checkError("Example: get version failed", err)
	dot.Logger.Info().Msgf("Version get suppose: %s, actually: %s\n", versionControlDemoVersionOne, v)

	// set(update) version
	checkError("Example: set version failed", c.RedisClient.SetVersion(versionControlDemoKey, versionControlDemoVersionTwo))

	// get version
	v, err = c.RedisClient.GetVersion(versionControlDemoKey)
	checkError("Example: get version failed", err)
	dot.Logger.Info().Msgf("Version get suppose: %s, actually: %s\n", versionControlDemoVersionTwo, v)

	// set version 3 times
	for i := 3; i < 6; i++ {
		checkError("Example: set version failed", c.RedisClient.SetVersion(versionControlDemoKey, strconv.Itoa(i)))
	}

	// get all versions
	versions, err := c.RedisClient.GetVersions(versionControlDemoKey)
	checkError("Example: get all versions failed", err)

	dot.Logger.Info().Msg("All versions: 5 <- 4 <- 3 <- 2| <- suppose |")
	dot.Logger.Info().Msgf("All versions: %v", versions)

	dot.Logger.Info().Msg("Node: version control demo passed.")
	dot.Logger.Info().Msg("-------")

}

func checkError(msg string, err error) {
	if err != nil {
		dot.Logger.Error().AnErr(msg, err).Send()
		os.Exit(-1)
	}
}
