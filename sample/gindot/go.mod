module github.com/scryinfo/dot/sample/gindot

go 1.12

require (
	github.com/gin-gonic/gin v1.4.0
	github.com/scryinfo/dot v0.1.3-0.20190531094641-c6405bc7be23
	github.com/scryinfo/dot/dots/gindot v0.0.0
	github.com/scryinfo/scryg v0.1.3-0.20190523074957-3a6377ac45ea
)

replace (
	github.com/scryinfo/dot v0.1.3-0.20190530060438-14a6f5f91e65 => ../../
	github.com/scryinfo/dot/dots/gindot v0.0.0 => ../../dots/gindot
)
