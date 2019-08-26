// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"fmt"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/gindot"
	"github.com/scryinfo/dot/dots/grpc/gserver"
	"github.com/scryinfo/dot/dots/line"
	"github.com/scryinfo/dot/tools/config/data/nobl"
	"github.com/scryinfo/scryg/sutils/ssignal"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	l, err := line.BuildAndStartBy(&dot.Builder{
		Add: add,
		AfterStart: func(l dot.Line) {
			go func() {
				switch runtime.GOOS {
				case "windows":
					windowsBrowser()
				case "linux":
					linuxBrowser()
				default:
					dot.Logger().Fatalln("无法识别的操作系统")
				}
			}()
		},
	})
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
	lives := nobl.RpcImplementTypeLives()
	lives = append(lives, gserver.GinNoblTypeLives()...)
	l.ToDotEventer().AddLiveEvents(dot.LiveId(gindot.EngineLiveId), &dot.LiveEvents{
		AfterCreate: func(live *dot.Live, l dot.Line) {
			if g, ok := live.Dot.(*gindot.Engine); ok {
				g.GinEngine().StaticFS("/", http.Dir("../../app/dist"))
			}
		},
	})

	//4943e959-7ad7-42c6-84dd-8b24e9ed30bb

	return l.PreAdd(lives...)
}
func linuxBrowser() {
	err := exec.Command("x-www-browser", "http://localhost:8080").Run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "please click ->   http://localhost:8080\n")
	}
}
func windowsBrowser() {
	err := exec.Command("cmd", "/C", "start", "http://localhost:8080").Run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "please click ->   http://localhost:8080\n")
	}
}
