// Scry Info.  All rights reserved.
// license that can be found in the license file.

package line

import (
	"github.com/scryinfo/dot/dot"
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

	err = nil

	if len(builder.LineLiveId) < 1 {
		builder.LineLiveId = "default"
	}

	l = NewLine(builder)

	if builder.BeforeCreate != nil {
		builder.BeforeCreate(l)
	}
	{
		err = l.ToLifer().Create(nil)

		if err != nil {
			return
		}

		if builder.Add != nil {
			err = builder.Add(l)
			if err != nil {
				return
			}
		}

		err = l.Rely()
		if err != nil {
			return
		}

		err = l.CreateDots()
	}
	if err != nil {
		return
	}
	if builder.AfterCreate != nil {
		builder.AfterCreate(l)
	}

	dot.Logger().Infoln("dots create")

	if builder.BeforeStart != nil {
		builder.BeforeStart(l)
	}
	err = l.ToLifer().Start(false)
	if builder.AfterStart != nil {
		builder.AfterStart(l)
	}

	if err != nil {
		return
	}

	dot.Logger().Infoln("dots start")

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
	dot.Logger().Infoln("dots stop")

	if builder.BeforeDestroy != nil {
		builder.BeforeDestroy(l)
	}
	_ = l.ToLifer().Destroy(ignore)
	if builder.AfterDestroy != nil {
		builder.AfterDestroy(l)
	}

	dot.Logger().Infoln("dots destroy")
}
