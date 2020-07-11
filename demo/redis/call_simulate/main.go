package main

import (
    "github.com/scryinfo/dot/dot"
    "github.com/scryinfo/dot/dots/grpc/gserver"
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
    lives := gserver.GinNoblTypeLives()
    //lives = append(lives, redis.)
    err := l.PreAdd(lives...)

    return err
}
