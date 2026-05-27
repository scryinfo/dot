package pebble2

import (
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
	dpPath, err := sconfig.FullPath(config.DbPath)
	if err != nil {
		logger.Error().AnErr("cant get full db path", err).Send()
		return nil, nil, err
	}
	db, err := pebble.Open(dpPath, &pebble.Options{})
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
