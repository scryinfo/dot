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

//SConfig 配置属于一个组件Dot,但它太基础了，每一个Dot都需要，所以定义到 dot.go文件中
//S表示scryinfo config 这个名字用的地方太多，加一个s以示区别
type SConfiger interface {
	Lifer

	//配置文件路径
	ConfigRootPath() string
	//不包含路径，只有文件名
	ConfigFile() string
	//key 是否存在
	Key(key string) bool
	//如果没有配置或配置为空，返回nil
	ToJson() (m map[string]interface{}, err error)

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

//SConfig
//执行文件目录 expath，执行文件名 exname (不包含扩展名),expath的上层目录下的conf目录 exconf 配置文件目录 confpath
//配置文件目录查找过程为：
// 1，命令行参数 confpath
// 2, expath下的 exname_conf
// 3, expath下的 conf
// 4，exconf下的 exname_conf
// 5，exconf的 conf
// 6，以上目录都不存在，则使用 expath 作为 confpath
//配置文件查找过程
//1，查找confpath下的 exname.json
//2，查找confpath下的 conf.json
//3，以上文件都不存在，则没有配置文件
type SConfig struct {
	RootPath string
	File     string
	json     *simplejson.Json
}

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
		if sfile.ExitFile(Cmd.Confpath) {
			c.RootPath = Cmd.Confpath
		} else if configPath := filepath.Join(exPath, exName+"_conf"); sfile.ExitFile(configPath) {
			c.RootPath = configPath
		} else if configPath := filepath.Join(exPath, "conf"); sfile.ExitFile(configPath) {
			c.RootPath = configPath
		} else if configPath := filepath.Join(binPath, exName+"_conf"); sfile.ExitFile(configPath) { //prefer the path
			c.RootPath = configPath
		} else if configPath := filepath.Join(binPath, "conf"); sfile.ExitFile(configPath) {
			c.RootPath = configPath
		}

		if len(c.RootPath) < 1 {
			c.RootPath = exPath
		}

		if len(Cmd.Conffile) > 0 {
			c.File = Cmd.Conffile
		} else if file := filepath.Join(c.RootPath, exName+".json"); sfile.ExitFile(file) {
			c.File = exName + ".json"
		} else if file := filepath.Join(c.RootPath, "conf.json"); sfile.ExitFile(file) {
			c.File = "conf.json"
		}
	}

	if len(c.RootPath) > 0 && !sfile.ExitFile(c.RootPath) {
		os.MkdirAll(c.RootPath, os.ModePerm)
	}
}

func (c *SConfig) Create(conf SConfiger) error {

	var err error = nil
	//check or init config path
	c.rootPath()

	//read file and set config
	if sfile.ExitFile(c.File) {
		var data []byte
		if data, err = ioutil.ReadFile(c.File); err == nil {
			c.json, err = simplejson.NewJson(data)
		}
	}

	return err
}

func (c *SConfig) Start() error {
	return nil
}
func (c *SConfig) Stop() error {
	return nil
}

func (c *SConfig) Destroy() error {
	c.json = nil
	return nil
}

func (c *SConfig) ConfigRootPath() string {
	return c.RootPath
}

func (c *SConfig) ConfigFile() string {
	return c.File
}

func (c *SConfig) Key(key string) bool {

	re := false
	if c.json != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			_, re = c.json.CheckGet(key)
		} else if len(keys) > 1 {
			t := c.json.GetPath(keys...)
			if t != nil {
				re = true
			}
		}
	}

	return re
}

func (c *SConfig) ToJson() (m map[string]interface{}, err error) {
	return c.json.Map()
}

func (c *SConfig) DefInterface(key string, def interface{}) interface{} {

	re := def
	if c.json != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.json.CheckGet(key); ok {
				re = t.Interface()
			}
		} else if len(keys) > 1 {
			t := c.json.GetPath(keys...)
			if t != nil {
				re = t.Interface()
			}
		}
	}

	return re
}

