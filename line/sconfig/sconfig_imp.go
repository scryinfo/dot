// Scry Info.  All rights reserved.
// license that can be found in the license file.

package sconfig

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/lib/kits"
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
	wdPath     string
	exePath    string //Executable file path 可执行文件路径
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
	fmt.Println("initing conifg")
	err := conf.LoadLocal()
	fmt.Printf("config file: %s/%s\nexe path: %s\nwd path: %s\n", conf.confPath, conf.file, conf.exePath, conf.wdPath)
	if err != nil {
		// the config is the first, and the logger is not initialized, so use fmt.Printf
		fmt.Printf("cant read the config: %+v\n", err)
	}

	return conf, err
}

// NewConfig new sConfig
func NewConfigFromGetter(getter dot.ConfigGetter) (*SConfig, error) {
	conf := &SConfig{
		simpleConf: viper.New(),
		fileType:   getter.FileType(),
	}
	fmt.Println("initing conifg")
	err := conf.Load(getter)
	fmt.Printf("config file: %s/%s\nexe path: %s\nwd path: %s\n", conf.confPath, conf.file, conf.exePath, conf.wdPath)
	if err != nil {
		// the config is the first, and the logger is not initialized, so use fmt.Printf
		fmt.Printf("cant read the config: %+v\n", err)
	}

	return conf, err
}

func NewLineConfig[T any](config *SConfig) (*T, error) {
	conf, err := Unmarshal[T](config)
	if err != nil {
		log.Fatalf("%v", err)
		return nil, err
	}
	return &conf, nil
}

