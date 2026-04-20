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
func InitLogger(conf *LogConfig) *LoggerType {
	if len(conf.FileName) < 1 {
		conf.FileName = "logs/log.log"
	}
	rotator := &lumberjack.Logger{
		Filename:   conf.FileName,
		MaxSize:    conf.MaxSize,
		MaxBackups: conf.MaxBackups,
		MaxAge:     conf.MaxAge,
		Compress:   conf.Compress,
	}
	level, err := zerolog.ParseLevel(conf.Level)
	if err != nil {
		level = zerolog.DebugLevel
	}

	if IsDebug {
		var writer io.Writer = rotator
		if conf.AddStdOut {
			consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
			writer = zerolog.MultiLevelWriter(consoleWriter, rotator)
		}
		Logger = zerolog.New(writer).With().Caller().Logger().Level(level)

	} else {
		var writer io.Writer = rotator
		if conf.AddStdOut {
			consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
			writer = zerolog.MultiLevelWriter(consoleWriter, rotator)
		}
		Logger = zerolog.New(writer).With().Timestamp().Caller().Logger().Level(level)
	}
	log.Logger = Logger
	if conf.SetSlog {
		slog.SetDefault(slog.New(&FastZerologHandler{logger: Logger}))
	}
	return &Logger
}

type LogConfig struct {
	FileName   string `json:"fileName" toml:"fileName" yaml:"fileName"`
	MaxSize    int    `json:"maxSize" toml:"maxSize" yaml:"maxSize"`
	MaxBackups int    `json:"maxBackups" toml:"maxBackups" yaml:"maxBackups"`
	// days
	MaxAge    int    `json:"maxAge" toml:"maxAge" yaml:"maxAge"`
	Compress  bool   `json:"compress" toml:"compress" yaml:"compress"`
	Level     string `json:"level" toml:"level" yaml:"level"`
	AddStdOut bool   `json:"addStdOut" toml:"addStdOut" yaml:"addStdOut"`
	// 是否把zerolog设置为系统日志
	SetSlog bool `json:"setSlog" toml:"setSlog" yaml:"setSlog"`
}

func TestLogConfig() LogConfig {
	return LogConfig{
		FileName:   "logs/log.log",
		MaxSize:    10 << 20,
		MaxBackups: 2,
		MaxAge:     2, // days
		Compress:   false,
		Level:      "debug",
		AddStdOut:  true,
		SetSlog:    true,
	}
}
