module github.com/scryinfo/dot/dots/db/gorms

go 1.13

require (
	github.com/jinzhu/gorm v1.9.11
	github.com/lib/pq v1.2.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/scryinfo/dot v0.1.3
	go.uber.org/zap v1.13.0
	google.golang.org/appengine v1.6.5 // indirect
)

replace github.com/scryinfo/dot v0.1.3 => ../../../
