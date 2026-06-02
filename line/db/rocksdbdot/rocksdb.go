package rocksdbdot

import (
	"path/filepath"

	"github.com/linxGnu/grocksdb"
	"github.com/scryinfo/dot/dot"
)

type RocksDbDot struct {
	logger *dot.LoggerType
	config *RocksdbDotConfig
	db     *grocksdb.DB
}

type RocksdbDotConfig struct {
	DbPath string `json:"db_path" toml:"db_path" yaml:"db_path" mapstructure:"db_path"`
}

func NewRocksDbDot(config *RocksdbDotConfig, sconfig dot.SConfig, logger *dot.LoggerType) (*RocksDbDot, func(), error) {
	if config.DbPath == "" {
		config.DbPath = "data/rkv"
	}
	{
		dpPath, err := sconfig.FullPath(config.DbPath)
		if err != nil {
			config.DbPath = filepath.Join(sconfig.WdPath(), config.DbPath)
		} else {
			config.DbPath = dpPath
		}
	}
	opts := grocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	defer opts.Destroy()

	logger.Info().Msgf("rocksdb path: %s", config.DbPath)
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
