package conns

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/pkg/errors"
	"time"

	"go.etcd.io/etcd/v3/clientv3"
	etcdnaming "go.etcd.io/etcd/v3/clientv3/naming"
	"google.golang.org/grpc"
	"google.golang.org/grpc/naming"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/utils"
)

const EtcdConnsTypeID = "c91c8c68-281a-4949-8e55-3516664e80c7"

//see https://github.com/etcd-io/etcd/blob/master/Documentation/dev-guide/grpc_naming.md

type configEtcdConns struct {
	Endpoints   []string        `json:"endpoints"`
	DialTimeout int64           `json:"dialTimeout"`
	Tls         utils.TlsConfig `json:"tls"` //not sure about the Mutual Authentication can work

	Names []string `json:"names"` //服务名字
}

//etcd client
type EtcdConns struct {
	etcdClient   *clientv3.Client
	grpcResolver *etcdnaming.GRPCResolver //ref the etcdClient
	conns        map[string]*grpc.ClientConn
	conf         configEtcdConns
	ctx          context.Context
	cancelFun    context.CancelFunc
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

func (c *EtcdConns) Context() context.Context {
	return c.ctx
}

func (c *EtcdConns) CancelFun() context.CancelFunc {
	return c.cancelFun
}

func (c *EtcdConns) GRPCResolver() *etcdnaming.GRPCResolver {
	return c.grpcResolver
}
func (c *EtcdConns) RegisterServer(ctx context.Context, name string, addr string) error {
	if c.grpcResolver != nil {
		return c.grpcResolver.Update(ctx, name, naming.Update{Op: naming.Add, Addr: addr})
	}
	return errors.New("GRPC Resolver is null")
}
func (c *EtcdConns) UnRegisterServer(ctx context.Context, name string, addr string) error {
	if c.grpcResolver != nil {
		return c.grpcResolver.Update(ctx, name, naming.Update{Op: naming.Delete, Addr: addr})
	}
	return errors.New("GRPC Resolver is null")
}

//func (c *EtcdConns) Create(l dot.Line) error {
//
//}
//func (c *EtcdConns) Injected(l dot.Line) error {
//
//}
//func (c *EtcdConns) AfterAllInject(l dot.Line) {
//
//}
//
//func (c *EtcdConns) Start(ignore bool) error {
//
//}

func (c *EtcdConns) Stop(ignore bool) error {
	if c.cancelFun != nil {
		c.cancelFun()
	}
	if c.etcdClient != nil {
		c.etcdClient.Close()
		//c.etcdClient = nil // maybe somewhere use the client, so do not set nil
	}
	return nil
}

//func (c *EtcdConns) Destroy(ignore bool) error {
//
//}

//construct dot
func newEtcdConns(conf []byte) (dot.Dot, error) {
	dconf := &configEtcdConns{}
	err := dot.UnMarshalConfig(conf, dconf)
	if err != nil {
		return nil, err
	}

	d := &EtcdConns{conf: *dconf, conns: make(map[string]*grpc.ClientConn)}
	d.ctx, d.cancelFun = context.WithCancel(context.Background())
	{
		var tlsConfig *tls.Config
		{
			tlsConfig, err = utils.GetTlsConfig(&d.conf.Tls)
			if err != nil {
				return nil, err
			}
		}
		if d.conf.DialTimeout < 1 {
			d.conf.DialTimeout = 6
		}

		d.etcdClient, err = clientv3.New(clientv3.Config{
			Endpoints:   d.conf.Endpoints,
			DialTimeout: time.Duration(d.conf.DialTimeout) * time.Second,
			TLS:         tlsConfig,
			Context:     d.ctx,
		})
		if err != nil {
			return nil, err
		}
	}
	d.grpcResolver = &etcdnaming.GRPCResolver{
		Client: d.etcdClient,
	}
	if len(d.conf.Names) > 0 {
		b := grpc.RoundRobin(d.grpcResolver)

		for _, n := range d.conf.Names {
			conn, err := grpc.Dial(n, grpc.WithBalancer(b), grpc.WithInsecure())
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
		Meta: dot.Metadata{TypeID: EtcdConnsTypeID, NewDoter: func(conf []byte) (dot.Dot, error) {
			return newEtcdConns(conf)
		}},
		//Lives: []dot.Live{
		//	{
		//		LiveID:    EtcdConnsTypeID,
		//		RelyLives: map[string]dot.LiveID{"some field": "some id"},
		//	},
		//},
	}

	lives := []*dot.TypeLives{tl}

	return lives
}

//EtcdConnsConfigTypeLive
func EtcdConnsConfigTypeLive() *dot.ConfigTypeLive {
	paths := make([]string, 0)
	paths = append(paths, "")
	return &dot.ConfigTypeLive{
		TypeIDConfig: EtcdConnsTypeID,
		ConfigInfo:   &configEtcdConns{},
	}
}

func NewEtcd(endpoints []string, names []string) *EtcdConns {
	conf := &configEtcdConns{
		Endpoints:   endpoints,
		DialTimeout: 10,
		Tls:         utils.TlsConfig{},
		Names:       names,
	}
	bs, _ := json.Marshal(conf)
	re, _ := newEtcdConns(bs)
	etcdConns, _ := re.(*EtcdConns)
	return etcdConns
}
