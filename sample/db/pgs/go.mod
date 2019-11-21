module github.com/scryinfo/dot/sample/db/pgs

go 1.13

require (
	github.com/go-pg/pg v8.0.6+incompatible
	github.com/scryinfo/dot v0.1.3
	github.com/scryinfo/dot/dots/db/pgs v0.0.0
)

replace (
	github.com/scryinfo/dot v0.1.3 => ../../../
	github.com/scryinfo/dot/dots/db/pgs v0.0.0 => ../../../dots/db/pgs/
)
