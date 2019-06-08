module github.com/scryinfo/dot/sample/gindot

go 1.12

require (
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.4.0
	github.com/scryinfo/dot v0.1.3-0.20190608070725-56a7ef82fed1
	github.com/scryinfo/dot/dots/gindot v0.0.0-20190608071040-ca1e16c496ad
	github.com/scryinfo/scryg v0.1.3-0.20190608053141-a292b801bfd6
	golang.org/x/net v0.0.0-20190607181551-461777fb6f67 // indirect
	golang.org/x/sys v0.0.0-20190608050228-5b15430b70e3 // indirect
)

replace (
	github.com/scryinfo/dot v0.1.3-0.20190530060438-14a6f5f91e65 => ../../
	github.com/scryinfo/dot/dots/gindot v0.0.0 => ../../dots/gindot
)
