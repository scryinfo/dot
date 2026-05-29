package pebble2dot

import (
	"path/filepath"

	"github.com/cockroachdb/pebble/v2"
	"github.com/scryinfo/dot/dot"
)

type Pebble2 struct {
	db *pebble.DB
}

type Pebble2Config struct {
	DbPath string `json:"dbPath" toml:"dbPath" yaml:"dbPath"`
}

func NewPebble2(config *Pebble2Config, sconfig dot.SConfig, logger *dot.LoggerType) (*Pebble2, func(), error) {
	{
		dpPath, err := sconfig.FullPath(config.DbPath)
		if err != nil {
			config.DbPath = filepath.Join(sconfig.WdPath(), config.DbPath)
		} else {
			config.DbPath = dpPath
		}
	}
	logger.Info().Msgf("pebble path: %s", config.DbPath)
	opt := pebble.Options{
		Logger: &pebbleZerologAdapter{logger: logger},
	}
	db, err := pebble.Open(config.DbPath, &opt)
	if err != nil {
		return nil, nil, err
	}
	return &Pebble2{db: db}, func() {
		err := db.Close()
		if err != nil {
			logger.Error().AnErr("cant close db", err).Send()
		}
	}, nil
}

func (p *Pebble2) Db() *pebble.DB {
	return p.db
}

type pebbleZerologAdapter struct {
	logger *dot.LoggerType
}

func (a *pebbleZerologAdapter) Infof(format string, args ...any) {
	a.logger.Info().Msgf(format, args...)
}
func (a *pebbleZerologAdapter) Errorf(format string, args ...any) {
	a.logger.Error().Msgf(format, args...)
}

func (a *pebbleZerologAdapter) Fatalf(format string, args ...any) {
	a.logger.Fatal().Msgf(format, args...)
}
