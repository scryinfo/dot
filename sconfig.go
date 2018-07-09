package dot

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitly/go-simplejson"

	"github.com/scryinfo/scryg/sutils/sfile"
)

var (
	_ SConfiger = (*SConfig)(nil) //just static check implemet the interface
)

//SConfiger 配置属于一个组件Dot,但它太基础了，每一个Dot都需要，所以定义到 dot.go文件中
//S表示scryinfo config 这个名字用的地方太多，加一个s以示区别
type SConfiger interface {
	Lifer

	//配置文件路径
	ConfigPath() string
	//不包含路径，只有文件名
	ConfigFile() string
	//key 是否存在
	Key(key string) bool
	//如果没有配置或配置为空，返回nil
	Map() (m map[string]interface{}, err error)

	//如果没有找到对应的key或数据类型不能转换，一定要特别小心默认值问题，所以在函数前面特别增加“Def”，以提示默认值
	DefInterface(key string, def interface{}) interface{}
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

//SConfig implement SConfiger
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
type SConfig struct {
	ConfPath string
	File     string
	Json     *simplejson.Json
}

const (
	extensionName = ".json"
	separator     = "_"
	conf          = "conf"
)

//NewSConfig new SConfig
func NewSConfig() *SConfig {
	return &SConfig{}
}

func (c *SConfig) rootPath() {

	if ex, err := os.Executable(); err == nil {
		exPath := filepath.Dir(ex)
		binPath := filepath.Dir(exPath)
		exName := filepath.Base(ex)
		ext := filepath.Ext(ex)
		exName = exName[0 : len(exName)-len(ext)]
		if sfile.ExitFile(GCmd.ConfigPath) {
			c.ConfPath = GCmd.ConfigPath
		} else if configPath := filepath.Join(exPath, exName+separator+conf); sfile.ExitFile(configPath) {
			c.ConfPath = configPath
		} else if configPath := filepath.Join(exPath, conf); sfile.ExitFile(configPath) {
			c.ConfPath = configPath
		} else if configPath := filepath.Join(binPath, exName+separator+conf); sfile.ExitFile(configPath) { //prefer the path
			c.ConfPath = configPath
		} else if configPath := filepath.Join(binPath, conf); sfile.ExitFile(configPath) {
			c.ConfPath = configPath
		}

		if len(c.ConfPath) < 1 {
			c.ConfPath = exPath
		}

		if file:= filepath.Join(c.ConfPath, GCmd.ConfigFile); len(GCmd.ConfigFile) > 0 && sfile.ExitFile(file) {
			c.File = file
		} else if file := filepath.Join(c.ConfPath, exName+extensionName); sfile.ExitFile(file) {
			c.File = file
		} else if file := filepath.Join(c.ConfPath, conf +extensionName); sfile.ExitFile(file) {
			c.File = file
		}
	}

	if len(c.ConfPath) > 0 && !sfile.ExitFile(c.ConfPath) {
		os.MkdirAll(c.ConfPath, os.ModePerm)
	}
}

//Create implement
func (c *SConfig) Create(conf SConfiger) error {

	var err error
	//check or init config path
	c.rootPath()

	//read file and set config
	if sfile.ExitFile(c.File) {
		var data []byte
		if data, err = ioutil.ReadFile(c.File); err == nil {
			c.Json, err = simplejson.NewJson(data)
		}
	}

	return err
}

//Start  implement
func (c *SConfig) Start() error {
	return nil
}

//Stop  implement
func (c *SConfig) Stop() error {
	return nil
}

//Destroy  implement
func (c *SConfig) Destroy() error {
	c.Json = nil
	return nil
}

//ConfigPath  implement
func (c *SConfig) ConfigPath() string {
	return c.ConfPath
}

//ConfigFile  implement
func (c *SConfig) ConfigFile() string {
	return c.File
}

//Key  implement
func (c *SConfig) Key(key string) bool {

	re := false
	if c.Json != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			_, re = c.Json.CheckGet(key)
		} else if len(keys) > 1 {
			t := c.Json.GetPath(keys...)
			if t != nil {
				re = true
			}
		}
	}

	return re
}

