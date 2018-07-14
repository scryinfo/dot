package dot

import (
	"reflect"
)

//TypeId dot 的类型唯一id
type TypeId string

//LiveId dot 的实例唯一id
type LiveId string

func (c *TypeId) String() string {
	return string(*c)
}

func (c *LiveId) String() string {
	return string(*c)
}

//MetaData dot的元数据
type MetaData struct {
	TypeId      TypeId
	Version     string
	Name        string
	ShowName    string
	Single      bool
	RelyTypeIds []TypeId
	NewDoter    Newer        `json:"-"`
	RefType     reflect.Type `json:"-"`
}

//Live 依赖的实例
type Live struct {
	TypeId    TypeId
	LiveId    LiveId
	RelyLives []LiveId
	Dot       Dot
}

//NewMetaData @dot.MetaData 的构造函数
func NewMetaData() *MetaData {
	m := &MetaData{}
	return m
}

func (m *MetaData) Clone() *MetaData {
	c := *m
	c.RelyTypeIds = make([]TypeId, len(m.RelyTypeIds))
	copy(c.RelyTypeIds, m.RelyTypeIds)
	return &c
}

//NewDot 构造一个 dot
func (m *MetaData) NewDot(args interface{}) (dot Dot, err error) {

	dot = nil
	err = nil

	if m.NewDoter != nil {
		dot, err = m.NewDoter.New(args)
	} else if m.RefType != nil {
		v := reflect.New(m.RefType)
		dot = v.Interface()
	}

	return
}

//Newer 创建
type Newer interface {
	New(args interface{}) (dot Dot, err error)
}

//Dot 组件
type Dot interface{}

//Lifer 生命周期过程为：
// Create, Start,Stop,Destroy
// Create 与 Start是分开的， 为了解决不同dot实例之间的依赖， 如果依赖没有问题，那么可以直接在Create中创建并开始，把Start定为空
type Lifer interface {
	//Create 创建 dot， 在这个方法在进行初始，也运行或监听相同内容，最好放在Start方法中实现
	Create(conf SConfig) error
	//Start
	//ignore 在调用其它Lifer时，true 出错出后继续，false 出现一个错误直接返回
	Start(ignore bool) error
	//Stop
	//ignore 在调用其它Lifer时，true 出错出后继续，false 出现一个错误直接返回
	Stop(ignore bool) error
	//Destroy 销毁 Dot
	//ignore 在调用其它Lifer时，true 出错出后继续，false 出现一个错误直接返回
	Destroy(ignore bool) error
}

//Tager dot自己的标签数据，dot自己使用
type Tager interface {
	//SetTag set tag
	SetTag(tag interface{})
	//GetTag get tag
	GetTag() (tag interface{})
}

//StatusType status type
type StatusType int

//Statuser Status
type Statuser interface {
	Status() StatusType
}

//HotConfiger hot change config
type HotConfiger interface {
	//Update 更新配置信息， 返回true表示成功
	HotConfig(newConf SConfig) bool
}

//Checker 检测dot，运行一些验证或测试数据，返回对应的结果
type Checker interface {
	Check(args interface{}) interface{}
}

const (
	//TagDot tag dot
	TagDot = "dot"
)
