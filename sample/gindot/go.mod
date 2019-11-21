module github.com/scryinfo/dot/sample/gindot

go 1.12

require (
	github.com/gin-gonic/gin v1.4.0
	github.com/scryinfo/dot v0.1.3
	github.com/scryinfo/dot/dots/gindot v0.0.0-20191121022614-959828ad21d4
	github.com/scryinfo/scryg v0.1.3-0.20190608053141-a292b801bfd6
	go.uber.org/zap v1.13.0
)

replace (
	github.com/scryinfo/dot v0.1.3 => ../../
	github.com/scryinfo/dot/dots/gindot v0.0.0-20190705064650-8b2f44b376f8 => ../../dots/gindot/
)
