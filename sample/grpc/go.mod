module github.com/scryinfo/dot/sample/grpc

go 1.12

require (
	github.com/golang/protobuf v1.3.1
	github.com/scryinfo/dot v0.1.3-0.20190606035354-2f67d9575b1d
	github.com/scryinfo/dot/dots/grpc v0.0.0-20190605102428-60c55a6bb7e9
	github.com/scryinfo/scryg v0.1.3-0.20190523074957-3a6377ac45ea
	golang.org/x/net v0.0.0-20190603091049-60506f45cf65 // indirect
	golang.org/x/sys v0.0.0-20190602015325-4c4f7f33c9ed // indirect
	google.golang.org/grpc v1.21.0
)

replace (
	github.com/scryinfo/dot v0.0.0 => ../../
	github.com/scryinfo/dot/dots/grpc v0.0.0 => ../../dots/grpc
)
