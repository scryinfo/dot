// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"context"
	"fmt"
	"github.com/scryinfo/dot/sample/grpc/go_out/hidot"
	"github.com/scryinfo/dot/sample/grpc/nobl"
	"go.uber.org/zap"
	"os"
	"time"

	"github.com/scryinfo/dot/dot"
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

	go func() {
		dd, _ := l.ToInjecter().GetByLiveId(dot.LiveId(nobl.HiClientTypeId))
		if hiclient, ok := dd.(*nobl.HiClient); ok {
			logger := dot.Logger()
			for i := 0; i < 1000; i++ {
				name := fmt.Sprintf("client: %d", i)
				req, err := hiclient.HiClient().Hi(context.Background(), &hidot.HiReq{Name: name})
				if err != nil {
					logger.Infoln(err.Error())
				} else {
					logger.Infoln(fmt.Sprintf("%s %s", name, req.Name))
				}
				time.Sleep(100 * time.Millisecond)
			}
			logger.Infoln("done")
		}
	}()

	ssignal.WaitCtrlC(func(s os.Signal) bool { //third wait for exit
		return false
	})

}

func add(l dot.Line) error {
	return l.PreAdd(nobl.HiClientTypeLives()...)
}
