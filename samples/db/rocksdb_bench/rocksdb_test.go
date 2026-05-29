package rocksdbbench

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"github.com/cockroachdb/pebble/v2"
	"github.com/dgraph-io/badger/v4"
	"github.com/linxGnu/grocksdb"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/lib/kits"
	"github.com/scryinfo/dot/line/db/baderdot"
	"github.com/scryinfo/dot/line/db/pebble2dot"
	"github.com/scryinfo/dot/line/db/rocksdbdot"
	"github.com/scryinfo/dot/line/sconfig"
	"github.com/scryinfo/scryg/sutils/sfile"
)

func newLogger() *dot.LoggerType {
	conf := dot.TestLogConfig()
	conf.AddStdOut = false
	conf.SetSlog = false
	conf.Level = "error"
	return dot.NewLogger(&conf)
}

func BenchmarkGorocksdb(b *testing.B) {
	sourcePath := filepath.Dir(kits.Config.GetCallSourceFile())
	logger := newLogger()
	logger.Info().Msgf("rocksdb path: %s", filepath.Join(sourcePath, "temp/gorocksdb"))

	config := rocksdbdot.RocksDbDotConfig{
		DbPath: filepath.Join(sourcePath, "temp/gorocksdb"),
	}
	if !sfile.ExistDir(config.DbPath) {
		err := os.MkdirAll(config.DbPath, 0755)
		if err != nil {
			b.Fatal(err)
		}
	}
	db, cleaner, err := rocksdbdot.NewRocksDbDot(&config, sconfig.NewTestSConfig(sourcePath, sourcePath, sourcePath), logger)
	if err != nil {
		b.Fatal(err)
	}
	defer cleaner()

	b.ResetTimer()
	opt := grocksdb.NewDefaultWriteOptions()

	for i := 0; i < b.N; i++ {
		err := db.Db().Put(opt, binary.LittleEndian.AppendUint32(nil, uint32(i)), []byte("value"))
		if err != nil {
			b.Fatal(err)
		}
	}
	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "tps")
}

func BenchmarkBadgerdb(b *testing.B) {
	sourcePath := filepath.Dir(kits.Config.GetCallSourceFile())
	logger := newLogger()
	config := baderdot.BaderDbDotConfig{
		DbPath:   filepath.Join(sourcePath, "temp/badgerdb"),
		Loglevel: "error",
	}
	if !sfile.ExistDir(config.DbPath) {
		err := os.MkdirAll(config.DbPath, 0755)
		if err != nil {
			b.Fatal(err)
		}
	}
	db, cleaner, err := baderdot.NewBaderDot(&config, sconfig.NewTestSConfig(sourcePath, sourcePath, sourcePath), logger)
	if err != nil {
		b.Fatal(err)
	}
	defer cleaner()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := db.Db().Update(func(txn *badger.Txn) error {
			return txn.Set(binary.LittleEndian.AppendUint32([]byte{}, uint32(i)), []byte("value"))
		})
		if err != nil {
			b.Fatal(err)
		}

	}
	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "tps")
}

func BenchmarkPebble(b *testing.B) {
	sourcePath := filepath.Dir(kits.Config.GetCallSourceFile())
	logger := newLogger()
	config := pebble2dot.Pebble2Config{
		DbPath: filepath.Join(sourcePath, "temp/pebble"),
	}
	if !sfile.ExistDir(config.DbPath) {
		err := os.MkdirAll(config.DbPath, 0755)
		if err != nil {
			b.Fatal(err)
		}
	}
	db, cleaner, err := pebble2dot.NewPebble2(&config, sconfig.NewTestSConfig(sourcePath, sourcePath, sourcePath), logger)
	if err != nil {
		b.Fatal(err)
	}
	defer cleaner()

	b.ResetTimer()
	opt := &pebble.WriteOptions{}

	for i := 0; i < b.N; i++ {
		err := db.Db().Set(binary.LittleEndian.AppendUint32(nil, uint32(i)), []byte("value"), opt)
		if err != nil {
			b.Fatal(err)
		}

	}
	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "tps")
}
