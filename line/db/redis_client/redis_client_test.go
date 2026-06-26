package redis_client //nolint:golint

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func TestRedisClient_GetVersions(t *testing.T) {
	mRedis := miniredis.RunT(t)
	redisClient := makeTestRedisClientFalse(mRedis.Addr())
	assert.NotNil(t, redisClient)
	category := "test"
	version := "1.0"
	itemKey1 := "tkey"
	value1 := "value1"

	{
		redisClient.SetVersion(category, version)
		getVersion, err := redisClient.GetVersion(category)
		assert.Nil(t, err)
		assert.Equal(t, version, getVersion)
		getVersions, err := redisClient.GetVersions(category)
		assert.Nil(t, err)
		assert.Equal(t, version, getVersions[0])
		time.Sleep(2 * time.Second)
	}

	{
		getValue, err := redisClient.Get(category, version, itemKey1)
		assert.NotNil(t, err)
		assert.Equal(t, "", getValue)
		_, err = redisClient.Set(category, version, itemKey1, value1, 0)
		assert.Nil(t, err)

		getValue, err = redisClient.Get(category, version, itemKey1)
		assert.Nil(t, err)
		assert.Equal(t, value1, getValue)

		count, err := redisClient.Del(category, version, itemKey1)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), count)
		getValue, err = redisClient.Get(category, version, itemKey1)
		assert.NotNil(t, err)
		assert.Equal(t, "", getValue)
	}

	{
		redisClient.Set(category, version, itemKey1, value1, 0)
		err := redisClient.CleanVersion(category, version)
		assert.Nil(t, err)
		getValue, err := redisClient.Get(category, version, itemKey1)
		assert.NotNil(t, err)
		assert.Equal(t, "", getValue)

		getVersion, err := redisClient.GetVersion(category)
		assert.NotNil(t, err)
		assert.Equal(t, "", getVersion)
		getVersions, err := redisClient.GetVersions(category)
		assert.Nil(t, err)
		assert.Equal(t, []string{}, getVersions)
	}

	redisClient = makeTestRedisClientTrue(mRedis.Addr())
	assert.NotEqual(t, nil, redisClient)
	{
		redisClient.SetVersion(category, version)
		getVersion, err := redisClient.GetVersion(category)
		assert.Nil(t, err)
		assert.Equal(t, version, getVersion)
		getVersions, err := redisClient.GetVersions(category)
		assert.Nil(t, err)
		assert.Equal(t, version, getVersions[0])
	}

	{
		getValue, err := redisClient.Get(category, version, itemKey1)
		assert.NotNil(t, err)
		assert.Equal(t, "", getValue)
		_, err = redisClient.Set(category, version, itemKey1, value1, 0)
		assert.Nil(t, err)

		getValue, err = redisClient.Get(category, version, itemKey1)
		assert.Nil(t, err)
		assert.Equal(t, value1, getValue)

		count, err := redisClient.Del(category, version, itemKey1)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), count)
		getValue, err = redisClient.Get(category, version, itemKey1)
		assert.NotNil(t, err)
		assert.Equal(t, "", getValue)
	}

	{
		redisClient.Set(category, version, itemKey1, value1, 0)
		err := redisClient.CleanVersion(category, version)
		assert.Nil(t, err)
		getValue, err := redisClient.Get(category, version, itemKey1)
		assert.NotNil(t, err)
		assert.Equal(t, "", getValue)

		getVersion, err := redisClient.GetVersion(category)
		assert.NotNil(t, err)
		assert.Equal(t, "", getVersion)
		getVersions, err := redisClient.GetVersions(category)
		assert.Nil(t, err)
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

func makeTestRedisClientTrue(addr string) *RedisClient {
	if clientTrue != nil {
		return clientTrue
	}
	config := &RedisConfig{
		Addr:                addr,
		KeepMaxVersionCount: 3,
		VersionFromRedis:    true,
	}
	clientTrue, clientCancelTrue, _ = NewRedisClient(config)
	return clientTrue
}
func makeTestRedisClientFalse(addr string) *RedisClient {
	if clientFalse != nil {
		return clientFalse
	}
	config := &RedisConfig{
		Addr:                addr,
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
