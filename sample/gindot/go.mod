module github.com/scryinfo/dot/sample/gindot

go 1.12

require (
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.4.0
	github.com/scryinfo/dot v0.1.3-0.20190625102047-666d44ee7d72
	github.com/scryinfo/dot/dots/gindot v0.0.0-20190625102047-666d44ee7d72
	github.com/scryinfo/scryg v0.1.3-0.20190608053141-a292b801bfd6
	go.uber.org/zap v1.10.0
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859 // indirect
	golang.org/x/sys v0.0.0-20190624142023-c5567b49c5d0 // indirect
)

replace (
	github.com/scryinfo/dot v0.1.3-0.20190625102047-666d44ee7d72 => ../../
	github.com/scryinfo/dot/dots/gindot v0.0.0-20190625102047-666d44ee7d72 => ../../dots/gindot/
)
