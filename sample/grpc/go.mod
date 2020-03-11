module github.com/scryinfo/dot/sample/grpc

go 1.14

require (
	github.com/golang/protobuf v1.3.4
	github.com/json-iterator/go v1.1.8 // indirect
	github.com/mattn/go-isatty v0.0.10 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/scryinfo/dot v0.1.4
	github.com/scryinfo/dot/dots/grpc v0.0.0-20200311030916-18de37ac25e4
	github.com/scryinfo/scryg v0.1.3
	go.uber.org/zap v1.14.0
	golang.org/x/sys v0.0.0-20191120155948-bd437916bb0e // indirect
	google.golang.org/grpc v1.28.0
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
	gopkg.in/yaml.v2 v2.2.7 // indirect
)

replace (
	github.com/scryinfo/dot v0.1.4 => ../../
	github.com/scryinfo/dot/dots/gindot v0.0.0-20190622091252-bab0929bd7e7 => ../../dots/gindot
	github.com/scryinfo/dot/dots/gindot v0.0.0-20191026032307-4fe8cc8e04c9 => ../../dots/gindot
	github.com/scryinfo/dot/dots/grpc v0.0.0-20190705064910-5975ec5bbacd => ../../dots/grpc
)
