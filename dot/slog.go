package dot

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
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

func (log *ULog) SetLogFile(file string) {
	log.OutputPath = file
	log.Create(nil)
}

//NewConfiger new sConfig
func NewLoger(lv int8, file string) *ULog {
	return &ULog{
		LogLevel:   lv,
		OutputPath: file,
	}
}

type ULog struct {
	level      zap.AtomicLevel
	OutputPath string
	Logger     *zap.Logger
	LogLevel   int8
}

type DotLog struct {
	Log SLogger `dot:"d8299d21-4f43-48bd-9a5c-654c4395ea17"`
}

func (log *ULog) GetLevel() int8 {
	return log.LogLevel
}

//SetLevel set level
func (log *ULog) SetLevel(levels int8) {
	log.level.SetLevel(zapcore.Level(levels))
}

//Debugln debug
func (log *ULog) Debugln(msg string, fields ...zap.Field) {
	log.Logger.Debug(msg, fields...)
}

//Debug debug
func (log *ULog) Debug(mstr MakeStringer) {
	if ce := log.Logger.Check(DebugLevel, mstr()); ce != nil {
		log.Logger.Debug(mstr())
	}
}

//Infoln info
func (log *ULog) Infoln(msg string, fields ...zap.Field) {
	log.Logger.Info(msg, fields...)
}

//Info info
func (log *ULog) Info(mstr MakeStringer) {
	if ce := log.Logger.Check(DebugLevel, mstr()); ce != nil {
		log.Logger.Info(mstr())
	}
}

////Warnln warn
func (log *ULog) Warnln(msg string, fields ...zap.Field) {
	log.Logger.Warn(msg, fields...)
}

//Warn warn
func (log *ULog) Warn(mstr MakeStringer) {
	if ce := log.Logger.Check(DebugLevel, mstr()); ce != nil {
		log.Logger.Warn(mstr())
	}
}

////Errorln error
func (log *ULog) Errorln(msg string, fields ...zap.Field) {
	log.Logger.Error(msg, fields...)
}

////Error error
func (log *ULog) Error(mstr MakeStringer) {
	if ce := log.Logger.Check(DebugLevel, mstr()); ce != nil {
		log.Logger.Error(mstr())
	}
}

////Fatalln fatal
func (log *ULog) Fatalln(msg string, fields ...zap.Field) {
	log.Logger.Fatal(msg, fields...)
}

////Fatal fatal
func (log *ULog) Fatal(mstr MakeStringer) {
	if ce := log.Logger.Check(DebugLevel, mstr()); ce != nil {
		log.Logger.Fatal(mstr())
	}
}

//func Add(l line.Line)  {
//	l.AddNewerByLiveId(LiveId("d8299d21-4f43-48bd-9a5c-654c4395ea17"), func(conf interface{}) (d Dot, err error) {
//		d = &ULog{
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

func (log *ULog) Create(conf SConfig) (err error) {

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
func (log *ULog) Start(ignore bool) error {
	return nil
}

//Stop
//ignore 在调用其它Lifer时，true 出错出后继续，false 出现一个错误直接返回
func (log *ULog) Stop(ignore bool) error {
	return nil
}

//Destroy 销毁 Dot
//ignore 在调用其它Lifer时，true 出错出后继续，false 出现一个错误直接返回
func (log *ULog) Destroy(ignore bool) error {
	defer log.Logger.Sync()
	return nil
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
}
