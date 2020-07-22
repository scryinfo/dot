module github.com/scryinfo/dot/dots/db/redis_client

go 1.14

require (
	github.com/go-redis/redis/v8 v8.0.0-beta.6
	github.com/pkg/errors v0.9.1
	github.com/scryinfo/dot v0.1.4
	github.com/stretchr/testify v1.6.1
	go.uber.org/zap v1.14.0
)

replace github.com/scryinfo/dot v0.1.4 => ../../../
