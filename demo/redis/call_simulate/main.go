package main

import (
    "github.com/scryinfo/dot/dot"
    dot_redis "github.com/scryinfo/dot/dots/db/redis"
    "github.com/scryinfo/dot/dots/line"
    "github.com/scryinfo/scryg/sutils/ssignal"
    "go.uber.org/zap"
    "os"
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
    lives := dot_redis.RedisTypeLives()
    err := l.PreAdd(lives...)

    return err
}
