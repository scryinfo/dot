package main

import (
	"os"

	"go.uber.org/zap"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/line"
	"github.com/scryinfo/dot/sample/db/redis_client/call"
	"github.com/scryinfo/scryg/sutils/ssignal"
)

func main() {
	l, err := line.BuildAndStart(add)
	if err != nil {
		dot.Logger().Errorln("", zap.Error(err))
	}
	defer line.StopAndDestroy(l, true)
	dot.Logger().Infoln("dot ok")

	ssignal.WaitCtrlC(func(s os.Signal) bool {
		return false
	})
}
func add(l dot.Line) error {
	lives := call.RedisDemoTypeLives()
	err := l.PreAdd(lives...)

	return err
}
