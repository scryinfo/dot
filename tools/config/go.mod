module github.com/scryinfo/dot/tools/config

go 1.14

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/golang/protobuf v1.3.4
	github.com/gookit/config v1.1.0
	github.com/kr/pty v1.1.4 // indirect
	github.com/scryinfo/dot v0.1.5-0.20200711025551-7ba9a5161bd4
	github.com/scryinfo/dot/dots/gindot v0.0.0-20200520093457-f8a16513551b
	github.com/scryinfo/dot/dots/grpc v0.0.0-20191121024911-104dd77f0dab
	github.com/scryinfo/scryg v0.1.3
	github.com/stretchr/objx v0.2.0 // indirect
	go.uber.org/zap v1.14.0
	golang.org/x/tools v0.0.0-20200310231627-71bfc1b943ce
	google.golang.org/grpc v1.28.0
	gopkg.in/yaml.v2 v2.2.7
)

replace (
	github.com/scryinfo/dot => ../../
	github.com/scryinfo/dot/dots/gindot => ../../dots/gindot
	github.com/scryinfo/dot/dots/grpc => ../../dots/grpc
)
