package rocksdbdot

import (
	"github.com/linxGnu/grocksdb"
	"github.com/scryinfo/dot/dot"
)

type RocksDbDot struct {
	logger *dot.LoggerType
	config *RocksDbDotConfig
	db     *grocksdb.DB
}

type RocksDbDotConfig struct {
	DbPath string `json:"dbPath" toml:"dbPath" yaml:"dbPath"`
}

func NewRocksDbDot(config *RocksDbDotConfig, sconfig dot.SConfig, logger *dot.LoggerType) (*RocksDbDot, func(), error) {
	if config.DbPath == "" {
		config.DbPath = "data/rkv"
	}
	var err error
	config.DbPath, err = sconfig.FullPath(config.DbPath)
	if err != nil {
		logger.Error().AnErr("FullPath failed", err).Send()
		return nil, nil, err
	}
	opts := grocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	defer opts.Destroy()

	db, err := grocksdb.OpenDb(opts, config.DbPath)
	if err != nil {
		logger.Error().AnErr("OpenDb failed", err).Send()
		return nil, nil, err
	}
	logger.Info().Msgf("DB opened: %s", config.DbPath)

	return &RocksDbDot{
			logger: logger,
			config: config,
			db:     db,
		}, func() {
			defer db.Close()
		}, nil
}

func (r *RocksDbDot) Db() *grocksdb.DB {
	return r.db
}
