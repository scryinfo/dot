// Scry Info.  All rights reserved.
// license that can be found in the license file.

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

//Default log, if not created, then returned log will be output to “before.log” file，all log will be returned
var logger SLogger = nil

//Return default log, this API is used to call log easily
//This method does not consider thread security, Adjusting value is not suggested after program initialization
//Note: Default log, if log is not created, then returned log will be output to control panel, all log will be output
func Logger() SLogger {
	return logger
}

//Set default log,
//This method does not consider thread security, Adjusting value is not suggested after program initialization
func SetLogger(log SLogger) {
	if logger != nil {
		if d, ok := logger.(Destroyer); ok {
			_ = d.Destroy(true)
		}
		logger = nil
	}
	logger = log
}

//MakeStringer Generate log string
type MakeStringer func() string

//SLogger log belongs to one component Dot, but it is too basic, most Dot need it, so defined it to dot.go file
//All log calling should not call function in parameters, function run priorly than log, if must call function, you should use callback(must run normally)
//S represents scryinfo, log name used frequently so add s to distinguish it
type SLogger interface {
	//GetLevel get level
	GetLevel() Level
	//SetLevel set level
	SetLevel(level Level)

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
	//NewLogger return new logger
	NewLogger(callerSkip int) SLogger
}

type LogConfig struct {
	File  string `json:"file"`
	Level string `json:"level"`
}

//Initialize one default log, let program use log at first, output to “before.log” file, all log will be output
func init() {
	SetLogger(newBlog())
}
