// Scry Info.  All rights reserved.
// license that can be found in the license file.

package sconfig

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/scryg/sutils/sfile"
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
	simpleJson *simplejson.Json //All config
}

const (
	extensionName = ".json" //extension name of config file 配置文件的扩展名
	separator     = "_"     //Separator
	conf          = "conf"
)

//NewConfiger new sConfig
func NewConfiger() *sConfig {
	return &sConfig{}
}

func (c *sConfig) RootPath() {

	if ex, err := os.Executable(); err == nil {
		exPath := filepath.Dir(ex)
		binPath := filepath.Dir(exPath)
		exName := filepath.Base(ex)
		ext := filepath.Ext(ex)
		exName = exName[0 : len(exName)-len(ext)]
		if sfile.ExitFile(dot.GCmd.ConfigPath) {
			c.confPath = dot.GCmd.ConfigPath
		} else if configPath := filepath.Join(exPath, exName+separator+conf); sfile.ExitFile(configPath) {
			c.confPath = configPath
		} else if configPath := filepath.Join(exPath, conf); sfile.ExitFile(configPath) {
			c.confPath = configPath
		} else if configPath := filepath.Join(binPath, exName+separator+conf); sfile.ExitFile(configPath) { //prefer the path
			c.confPath = configPath
		} else if configPath := filepath.Join(binPath, conf); sfile.ExitFile(configPath) {
			c.confPath = configPath
		}

		if len(c.confPath) < 1 {
			c.confPath = exPath
		}

		if file := filepath.Join(c.confPath, dot.GCmd.ConfigFile); len(dot.GCmd.ConfigFile) > 0 && sfile.ExitFile(file) {
			c.file = dot.GCmd.ConfigFile
		} else if file := filepath.Join(c.confPath, exName+extensionName); sfile.ExitFile(file) {
			c.file = exName + extensionName
		} else if file := filepath.Join(c.confPath, conf+extensionName); sfile.ExitFile(file) {
			c.file = conf + extensionName
		}
	}

	if len(c.confPath) > 0 && !sfile.ExitFile(c.confPath) {
		os.MkdirAll(c.confPath, os.ModePerm)
	}
}

//Create implement
func (c *sConfig) Create(l dot.Line) error {

	fname := filepath.Join(c.ConfigPath(), c.ConfigFile())
	if len(c.ConfigFile()) < 1 || !sfile.ExitFile(fname) {
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
	c.simpleJson, err = simplejson.NewFromReader(f)
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
	c.simpleJson = nil
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
	if c.simpleJson != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			_, re = c.simpleJson.CheckGet(key)
		} else if len(keys) > 1 {
			t := c.simpleJson.GetPath(keys...)
			if t != nil {
				re = true
			}
		}
	}

	return re
}

//Map  implement
func (c *sConfig) Map() (m map[string]interface{}, err error) {
	return c.simpleJson.Map()
}

//Unmarshal implement
func (c *sConfig) Unmarshal(s interface{}) error {
	f := filepath.Join(c.ConfigPath(), c.ConfigFile())
	var data []byte
	var err error
	if c.simpleJson != nil {
		data, err = c.simpleJson.MarshalJSON()
	} else if sfile.ExitFile(f) {
		data, err = ioutil.ReadFile(filepath.Join(c.ConfigPath(), c.ConfigFile()))
	}

	if err == nil {
		err = json.Unmarshal(data, s)
	}
	return err
}

func (c *sConfig) Marshal(data []byte) error {
	var err error
	c.simpleJson, err = simplejson.NewJson(data)
	return err
}

//DefInterface  implement
func (c *sConfig) DefInterface(key string, def interface{}) interface{} {

	re := def
	if c.simpleJson != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.simpleJson.CheckGet(key); ok {
				re = t.Interface()
			}
		} else if len(keys) > 1 {
			t := c.simpleJson.GetPath(keys...)
			if t != nil {
				re = t.Interface()
			}
		}
	}

	return re
}

func (c *sConfig) UnmarshalKey(key string, obj interface{}) error {

	var err error = nil
	if c.simpleJson != nil {
		var bs []byte = nil
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.simpleJson.CheckGet(key); ok {
				bs, err = json.Marshal(t)
			}
		} else if len(keys) > 1 {
			t := c.simpleJson.GetPath(keys...)
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
	if c.simpleJson != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.simpleJson.CheckGet(key); ok {
				if t2, err := t.Array(); err == nil {
					re = t2
				}
			}
		} else if len(keys) > 1 {
			t := c.simpleJson.GetPath(keys...)
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
	if c.simpleJson != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.simpleJson.CheckGet(key); ok {
				if t2, err := t.Map(); err == nil {
					re = t2
				}
			}
		} else if len(keys) > 1 {
			t := c.simpleJson.GetPath(keys...)
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
	if c.simpleJson != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.simpleJson.CheckGet(key); ok {
				if t2, err := t.String(); err == nil {
					re = t2
				}
			}
		} else if len(keys) > 1 {
			t := c.simpleJson.GetPath(keys...)
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
	if c.simpleJson != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.simpleJson.CheckGet(key); ok {
				if t2, err := t.Int(); err == nil {
					re = int32(t2)
				}
			}
		} else if len(keys) > 1 {
			t := c.simpleJson.GetPath(keys...)
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
	if c.simpleJson != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.simpleJson.CheckGet(key); ok {
				if t2, err := t.Uint64(); err == nil {
					re = uint32(t2)
				}
			}
		} else if len(keys) > 1 {
			t := c.simpleJson.GetPath(keys...)
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
	if c.simpleJson != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.simpleJson.CheckGet(key); ok {
				if t2, err := t.Int64(); err == nil {
					re = t2
				}
			}
		} else if len(keys) > 1 {
			t := c.simpleJson.GetPath(keys...)
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
	if c.simpleJson != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.simpleJson.CheckGet(key); ok {
				if t2, err := t.Uint64(); err == nil {
					re = t2
				}
			}
		} else if len(keys) > 1 {
			t := c.simpleJson.GetPath(keys...)
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
	if c.simpleJson != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.simpleJson.CheckGet(key); ok {
				if t2, err := t.Bool(); err == nil {
					re = t2
				}
			}
		} else if len(keys) > 1 {
			t := c.simpleJson.GetPath(keys...)
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
	if c.simpleJson != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.simpleJson.CheckGet(key); ok {
				if t2, err := t.Float64(); err == nil {
					re = float32(t2)
				}
			}
		} else if len(keys) > 1 {
			t := c.simpleJson.GetPath(keys...)
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
	if c.simpleJson != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.simpleJson.CheckGet(key); ok {
				if t2, err := t.Float64(); err == nil {
					re = t2
				}
			}
		} else if len(keys) > 1 {
			t := c.simpleJson.GetPath(keys...)
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
