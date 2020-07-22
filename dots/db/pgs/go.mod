module github.com/scryinfo/dot/dots/db/pgs

go 1.14

require (
	github.com/go-pg/pg/v9 v9.1.3
	github.com/pkg/errors v0.9.1 // indirect
	github.com/scryinfo/dot v0.1.4
	github.com/scryinfo/scryg v0.1.3
	github.com/vmihailenco/msgpack/v4 v4.3.9 // indirect
	go.uber.org/multierr v1.5.0 // indirect
	go.uber.org/zap v1.14.0 // indirect
	golang.org/x/crypto v0.0.0-20200302210943-78000ba7a073 // indirect
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/tools v0.0.0-20200310231627-71bfc1b943ce // indirect
)

replace (
	github.com/scryinfo/dot v0.1.4 => ../../../
)
