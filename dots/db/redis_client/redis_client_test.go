package redis_client //nolint:golint

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRedisClient_GetVersions(t *testing.T) {
	redisClient := makeTestRedisClientFalse()
	assert.NotEqual(t, nil, redisClient)
	category := "test"
	version := "1.0"
	itemKey1 := "tkey"
	value1 := "value1"

	{
		redisClient.SetVersion(category, version)
		getVersion, err := redisClient.GetVersion(category)
		assert.Equal(t, nil, err)
		assert.Equal(t, version, getVersion)
		getVersions, err := redisClient.GetVersions(category)
		assert.Equal(t, nil, err)
		assert.Equal(t, version, getVersions[0])
		time.Sleep(2 * time.Second)
	}

	{
		getValue, err := redisClient.Get(category, version, itemKey1)
		assert.NotNil(t, err)
		assert.Equal(t, "", getValue)
		_, err = redisClient.Set(category, version, itemKey1, value1, 0)
		assert.Equal(t, nil, err)

		getValue, err = redisClient.Get(category, version, itemKey1)
		assert.Equal(t, nil, err)
		assert.Equal(t, value1, getValue)

		count, err := redisClient.Del(category, version, itemKey1)
		assert.Equal(t, nil, err)
		assert.Equal(t, int64(1), count)
		getValue, err = redisClient.Get(category, version, itemKey1)
		assert.NotNil(t, err)
		assert.Equal(t, "", getValue)
	}

	{
		redisClient.Set(category, version, itemKey1, value1, 0)
		err := redisClient.CleanVersion(category, version)
		assert.Equal(t, nil, err)
		getValue, err := redisClient.Get(category, version, itemKey1)
		assert.NotNil(t, err)
		assert.Equal(t, "", getValue)

		getVersion, err := redisClient.GetVersion(category)
		assert.NotNil(t, err)
		assert.Equal(t, "", getVersion)
		getVersions, err := redisClient.GetVersions(category)
		assert.Equal(t, nil, err)
		assert.Equal(t, []string{}, getVersions)
	}

	redisClient = makeTestRedisClientTrue()
	assert.NotEqual(t, nil, redisClient)
	{
		redisClient.SetVersion(category, version)
		getVersion, err := redisClient.GetVersion(category)
		assert.Equal(t, nil, err)
		assert.Equal(t, version, getVersion)
		getVersions, err := redisClient.GetVersions(category)
		assert.Equal(t, nil, err)
		assert.Equal(t, version, getVersions[0])
	}

	{
		getValue, err := redisClient.Get(category, version, itemKey1)
		assert.NotNil(t, err)
		assert.Equal(t, "", getValue)
		_, err = redisClient.Set(category, version, itemKey1, value1, 0)
		assert.Equal(t, nil, err)

		getValue, err = redisClient.Get(category, version, itemKey1)
		assert.Equal(t, nil, err)
		assert.Equal(t, value1, getValue)

		count, err := redisClient.Del(category, version, itemKey1)
		assert.Equal(t, nil, err)
		assert.Equal(t, int64(1), count)
		getValue, err = redisClient.Get(category, version, itemKey1)
		assert.NotNil(t, err)
		assert.Equal(t, "", getValue)
	}

	{
		redisClient.Set(category, version, itemKey1, value1, 0)
		err := redisClient.CleanVersion(category, version)
		assert.Equal(t, nil, err)
		getValue, err := redisClient.Get(category, version, itemKey1)
		assert.NotNil(t, err)
		assert.Equal(t, "", getValue)

		getVersion, err := redisClient.GetVersion(category)
		assert.NotNil(t, err)
		assert.Equal(t, "", getVersion)
		getVersions, err := redisClient.GetVersions(category)
		assert.Equal(t, nil, err)
		assert.Equal(t, []string{}, getVersions)
	}
	cleanup()
}

var (
	clientTrue        *RedisClient
	clientFalse       *RedisClient
	clientCancelTrue  func()
	clientCancelFalse func()
)

func makeTestRedisClientTrue() *RedisClient {
	if clientTrue != nil {
		return clientTrue
	}
	config := &RedisConfig{
		Addr:                "127.0.0.1:6379",
		KeepMaxVersionCount: 3,
		VersionFromRedis:    true,
	}
	clientTrue, clientCancelTrue, _ = NewRedisClient(config)
	return clientTrue
}
func makeTestRedisClientFalse() *RedisClient {
	if clientFalse != nil {
		return clientFalse
	}
	config := &RedisConfig{
		Addr:                "127.0.0.1:6379",
		KeepMaxVersionCount: 3,
		VersionFromRedis:    false,
	}
	clientFalse, clientCancelFalse, _ = NewRedisClient(config)
	return clientFalse
}

func cleanup() {
	clientCancelTrue()
	clientCancelFalse()
}
