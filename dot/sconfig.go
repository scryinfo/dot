// Scry Info.  All rights reserved.
// license that can be found in the license file.

package dot

import (
	"encoding/json"

	"github.com/scryinfo/scryg/sutils/skit"
)

const (
	SconfigLiveId = ""
)

//SConfig 配置属于一个组件Dot,但它太基础了，每一个Dot都需要，所以定义到 dot.go文件中
//S表示scryinfo config 这个名字用的地方太多，加一个s以示区别
type SConfig interface {
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

	//优先从内存中取数据出来操作， 如果内存为nil查看是否有对就原配置文件
	Unmarshal(s interface{}) error
	//解析key为对应的类型
	UnmarshalKey(key string, obj interface{}) error

	Marshal(data []byte) error

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

func UnMarshalConfig(conf []byte, obj interface{}) (err error) {
	err = nil
	if conf != nil {
		err = json.Unmarshal(conf, obj)
	} else {
		err = SError.Parameter
	}
	return err
}

func MarshalConfig(lconf *LiveConfig) (conf []byte, err error) {
	conf = nil
	err = nil

	if lconf != nil {
		if !skit.IsNil(lconf.Json) {
			conf, err = lconf.Json.MarshalJSON()
		}
	}
	return conf, err
}
