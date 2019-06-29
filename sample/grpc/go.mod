module github.com/scryinfo/dot/sample/grpc

go 1.12

require (
	github.com/golang/protobuf v1.3.1
	github.com/scryinfo/dot v0.1.3-0.20190629071017-dbfa9f04a27e
	github.com/scryinfo/dot/dots/grpc v0.0.0-20190629072008-57d44dc20b86
	github.com/scryinfo/scryg v0.1.3-0.20190608053141-a292b801bfd6
	go.uber.org/zap v1.10.0
	google.golang.org/appengine v1.4.0 // indirect
	google.golang.org/genproto v0.0.0-20190620144150-6af8c5fc6601 // indirect
	google.golang.org/grpc v1.21.1
)

replace (
	github.com/scryinfo/dot v0.1.3-0.20190625101940-1336d6ee5a13 => ../../
	github.com/scryinfo/dot/dots/grpc v0.0.0-20190619080120-397c86d1d25b => ../../dots/grpc
)