func (c *SConfig) DefArray(key string, def []interface{}) []interface{} {

	re := def
	if c.json != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.json.CheckGet(key); ok {
				if t2, err := t.Array(); err == nil {
					re = t2
				}
			}
		} else if len(keys) > 1 {
			t := c.json.GetPath(keys...)
			if t != nil {
				if t2, err := t.Array(); err == nil {
					re = t2
				}
			}
		}
	}
	return re
}

func (c *SConfig) DefMap(key string, def map[string]interface{}) map[string]interface{} {

	re := def
	if c.json != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.json.CheckGet(key); ok {
				if t2, err := t.Map(); err == nil {
					re = t2
				}
			}
		} else if len(keys) > 1 {
			t := c.json.GetPath(keys...)
			if t != nil {
				if t2, err := t.Map(); err == nil {
					re = t2
				}
			}
		}
	}

	return re
}

func (c *SConfig) DefString(key string, def string) string {

	re := def
	if c.json != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.json.CheckGet(key); ok {
				if t2, err := t.String(); err == nil {
					re = t2
				}
			}
		} else if len(keys) > 1 {
			t := c.json.GetPath(keys...)
			if t != nil {
				if t2, err := t.String(); err == nil {
					re = t2
				}
			}
		}
	}

	return re
}

func (c *SConfig) DefInt32(key string, def int32) int32 {

	re := def
	if c.json != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.json.CheckGet(key); ok {
				if t2, err := t.Int(); err == nil {
					re = int32(t2)
				}
			}
		} else if len(keys) > 1 {
			t := c.json.GetPath(keys...)
			if t != nil {
				if t2, err := t.Int(); err == nil {
					re = int32(t2)
				}
			}
		}
	}

	return re
}

func (c *SConfig) DefUint32(key string, def uint32) uint32 {

	re := def
	if c.json != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.json.CheckGet(key); ok {
				if t2, err := t.Uint64(); err == nil {
					re = uint32(t2)
				}
			}
		} else if len(keys) > 1 {
			t := c.json.GetPath(keys...)
			if t != nil {
				if t2, err := t.Uint64(); err == nil {
					re = uint32(t2)
				}
			}
		}
	}

	return re
}

func (c *SConfig) DefInt64(key string, def int64) int64 {

	re := def
	if c.json != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.json.CheckGet(key); ok {
				if t2, err := t.Int64(); err == nil {
					re = t2
				}
			}
		} else if len(keys) > 1 {
			t := c.json.GetPath(keys...)
			if t != nil {
				if t2, err := t.Int64(); err == nil {
					re = t2
				}
			}
		}
	}

	return re

}

func (c *SConfig) DefUint64(key string, def uint64) uint64 {
	re := def
	if c.json != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.json.CheckGet(key); ok {
				if t2, err := t.Uint64(); err == nil {
					re = t2
				}
			}
		} else if len(keys) > 1 {
			t := c.json.GetPath(keys...)
			if t != nil {
				if t2, err := t.Uint64(); err == nil {
					re = t2
				}
			}
		}
	}

	return re
}

func (c *SConfig) DefBool(key string, def bool) bool {
	re := def
	if c.json != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.json.CheckGet(key); ok {
				if t2, err := t.Bool(); err == nil {
					re = t2
				}
			}
		} else if len(keys) > 1 {
			t := c.json.GetPath(keys...)
			if t != nil {
				if t2, err := t.Bool(); err == nil {
					re = t2
				}
			}
		}
	}

	return re
}

func (c *SConfig) DefFloat32(key string, def float32) float32 {
	re := def
	if c.json != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.json.CheckGet(key); ok {
				if t2, err := t.Float64(); err == nil {
					re = float32(t2)
				}
			}
		} else if len(keys) > 1 {
			t := c.json.GetPath(keys...)
			if t != nil {
				if t2, err := t.Float64(); err == nil {
					re = float32(t2)
				}
			}
		}
	}

	return re
}

func (c *SConfig) DefFloat64(key string, def float64) float64 {
	re := def
	if c.json != nil {
		keys := c.keys(key)
		if len(keys) == 1 {
			if t, ok := c.json.CheckGet(key); ok {
				if t2, err := t.Float64(); err == nil {
					re = t2
				}
			}
		} else if len(keys) > 1 {
			t := c.json.GetPath(keys...)
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
