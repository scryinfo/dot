// Scry Info.  All rights reserved.
// license that can be found in the license file.

package sconfig

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/scryg/sutils/sfile"
)

var (
	_ dot.SConfig = (*SConfig)(nil) //just static check implemet the interface
)

// SConfig implement SConfig
// Run executable file content expath，xecutable file name exname (without extension name),expath same content, conf content exconf， config file content confpath
// The process for searching config file content:
// 1，Command line parameter confpath
// 2, exname_conf under expath
// 3, conf under expath
// 4，exname_conf under exconf
// 5，conf in exconf
// 6，If the content above does not existed, then use expath as confpath
// Note: Check whether content existing rather than check whether corrsponding parameters existing
// Config file searching process
// 1，Command line parameter conffile
// 2，Search exname.json under confpath
// 3，Search conf.json under confpath
// 4，If file above do not existing, then no config file
// Note: Check whether file existing
type SConfig struct {
	confPath   string //Config path
	file       string //File name
	fileType   string //json,yaml,toml
	simpleConf *viper.Viper
}

const (
	extensionNameJson = ".json" //extension name of config file 配置文件的扩展名
	extensionNameYaml = ".yaml" //extension name of config file 配置文件的扩展名
	extensionNameToml = ".toml" //extension name of config file 配置文件的扩展名
	separator         = "_"     //Separator
	conf              = "conf"
)

// NewConfig new sConfig
func NewConfig() (*SConfig, error) {
	conf := &SConfig{
		simpleConf: viper.New(),
	}
	conf.RootPath()
	err := conf.create()
	if err != nil {
		dot.Logger.Error().AnErr("cant read the config", err).Send()
	}

	return conf, err
}

func (c *SConfig) RootPath() error {

	ex, err := os.Executable()
	if err != nil {
		return err
	}
	{
		exPath := filepath.Dir(ex)
		binPath := filepath.Dir(exPath)
		exName := filepath.Base(ex)
		ext := filepath.Ext(ex)
		exName = exName[0 : len(exName)-len(ext)]
		if sfile.ExistFile(dot.GCmd.ConfigPath) {
			c.confPath = dot.GCmd.ConfigPath
		} else if configPath := filepath.Join(exPath, exName+separator+conf); sfile.ExistFile(configPath) {
			c.confPath = configPath
		} else if configPath := filepath.Join(exPath, conf); sfile.ExistFile(configPath) {
			c.confPath = configPath
		} else if configPath := filepath.Join(binPath, exName+separator+conf); sfile.ExistFile(configPath) { //prefer the path
			c.confPath = configPath
		} else if configPath := filepath.Join(binPath, conf); sfile.ExistFile(configPath) {
			c.confPath = configPath
		}

		if len(c.confPath) < 1 {
			c.confPath = exPath
		}

		if file := filepath.Join(c.confPath, dot.GCmd.ConfigFile); len(dot.GCmd.ConfigFile) > 0 && sfile.ExistFile(file) {
			c.file = dot.GCmd.ConfigFile
			c.getFileType()
		} else if file := filepath.Join(c.confPath, exName+extensionNameJson); sfile.ExistFile(file) {
			c.file = exName + extensionNameJson
			c.fileType = extensionNameJson
		} else if file := filepath.Join(c.confPath, conf+extensionNameJson); sfile.ExistFile(file) {
			c.file = conf + extensionNameJson
			c.fileType = extensionNameJson
		} else if file := filepath.Join(c.confPath, exName+extensionNameToml); sfile.ExistFile(file) {
			c.file = exName + extensionNameToml
			c.fileType = extensionNameToml
		} else if file := filepath.Join(c.confPath, conf+extensionNameToml); sfile.ExistFile(file) {
			c.file = conf + extensionNameToml
			c.fileType = extensionNameToml
		} else if file := filepath.Join(c.confPath, exName+extensionNameYaml); sfile.ExistFile(file) {
			c.file = exName + extensionNameYaml
			c.fileType = extensionNameYaml
		} else if file := filepath.Join(c.confPath, conf+extensionNameYaml); sfile.ExistFile(file) {
			c.file = conf + extensionNameYaml
			c.fileType = extensionNameYaml
		}
		c.simpleConf.SetConfigFile(c.file)
		//c.simpleConf.SetConfigType(c.fileType)
		c.simpleConf.AddConfigPath(c.confPath)
	}

	if len(c.confPath) > 0 && !sfile.ExistFile(c.confPath) {
		err := os.MkdirAll(c.confPath, os.ModePerm)
		if err != nil {
			dot.Logger.Debug().Err(err).Send()
			return err
		}
	}
	return nil
}