func (p *SConfig) LoadLocal() error {
	if p.wdPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		p.wdPath = filepath.ToSlash(wd)
	}

	{
		exeFile, err := os.Executable()
		if err != nil {
			return err
		}
		fmt.Printf("exe file: %s\n", exeFile)
		p.exePath = filepath.ToSlash(filepath.Dir(exeFile))
		checkPath := p.exePath
		binPath := filepath.Dir(p.exePath)
		exeName := filepath.Base(exeFile)
		ext := filepath.Ext(exeFile)
		exeName = exeName[0 : len(exeName)-len(ext)]
		if dot.IsDebug {
			mainFile := kits.Config.GetMainPackageDir()
			if len(mainFile) > 0 {
				mainFile = filepath.ToSlash(mainFile)
				checkPath = filepath.Dir(mainFile)
				exeName = filepath.Base(mainFile)
				exeName = exeName[0 : len(exeName)-len(".go")]
				fmt.Printf("check path is from main file: %s\n", mainFile)
			}
		}
		if p.confPath == "" {
			if sfile.ExistDir(dot.GCmd.ConfigPath) {
				var err error
				p.confPath, err = filepath.Abs(dot.GCmd.ConfigPath)
				if err != nil {
					return err
				}
				fmt.Printf("get conf path from cmd parameter, %s\n", dot.GCmd.ConfigPath)
			} else if configPath := filepath.Join(checkPath, exeName+extensionNameToml); sfile.ExistFile(configPath) {
				p.confPath = checkPath
				p.file = exeName + extensionNameToml
				p.fileType = extensionNameToml[1:]
				fmt.Printf("get conf file from check path, config path: %s, file: %s, file type: %s\n", p.confPath, p.file, p.fileType)
			} else if configPath := filepath.Join(checkPath, exeName+extensionNameYaml); sfile.ExistFile(configPath) {
				p.confPath = checkPath
				p.file = exeName + extensionNameYaml
				p.fileType = extensionNameYaml[1:]
				fmt.Printf("get conf file from check path, config path: %s, file: %s, file type: %s\n", p.confPath, p.file, p.fileType)
			} else if configPath := filepath.Join(checkPath, exeName+extensionNameJson); sfile.ExistFile(configPath) {
				p.confPath = checkPath
				p.file = exeName + extensionNameJson
				p.fileType = extensionNameJson[1:]
				fmt.Printf("get conf file from check path, config path: %s, file: %s, file type: %s\n", p.confPath, p.file, p.fileType)
			} else if configPath := filepath.Join(checkPath, conf+extensionNameToml); sfile.ExistFile(configPath) {
				p.confPath = checkPath
				p.file = conf + extensionNameToml
				p.fileType = extensionNameToml[1:]
				fmt.Printf("get conf file from check path, config path: %s, file: %s, file type: %s\n", p.confPath, p.file, p.fileType)
			} else if configPath := filepath.Join(checkPath, conf+extensionNameYaml); sfile.ExistFile(configPath) {
				p.confPath = checkPath
				p.file = conf + extensionNameYaml
				p.fileType = extensionNameYaml[1:]
				fmt.Printf("get conf file from check path, config path: %s, file: %s, file type: %s\n", p.confPath, p.file, p.fileType)
			} else if configPath := filepath.Join(checkPath, conf+extensionNameJson); sfile.ExistFile(configPath) {
				p.confPath = checkPath
				p.file = conf + extensionNameJson
				p.fileType = extensionNameJson[1:]
				fmt.Printf("get conf file from check path, config path: %s, file: %s, file type: %s\n", p.confPath, p.file, p.fileType)
			} else if configPath := filepath.Join(binPath, exeName+extensionNameToml); sfile.ExistFile(configPath) {
				p.confPath = binPath
				p.file = exeName + extensionNameToml
				p.fileType = extensionNameToml[1:]
				fmt.Printf("get conf file from bin path, config path: %s, file: %s, file type: %s\n", p.confPath, p.file, p.fileType)
			} else if configPath := filepath.Join(binPath, exeName+extensionNameYaml); sfile.ExistFile(configPath) {
				p.confPath = binPath
				p.file = exeName + extensionNameYaml
				p.fileType = extensionNameYaml[1:]
				fmt.Printf("get conf file from bin path, config path: %s, file: %s, file type: %s\n", p.confPath, p.file, p.fileType)
			} else if configPath := filepath.Join(binPath, exeName+extensionNameJson); sfile.ExistFile(configPath) {
				p.confPath = binPath
				p.file = exeName + extensionNameJson
				p.fileType = extensionNameJson[1:]
				fmt.Printf("get conf file from bin path, config path: %s, file: %s, file type: %s\n", p.confPath, p.file, p.fileType)
			} else if configPath := filepath.Join(p.wdPath, exeName+extensionNameToml); sfile.ExistFile(configPath) {
				p.confPath = p.wdPath
				p.file = exeName + extensionNameToml
				p.fileType = extensionNameToml[1:]
				fmt.Printf("get conf file from wd path, config path: %s, file: %s, file type: %s\n", p.confPath, p.file, p.fileType)
			} else if configPath := filepath.Join(p.wdPath, exeName+extensionNameYaml); sfile.ExistFile(configPath) {
				p.confPath = p.wdPath
				p.file = exeName + extensionNameYaml
				p.fileType = extensionNameYaml[1:]
				fmt.Printf("get conf file from wd path, config path: %s, file: %s, file type: %s\n", p.confPath, p.file, p.fileType)
			} else if configPath := filepath.Join(p.wdPath, exeName+extensionNameJson); sfile.ExistFile(configPath) {
				p.confPath = p.wdPath
				p.file = exeName + extensionNameJson
				p.fileType = extensionNameJson[1:]
				fmt.Printf("get conf file from wd path, config path: %s, file: %s, file type: %s\n", p.confPath, p.file, p.fileType)
			}
			if len(p.confPath) < 1 {
				p.confPath = p.exePath
			}
		}

		if file := filepath.Join(p.confPath, dot.GCmd.ConfigFile); len(dot.GCmd.ConfigFile) > 0 && sfile.ExistFile(file) {
			p.file = dot.GCmd.ConfigFile
			p.getFileType()
			fmt.Printf("get conf file from cmd config file: %s\n", dot.GCmd.ConfigFile)
		} else if p.file == "" {
			err := fmt.Errorf("config file not found")
			return err
		}
		p.simpleConf.SetConfigFile(p.file)
		p.simpleConf.SetConfigType(p.fileType)
	}
	f, err := os.Open(filepath.Join(p.confPath, p.file))
	if err != nil {
		return err
	}
	defer func() {
		if f != nil {
			_ = f.Close()
		}
	}()
	return p.simpleConf.ReadConfig(f)
}

func (p *SConfig) Load(getter dot.ConfigGetter) error {
	// p.fileType = getter.FileType()
	return p.simpleConf.ReadConfig(getter)
}

// ConfigPath  implement
func (p *SConfig) ConfigPath() string {
	return p.confPath
}

// ExePath implements [dot.SConfig].
func (p *SConfig) ExePath() string {
	return p.exePath
}

// WdPath implements [dot.SConfig].
func (p *SConfig) WdPath() string {
	return p.wdPath
}

// ConfigFile  implement
func (p *SConfig) ConfigFile() string {
	return p.file
}

