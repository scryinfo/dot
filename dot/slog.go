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
func NewLogger(conf *LogConfig) *LoggerType {
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
		if conf.AddStdout {
			consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
			writer = zerolog.MultiLevelWriter(consoleWriter, rotator)
		}
		Logger = zerolog.New(writer).With().Caller().Logger().Level(level)

	} else {
		var writer io.Writer = rotator
		if conf.AddStdout {
			consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
			writer = zerolog.MultiLevelWriter(consoleWriter, rotator)
		}
		Logger = zerolog.New(writer).With().Timestamp().Caller().Logger().Level(level)
	}
	log.Logger = Logger
	if conf.SetSlog {
		slog.SetDefault(slog.New(&FastZerologHandler{logger: Logger}))
	}
	Logger.Info().Msgf("log created")
	return &Logger
}

type LogConfig struct {
	FileName   string `json:"file_name" toml:"file_name" yaml:"file_name" mapstructure:"file_name"`
	MaxSize    int    `json:"max_size" toml:"max_size" yaml:"max_size" mapstructure:"max_size"`
	MaxBackups int    `json:"max_backups" toml:"max_backups" yaml:"max_backups" mapstructure:"max_backups"`
	// days
	MaxAge    int    `json:"max_age" toml:"max_age" yaml:"max_age" mapstructure:"max_age"`
	Compress  bool   `json:"compress" toml:"compress" yaml:"compress" mapstructure:"compress"`
	Level     string `json:"level" toml:"level" yaml:"level" mapstructure:"level"`
	AddStdout bool   `json:"add_stdout" toml:"add_stdout" yaml:"add_stdout" mapstructure:"add_stdout"`
	// 是否把zerolog设置为系统日志
	SetSlog bool `json:"set_slog" toml:"set_slog" yaml:"set_slog" mapstructure:"set_slog"`
}

func TestLogConfig() LogConfig {
	return LogConfig{
		FileName:   "logs/log.log",
		MaxSize:    10 << 20,
		MaxBackups: 2,
		MaxAge:     2, // days
		Compress:   false,
		Level:      "debug",
		AddStdout:  true,
		SetSlog:    true,
	}
}

func NewTestLogger() *zerolog.Logger {
	conf := TestLogConfig()
	return NewLogger(&conf)
}
