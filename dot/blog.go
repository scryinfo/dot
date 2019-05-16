package dot

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	blogName = "before.log"
)

//这个是为了解决在日志还没有被初化前， 日志输出的问题
//blog == before log
type blog struct {
	logger *zap.Logger
}

func (c *blog) Destroy(ignore bool) error {
	if c.logger != nil {
		c.logger.Sync()
		c.logger = nil
	}
	return nil
}

func (c *blog) GetLevel() Level {
	return zap.DebugLevel
}

func (c *blog) SetLevel(level Level) {

}

func (c *blog) Debugln(msg string, fields ...zap.Field) {
	c.logger.Debug(msg, fields...)
}

func (c *blog) Debug(mstr MakeStringer) {
	c.logger.Debug(mstr())
}

func (c *blog) Infoln(msg string, fields ...zap.Field) {
	c.logger.Debug(msg, fields...)
}

func (c *blog) Info(mstr MakeStringer) {
	c.logger.Debug(mstr())
}

func (c *blog) Warnln(msg string, fields ...zap.Field) {
	c.logger.Debug(msg, fields...)
}

func (c *blog) Warn(mstr MakeStringer) {
	c.logger.Debug(mstr())
}

func (c *blog) Errorln(msg string, fields ...zap.Field) {
	c.logger.Debug(msg, fields...)
}

func (c *blog) Error(mstr MakeStringer) {
	c.logger.Debug(mstr())
}

func (c *blog) Fatalln(msg string, fields ...zap.Field) {
	c.logger.Debug(msg, fields...)
}

func (c *blog) Fatal(mstr MakeStringer) {
	c.logger.Debug(mstr())
}

func newBlog() *blog {

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
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	atom := zap.NewAtomicLevel()

	atom.SetLevel(zap.DebugLevel)

	customCfg := zap.Config{
		Level:            atom,
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    encoderCfg,
		OutputPaths:      []string{"stderr", blogName},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := customCfg.Build()

	if err != nil {
		//todo 这里还没有日志， 无法输出来日志
		fmt.Println(err)
	}

	return &blog{logger: logger}
}