// Create implement
func (c *SConfig) create() error {

	fname := filepath.Join(c.ConfigPath(), c.ConfigFile())
	if len(c.ConfigFile()) < 1 || !sfile.ExistFile(fname) {
		return nil
	}
	f, err := os.Open(fname)
	if err != nil {
		return err
	}

	if state, err := f.Stat(); err == nil && state.Size() < 1 {
		return nil
	}
	defer f.Close()

	err = c.simpleConf.ReadInConfig()

	return err
}

// ConfigPath  implement
func (c *SConfig) ConfigPath() string {
	return c.confPath
}

// ConfigFile  implement
func (c *SConfig) ConfigFile() string {
	return c.file
}

// Key  implement
func (c *SConfig) ExistKey(key string) bool {

	re := false
	if c.simpleConf != nil {
		re = c.simpleConf.InConfig(key)
	}
	return re
}

// Unmarshal implement
func Unmarshal[T any](c *SConfig) (T, error) {
	//f := filepath.Join(c.ConfigPath(), c.ConfigFile())
	//var data []byte
	var err error
	var t T
	err = c.simpleConf.Unmarshal(&t)
	return t, err
}

func UnmarshalKey[T any](c *SConfig, key string) (T, error) {
	var err error
	var t T
	if c.simpleConf != nil {
		err = c.simpleConf.UnmarshalKey(key, &t)
	}
	return t, err
}

// Map  implement
func (c *SConfig) Map() map[string]any {
	return c.simpleConf.AllSettings()
}

// DefMap  implement
func (c *SConfig) DefMap(key string, def map[string]any) map[string]any {
	re := def
	if c.simpleConf != nil {
		re = c.simpleConf.GetStringMap(key)
	}
	return re
}

// DefString  implement
func (c *SConfig) DefString(key string, def string) string {

	re := def
	if c.simpleConf != nil {
		re = c.simpleConf.GetString(key)
	}

	return re
}

// DefInt32  implement
func (c *SConfig) DefInt32(key string, def int32) int32 {

	re := def
	if c.simpleConf != nil {
		re = c.simpleConf.GetInt32(key)
	}

	return re
}

// DefUint32  implement
func (c *SConfig) DefUint32(key string, def uint32) uint32 {

	re := def
	if c.simpleConf != nil {
		re = c.simpleConf.GetUint32(key)
	}

	return re
}

// DefInt64  implement
func (c *SConfig) DefInt64(key string, def int64) int64 {

	re := def
	if c.simpleConf != nil {
		re = c.simpleConf.GetInt64(key)
	}

	return re

}

// DefUint64  implement
func (c *SConfig) DefUint64(key string, def uint64) uint64 {
	re := def
	if c.simpleConf != nil {
		re = c.simpleConf.GetUint64(key)
	}

	return re
}

// DefBool  implement
func (c *SConfig) DefBool(key string, def bool) bool {
	re := def
	if c.simpleConf != nil {
		re = c.simpleConf.GetBool(key)
	}

	return re
}

// DefFloat32  implement
func (c *SConfig) DefFloat32(key string, def float32) float32 {
	re := def
	if c.simpleConf != nil {
		re = float32(c.simpleConf.GetFloat64(key))
	}
	return re
}

// DefFloat64  implement
func (c *SConfig) DefFloat64(key string, def float64) float64 {
	re := def
	if c.simpleConf != nil {
		re = c.simpleConf.GetFloat64(key)
	}
	return re
}

func (c *SConfig) getFileType() {
	re := strings.Split(c.file, ".")
	if l := len(re); l >= 2 {
		switch re[l-1] {
		case "json":
			c.fileType = extensionNameJson
		case "yaml":
			c.fileType = extensionNameYaml
		case "toml":
			c.fileType = extensionNameToml
		}
	}
}
