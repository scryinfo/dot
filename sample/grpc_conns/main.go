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
	dot.Logger().Infoln("dot ok")
	//second step ....

	dd, _ := l.ToInjecter().GetByLiveId(dot.LiveId(conns.ConnNameTypeId))
	fmt.Println(dd)

	ssignal.WatiCtrlC(func(s os.Signal) bool { //third wait for exit
		return false
	})

	line.StopAndDestroy(l, true) //fourth step stop and destroy dots
}

func add(l dot.Line) error {
	var err error
	// 给typeid指定newer
	err = l.PreAdd(conns.TypeLiveConns())
	err = l.PreAdd(conns.TypeLiveConnName())
	return err
}
