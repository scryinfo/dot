package dot

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level = zapcore.Level

const (
	LogLiveId = "d8299d21-4f43-48bd-9a5c-654c4395ea17"

	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = zapcore.DebugLevel
	// InfoLevel is the default logging priority.
	InfoLevel = zapcore.InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = zapcore.WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = zapcore.ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	FatalLevel = zapcore.FatalLevel
)

//MakeStringer 生成日志字符串
type MakeStringer func() string

//SLogger 日志属于一个组件Dot,但它太基础了，大部分 Dot都需要，所以定义到 dot.go文件中
//所有日志调用时，不要在参数中调用函数，函数的运行是先于日志等级的，如果一定要调用函数，可以使用回调函数的方式（注回调函数一定要正常运行）
//S表示 scryInfo, log这个名字用的地方太多，加一个s以示区别
type SLogger interface {
	Lifer

	//GetLevel get level
	GetLevel() int8
	//SetLevel set level
	SetLevel(level int8)

	//Debugf debug
	//Debugf(format string, args ...interface{})
	//Debugln debug
	Debugln(msg string, fields ...zap.Field)
	//Debug debug
	Debug(mstr MakeStringer)

	//Infof info
	//Infof(format string, args ...interface{})
	//Infoln info
	Infoln(msg string, fields ...zap.Field)
	//Info info
	Info(mstr MakeStringer)

	//Warnf warn
	//Warnf(format string, args ...interface{})
	//Warnln warn
	Warnln(msg string, fields ...zap.Field)
	//Warn warn
	Warn(mstr MakeStringer)

	//Errorf error
	//Errorf(format string, args ...interface{})
	//Errorln error
	Errorln(msg string, fields ...zap.Field)
	//Error error
	Error(mstr MakeStringer)

	//Fatalf fatal
	//Fatalf(format string, args ...interface{})
	//Fatalln fatal
	Fatalln(msg string, fields ...zap.Field)
	//Fatal fatal
	Fatal(mstr MakeStringer)

	SetLogFile(fiel string)
}
