module github.com/scryinfo/dot/dots/db/gorms

go 1.14

require (
	cloud.google.com/go v0.37.4 // indirect
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/jinzhu/gorm v1.9.12
	github.com/lib/pq v1.3.0 // indirect
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/pkg/errors v0.9.1
	github.com/scryinfo/dot v0.1.3
	go.uber.org/zap v1.14.0
	google.golang.org/appengine v1.6.5 // indirect
)

replace github.com/scryinfo/dot v0.1.3 => ../../../
