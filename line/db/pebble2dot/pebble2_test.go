package pebble2dot

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"github.com/cockroachdb/pebble/v2"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/lib/kits"
	"github.com/scryinfo/dot/line/sconfig"
	"github.com/scryinfo/scryg/sutils/sfile"
	"github.com/stretchr/testify/assert"
)

func TestPebble2(t *testing.T) {
	sourcePath := filepath.Dir(kits.Config.GetCallSourceFile())
	logger := dot.NewTestLogger()
	config := Pebble2Config{
		DbPath: filepath.Join(sourcePath, "temp/pebble"),
	}
	if !sfile.ExistDir(config.DbPath) {
		err := os.MkdirAll(config.DbPath, 0755)
		assert.Nil(t, err)
	}
	db, fClear, err := NewPebble2(&config, sconfig.NewTestSConfig(sourcePath, sourcePath, sourcePath), logger)
	assert.Nil(t, err)
	defer fClear()
	opt := &pebble.WriteOptions{}
	err = db.Db().Set(binary.LittleEndian.AppendUint32(nil, uint32(10)), []byte("value"), opt)
	assert.Nil(t, err)
}
