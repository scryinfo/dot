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
//执行文件目录 expath，执行文件名 exname (不包含扩展名),expath的同层目录的conf目录 exconf， 配置文件目录 confpath
//配置文件目录查找过程为：
// 1，命令行参数 confpath
// 2, expath下的 exname_conf
// 3, expath下的 conf
// 4，exconf下的 exname_conf
// 5，exconf的 conf
// 6，以上目录都不存在，则使用 expath 作为 confpath
//注： 是检查目录是否存在，不是检测是否有对应的变量
//配置文件查找过程
//1，命令行参数 conffile
//2，查找confpath下的 exname.json
//3，查找confpath下的 conf.json
//4，以上文件都不存在，则没有配置文件
//注： 是检测文件是否存在
type sConfig struct {
	confPath   string           //配置路径
	file       string           //文件名
	simpleJson *simplejson.Json //整个配置
}

const (
	extensionName = ".json" //配置文件的扩展名
	separator     = "_"     //分隔符
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
