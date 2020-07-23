package redis_client //nolint:golint

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedisClient_GetVersions(t *testing.T) {
	redisClient := GetRedisClient(false)
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

	redisClient = GetRedisClient(true)
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
}

func GetRedisClient(versionFromRedis bool) *RedisClient {
	config := &configRedis{
		Addr:                "127.0.0.1:6379",
		KeepMaxVersionCount: 3,
		VersionFromRedis:    versionFromRedis,
		TrySeconds:          80,
	}
	bs, _ := json.Marshal(config)
	return RedisClientTest(string(bs))
}
