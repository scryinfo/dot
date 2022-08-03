// Scry Info.  All rights reserved.
// license that can be found in the license file.

package sconfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	"github.com/bitly/go-simplejson"
	"github.com/pelletier/go-toml"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/scryg/sutils/sfile"
	"github.com/scryinfo/yaml"
)

var (
	_ dot.SConfig = (*sConfig)(nil) //just static check implemet the interface
)

//sConfig implement SConfig
//Run executable file content expath，xecutable file name exname (without extension name),expath same content, conf content exconf， config file content confpath
//The process for searching config file content:
// 1，Command line parameter confpath
// 2, exname_conf under expath
// 3, conf under expath
// 4，exname_conf under exconf
// 5，conf in exconf
// 6，If the content above does not existed, then use expath as confpath
//Note: Check whether content existing rather than check whether corrsponding parameters existing
//Config file searching process
//1，Command line parameter conffile
//2，Search exname.json under confpath
//3，Search conf.json under confpath
//4，If file above do not existing, then no config file
//Note: Check whether file existing
type sConfig struct {
	confPath   string           //Config path
	file       string           //File name
	fileType   string           //json,yaml,toml
	simpleJSON *simplejson.Json //All config
	simpleConf *viper.Viper
}

const (
	extensionNameJson = ".json" //extension name of config file 配置文件的扩展名
	extensionNameYaml = ".yaml" //extension name of config file 配置文件的扩展名
	extensionNameToml = ".toml" //extension name of config file 配置文件的扩展名
	separator         = "_"     //Separator
	conf              = "conf"
)

//NewConfig new sConfig
func NewConfig() *sConfig {
	return &sConfig{
		simpleConf: viper.New(),
	}
}

func (c *sConfig) RootPath() {

	if ex, err := os.Executable(); err == nil {
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
			logger := dot.Logger()
			if logger != nil {
				logger.Debugln(fmt.Sprint(err))
			}
		}
	}
}

//Create implement
func (c *sConfig) Create(l dot.Line) error {

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
	switch c.fileType {
	case extensionNameJson:
		c.simpleJSON, err = simplejson.NewFromReader(f)
	case extensionNameToml:
		t, err := toml.LoadReader(f)
		if err == nil {
			jsonStr, err := json.Marshal(t.ToMap())
			if err == nil {
				c.simpleJSON, err = simplejson.NewJson(jsonStr)
			}
		}
	case extensionNameYaml:
		var yamlBytes, jsonBytes []byte
		yamlBytes, err = ioutil.ReadAll(f)
		jsonBytes, err = yaml.YAMLToJSON(yamlBytes)
		if err == nil {
			c.simpleJSON, err = simplejson.NewJson(jsonBytes)
		}
	}

	err = c.simpleConf.ReadInConfig()

	return err
}

////Start  implement
//func (c *sConfig) Start(ignore bool) error {
//	return nil
//}
//
////Stop  implement
//func (c *sConfig) Stop(ignore bool) error {
//	return nil
//}

//Destroy  implement
func (c *sConfig) Destroy(ignore bool) error {
	c.simpleJSON = nil
	c.simpleConf = nil
	return nil
}

//ConfigPath  implement
func (c *sConfig) ConfigPath() string {
	return c.confPath
}

//ConfigFile  implement
func (c *sConfig) ConfigFile() string {
	return c.file
}

//Key  implement
func (c *sConfig) Key(key string) bool {

	re := false
	if c.simpleJSON != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			_, re = c.simpleJSON.CheckGet(key)
		} else if len(keys) > 1 {
			t := c.simpleJSON.GetPath(keys...)
			if t != nil {
				re = true
			}
		}
	}

	return re
}

//Map  implement
func (c *sConfig) Map() (m map[string]interface{}, err error) {
	c.simpleConf.AllSettings()
	return c.simpleJSON.Map()
}

//Unmarshal implement
func (c *sConfig) Unmarshal(s interface{}) error {
	//f := filepath.Join(c.ConfigPath(), c.ConfigFile())
	//var data []byte
	var err error

	err = c.simpleConf.Unmarshal(s)

	//if sfile.ExistFile(f) {
	//	data, err = ioutil.ReadFile(filepath.Join(c.ConfigPath(), c.ConfigFile()))
	//	if err == nil {
	//		switch c.fileType {
	//		case extensionNameJson:
	//			err = json.Unmarshal(data, s)
	//		case extensionNameToml:
	//			err = _toml.Unmarshal(data, s)
	//		case extensionNameYaml:
	//			err = yaml.Unmarshal(data, s)
	//		}
	//	}
	//}

	return err
}

func (c *sConfig) Marshal(data []byte) error {
	var err error
	c.simpleJSON, err = simplejson.NewJson(data)
	return err
}

//DefInterface  implement
func (c *sConfig) DefInterface(key string, def interface{}) interface{} {

	re := def
	if c.simpleJSON != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.simpleJSON.CheckGet(key); ok {
				re = t.Interface()
			}
		} else if len(keys) > 1 {
			t := c.simpleJSON.GetPath(keys...)
			if t != nil {
				re = t.Interface()
			}
		}
	}

	return re
}

func (c *sConfig) UnmarshalKey(key string, obj interface{}) error {

	var err error = nil
	if c.simpleJSON != nil {
		var bs []byte = nil
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.simpleJSON.CheckGet(key); ok {
				bs, err = json.Marshal(t)
			}
		} else if len(keys) > 1 {
			t := c.simpleJSON.GetPath(keys...)
			if t != nil {
				bs, err = json.Marshal(t)
			}
		}
		if err == nil {
			err = json.Unmarshal(bs, obj)
			if err != nil {
				obj = nil
			}
		}
	}
	return err
}

