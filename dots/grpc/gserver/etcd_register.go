package gserver

//
//import (
//	"github.com/scryinfo/dot/dot"
//	"github.com/scryinfo/dot/dots/grpc/conns"
//	"go.uber.org/zap"
//	"runtime"
//)
//
//const EtcdRegisterTypeID = "c2c90392-ed06-4e7f-ac2a-662b0153ecb5"
//
//type ServerItem struct {
//	Name  string   `json:"name"`
//	Addrs []string `json:"addrs"` //可以对应多个地址
//}
//
////do not config, if the ServerNobl in the current process
//type ConfigEtcdRegister struct {
//	Items []ServerItem
//}
//type EtcdRegister struct {
//	EtcdConns *conns.EtcdConns `dot:""`
//	items     []ServerItem
//	conf      ConfigEtcdRegister
//}
//
////func (c *EtcdRegister) Create(l dot.Line) error {
////	//todo add
////}
////func (c *EtcdRegister) Injected(l dot.Line) error {
////	//todo add
////}
//func (c *EtcdRegister) AfterAllInject(l dot.Line) {
//	if l != nil {
//		var items []ServerItem
//		l.EachLives(func(live *dot.Live, metadata *dot.Metadata) bool {
//			if live == nil {
//				return true
//			}
//			if s, ok := live.Dot.(ServerNobl); ok {
//				items = append(items, s.ServerItem())
//			}
//			return true
//		})
//		nameItems := make(map[string]*ServerItem, len(items)+len(c.conf.Items))
//		putToMap := func(nameMap map[string]*ServerItem, serverItems []ServerItem) {
//			for i := range serverItems {
//				item := &serverItems[i]
//				mapItem := nameMap[item.Name]
//				if mapItem == nil {
//					mapItem = &ServerItem{Name: item.Name}
//					nameMap[item.Name] = mapItem
//				}
//				for _, addr := range item.Addrs {
//					exist := false
//					for _, old := range mapItem.Addrs { //find the new addr
//						if addr == old {
//							exist = true
//							break
//						}
//					}
//					if !exist { //if it new, append it
//						mapItem.Addrs = append(mapItem.Addrs, addr)
//					}
//				}
//			}
//		}
//		putToMap(nameItems, c.conf.Items)
//		putToMap(nameItems, items)
//
//		c.items = make([]ServerItem, 0, len(nameItems))
//		for _, item := range nameItems {
//			c.items = append(c.items, *item)
//		}
//	}
//}
//
////
////func (c *EtcdRegister) Start(ignore bool) error {
////
////}
//
//func (c *EtcdRegister) AfterAllStart() {
//	logger := dot.Logger()
//	if c.EtcdConns != nil {
//		go func() {
//			for i := range c.items {
//				item := &c.items[i]
//				for _, addr := range item.Addrs {
//					err := c.EtcdConns.RegisterServer(c.EtcdConns.Context(), item.Name, addr)
//					if err != nil {
//						logger.Infoln("EtcdRegister", zap.Error(err))
//					}
//				}
//			}
//		}()
//
//		runtime.Gosched() //给goroutine运行机会
//	}
//}
//
//func (c *EtcdRegister) Stop(ignore bool) error {
//	logger := dot.Logger()
//	var err error
//	if c.EtcdConns != nil {
//		for i := range c.items {
//			item := &c.items[i]
//			for _, addr := range item.Addrs {
//				err2 := c.EtcdConns.UnRegisterServer(c.EtcdConns.Context(), item.Name, addr)
//				if err2 != nil {
//					err = err2
//					logger.Infoln("EtcdRegister", zap.Error(err))
//				}
//			}
//		}
//	}
//	return err
//}
//
////func (c *EtcdRegister) Destroy(ignore bool) error {
////
////}
//
////construct dot
//func newEtcdRegister(conf []byte) (dot.Dot, error) {
//	dconf := &ConfigEtcdRegister{}
//	_ = dot.UnMarshalConfig(conf, dconf)
//	//if err != nil {
//	//	return nil, err
//	//}
//
//	d := &EtcdRegister{conf: *dconf}
//	return d, nil
//}
//
////NewEctcRegisterTest for test
//func NewEctcRegisterTest(config *ConfigEtcdRegister) *EtcdRegister {
//	if config == nil {
//		config = &ConfigEtcdRegister{}
//	}
//	d := &EtcdRegister{conf: *config}
//	return d
//}
//
////EtcdRegisterTypeLives
//func EtcdRegisterTypeLives() []*dot.TypeLives {
//	tl := &dot.TypeLives{
//		Meta: dot.Metadata{TypeID: EtcdRegisterTypeID, NewDoter: func(conf []byte) (dot.Dot, error) {
//			return newEtcdRegister(conf)
//		}},
//		Lives: []dot.Live{
//			{
//				LiveID:    EtcdRegisterTypeID,
//				RelyLives: map[string]dot.LiveID{"EtcdConns": conns.EtcdConnsTypeID},
//			},
//		},
//	}
//
//	lives := []*dot.TypeLives{tl}
//	lives = append(lives, conns.EtcdConnsTypeLives()...)
//
//	return lives
//}
//
////EtcdRegisterConfigTypeLive
//func EtcdRegisterConfigTypeLive() *dot.ConfigTypeLive {
//	paths := make([]string, 0)
//	paths = append(paths, "")
//	return &dot.ConfigTypeLive{
//		TypeIDConfig: EtcdRegisterTypeID,
//		ConfigInfo: &ConfigEtcdRegister{
//			Items: []ServerItem{},
//		},
//	}
//}