// Key  implement
func (p *SConfig) ExistKey(key string) bool {

	re := false
	if p.simpleConf != nil {
		re = p.simpleConf.InConfig(key)
	}
	return re
}

// Unmarshal implement
func Unmarshal[T any](c *SConfig) (T, error) {
	//f := filepath.Join(c.ConfigPath(), c.ConfigFile())
	//var data []byte
	var err error
	var t T
	err = c.simpleConf.Unmarshal(&t, func(dc *mapstructure.DecoderConfig) {
		dc.ErrorUnused = true
		dc.ErrorUnset = true
		dc.TagName = c.fileType
	})
	return t, err
}

func UnmarshalKey[T any](c *SConfig, key string) (T, error) {
	var err error
	var t T
	if c.simpleConf != nil {
		err = c.simpleConf.UnmarshalKey(key, &t, func(dc *mapstructure.DecoderConfig) {
			dc.ErrorUnused = true
			dc.ErrorUnset = true
			dc.TagName = c.fileType
		})
	}
	return t, err
}

// Map  implement
func (p *SConfig) Map() map[string]any {
	return p.simpleConf.AllSettings()
}

// DefMap  implement
func (p *SConfig) DefMap(key string, def map[string]any) map[string]any {
	re := def
	if p.simpleConf != nil {
		re = p.simpleConf.GetStringMap(key)
	}
	return re
}

// DefString  implement
func (p *SConfig) DefString(key string, def string) string {

	re := def
	if p.simpleConf != nil {
		re = p.simpleConf.GetString(key)
	}

	return re
}

// DefInt32  implement
func (p *SConfig) DefInt32(key string, def int32) int32 {

	re := def
	if p.simpleConf != nil {
		re = p.simpleConf.GetInt32(key)
	}

	return re
}

// DefUint32  implement
func (p *SConfig) DefUint32(key string, def uint32) uint32 {

	re := def
	if p.simpleConf != nil {
		re = p.simpleConf.GetUint32(key)
	}

	return re
}

// DefInt64  implement
func (p *SConfig) DefInt64(key string, def int64) int64 {

	re := def
	if p.simpleConf != nil {
		re = p.simpleConf.GetInt64(key)
	}

	return re

}

// DefUint64  implement
func (p *SConfig) DefUint64(key string, def uint64) uint64 {
	re := def
	if p.simpleConf != nil {
		re = p.simpleConf.GetUint64(key)
	}

	return re
}

// DefBool  implement
func (p *SConfig) DefBool(key string, def bool) bool {
	re := def
	if p.simpleConf != nil {
		re = p.simpleConf.GetBool(key)
	}

	return re
}

// DefFloat32  implement
func (p *SConfig) DefFloat32(key string, def float32) float32 {
	re := def
	if p.simpleConf != nil {
		re = float32(p.simpleConf.GetFloat64(key))
	}
	return re
}

// DefFloat64  implement
func (p *SConfig) DefFloat64(key string, def float64) float64 {
	re := def
	if p.simpleConf != nil {
		re = p.simpleConf.GetFloat64(key)
	}
	return re
}

func (p *SConfig) getFileType() {
	re := strings.Split(p.file, ".")
	if l := len(re); l >= 2 {
		switch re[l-1] {
		case "json":
			p.fileType = extensionNameJson[1:]
		case "yaml":
			p.fileType = extensionNameYaml[1:]
		case "toml":
			p.fileType = extensionNameToml[1:]
		}
	}
}

// FullPath implements [dot.SConfig].
// if file exists, return the full path of the file,
// else join config path, if the file exists, return the full path,
// else join wd path, if the file exists, return the full path,
// else join exe path, if the file exists, return the full path,
// otherwise return an error
func (p *SConfig) FullPath(file string) (string, error) {
	if sfile.Exist(file) {
		return filepath.Abs(file)
	}
	if sfile.Exist(filepath.Join(p.confPath, file)) {
		return filepath.Abs(filepath.Join(p.confPath, file))
	}
	if sfile.Exist(filepath.Join(p.wdPath, file)) {
		return filepath.Abs(filepath.Join(p.wdPath, file))
	}
	if sfile.Exist(filepath.Join(p.exePath, file)) {
		return filepath.Abs(filepath.Join(p.exePath, file))
	}
	return "", fmt.Errorf("file %s not found", file)
}

func NewTestSConfig(confPath, wdPath, exePath string) *SConfig {
	return &SConfig{
		confPath:   confPath,
		wdPath:     wdPath,
		exePath:    exePath,
		file:       "config.toml",
		fileType:   "toml",
		simpleConf: nil,
	}
}
