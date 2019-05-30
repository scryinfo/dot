module github.com/scryinfo/dot/dots/grpc

go 1.12

require (
	github.com/scryinfo/dot v0.1.3-0.20190530023729-40528e80ddb2
	github.com/scryinfo/scryg v0.1.3-0.20190523074957-3a6377ac45ea
	golang.org/x/net v0.0.0-20190522155817-f3200d17e092
	google.golang.org/grpc v1.21.0
)

replace (
	github.com/scryinfo/dot v0.0.0 => ../../
	github.com/scryinfo/dot/dots/grpc v0.0.0 => ../../dots/gindot
)
