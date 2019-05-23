// Scry Info.  All rights reserved.
// license that can be found in the license file.

package slog

import (
	"time"

	"github.com/scryinfo/dot/dot"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	_ dot.SLogger = (*sLogger)(nil)
)

//NewConfiger new sConfig
func NewSLogger(conf *dot.LogConfig, l dot.Line) *sLogger {
	if conf == nil {
		conf = &dot.LogConfig{
			File:  "log.log",
			Level: zapcore.InfoLevel.String(),
		}
	}

	if len(conf.Level) < 1 {
		conf.Level = zapcore.InfoLevel.String()
	}
	if len(conf.File) < 1 {
		conf.File = "log.log"
	}
	re := &sLogger{
		conf: *conf,
	}
	re.Create(l)
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

//SetLevel set level
func (log *sLogger) SetLevel(levels dot.Level) {
	log.conf.Level = levels.String()
	log.level.SetLevel(levels)
}

//Debugln debug
func (log *sLogger) Debugln(msg string, fields ...zap.Field) {
	log.Logger.Debug(msg, fields...)
}

//Debug debug
func (log *sLogger) Debug(mstr dot.MakeStringer) {
	if ce := log.Logger.Check(dot.DebugLevel, mstr()); ce != nil {
		log.Logger.Debug(mstr())
	}
}

//Infoln info
func (log *sLogger) Infoln(msg string, fields ...zap.Field) {
	log.Logger.Info(msg, fields...)
}

//Info info
func (log *sLogger) Info(mstr dot.MakeStringer) {
	if ce := log.Logger.Check(dot.DebugLevel, mstr()); ce != nil {
		log.Logger.Info(mstr())
	}
}

////Warnln warn
func (log *sLogger) Warnln(msg string, fields ...zap.Field) {
	log.Logger.Warn(msg, fields...)
}

//Warn warn
func (log *sLogger) Warn(mstr dot.MakeStringer) {
	if ce := log.Logger.Check(dot.DebugLevel, mstr()); ce != nil {
		log.Logger.Warn(mstr())
	}
}

////Errorln error
func (log *sLogger) Errorln(msg string, fields ...zap.Field) {
	log.Logger.Error(msg, fields...)
}

////Error error
func (log *sLogger) Error(mstr dot.MakeStringer) {
	if ce := log.Logger.Check(dot.DebugLevel, mstr()); ce != nil {
		log.Logger.Error(mstr())
	}
}

////Fatalln fatal
func (log *sLogger) Fatalln(msg string, fields ...zap.Field) {
	log.Logger.Fatal(msg, fields...)
}

////Fatal fatal
func (log *sLogger) Fatal(mstr dot.MakeStringer) {
	if ce := log.Logger.Check(dot.DebugLevel, mstr()); ce != nil {
		log.Logger.Fatal(mstr())
	}
}

func (log *sLogger) Create(l dot.Line) (err error) {

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
		EncodeTime:     TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	atom := zap.NewAtomicLevel()

	atom.SetLevel(log.GetLevel())

	log.level = atom

	customCfg := zap.Config{
		Level:            log.level,
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    encoderCfg,
		OutputPaths:      []string{"stderr", log.conf.File},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := customCfg.Build()

	if err != nil {

	} else {
		log.Logger = logger
	}

	return err
}

////start
//func (log *sLogger) Start(ignore bool) error {
//	return nil
//}
//
////Stop
////ignore When calling other Lifer, if true erred then continue, if false erred then return directly
//func (log *sLogger) Stop(ignore bool) error {
//	return nil
//}

//Destroy Destroy Dot
//ignore When calling other Lifer, if true erred then continue, if false erred then return directly
func (log *sLogger) Destroy(ignore bool) error {
	if log.Logger != nil {
		log.Logger.Sync()
		log.Logger = nil
	}
	return nil
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
}
