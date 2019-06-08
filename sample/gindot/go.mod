module github.com/scryinfo/dot/sample/gindot

go 1.12

require (
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.4.0
	github.com/scryinfo/dot v0.1.3-0.20190608033438-4c5a4d63587d
	github.com/scryinfo/dot/dots/gindot v0.0.0-20190608033624-4e8973686871
	github.com/scryinfo/scryg v0.1.3-0.20190608032618-f4f2c5103cd2
	golang.org/x/net v0.0.0-20190607181551-461777fb6f67 // indirect
	golang.org/x/sys v0.0.0-20190606203320-7fc4e5ec1444 // indirect
)

replace (
	github.com/scryinfo/dot v0.1.3-0.20190530060438-14a6f5f91e65 => ../../
	github.com/scryinfo/dot/dots/gindot v0.0.0 => ../../dots/gindot
)
