package baderdot

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/scryinfo/dot/dot"
)

type BaderDbDot struct {
	db *badger.DB
}
type BaderDbDotConfig struct {
	DbPath     string `json:"dbPath" toml:"dbPath" yaml:"dbPath" `
	BackupPath string `json:"backupPath" toml:"backupPath" yaml:"backupPath" `
	// info warning debug error
	Loglevel string `json:"logLevel" toml:"logLevel" yaml:"logLevel" `
}

func NewBaderDot(config *BaderDbDotConfig, sconfig dot.SConfig, logger *dot.LoggerType) (*BaderDbDot, error) {
	var err error
	var dbPath string
	dbPath, err = sconfig.FullPath(config.DbPath)
	if err != nil {
		logger.Error().Err(err).Send()
		return nil, err
	}
	logger.Info().Msgf("full bader db path: %s", dbPath)
	logLevel := badger.INFO
	switch config.Loglevel {
	case "debug":
		logLevel = badger.DEBUG
	case "warning":
		logLevel = badger.WARNING
	case "error":
		logLevel = badger.ERROR
	}
	dbBadger, err := badger.Open(badger.DefaultOptions(dbPath).WithLogger(&dblogger{Logger: logger}).WithLoggingLevel(logLevel))
	if err != nil {
		logger.Error().Err(err).Send()
		return nil, err
	}
	return &BaderDbDot{db: dbBadger}, nil
}

func (p *BaderDbDot) Db() *badger.DB {
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