//Map  implement
func (c *SConfig) Map() (m map[string]interface{}, err error) {
	return c.Json.Map()
}

//DefInterface  implement
func (c *SConfig) DefInterface(key string, def interface{}) interface{} {

	re := def
	if c.Json != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.Json.CheckGet(key); ok {
				re = t.Interface()
			}
		} else if len(keys) > 1 {
			t := c.Json.GetPath(keys...)
			if t != nil {
				re = t.Interface()
			}
		}
	}

	return re
}

//DefArray  implement
func (c *SConfig) DefArray(key string, def []interface{}) []interface{} {

	re := def
	if c.Json != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.Json.CheckGet(key); ok {
				if t2, err := t.Array(); err == nil {
					re = t2
				}
			}
		} else if len(keys) > 1 {
			t := c.Json.GetPath(keys...)
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
func (c *SConfig) DefMap(key string, def map[string]interface{}) map[string]interface{} {

	re := def
	if c.Json != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.Json.CheckGet(key); ok {
				if t2, err := t.Map(); err == nil {
					re = t2
				}
			}
		} else if len(keys) > 1 {
			t := c.Json.GetPath(keys...)
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
func (c *SConfig) DefString(key string, def string) string {

	re := def
	if c.Json != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.Json.CheckGet(key); ok {
				if t2, err := t.String(); err == nil {
					re = t2
				}
			}
		} else if len(keys) > 1 {
			t := c.Json.GetPath(keys...)
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
func (c *SConfig) DefInt32(key string, def int32) int32 {

	re := def
	if c.Json != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.Json.CheckGet(key); ok {
				if t2, err := t.Int(); err == nil {
					re = int32(t2)
				}
			}
		} else if len(keys) > 1 {
			t := c.Json.GetPath(keys...)
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
func (c *SConfig) DefUint32(key string, def uint32) uint32 {

	re := def
	if c.Json != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.Json.CheckGet(key); ok {
				if t2, err := t.Uint64(); err == nil {
					re = uint32(t2)
				}
			}
		} else if len(keys) > 1 {
			t := c.Json.GetPath(keys...)
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
func (c *SConfig) DefInt64(key string, def int64) int64 {

	re := def
	if c.Json != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.Json.CheckGet(key); ok {
				if t2, err := t.Int64(); err == nil {
					re = t2
				}
			}
		} else if len(keys) > 1 {
			t := c.Json.GetPath(keys...)
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
func (c *SConfig) DefUint64(key string, def uint64) uint64 {
	re := def
	if c.Json != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.Json.CheckGet(key); ok {
				if t2, err := t.Uint64(); err == nil {
					re = t2
				}
			}
		} else if len(keys) > 1 {
			t := c.Json.GetPath(keys...)
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
func (c *SConfig) DefBool(key string, def bool) bool {
	re := def
	if c.Json != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.Json.CheckGet(key); ok {
				if t2, err := t.Bool(); err == nil {
					re = t2
				}
			}
		} else if len(keys) > 1 {
			t := c.Json.GetPath(keys...)
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
func (c *SConfig) DefFloat32(key string, def float32) float32 {
	re := def
	if c.Json != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.Json.CheckGet(key); ok {
				if t2, err := t.Float64(); err == nil {
					re = float32(t2)
				}
			}
		} else if len(keys) > 1 {
			t := c.Json.GetPath(keys...)
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
func (c *SConfig) DefFloat64(key string, def float64) float64 {
	re := def
	if c.Json != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.Json.CheckGet(key); ok {
				if t2, err := t.Float64(); err == nil {
					re = t2
				}
			}
		} else if len(keys) > 1 {
			t := c.Json.GetPath(keys...)
			if t != nil {
				if t2, err := t.Float64(); err == nil {
					re = t2
				}
			}
		}
	}

	return re
}

func (c *SConfig) keys(k string) []string {
	re := strings.Split(k, ".")
	if re == nil {
		re = []string{}
	}
	return re
}
