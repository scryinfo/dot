package badgerdot

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"github.com/dgraph-io/badger/v4"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/lib/kits"
	"github.com/scryinfo/dot/line/sconfig"
	"github.com/scryinfo/scryg/sutils/sfile"
	"github.com/stretchr/testify/assert"
)

func TestBadger(t *testing.T) {
	db, fClear := newTestBadger(t)
	defer fClear()
	err := db.Db().Update(func(txn *badger.Txn) error {
		return txn.Set(binary.LittleEndian.AppendUint32(nil, uint32(10)), []byte("value"))
	})
	assert.Nil(t, err)
}

func newTestBadger(t *testing.T) (*BadgerDbDot, func()) {
	sourcePath := filepath.Dir(kits.Config.GetCallSourceFile())
	logger := dot.NewTestLogger()
	config := BadgerDbDotConfig{
		DbPath: filepath.Join(sourcePath, "temp/badger"),
	}
	if !sfile.ExistDir(config.DbPath) {
		err := os.MkdirAll(config.DbPath, 0755)
		assert.Nil(t, err)
	}
	db, fClear, err := NewBadgerDot(&config, sconfig.NewTestSConfig(sourcePath, sourcePath, sourcePath), logger)
	assert.Nil(t, err)
	return db, fClear
}
