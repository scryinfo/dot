module github.com/scryinfo/dot/sample/grpc

go 1.12

require (
	github.com/golang/protobuf v1.3.1
	github.com/scryinfo/dot v0.1.3-0.20190625101940-1336d6ee5a13
	github.com/scryinfo/dot/dots/gindot v0.0.0-20190625102047-666d44ee7d72 // indirect
	github.com/scryinfo/dot/dots/grpc v0.0.0-20190625102047-666d44ee7d72
	github.com/scryinfo/scryg v0.1.3-0.20190608053141-a292b801bfd6
	go.uber.org/zap v1.10.0
	golang.org/x/sync v0.0.0-20181108010431-42b317875d0f // indirect
	google.golang.org/appengine v1.4.0 // indirect
	google.golang.org/genproto v0.0.0-20190620144150-6af8c5fc6601 // indirect
	google.golang.org/grpc v1.21.1
)

replace (
	github.com/scryinfo/dot v0.1.3-0.20190625101940-1336d6ee5a13 => ../../
	github.com/scryinfo/dot/dots/grpc v0.0.0-20190619080120-397c86d1d25b => ../../dots/grpc
)
