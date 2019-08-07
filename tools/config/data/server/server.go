// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/grpc/gserver"
	"github.com/scryinfo/dot/dots/line"
	"github.com/scryinfo/dot/tools/config/data/nobl"
	"github.com/scryinfo/scryg/sutils/ssignal"
	"go.uber.org/zap"
	"os"
)

func main() {
	l, err := line.BuildAndStart(add) //first step create line and dots
	if err != nil {
		dot.Logger().Errorln("", zap.Error(err))
		return
	}
	defer line.StopAndDestroy(l, true) //fourth step stop and destroy dots

	dot.Logger().Infoln("dot ok")
	//second step ....

	ssignal.WaitCtrlC(func(s os.Signal) bool { //third wait for exit
		return false
	})

}

func add(l dot.Line) error {
	lives := nobl.HiServerTypeLives()
	lives = append(lives, gserver.HttpNoblTypeLives()...)
	return l.PreAdd(lives...)
}
