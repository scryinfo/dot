// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"go.uber.org/zap"
	"os"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/grpc/conns"
	"github.com/scryinfo/dot/dots/line"
	"github.com/scryinfo/scryg/sutils/ssignal"
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

	dd, _ := l.ToInjecter().GetByLiveID(dot.LiveID(conns.ConnNameTypeID))
	if dd != nil {

	}

	ssignal.WaitCtrlC(func(s os.Signal) bool { //third wait for exit
		return false
	})

}

func add(l dot.Line) error {
	var err error
	err = l.PreAdd(conns.ConnNameTypeLives()...)
	return err
}
