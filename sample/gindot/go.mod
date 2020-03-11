module github.com/scryinfo/dot/sample/gindot

go 1.14

require (
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/gin-gonic/gin v1.5.0
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/scryinfo/dot v0.1.4
	github.com/scryinfo/dot/dots/gindot v0.0.0-20200311030916-18de37ac25e4
	github.com/scryinfo/scryg v0.1.3
	go.uber.org/zap v1.14.0
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
)

replace (
	github.com/scryinfo/dot v0.1.4 => ../../
	github.com/scryinfo/dot/dots/gindot v0.0.0-20190705064650-8b2f44b376f8 => ../../dots/gindot/
)
