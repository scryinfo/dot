module github.com/scryinfo/dot/sample/grpc

go 1.12

require (
	github.com/golang/protobuf v1.3.1
	github.com/scryinfo/dot v0.1.3-0.20190530023729-40528e80ddb2
	github.com/scryinfo/dot/dots/grpc v0.0.0-20190530023729-40528e80ddb2
	github.com/scryinfo/scryg v0.1.3-0.20190523074957-3a6377ac45ea
	google.golang.org/grpc v1.21.0
)

replace (
	github.com/scryinfo/dot v0.0.0 => ../../
	github.com/scryinfo/dot/dots/grpc v0.0.0 => ../../dots/grpc
)
