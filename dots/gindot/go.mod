module github.com/scryinfo/dot/dots/gindot

go 1.14

require (
	github.com/gin-gonic/gin v1.5.0
	github.com/golang/gddo v0.0.0-20191216155521-fbfc0f5e7810
	github.com/google/go-cmp v0.4.0 // indirect
	github.com/scryinfo/dot v0.1.4
	github.com/scryinfo/scryg v0.1.3
	go.uber.org/zap v1.14.0
)

replace (
	github.com/scryinfo/dot v0.1.4 => ../../
)