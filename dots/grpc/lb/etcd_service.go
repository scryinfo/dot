package lb
//
//import (
//	"go.etcd.io/etcd/client/v3/embed"
//
//	"github.com/scryinfo/dot/dot"
//	"github.com/scryinfo/dot/utils"
//)
//
//const EtcdServiceTypeID = "59b10f33-4f8b-40b9-9912-02068f78c228"
//
//type configEtcdService struct {
//	Pears     []string        `json:"endpoints"`
//	TlsClient utils.TlsConfig `json:"tls"`
//}
//type EtcdService struct {
//	etcd *embed.Etcd
//	conf configEtcdService
//	//todo add
//}
//
////func (c *EtcdService) Create(l dot.Line) error {
////	//todo add
////}
////func (c *EtcdService) Injected(l dot.Line) error {
////	//todo add
////}
////func (c *EtcdService) AfterAllInject(l dot.Line) {
////	//todo add
////}
////
////func (c *EtcdService) Start(ignore bool) error {
////	//todo add
////}
////
////func (c *EtcdService) Stop(ignore bool) error {
////	//todo add
////}
//
//func (c *EtcdService) Destroy(ignore bool) error {
//
//	if c.etcd != nil {
//		c.etcd.Close()
//	}
//
//	return nil
//}
//
////construct dot
//func newEtcdService(conf []byte) (dot.Dot, error) {
//	dconf := &configEtcdService{}
//
//	err := dot.UnMarshalConfig(conf, dconf)
//	if err != nil {
//		return nil, err
//	}
//
//	d := &EtcdService{conf: *dconf}
//
//	cfg := embed.NewConfig()
//	cfg.Dir = "etcd"
//	d.etcd, err = embed.StartEtcd(cfg) //todo how to set the logger
//	if err != nil {
//		d.etcd = nil
//		return nil, err
//	}
//
//	return d, nil
//}
//
////EtcdServiceTypeLives
//func EtcdServiceTypeLives() []*dot.TypeLives {
//	tl := &dot.TypeLives{
//		Meta: dot.Metadata{TypeID: EtcdServiceTypeID, NewDoter: func(conf []byte) (dot.Dot, error) {
//			return newEtcdService(conf)
//		}},
//		//Lives: []dot.Live{
//		//	{
//		//		LiveID:    EtcdServiceTypeID,
//		//		RelyLives: map[string]dot.LiveID{"some field": "some id"},
//		//	},
//		//},
//	}
//
//	lives := []*dot.TypeLives{tl}
//
//	return lives
//}
//
////EtcdServiceConfigTypeLive
//func EtcdServiceConfigTypeLive() *dot.ConfigTypeLive {
//	paths := make([]string, 0)
//	paths = append(paths, "")
//	return &dot.ConfigTypeLive{
//		TypeIDConfig: EtcdServiceTypeID,
//		ConfigInfo:   &configEtcdService{
//			//todo
//		},
//	}
//}
