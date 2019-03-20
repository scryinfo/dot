package dot

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitly/go-simplejson"

	"github.com/scryInfo/scryg/sutils/sfile"
)

var (
	_ SConfig = (*sConfig)(nil) //just static check implemet the interface
)

//SConfig 配置属于一个组件Dot,但它太基础了，每一个Dot都需要，所以定义到 dot.go文件中
//S表示scryInfo config 这个名字用的地方太多，加一个s以示区别
type SConfig interface {
	Lifer

	//RootPath root path
	RootPath()
	//配置文件路径
	ConfigPath() string
	//不包含路径，只有文件名
	ConfigFile() string
	//key 是否存在
	Key(key string) bool
	//如果没有配置或配置为空，返回nil
	Map() (m map[string]interface{}, err error)

	Unmarshal(s interface{}) error

	Marshal(data []byte) error

	//如果没有找到对应的key或数据类型不能转换，一定要特别小心默认值问题，所以在函数前面特别增加“Def”，以提示默认值
	DefInterface(key string, def interface{}) interface{}
	DefJson(key string, def interface{}) interface{}
	DefArray(key string, def []interface{}) []interface{}
	DefMap(key string, def map[string]interface{}) map[string]interface{}
	DefString(key string, def string) string
	DefInt32(key string, def int32) int32
	DefUint32(key string, def uint32) uint32
	DefInt64(key string, def int64) int64
	DefUint64(key string, def uint64) uint64
	DefBool(key string, def bool) bool
	DefFloat32(key string, def float32) float32
	DefFloat64(key string, def float64) float64
}

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
	confPath   string
	file       string
	simpleJson *simplejson.Json
}

const (
	extensionName = ".json"
	separator     = "_"
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
		if sfile.ExitFile(GCmd.ConfigPath) {
			c.confPath = GCmd.ConfigPath
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

		if file := filepath.Join(c.confPath, GCmd.ConfigFile); len(GCmd.ConfigFile) > 0 && sfile.ExitFile(file) {
			c.file = GCmd.ConfigFile
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
func (c *sConfig) Create(conf SConfig) error {

	f, err := os.Open(filepath.Join(c.ConfigPath(), c.ConfigFile()))
	if err != nil {
		return err
	}
	defer f.Close()
	c.simpleJson, err = simplejson.NewFromReader(f)
	return err
}

//Start  implement
func (c *sConfig) Start(ignore bool) error {
	return nil
}

//Stop  implement
func (c *sConfig) Stop(ignore bool) error {
	return nil
}

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
	if sfile.ExitFile(f) {
		data, err = ioutil.ReadFile(filepath.Join(c.ConfigPath(), c.ConfigFile()))
	} else if c.simpleJson != nil {
		data, err = c.simpleJson.MarshalJSON()
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

func (c *sConfig) DefJson(key string, def interface{}) interface{} {

	//re := def
	if c.simpleJson != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.simpleJson.CheckGet(key); ok {
				jsonStr, _ := json.Marshal(t)
				json.Unmarshal(jsonStr,def)
			}
		} else if len(keys) > 1 {
			t := c.simpleJson.GetPath(keys...)
			if t != nil {
				jsonStr, _ := json.Marshal(t)
				json.Unmarshal(jsonStr,def)
			}
		}
	}
	return def
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
