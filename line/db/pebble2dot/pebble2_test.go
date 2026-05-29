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
)

func TestPebble2(t *testing.T) {
	sourcePath := filepath.Dir(kits.Config.GetCallSourceFile())
	logger := dot.NewTestLogger()
	config := Pebble2Config{
		DbPath: filepath.Join(sourcePath, "temp/pebble"),
	}
	if !sfile.ExistDir(config.DbPath) {
		err := os.MkdirAll(config.DbPath, 0755)
		if err != nil {
			t.Fatal(err)
		}
	}
	db, cleaner, err := NewPebble2(&config, sconfig.NewTestSConfig(sourcePath, sourcePath, sourcePath), logger)
	if err != nil {
		t.Fatal(err)
	}
	defer cleaner()
	opt := &pebble.WriteOptions{}
	err = db.Db().Set(binary.LittleEndian.AppendUint32(nil, uint32(10)), []byte("value"), opt)
	if err != nil {
		t.Fatal(err)
	}
}
