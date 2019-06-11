module github.com/scryinfo/dot/sample/grpc

go 1.12

require (
	github.com/golang/protobuf v1.3.1
	github.com/scryinfo/dot v0.1.3-0.20190611102425-8704955cbba0
	github.com/scryinfo/dot/dots/grpc v0.0.0-20190611102605-52d9c61cc213
	github.com/scryinfo/scryg v0.1.3-0.20190608053141-a292b801bfd6
	go.uber.org/zap v1.10.0
	golang.org/x/net v0.0.0-20190607181551-461777fb6f67 // indirect
	golang.org/x/sys v0.0.0-20190608050228-5b15430b70e3 // indirect
	google.golang.org/genproto v0.0.0-20190605220351-eb0b1bdb6ae6 // indirect
	google.golang.org/grpc v1.21.1
)

replace (
	github.com/scryinfo/dot v0.0.0 => ../../
	github.com/scryinfo/dot/dots/grpc v0.0.0 => ../../dots/grpc
)
