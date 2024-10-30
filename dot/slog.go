// Scry Info.  All rights reserved.
// license that can be found in the license file.

package dot

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Level level of log
type Level = zapcore.Level

const (
	//LogLiveID log dot live id
	LogLiveID = "d8299d21-4f43-48bd-9a5c-654c4395ea17"

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
	//FatalLevel DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	FatalLevel = zapcore.FatalLevel
)

// Default log, if not created, then returned log will be output to “before.log” file，all log will be returned
var logger SLogger = nil

// Logger Return default log, this API is used to call log easily
// This method does not consider thread security, Adjusting value is not suggested after program initialization
// Note: Default log, if log is not created, then returned log will be output to control panel, all log will be output
func Logger() SLogger {
	return logger
}

// SetLogger Set default log,
// This method does not consider thread security, Adjusting value is not suggested after program initialization
func SetLogger(log SLogger) {
	if logger != nil {
		if d, ok := logger.(Destroyer); ok {
			_ = d.Destroy(true)
		}
		logger = nil
	}
	logger = log
}

// MakeStringer Generate log string
type MakeStringer func() string

// SLogger log belongs to one component Dot, but it is too basic, most Dot need it, so defined it to dot.go file
// All log calling should not call function in parameters, function run priorly than log, if must call function, you should use callback(must run normally)
// S represents scry info, log name used frequently so add s to distinguish it
type SLogger interface {
	//GetLevel get level
	GetLevel() Level
	//SetLevel set level
	SetLevel(level Level)
	//Debugln debug
	Debugln(msg string, fields ...zap.Field)
	//Debug debug
	Debug(mstr MakeStringer)
	//Infoln info
	Infoln(msg string, fields ...zap.Field)
	//Info info
	Info(mstr MakeStringer)
	//Warnln warn
	Warnln(msg string, fields ...zap.Field)
	//Warn warn
	Warn(mstr MakeStringer)
	//Errorln error
	Errorln(msg string, fields ...zap.Field)
	//Error error
	Error(mstr MakeStringer)
	//Fatalln fatal
	Fatalln(msg string, fields ...zap.Field)
	//Fatal fatal
	Fatal(mstr MakeStringer)
	//NewLogger return new logger
	NewLogger(callerSkip int) SLogger
}

// LogConfig log config
type LogConfig struct {
	DirPath       string `json:"dirPath" toml:"dirPath" yaml:"dirPath"` // 文件夹路径
	File          string `json:"file" toml:"file" yaml:"file"`          // 日志文件名
	Level         string `json:"level" toml:"level" yaml:"level"`
	Time          string `json:"time" toml:"time" yaml:"time"`                            //hour,h,day,d,month,m,y,year
	MaxSize       int    `json:"maxSize" toml:"maxSize" yaml:"maxSize"`                   // 日志单文件,大小上限     最少：120M
	MaxBackups    int    `json:"maxBackups" toml:"maxBackups" yaml:"maxBackups"`          // 日志文件，最多有多少个   default：10
	MaxAge        int    `json:"maxAge" toml:"maxAge" yaml:"maxAge"`                      // 日志文件保存时长        最少：30days
	IsOpenConsole bool   `json:"isOpenConsole" toml:"isOpenConsole" yaml:"isOpenConsole"` // 是否输出到控制台
}

// Initialize one default log, let program use log at first, output to “before.log” file, all log will be output
func init() {
	//SetLogger(newBlog())
}
