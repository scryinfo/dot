module github.com/scryinfo/dot/sample/grpc

go 1.12

require (
	github.com/golang/protobuf v1.3.1
	github.com/scryinfo/dot v0.1.3-0.20190608033438-4c5a4d63587d
	github.com/scryinfo/dot/dots/grpc v0.0.0-20190608033624-4e8973686871
	github.com/scryinfo/scryg v0.1.3-0.20190608032618-f4f2c5103cd2
	golang.org/x/net v0.0.0-20190607181551-461777fb6f67 // indirect
	golang.org/x/sys v0.0.0-20190606203320-7fc4e5ec1444 // indirect
	google.golang.org/genproto v0.0.0-20190605220351-eb0b1bdb6ae6 // indirect
	google.golang.org/grpc v1.21.1
)

replace (
	github.com/scryinfo/dot v0.0.0 => ../../
	github.com/scryinfo/dot/dots/grpc v0.0.0 => ../../dots/grpc
)
