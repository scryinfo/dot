module github.com/scryinfo/dot/tools/config

go 1.12

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/gin-gonic/gin v1.4.0
	github.com/golang/protobuf v1.3.2
	github.com/gookit/config v1.1.0
	github.com/scryinfo/dot v0.1.3-0.20190705064446-6614e45bf155
	github.com/scryinfo/dot/dots/grpc v0.0.0-20190823030926-b7234e66ebf4
	github.com/scryinfo/scryg v0.1.3-0.20190608053141-a292b801bfd6
	github.com/stretchr/objx v0.2.0 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0
	golang.org/x/tools v0.0.0-20190809145639-6d4652c779c4
	google.golang.org/grpc v1.22.1
	gopkg.in/yaml.v2 v2.2.2
)

replace (
	github.com/scryinfo/dot v0.1.3-0.20190705064446-6614e45bf155 => ../../
	github.com/scryinfo/dot/dots/gindot v0.0.0-20190622091252-bab0929bd7e7 => ../gindot
)
