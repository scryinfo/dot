module github.com/scryinfo/dot/tools/config

go 1.14

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/golang/protobuf v1.4.2
	github.com/gookit/config v1.1.0
	github.com/scryinfo/dot v0.1.5-0.20200711025551-7ba9a5161bd4
	github.com/scryinfo/scryg v0.1.3
	go.uber.org/zap v1.15.0
	golang.org/x/tools v0.0.0-20200207183749-b753a1ba74fa
	google.golang.org/grpc v1.29.1
	gopkg.in/yaml.v2 v2.2.8
)

replace github.com/scryinfo/dot => ../
