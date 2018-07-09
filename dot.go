package dot

import (
	"reflect"
)

type TypeId string
type InstanceId string

func (c *TypeId) String() string {
	return string(*c)
}

func (c *InstanceId) String() string {
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
	NewDoter    NewDoter
	RefType     reflect.Type
}

type RelyInstance struct {
	InstId    InstanceId
	RelyInsts []InstanceId
}

func NewMetaData() *MetaData {
	m := &MetaData{}
	return m
}

func (m *MetaData) NewDot(args interface{}) Dot {

	if m.NewDoter != nil {
		return m.NewDoter.NewDot(args)
	} else {
		return reflect.New(m.RefType)
	}
}

//NewDot 创建
type NewDoter interface {
	NewDot(args interface{}) Dot
}

//
type Dot interface {
}

//生命周期过程为：
// Create, Start,Stop,Destroy
// Create 与 Start是分开的， 为了解决不同dot实例之间的依赖， 如果依赖没有问题，那么可以直接在Create中创建并开始，把Start定为空
type Lifer interface {
	//Create 创建 dot， 在这个方法在进行初始，也运行或监听相同内容，最好放在Start方法中实现
	Create(conf SConfiger) error
	//Start
	Start() error
	//Stop
	Stop() error
	//Destroy 销毁 Dot
	Destroy() error
}

//Tag dot自己的标签数据，dot自己使用
type Tager interface {
	//
	SetTag(tag interface{})
	//
	GetTag() (tag interface{})
}

type StatusType int

//Status Status
type Statuser interface {
	Status() StatusType
}

//HotConfig
type HotConfiger interface {
	//Update 更新配置信息， 返回true表示成功
	HotConfig(newConf SConfiger) bool
}

//Check 检测dot，运行一些验证或测试数据，返回对应的结果
type Checker interface {
	Check(args interface{}) interface{}
}

//命令行参数
const (
	//配置文件路径
	Cmd_Confpath = "confpath"
)
