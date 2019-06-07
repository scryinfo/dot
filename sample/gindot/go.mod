module github.com/scryinfo/dot/sample/gindot

go 1.12

require (
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.4.0
	github.com/scryinfo/dot v0.1.3-0.20190607005633-fbeee0d18475
	github.com/scryinfo/dot/dots/gindot v0.0.0-20190607005812-02c9f0eaee8f
	github.com/scryinfo/scryg v0.1.3-0.20190523074957-3a6377ac45ea
	golang.org/x/net v0.0.0-20190606173856-1492cefac77f // indirect
	golang.org/x/sys v0.0.0-20190606203320-7fc4e5ec1444 // indirect
)

replace (
	github.com/scryinfo/dot v0.1.3-0.20190530060438-14a6f5f91e65 => ../../
	github.com/scryinfo/dot/dots/gindot v0.0.0 => ../../dots/gindot
)
