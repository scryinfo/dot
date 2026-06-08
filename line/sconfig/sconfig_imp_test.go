// Scry Info.  All rights reserved.
// license that can be found in the license file.

package sconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/scryinfo/dot/lib/kits"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	// const TEMP = "temp"
	exeFile, err := os.Executable()
	assert.Nil(t, err)
	exeDir := filepath.Dir(exeFile)

	configFile := filepath.Join(exeDir, filepath.Base(exeFile[:len(exeFile)-len(filepath.Ext(exeFile))]))
	configFileExt := configFile + extensionNameToml
	err = os.WriteFile(configFileExt, []byte("#"), 0644)
	assert.Nil(t, err)
	fmt.Println("created configFileExt:", configFileExt)
	defer func() {
		os.Remove(configFileExt)
	}()

	conf, err := NewConfig()
	if err != nil {
		fmt.Printf("err: %+v\n", err)
	}
	time.Sleep(1 * time.Second)
	assert.Nil(t, err)
	assert.NotNil(t, conf)
	assert.Equal(t, conf.confPath, filepath.ToSlash(exeDir), "confPath should be exeDir")
	assert.Equal(t, conf.wdPath, filepath.ToSlash(filepath.Dir(kits.Config.GetCallSourceFile())), "wdPath should be call source file dir")
	assert.Equal(t, conf.exePath, filepath.ToSlash(exeDir), "exePath should be exeDir")
	assert.Equal(t, conf.file, filepath.Base(configFileExt), "file should be config file name")
	assert.Equal(t, conf.fileType, extensionNameToml[1:], "fileType should be extension name without dot")

}
