module github.com/scryinfo/dot/tools/config

go 1.12

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/golang/protobuf v1.3.2
	github.com/gookit/config v1.1.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/scryinfo/dot v0.1.3-0.20190705064446-6614e45bf155
	github.com/scryinfo/dot/dots/gindot v0.0.0-20190827012258-f2d81985edf7
	github.com/scryinfo/dot/dots/grpc v0.0.0-20190827012258-f2d81985edf7
	github.com/scryinfo/scryg v0.1.3-0.20190608053141-a292b801bfd6
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0
	golang.org/x/tools v0.0.0-20190809145639-6d4652c779c4
	google.golang.org/grpc v1.22.1
	gopkg.in/yaml.v2 v2.2.2
)

replace (
	github.com/scryinfo/dot v0.1.3-0.20190705064446-6614e45bf155 => ../../
	github.com/scryinfo/dot/dots/grpc v0.0.0-20190823030926-b7234e66ebf4 => ../../dots/grpc
)
