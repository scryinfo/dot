package badgerdot

import (
	"path/filepath"

	"github.com/dgraph-io/badger/v4"
	"github.com/scryinfo/dot/dot"
)

type BadgerDbDot struct {
	db *badger.DB
}
type BadgerDbDotConfig struct {
	DbPath     string `json:"db_path" toml:"db_path" yaml:"db_path" mapstructure:"db_path"`
	BackupPath string `json:"backup_path" toml:"backup_path" yaml:"backup_path" mapstructure:"backup_path"`
	// info warning debug error
	Loglevel string `json:"log_level" toml:"log_level" yaml:"log_level" mapstructure:"log_level"`
}

func NewBadgerDot(config *BadgerDbDotConfig, sconfig dot.SConfig, logger *dot.LoggerType) (*BadgerDbDot, func(), error) {

	{
		dpPath, err := sconfig.FullPath(config.DbPath)
		if err != nil {
			config.DbPath = filepath.Join(sconfig.WdPath(), config.DbPath)
		} else {
			config.DbPath = dpPath
		}
	}
	logger.Info().Msgf("full badger db path: %s", config.DbPath)
	logLevel := badger.INFO
	switch config.Loglevel {
	case "debug":
		logLevel = badger.DEBUG
	case "warning":
		logLevel = badger.WARNING
	case "error":
		logLevel = badger.ERROR
	}
	logger.Info().Msgf("badger db path: %s", config.DbPath)
	dbBadger, err := badger.Open(badger.DefaultOptions(config.DbPath).WithLogger(&dblogger{Logger: logger}).WithLoggingLevel(logLevel))
	if err != nil {
		logger.Error().Err(err).Send()
		return nil, nil, err
	}
	return &BadgerDbDot{db: dbBadger}, func() {
		err := dbBadger.Close()
		if err != nil {
			logger.Error().Err(err).Send()
		}
	}, nil
}

func (p *BadgerDbDot) Db() *badger.DB {
	return p.db
}

type dblogger struct {
	Logger *dot.LoggerType
}

func (c *dblogger) Errorf(format string, v ...any) {
	c.Logger.Error().Msgf(format, v...)
}

// Infof logs an INFO message to the logger specified in opts.
func (c *dblogger) Infof(format string, v ...any) {
	c.Logger.Info().Msgf(format, v...)
}

// Warningf logs a WARNING message to the logger specified in opts.
func (c *dblogger) Warningf(format string, v ...any) {
	c.Logger.Warn().Msgf(format, v...)
}

// Debugf logs a DEBUG message to the logger specified in opts.
func (c *dblogger) Debugf(format string, v ...any) {
	c.Logger.Debug().Msgf(format, v...)
}
