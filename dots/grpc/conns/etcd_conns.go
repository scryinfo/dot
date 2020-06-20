package conns

import (
	"crypto/tls"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/naming"
	"google.golang.org/grpc"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/utils"
)

const EtcdConnsTypeId = "c91c8c68-281a-4949-8e55-3516664e80c7"

type configEtcdConns struct {
	Endpoints   []string        `json:"endpoints"`
	DialTimeout int64           `json:"dialTimeout"`
	Names       []string        `json:"names"`
	Tls         utils.TlsConfig `json:"tls"`
}

type EtcdConns struct {
	etcdClient *clientv3.Client
	conns      map[string]*grpc.ClientConn
	conf       configEtcdConns
}

func (c *EtcdConns) ClientConn(serviceName string) *grpc.ClientConn {
	var conn *grpc.ClientConn = nil
	if len(c.conns) > 0 {
		if c, ok := c.conns[serviceName]; ok {
			conn = c
		}
	}
	return conn
}

func (c *EtcdConns) EtcdClient() *clientv3.Client {
	return c.etcdClient
}

//func (c *EtcdConns) Create(l dot.Line) error {
//	//todo add
//}
//func (c *EtcdConns) Injected(l dot.Line) error {
//	//todo add
//}
//func (c *EtcdConns) AfterAllInject(l dot.Line) {
//	//todo add
//}
//
//func (c *EtcdConns) Start(ignore bool) error {
//	//todo add
//}

func (c *EtcdConns) Stop(ignore bool) error {
	if c.etcdClient != nil {
		c.etcdClient.Close()
	}
	return nil
}

//func (c *EtcdConns) Destroy(ignore bool) error {
//	//todo add
//}

//construct dot
func newEtcdConns(conf []byte) (dot.Dot, error) {
	dconf := &configEtcdConns{}
	err := dot.UnMarshalConfig(conf, dconf)
	if err != nil {
		return nil, err
	}

	d := &EtcdConns{conf: *dconf, conns: make(map[string]*grpc.ClientConn)}
	{
		var tlsConfig *tls.Config
		{
			tlsConfig, err = utils.GetTlsConfig(&d.conf.Tls)
			if err != nil {
				return nil, err
			}
		}

		d.etcdClient, err = clientv3.New(clientv3.Config{
			Endpoints:   d.conf.Endpoints,
			DialTimeout: time.Duration(d.conf.DialTimeout) * time.Second,
			TLS:         tlsConfig,
		})
		if err != nil {
			return nil, err
		}
	}

	if len(d.conf.Names) > 0 {
		r := naming.GRPCResolver{Client: d.etcdClient}
		b := grpc.RoundRobin(&r) //todo 改为新版实现

		for _, n := range d.conf.Names {
			conn, err := grpc.Dial(n, grpc.WithBalancer(b), grpc.WithBlock())
			if err != nil {
				return nil, err
			}
			d.conns[n] = conn
		}
	}

	return d, nil
}

//EtcdConnsTypeLives
func EtcdConnsTypeLives() []*dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{TypeId: EtcdConnsTypeId, NewDoter: func(conf []byte) (dot.Dot, error) {
			return newEtcdConns(conf)
		}},
		//Lives: []dot.Live{
		//	{
		//		LiveId:    EtcdConnsTypeId,
		//		RelyLives: map[string]dot.LiveId{"some field": "some id"},
		//	},
		//},
	}

	lives := []*dot.TypeLives{tl}

	return lives
}

//EtcdConnsConfigTypeLive
func EtcdConnsConfigTypeLive() *dot.ConfigTypeLives {
	paths := make([]string, 0)
	paths = append(paths, "")
	return &dot.ConfigTypeLives{
		TypeIdConfig: EtcdConnsTypeId,
		ConfigInfo:   &configEtcdConns{},
	}
}
