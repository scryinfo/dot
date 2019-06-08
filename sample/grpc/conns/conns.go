// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"fmt"
	"os"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/grpc/conns"
	"github.com/scryinfo/dot/dots/line"
	"github.com/scryinfo/scryg/sutils/ssignal"
)

func main() {
	l, err := line.BuildAndStart(add) //first step create line and dots
	if err != nil {
		fmt.Println(err)
		return
	}
	defer line.StopAndDestroy(l, true) //fourth step stop and destroy dots

	dot.Logger().Infoln("dot ok")
	//second step ....

	dd, _ := l.ToInjecter().GetByLiveId(dot.LiveId(conns.ConnNameTypeId))
	fmt.Println(dd)

	ssignal.WaitCtrlC(func(s os.Signal) bool { //third wait for exit
		return false
	})

}

func add(l dot.Line) error {
	var err error
	// Point newer for typeid
	//err = l.PreAdd(conns.TypeLiveConns())
	err = l.PreAdd(conns.TypeLiveConnName()...)
	return err
}
