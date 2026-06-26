package dot

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/go-viper/mapstructure/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/scryinfo/dot/lib/kits"
	"github.com/spf13/viper"
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

	{
		var writer io.Writer = rotator
		if conf.AddStdout || (IsDebug && conf.AddStdoutInDebug) {
			consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
			writer = zerolog.MultiLevelWriter(consoleWriter, rotator)
		}
		Logger = zerolog.New(writer).With().Caller().Logger().Level(level)
		log.Logger = Logger
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
	MaxAge   int    `json:"max_age" toml:"max_age" yaml:"max_age" mapstructure:"max_age"`
	Compress bool   `json:"compress" toml:"compress" yaml:"compress" mapstructure:"compress"`
	Level    string `json:"level" toml:"level" yaml:"level" mapstructure:"level"`
	// when debug is true, add stdout to the log output
	AddStdoutInDebug bool `json:"add_stdout_in_debug" toml:"add_stdout_in_debug" yaml:"add_stdout_in_debug" mapstructure:"add_stdout_in_debug"`
	AddStdout        bool `json:"add_stdout" toml:"add_stdout" yaml:"add_stdout" mapstructure:"add_stdout"`
	// 是否把zerolog设置为系统日志
	SetSlog bool `json:"set_slog" toml:"set_slog" yaml:"set_slog" mapstructure:"set_slog"`
}

// log file name is log.toml
// read the log.toml in cmd parameter
// if no parameter, read the log.toml file from wd path, then exe path, then main source path
// if log.toml not found in main source path, read the log file name = "exe name"_log.toml
// read from wd path, then read from the exe path, then read from the main source path
func NewLogConfig() (*LogConfig, error) {
	logName := "log.toml"
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	logConfig, err := _readLogConfigFromPath(logName, wd)
	if err == nil {
		return logConfig, nil
	}
	exeFile, err := os.Executable()
	if err != nil {
		return nil, err
	}
	exePath := filepath.Dir(exeFile)
	logConfig, err = _readLogConfigFromPath(logName, exePath)
	if err == nil {
		return logConfig, nil
	}
	mainFile := kits.Config.GetMainPackageFile()
	mainPath := filepath.Dir(mainFile)
	logConfig, err = _readLogConfigFromPath(logName, mainPath)
	if err == nil {
		return logConfig, nil
	}

	{
		exeName := filepath.Base(exeFile)
		extName := filepath.Ext(exeName)
		exeName = exeName[:len(exeName)-len(extName)]
		logName = exeName + "_" + logName
		logConfig, err = _readLogConfigFromPath(logName, wd)
		if err == nil {
			return logConfig, nil
		}
		logConfig, err = _readLogConfigFromPath(logName, exePath)
		if err == nil {
			return logConfig, nil
		}
		logConfig, err = _readLogConfigFromPath(logName, mainPath)
		if err == nil {
			return logConfig, nil
		}
	}

	{
		mainName := filepath.Base(mainFile)
		extName := filepath.Ext(mainName)
		mainName = mainName[:len(mainName)-len(extName)]
		logName = mainName + "_" + logName
		logConfig, err = _readLogConfigFromPath(logName, wd)
		if err == nil {
			return logConfig, nil
		}
		logConfig, err = _readLogConfigFromPath(logName, exePath)
		if err == nil {
			return logConfig, nil
		}
		logConfig, err = _readLogConfigFromPath(logName, mainPath)
		if err == nil {
			return logConfig, nil
		}
	}
	return nil, fmt.Errorf("cant find log.toml(exename_log.toml, mainsource_log.toml) in gw/exe/mainsource")
}

func _readLogConfigFromPath(logName string, path string) (*LogConfig, error) {
	config, err := os.ReadFile(filepath.Join(path, logName))
	if err != nil {
		return nil, err
	}
	viperLog := viper.New()
	viperLog.SetConfigType("toml")
	if err := viperLog.ReadConfig(bytes.NewBuffer(config)); err != nil {
		return nil, err
	}
	var logConfig LogConfig
	err = viperLog.Unmarshal(&logConfig, func(dc *mapstructure.DecoderConfig) {
		dc.ErrorUnused = true
		dc.ErrorUnset = true
		dc.TagName = "toml"
	})
	if err != nil {
		return nil, err
	}
	return &logConfig, nil

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

// set the default logger for stdout
func init() {
	Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Caller().Logger().Level(zerolog.DebugLevel)
	log.Logger = Logger
	slog.SetDefault(slog.New(&FastZerologHandler{logger: Logger}))
}
