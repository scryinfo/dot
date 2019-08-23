// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/grpc/gserver"
	"github.com/scryinfo/dot/dots/line"
	"github.com/scryinfo/dot/tools/config/data/nobl"
	"github.com/scryinfo/scryg/sutils/ssignal"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	{
		ch := make(chan byte, 2)
		go func() {
			gin.SetMode(gin.ReleaseMode)
			router := gin.Default()
			router.StaticFS("/", http.Dir("../../app/dist"))
			ch <- 1
			_ = router.Run(":9090")

		}()

		switch runtime.GOOS {
		case "windows":
			windowsBrowser(ch)
		case "linux":
			linuxBrowser(ch)
		default:
			log.Fatal("无法识别的操作系统")
		}
		<-ch
		<-ch
	}

	{
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
}

func add(l dot.Line) error {
	lives := nobl.HiServerTypeLives()
	lives = append(lives, gserver.HttpNoblTypeLives()...)
	return l.PreAdd(lives...)
}
func linuxBrowser(ch chan byte) {
	go func() {
		ch <- 1
		err := exec.Command("x-www-browser", "http://localhost:9090").Run()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "please click ->   http://localhost:9090\n")
		}
	}()
}
func windowsBrowser(ch chan byte) {
	go func() {
		ch <- 1
		err := exec.Command("cmd", "/C", "start", "http://localhost:9090").Run()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "please click ->   http://localhost:9090\n")
		}
	}()
}
