package main

import (
	"fmt"
	"github.com/scryInfo/dot/dots/line"
)

func main() {

	//if sv,ok := slog.sLogger.(slog.SLogger); ok {
	//	fmt.Printf("v implements String(): %s\n", sv)
	//}

	l, _ := line.BuildAndStart(nil)

	log := l.SLogger()

	log.Debug(func() string {
		return "6666666666"
	})

	log.Debugln("ssssssssssssssssssssss")

	fmt.Println(log.GetLevel())

	line.StopAndDestroy(l, true)

}
