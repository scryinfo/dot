// Scry Info.  All rights reserved.
// license that can be found in the license file.

// 文件定时拆分处理
package demo

import (
	"fmt"
	"sync"
	"time"

	"github.com/scryinfo/dot/dot"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	_ dot.SLogger = (*sLogger)(nil)
)

func (log *sLogger) updateLogFile() {
	timer := time.NewTimer(log.interval)
	for {

		select {
		case <-log.cancel:
			//退出清理
			timer.Stop() //尽快清理timer
			return
		case <-timer.C:
			//do something
			m := sync.Mutex{}
			m.Lock()
			{
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

				atom := zap.NewAtomicLevel()

				atom.SetLevel(log.GetLevel())

				log.level = atom

				customCfg := zap.Config{
					Level:            log.level,
					Development:      true,
					Encoding:         "console",
					EncoderConfig:    encoderCfg,
					OutputPaths:      []string{"stderr", time.Now().Format(log.formatPrefix) + log.conf.File},
					ErrorOutputPaths: []string{"stderr"},
				}

				logger, err := customCfg.Build(zap.AddCallerSkip(1))

				if err == nil {
					log.Logger = logger
				} else {
					fmt.Print(err)
				}
			}
			m.Unlock()
		}
		timer.Reset(log.interval)
	}

}

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
	switch conf.Time {
	case "h", "H", "hour", "HOUR":
		re.interval = time.Hour
		re.formatPrefix = "2006-01-02 15 "
	case "d", "D", "day", "DAY":
		re.interval = time.Hour * 24
		re.formatPrefix = "2006-01-02 "
	case "m", "M", "month", "MONTH":
		re.interval = time.Hour * 24 * 31
		re.formatPrefix = "2006-01 "
	case "y", "Y", "year", "YEAR":
		re.interval = time.Hour * 24 * 31 * 12
		re.formatPrefix = "2006 "
	default:
		re.interval = time.Hour * 24 * 31 * 12
		re.formatPrefix = "2006 "
	}
	_ = re.Create(l)
	return re
}

type sLogger struct {
	cancel       chan bool
	level        zap.AtomicLevel
	Logger       *zap.Logger
	conf         dot.LogConfig
	interval     time.Duration
	formatPrefix string
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
	//top
	//n := &sLogger{
	//	conf:   log.conf,
	//	level:  log.level,
	//	Logger: log.Logger.WithOptions(zap.AddCallerSkip(callerSkip)),
	//}
	//return n

	//bottom
	log.Logger = log.Logger.WithOptions(zap.AddCallerSkip(callerSkip))
	return log
	/*
			top result:
			INFO	sconfig/sconfig.go:23	dot ok
			bottom result:
			INFO	runtime/proc.go:225	dot ok
		Other logs are printed the same,but in the top the output log is in the same file ,
		because n is not log,so log.updateLogFile()!=n.updateLogFile()
	*/
}

// Create
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
		EncodeTime:     timeEncoder,
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
		OutputPaths:      []string{"stderr", time.Now().Format(log.formatPrefix) + log.conf.File},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := customCfg.Build(zap.AddCallerSkip(1))

	if err == nil {
		log.Logger = logger
	}
	go log.updateLogFile()
	return err
}

// Destroy Destroy Dot
// ignore When calling other Lifer, if true erred then continue, if false erred then return directly
func (log *sLogger) Destroy(ignore bool) error {
	if log.Logger != nil {
		_ = log.Logger.Sync() //no log
		log.Logger = nil
		close(log.cancel)
	}
	return nil
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
}
