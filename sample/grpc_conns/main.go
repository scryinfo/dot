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

	ssignal.WatiCtrlC(func(s os.Signal) bool { //third wait for exit
		return false
	})

	line.StopAndDestroy(l, true) //fourth step stop and destroy dots
}

func add(l dot.Line) error {
	var err error
	// 给typeid指定newer
	err = l.PreAdd(&dot.TypeLives{
		Meta: dot.Metadata{TypeId: conns.DotTypeId, NewDoter: func(conf interface{}) (dot dot.Dot, err error) {
			return conns.NewDailConnections(conf)
		}},
	})
	return err
}
