// Scry Info.  All rights reserved.
// license that can be found in the license file.

package slog

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/scryinfo/dot/dot"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	_ dot.SLogger = (*sLogger)(nil)
)

const funLog = "Fun Log"

// NewSLogger new sConfig
func NewSLogger(conf *dot.LogConfig, l dot.Line) dot.SLogger {
	if conf == nil {
		conf = &dot.LogConfig{
			File: "log.log",
		}
	}

	if len(conf.Level) < 1 {
		conf.Level = zapcore.DebugLevel.String()
	}
	if len(conf.File) < 1 {
		conf.File = "log.log"
	}
	re := &sLogger{
		conf: *conf,
	}
	_ = re.Create(l)
	return re
}

type sLogger struct {
	level  zap.AtomicLevel
	Logger *zap.Logger
	conf   dot.LogConfig
}

func (log *sLogger) GetLevel() dot.Level {
	l := zap.InfoLevel
	_ = (&l).UnmarshalText([]byte(log.conf.Level))
	return l
}

// SetLevel set level
func (log *sLogger) SetLevel(levels dot.Level) {
	log.conf.Level = levels.String()
	log.level.SetLevel(levels)
}

// Debugln debug
func (log *sLogger) Debugln(msg string, fields ...zap.Field) {
	log.Logger.Debug(msg, fields...)
}

// Debug debug
func (log *sLogger) Debug(mstr dot.MakeStringer) {
	if dot.DebugLevel <= log.GetLevel() {
		log.Logger.Debug(funLog, zap.String("", mstr()))
	}
}

// Infoln info
func (log *sLogger) Infoln(msg string, fields ...zap.Field) {
	log.Logger.Info(msg, fields...)
}

// Info info
func (log *sLogger) Info(mstr dot.MakeStringer) {
	if dot.InfoLevel <= log.GetLevel() {
		log.Logger.Info(funLog, zap.String("", mstr()))
	}
}

// //Warnln warn
func (log *sLogger) Warnln(msg string, fields ...zap.Field) {
	log.Logger.Warn(msg, fields...)
}

// Warn warn
func (log *sLogger) Warn(mstr dot.MakeStringer) {
	if dot.WarnLevel <= log.GetLevel() {
		log.Logger.Warn(funLog, zap.String("", mstr()))
	}
}

// //Errorln error
func (log *sLogger) Errorln(msg string, fields ...zap.Field) {
	log.Logger.Error(msg, fields...)
}

// //Error error
func (log *sLogger) Error(mstr dot.MakeStringer) {

	if dot.ErrorLevel <= log.GetLevel() {
		log.Logger.Error(funLog, zap.String("", mstr()))
	}
}

// //Fatalln fatal
func (log *sLogger) Fatalln(msg string, fields ...zap.Field) {
	log.Logger.Fatal(msg, fields...)
}

// //Fatal fatal
func (log *sLogger) Fatal(mstr dot.MakeStringer) {
	if dot.FatalLevel <= log.GetLevel() {
		log.Logger.Fatal(funLog, zap.String("", mstr()))
	}
}

// NewLogger return new logger
func (log *sLogger) NewLogger(callerSkip int) dot.SLogger {
	n := &sLogger{
		conf:   log.conf,
		level:  log.level,
		Logger: log.Logger.WithOptions(zap.AddCallerSkip(callerSkip)),
	}
	return n
}

// Create
func (log *sLogger) Create(l dot.Line) (err error) {
	var maxAge = 0 //
	if log.conf.MaxAge > 0 {
		maxAge = log.conf.MaxAge
	}
	var maxSize = 10 // default : 10M
	if log.conf.MaxSize > 0 {
		maxSize = log.conf.MaxSize
	}
	var maxBackups = 0
	if log.conf.MaxBackups > 0 {
		maxBackups = log.conf.MaxBackups
	}
	hook := lumberjack.Logger{
		Filename:   filepath.Join(log.conf.DirPath, log.conf.File), // 日志文件路径
		MaxSize:    maxSize,                                        // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: maxBackups,                                     // 日志文件最多保存多少个备份
		MaxAge:     maxAge,                                         // 文件最多保存多少天
		Compress:   true,                                           // 是否压缩
	}

	encoderCfg := zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	logLevel := zap.NewAtomicLevel()
	logLevel.SetLevel(log.GetLevel())
	log.level = logLevel

	var core zapcore.Core
	if log.conf.IsOpenConsole { // 打印到控制台和文件
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),                                              // 编码器配置
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
			logLevel, // 日志级别
		)
	} else { // 打印到 文件
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),                  // 编码器配置
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(&hook)), // 打印到 文件
			logLevel, // 日志级别
		)
	}

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	logger := zap.New(core, caller, zap.AddCallerSkip(1), development)
	log.Logger = logger

	return err
}

// Destroy Destroy Dot
// ignore When calling other Lifer, if true erred then continue, if false erred then return directly
func (log *sLogger) Destroy(ignore bool) error {
	if log.Logger != nil {
		_ = log.Logger.Sync() //no log
		log.Logger = nil
	}
	return nil
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
}
