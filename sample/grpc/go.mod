module github.com/scryinfo/dot/sample/grpc

go 1.14

require (
	github.com/golang/protobuf v1.3.4
	github.com/scryinfo/dot v0.1.3
	github.com/scryinfo/dot/dots/grpc v0.0.0-20191121024911-104dd77f0dab
	github.com/scryinfo/scryg v0.1.3
	go.uber.org/zap v1.14.0
	google.golang.org/genproto v0.0.0-20190620144150-6af8c5fc6601 // indirect
	google.golang.org/grpc v1.21.1
)

replace (
	github.com/scryinfo/dot v0.1.3 => ../../
	github.com/scryinfo/dot/dots/gindot v0.0.0-20190622091252-bab0929bd7e7 => ../../dots/gindot
	github.com/scryinfo/dot/dots/gindot v0.0.0-20191026032307-4fe8cc8e04c9 => ../../dots/gindot
	github.com/scryinfo/dot/dots/grpc v0.0.0-20190705064910-5975ec5bbacd => ../../dots/grpc
)
