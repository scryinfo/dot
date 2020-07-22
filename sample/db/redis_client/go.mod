module github.com/scryinfo/dot/sample/db/redis_client

go 1.14

require (
	github.com/go-redis/redis/v8 v8.0.0-beta.6
	github.com/pkg/errors v0.9.1
	github.com/scryinfo/dot v0.1.5-0.20200711025551-7ba9a5161bd4
	github.com/scryinfo/dot/dots/db/redis_client v0.0.0-20200711033836-fdd979f912ac
	github.com/scryinfo/scryg v0.1.3
	go.uber.org/zap v1.15.0
)

replace github.com/scryinfo/dot/dots/db/redis_client => ./../../../dots/db/redis_client
