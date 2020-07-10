package gserver

import "github.com/scryinfo/dot/dot"

const EtcdRegisterTypeId = "c2c90392-ed06-4e7f-ac2a-662b0153ecb5"

type ServerItem struct {
	Name string   `json:"name"`
	Addr []string `json:"addr"` //可以对应多个地址
}

//do not config, if the ServerNobl in the current process
type ConfigEtcdRegister struct {
	Items []ServerItem
}
type EtcdRegister struct {
	serverNobls []ServerNobl
	conf        ConfigEtcdRegister
}

//func (c *EtcdRegister) Create(l dot.Line) error {
//	//todo add
//}
//func (c *EtcdRegister) Injected(l dot.Line) error {
//	//todo add
//}
func (c *EtcdRegister) AfterAllInject(l dot.Line) {

}

//
//func (c *EtcdRegister) Start(ignore bool) error {
//	//todo add
//}
//
//func (c *EtcdRegister) Stop(ignore bool) error {
//	//todo add
//}
//
//func (c *EtcdRegister) Destroy(ignore bool) error {
//	//todo add
//}

//construct dot
func newEtcdRegister(conf []byte) (dot.Dot, error) {
	dconf := &ConfigEtcdRegister{}
	err := dot.UnMarshalConfig(conf, dconf)
	if err != nil {
		return nil, err
	}

	d := &EtcdRegister{conf: *dconf}

	return d, nil
}

//EtcdRegisterTypeLives
func EtcdRegisterTypeLives() []*dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{TypeId: EtcdRegisterTypeId, NewDoter: func(conf []byte) (dot.Dot, error) {
			return newEtcdRegister(conf)
		}},
		Lives: []dot.Live{
			{
				LiveId:    EtcdRegisterTypeId,
				RelyLives: map[string]dot.LiveId{"_": ServerNoblTypeId},
			},
		},
	}

	lives := []*dot.TypeLives{tl}

	return lives
}

//EtcdRegisterConfigTypeLive
func EtcdRegisterConfigTypeLive() *dot.ConfigTypeLives {
	paths := make([]string, 0)
	paths = append(paths, "")
	return &dot.ConfigTypeLives{
		TypeIdConfig: EtcdRegisterTypeId,
		ConfigInfo: &ConfigEtcdRegister{
			Items: []ServerItem{},
		},
	}
}
