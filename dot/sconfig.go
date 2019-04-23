package dot

import ()

const (
	SconfigLiveId = ""
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
