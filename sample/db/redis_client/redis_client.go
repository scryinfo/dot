package main

import (
	"github.com/scryinfo/dot/dots/db/redis_client"
	"os"

	"go.uber.org/zap"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/line"
	"github.com/scryinfo/scryg/sutils/ssignal"
)

func main() {
	l, err := line.BuildAndStart(add)
	if err != nil {
		dot.Logger().Errorln("", zap.Error(err))
	}
	defer line.StopAndDestroy(l, true)
	dot.Logger().Infoln("dot ok")
	{ // get redisSample client do something
		redisSample := &RedisSample{}
		l.ToInjecter().Inject(redisSample)
		redisSample.Start(true)
	}

	ssignal.WaitCtrlC(func(s os.Signal) bool {
		return false
	})
}
func add(l dot.Line) error {
	lives := redis_client.RedisClientTypeLives()
	err := l.PreAdd(lives...)
	return err
}