//DefArray  implement
func (c *sConfig) DefArray(key string, def []interface{}) []interface{} {

	re := def
	if c.simpleJSON != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.simpleJSON.CheckGet(key); ok {
				if t2, err := t.Array(); err == nil {
					re = t2
				}
			}
		} else if len(keys) > 1 {
			t := c.simpleJSON.GetPath(keys...)
			if t != nil {
				if t2, err := t.Array(); err == nil {
					re = t2
				}
			}
		}
	}
	return re
}

//DefMap  implement
func (c *sConfig) DefMap(key string, def map[string]interface{}) map[string]interface{} {

	re := def
	if c.simpleJSON != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.simpleJSON.CheckGet(key); ok {
				if t2, err := t.Map(); err == nil {
					re = t2
				}
			}
		} else if len(keys) > 1 {
			t := c.simpleJSON.GetPath(keys...)
			if t != nil {
				if t2, err := t.Map(); err == nil {
					re = t2
				}
			}
		}
	}

	return re
}

//DefString  implement
func (c *sConfig) DefString(key string, def string) string {

	re := def
	if c.simpleJSON != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.simpleJSON.CheckGet(key); ok {
				if t2, err := t.String(); err == nil {
					re = t2
				}
			}
		} else if len(keys) > 1 {
			t := c.simpleJSON.GetPath(keys...)
			if t != nil {
				if t2, err := t.String(); err == nil {
					re = t2
				}
			}
		}
	}

	return re
}

//DefInt32  implement
func (c *sConfig) DefInt32(key string, def int32) int32 {

	re := def
	if c.simpleJSON != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.simpleJSON.CheckGet(key); ok {
				if t2, err := t.Int(); err == nil {
					re = int32(t2)
				}
			}
		} else if len(keys) > 1 {
			t := c.simpleJSON.GetPath(keys...)
			if t != nil {
				if t2, err := t.Int(); err == nil {
					re = int32(t2)
				}
			}
		}
	}

	return re
}

//DefUint32  implement
func (c *sConfig) DefUint32(key string, def uint32) uint32 {

	re := def
	if c.simpleJSON != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.simpleJSON.CheckGet(key); ok {
				if t2, err := t.Uint64(); err == nil {
					re = uint32(t2)
				}
			}
		} else if len(keys) > 1 {
			t := c.simpleJSON.GetPath(keys...)
			if t != nil {
				if t2, err := t.Uint64(); err == nil {
					re = uint32(t2)
				}
			}
		}
	}

	return re
}

//DefInt64  implement
func (c *sConfig) DefInt64(key string, def int64) int64 {

	re := def
	if c.simpleJSON != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.simpleJSON.CheckGet(key); ok {
				if t2, err := t.Int64(); err == nil {
					re = t2
				}
			}
		} else if len(keys) > 1 {
			t := c.simpleJSON.GetPath(keys...)
			if t != nil {
				if t2, err := t.Int64(); err == nil {
					re = t2
				}
			}
		}
	}

	return re

}

//DefUint64  implement
func (c *sConfig) DefUint64(key string, def uint64) uint64 {
	re := def
	if c.simpleJSON != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.simpleJSON.CheckGet(key); ok {
				if t2, err := t.Uint64(); err == nil {
					re = t2
				}
			}
		} else if len(keys) > 1 {
			t := c.simpleJSON.GetPath(keys...)
			if t != nil {
				if t2, err := t.Uint64(); err == nil {
					re = t2
				}
			}
		}
	}

	return re
}

//DefBool  implement
func (c *sConfig) DefBool(key string, def bool) bool {
	re := def
	if c.simpleJSON != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.simpleJSON.CheckGet(key); ok {
				if t2, err := t.Bool(); err == nil {
					re = t2
				}
			}
		} else if len(keys) > 1 {
			t := c.simpleJSON.GetPath(keys...)
			if t != nil {
				if t2, err := t.Bool(); err == nil {
					re = t2
				}
			}
		}
	}

	return re
}

//DefFloat32  implement
func (c *sConfig) DefFloat32(key string, def float32) float32 {
	re := def
	if c.simpleJSON != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.simpleJSON.CheckGet(key); ok {
				if t2, err := t.Float64(); err == nil {
					re = float32(t2)
				}
			}
		} else if len(keys) > 1 {
			t := c.simpleJSON.GetPath(keys...)
			if t != nil {
				if t2, err := t.Float64(); err == nil {
					re = float32(t2)
				}
			}
		}
	}
	return re
}

//DefFloat64  implement
func (c *sConfig) DefFloat64(key string, def float64) float64 {
	re := def
	if c.simpleJSON != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.simpleJSON.CheckGet(key); ok {
				if t2, err := t.Float64(); err == nil {
					re = t2
				}
			}
		} else if len(keys) > 1 {
			t := c.simpleJSON.GetPath(keys...)
			if t != nil {
				if t2, err := t.Float64(); err == nil {
					re = t2
				}
			}
		}
	}

	return re
}

func (c *sConfig) keys(k string) []string {
	re := strings.Split(k, ".")
	if re == nil {
		re = []string{}
	}
	return re
}

func (c *sConfig) getFileType() {
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
