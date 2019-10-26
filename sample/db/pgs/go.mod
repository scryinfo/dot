module github.com/scryinfo/dot/sample/db/pgs

go 1.13

require (
	github.com/go-pg/pg v8.0.6+incompatible
	github.com/scryinfo/dot v0.1.3-0.20191026032307-4fe8cc8e04c9
	github.com/scryinfo/dot/dots/db/pgs v0.0.0 // indirect
)

replace (
	github.com/scryinfo/dot v0.1.3-0.20190827105138-5c8f5fae41f0 => ../../../
	github.com/scryinfo/dot/dots/db/pgs v0.0.0 => ../../../dots/db/pgs/
)
