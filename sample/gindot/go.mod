module github.com/scryinfo/dot/sample/gindot

go 1.12

require (
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.4.0
	github.com/scryinfo/dot v0.1.3-0.20190606094239-93914ee47449
	github.com/scryinfo/dot/dots/gindot v0.0.0-20190606094545-eecf8fb80956
	github.com/scryinfo/scryg v0.1.3-0.20190523074957-3a6377ac45ea
	golang.org/x/net v0.0.0-20190603091049-60506f45cf65 // indirect
	golang.org/x/sys v0.0.0-20190602015325-4c4f7f33c9ed // indirect
)

replace (
	github.com/scryinfo/dot v0.1.3-0.20190530060438-14a6f5f91e65 => ../../
	github.com/scryinfo/dot/dots/gindot v0.0.0 => ../../dots/gindot
)
