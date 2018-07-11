package slog

import "github.com/scryinfo/dot"

type Level int

const (
	Debug Level = 1
	Info        = 2
	Warn        = 3
	Error       = 4
	Fatal       = 5
)

//MakeStringer 生成日志字符串
type MakeStringer = func() string

//SLogger 日志属于一个组件Dot,但它太基础了，大部分 Dot都需要，所以定义到 dot.go文件中
//所有日志调用时，不要在参数中调用函数，函数的运行是先于日志等级的，如果一定要调用函数，可以使用回调函数的方式（注回调函数一定要正常运行）
//S表示 scryinfo, log这个名字用的地方太多，加一个s以示区别
type SLogger interface {
	dot.Lifer

	//GetLevel get level
	GetLevel() Level
	//SetLevel set level
	SetLevel(level Level)

	//Debugf debug
	Debugf(format string, args ...interface{})
	//Debugln debug
	Debugln(args ...interface{})
	//Debug debug
	Debug(mstr MakeStringer)

	//Infof info
	Infof(format string, args ...interface{})
	//Infoln info
	Infoln(args ...interface{})
	//Info info
	Info(mstr MakeStringer)

	//Warnf warn
	Warnf(format string, args ...interface{})
	//Warnln warn
	Warnln(args ...interface{})
	//Warn warn
	Warn(mstr MakeStringer)

	//Errorf error
	Errorf(format string, args ...interface{})
	//Errorln error
	Errorln(args ...interface{})
	//Error error
	Error(mstr MakeStringer)

	//Fatalf fatal
	Fatalf(format string, args ...interface{})
	//Fatalln fatal
	Fatalln(args ...interface{})
	//Fatal fatal
	Fatal(mstr MakeStringer)
}
