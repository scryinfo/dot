package main

import (
	"fmt"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/line"
	"github.com/scryinfo/scryg/sutils/ssignal"
	"os"
)

func main() {
	l, err := line.BuildAndStart(add) //first step create line and dots
	if err != nil {
		fmt.Println(err)
		return
	}
	defer line.StopAndDestroy(l, true) //fourth step stop and destroy dots

	//second step ....

	ssignal.WatiCtrlC(func(s os.Signal) bool { //third wait for exit
		return false
	})
}

func add(l dot.Line) error {
	logger := dot.Logger()
	err := l.PreAdd(&dot.TypeLives{
		Meta: dot.Metadata{TypeId: "dotId",
			NewDoter: func(conf interface{}) (dot dot.Dot, err error) {
				d := 1
				return &d, nil
			},
		},
	})

	l.ToDotEventer().AddLiveEvents("dotId", &dot.LiveEvents{
		BeforeCreate: func(live *dot.Live, l dot.Line) {
			//do something
			logger.Infoln("BeforeCreate")
		},
		AfterCreate: func(live *dot.Live, l dot.Line) {
			logger.Infoln("AfterCreate")
			if d, ok := live.Dot.(*int); ok {
				logger.Infoln(fmt.Sprintf("dot: %d", *d))
			}
		},
		BeforeStart: func(live *dot.Live, l dot.Line) {
			logger.Infoln("BeforeStart")
		},
		AfterStart: func(live *dot.Live, l dot.Line) {
			logger.Infoln("AfterStart")
		},
		BeforeStop: func(live *dot.Live, l dot.Line) {
			logger.Infoln("BeforeStop")
		},
		AfterStop: func(live *dot.Live, l dot.Line) {
			logger.Infoln("AfterStop")
		},
		BeforeDestroy: func(live *dot.Live, l dot.Line) {
			logger.Infoln("BeforeDestroy")
		},
		AfterDestroy: func(live *dot.Live, l dot.Line) {
			logger.Infoln("AfterDestroy")
		},
	})

	return err
}
