module github.com/scryinfo/dot/sample/gindot

go 1.12

require (
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.4.0
	github.com/scryinfo/dot v0.1.3-0.20190615010906-d1f546a24302
	github.com/scryinfo/dot/dots/gindot v0.0.0-20190615013714-8fedb67ed253
	github.com/scryinfo/scryg v0.1.3-0.20190608053141-a292b801bfd6
	golang.org/x/net v0.0.0-20190613194153-d28f0bde5980 // indirect
	golang.org/x/sys v0.0.0-20190614160838-b47fdc937951 // indirect
)

replace (
	github.com/scryinfo/dot v0.1.3-0.20190530060438-14a6f5f91e65 => ../../
	github.com/scryinfo/dot/dots/gindot v0.0.0 => ../../dots/gindot
)
