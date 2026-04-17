package dot

import (
	"io"
	"log/slog"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerType = zerolog.Logger

var Logger LoggerType

// init log with zerolog and lumberjack
func InitLogger(conf LogConfig) *LoggerType {
	rotator := &lumberjack.Logger{
		Filename:   conf.FileName,
		MaxSize:    conf.MaxSize,
		MaxBackups: conf.MaxBackups,
		MaxAge:     conf.MaxAge,
		Compress:   conf.Compress,
	}

	if IsDebug {
		var writer io.Writer = rotator
		if conf.AddStdOut {
			consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
			writer = zerolog.MultiLevelWriter(consoleWriter, rotator)
		}
		Logger = zerolog.New(writer).With().Caller().Logger().Level(conf.Level)

	} else {
		var writer io.Writer = rotator
		if conf.AddStdOut {
			consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
			writer = zerolog.MultiLevelWriter(consoleWriter, rotator)
		}
		Logger = zerolog.New(writer).With().Timestamp().Caller().Logger().Level(conf.Level)
	}
	log.Logger = Logger
	if conf.SetSlog {
		slog.SetDefault(slog.New(&FastZerologHandler{logger: Logger}))
	}
	return &Logger
}

type LogConfig struct {
	FileName   string        `json:"fileName" toml:"fileName" yaml:"fileName"`
	MaxSize    int           `json:"maxSize" toml:"maxSize" yaml:"maxSize"`
	MaxBackups int           `json:"maxBackups" toml:"maxBackups" yaml:"maxBackups"`
	MaxAge     int           `json:"maxAge" toml:"maxAge" yaml:"maxAge"`
	Compress   bool          `json:"compress" toml:"compress" yaml:"compress"`
	Level      zerolog.Level `json:"level" toml:"level" yaml:"level"`
	AddStdOut  bool          `json:"addStdOut" toml:"addStdOut" yaml:"addStdOut"`
	// 是否把zerolog设置为系统日志
	SetSlog bool `json:"setSlog" toml:"setSlog" yaml:"setSlog"`
}
