package slog

import (
	"github.com/scryInfo/dot/dot"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func (log *sLogger) SetLogFile(file string) {
	log.OutputPath = file
	log.Create(nil)
}

//NewConfiger new sConfig
func NewSLogger(lv int8, file string) *sLogger {
	return &sLogger{
		LogLevel:   lv,
		OutputPath: file,
	}
}

type sLogger struct {
	level      zap.AtomicLevel
	OutputPath string
	Logger     *zap.Logger
	LogLevel   int8
}

func (log *sLogger) GetLevel() int8 {
	return log.LogLevel
}

//SetLevel set level
func (log *sLogger) SetLevel(levels int8) {
	log.level.SetLevel(zapcore.Level(levels))
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

//func Add(l line.Line)  {
//	l.AddNewerByLiveId(LiveId("d8299d21-4f43-48bd-9a5c-654c4395ea17"), func(conf interface{}) (d Dot, err error) {
//		d = &sLogger{
//		}
//		err = nil
//		t := reflect.ValueOf(conf)
//		if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
//			if t.Len() > 0 && t.Index(0).Kind() == reflect.Uint8 {
//				v := t.Slice(0, t.Len())
//				json.Unmarshal(v.Bytes(), d)
//			}
//		} else {
//			err = SError.Parameter
//		}
//		return
//	})
//}

func (log *sLogger) Create(l dot.Line) (err error) {

	//log.LogLevel = -1

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

	atom.SetLevel(zapcore.Level(log.LogLevel))

	log.level = atom

	customCfg := zap.Config{
		Level:            log.level,
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    encoderCfg,
		OutputPaths:      []string{"stderr", log.OutputPath},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := customCfg.Build()

	if err != nil {

	} else {
		log.Logger = logger
	}

	return err
}

//启动连接
func (log *sLogger) Start(ignore bool) error {
	return nil
}

//Stop
//ignore 在调用其它Lifer时，true 出错出后继续，false 出现一个错误直接返回
func (log *sLogger) Stop(ignore bool) error {
	return nil
}

//Destroy 销毁 Dot
//ignore 在调用其它Lifer时，true 出错出后继续，false 出现一个错误直接返回
func (log *sLogger) Destroy(ignore bool) error {
	defer log.Logger.Sync()
	return nil
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
}
