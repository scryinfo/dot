package dots

import (
	"github.com/scryInfo/dot/dot"
)

//  构造line并调用 create rely createdots start
//  add会在 create之后rely之前调用
func BuildAndStart(add dot.BuildNewer) (l dot.Line, err error)  {
	err = nil
	builder := &dot.Builder{Add: add}
	l,err = BuildAndStartBy(builder)
	return
}

//  构造line并调用 create rely createdots start
func BuildAndStartBy(builder *dot.Builder) (l dot.Line, err error) {

	err = nil

	l = New(builder)

	if builder.BeforeCeate != nil {
		builder.BeforeCeate(l)
	}
	err = l.ToLifer().Create(nil)
	if builder.AfterCreate != nil {
		builder.AfterCreate(l)
	}
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

	if builder.BeforeCreateDots != nil {
		builder.BeforeCreateDots(l)
	}
	err = l.CreateDots()
	if builder.AfterCreateDots != nil {
		builder.AfterCreateDots(l)
	}
	if err != nil {
		return
	}


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

	return
}

func StopAndDestroy(l dot.Line, ignore bool)  {
	builder := l.GetLineBuilder()

	if builder.BeforeStop != nil {
		builder.BeforeStop(l)
	}
	_ = l.ToLifer().Stop(ignore)
	if builder.AfterStop != nil {
		builder.AfterStop(l)
	}

	if builder.BeforeDestroy != nil {
		builder.BeforeDestroy(l)
	}
	_ = l.ToLifer().Destroy(ignore)
	if builder.AfterDestroy != nil {
		builder.AfterDestroy(l)
	}
}
