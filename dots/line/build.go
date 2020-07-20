// Scry Info.  All rights reserved.
// license that can be found in the license file.

package line

import (
	"encoding/json"
	"flag"
	"github.com/scryinfo/dot/dot"
	"go.uber.org/zap"
)

//  Construct line and call create rely createdots start
//  add will be called before reply, after create
func BuildAndStart(add dot.BuildNewer) (l dot.Line, err error) {
	err = nil
	builder := &dot.Builder{Add: add, LineLiveId: "default"}
	l, err = BuildAndStartBy(builder)
	return
}

//  Construct line and call create rely createdots start
func BuildAndStartBy(builder *dot.Builder) (l dot.Line, err error) {

	if !flag.Parsed() {
		flag.Parse()
	}

	err = nil

	if len(builder.LineLiveId) < 1 {
		builder.LineLiveId = "default"
	}

	line := newLine(builder)
	l = line
	if builder.BeforeCreate != nil {
		builder.BeforeCreate(line)
	}
	{
		err = line.Create(nil)

		if err != nil {
			return
		}

		if builder.Add != nil {
			err = builder.Add(line)
			if err != nil {
				return
			}
		}

		err = line.makeDotMetaFromConfig() // after the Add
		if err != nil {
			return
		}

		line.autoMakeLiveId() //issue #17

		line.makeRelays()

		dotOrder, circles := line.relyOrder() //do not care the error, it is circle dependency
		//circle dependency
		if len(circles) > 0 {
			bs, _ := json.Marshal(circles) //the %v just print the address of memory
			dot.Logger().Warnln("build", zap.String("", string(bs)))
		}

		err = line.CreateDots(dotOrder)
		if err != nil {
			return
		}
	}

	if builder.AfterCreate != nil {
		builder.AfterCreate(line)
	}

	dot.Logger().Infoln("dots Create")

	if builder.BeforeStart != nil {
		builder.BeforeStart(line)
	}
	err = line.Start(false)
	if builder.AfterStart != nil {
		builder.AfterStart(line)
	}

	if err != nil {
		return
	}

	dot.Logger().Infoln("dots Start")

	return
}

func StopAndDestroy(l dot.Line, ignore bool) {
	builder := l.GetLineBuilder()

	if builder.BeforeStop != nil {
		builder.BeforeStop(l)
	}
	_ = l.ToLifer().Stop(ignore)
	if builder.AfterStop != nil {
		builder.AfterStop(l)
	}
	dot.Logger().Infoln("dots Stop")

	if builder.BeforeDestroy != nil {
		builder.BeforeDestroy(l)
	}
	_ = l.ToLifer().Destroy(ignore)
	if builder.AfterDestroy != nil {
		builder.AfterDestroy(l)
	}
	dot.Logger().Infoln("dots Destroy")

	//dot.Logger().Infoln("dots destroy") maybe no logger
}
