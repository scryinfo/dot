module github.com/scryinfo/dot/dots/grpc

go 1.12

require (
	github.com/scryinfo/dot v0.1.3-0.20190607034741-c1f22eee2c61
	github.com/scryinfo/scryg v0.1.3-0.20190523074957-3a6377ac45ea
	golang.org/x/net v0.0.0-20190522155817-f3200d17e092
	golang.org/x/sync v0.0.0-20181108010431-42b317875d0f // indirect
	golang.org/x/sys v0.0.0-20190529164535-6a60838ec259 // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/appengine v1.4.0 // indirect
	google.golang.org/genproto v0.0.0-20190522204451-c2c4e71fbf69 // indirect
	google.golang.org/grpc v1.21.0
)

replace (
	github.com/scryinfo/dot v0.0.0 => ../../
	github.com/scryinfo/dot/dots/grpc v0.0.0 => ../../dots/gindot
)
