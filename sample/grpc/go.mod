module github.com/scryinfo/dot/sample/grpc

go 1.12

require (
	github.com/golang/protobuf v1.3.2
	github.com/scryinfo/dot v0.1.3-0.20190907084536-c60f8a67fccd
	github.com/scryinfo/dot/dots/grpc v0.0.0-20190705064910-5975ec5bbacd
	github.com/scryinfo/scryg v0.1.3-0.20190608053141-a292b801bfd6
	go.uber.org/zap v1.11.0
	google.golang.org/genproto v0.0.0-20190620144150-6af8c5fc6601 // indirect
	google.golang.org/grpc v1.21.1
)

replace (
	github.com/scryinfo/dot v0.1.3-0.20190705064446-6614e45bf155 => ../../
	github.com/scryinfo/dot/dots/grpc v0.0.0-20190705064910-5975ec5bbacd => ../../dots/grpc
)
