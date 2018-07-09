package dots

import "github.com/scryinfo/dot"

type SLogLevel int

const (
	DebugLog SLogLevel = 1
	InfoLog            = 2
	ErrorLog           = 3
)

//MakeLog 生成日志字符串
type MakeLog interface {
	MakeLog() string
}

//SLog 日志属于一个组件Dot,但它太基础了，大部分 Dot都需要，所以定义到 dot.go文件中
//所有日志调用时，不要在参数中调用函数，函数的运行是先于日志等级的，如果一定要调用函数，可以使用回调函数的方式（注回调函数一定要正常运行）
//S表示 scryinfo, log这个名字用的地方太多，加一个s以示区别
type SLoger interface {
	dot.Lifer

	Level() SLogLevel
	SetLevel(level SLogLevel)

	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})
	Debug(mstr MakeLog)

	Infof(format string, args ...interface{})
	Infoln(args ...interface{})
	Info(mstr MakeLog)

	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})
	Error(mstr MakeLog)
}
