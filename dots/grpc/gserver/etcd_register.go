package gserver

import (
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/grpc/conns"
	"go.uber.org/zap"
)

const EtcdRegisterTypeId = "c2c90392-ed06-4e7f-ac2a-662b0153ecb5"

type ServerItem struct {
	Name  string   `json:"name"`
	Addrs []string `json:"addrs"` //可以对应多个地址
}

//do not config, if the ServerNobl in the current process
type ConfigEtcdRegister struct {
	Items []ServerItem
}
type EtcdRegister struct {
	EtcdConns *conns.EtcdConns `dot:""`
	items     []ServerItem
	conf      ConfigEtcdRegister
}

//func (c *EtcdRegister) Create(l dot.Line) error {
//	//todo add
//}
//func (c *EtcdRegister) Injected(l dot.Line) error {
//	//todo add
//}
func (c *EtcdRegister) AfterAllInject(l dot.Line) {
	if l != nil {
		var items []ServerItem
		l.EachLives(func(live *dot.Live, metadata *dot.Metadata) bool {
			if live == nil {
				return true
			}
			if s, ok := live.Dot.(ServerNobl); ok {
				items = append(items, s.ServerItem())
			}
			return true
		})
		nameItems := make(map[string]*ServerItem, len(items)+len(c.conf.Items))
		for i := range c.conf.Items {
			it := &c.conf.Items[i]
			mit := nameItems[it.Name]
			if mit == nil {
				mit = &ServerItem{}
				nameItems[it.Name] = mit
			}
			for ai := range it.Addrs {
				mit.Addrs = append(mit.Addrs, it.Addrs[ai])
			}
		}

		for i := range items {
			it := &items[i]
			mit := nameItems[it.Name]
			if mit == nil {
				mit = &ServerItem{}
				nameItems[it.Name] = mit
			}
			for ai := range it.Addrs {
				mit.Addrs = append(mit.Addrs, it.Addrs[ai])
			}
		}

		c.items = make([]ServerItem, 0, len(nameItems))
		for _, item := range nameItems {
			c.items = append(c.items, *item)
		}
	}
}

//
//func (c *EtcdRegister) Start(ignore bool) error {
//	//todo add
//}

func (c *EtcdRegister) AfterAllStart() {
	logger := dot.Logger()
	if c.EtcdConns != nil {
		for i := range c.items {
			item := &c.items[i]
			for _, addr := range item.Addrs {
				err := c.EtcdConns.RegisterServer(item.Name, addr)
				if err != nil {
					logger.Infoln("EtcdRegister", zap.Error(err))
				}
			}
		}
	}
}

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
				RelyLives: map[string]dot.LiveId{"EtcdConns": conns.EtcdConnsTypeId},
			},
		},
	}

	lives := []*dot.TypeLives{tl}
	lives = append(lives, conns.EtcdConnsTypeLives()...)

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
