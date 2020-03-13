// Scry Info.  All rights reserved.
// license that can be found in the license file.

package conns

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/pkg/errors"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/grpc/lb"
	"github.com/scryinfo/dot/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/resolver"
	"io/ioutil"
)

const (
	ConnsTypeId = "7bf0a017-ef0c-496a-b04c-b1dc262abc8d"
)

//grpc connection, support one Scheme, multi services are below, every service can have multi address(ClientConn load balancing)
type Conns interface {
	//Return default connection, only one connection
	DefaultClientConn() *grpc.ClientConn
	//Return server name corresponding connection
	ClientConn(serviceName string) *grpc.ClientConn
	//Return ClientContext
	ClientContext(serviceName string) *ClientContext
	//return service name
	ServiceName() []string
	//return schemeName,  Scheme is defined at https://github.com/grpc/grpc/blob/master/doc/naming.md
	SchemeName() string
}

type ConnsConfig struct {
	Scheme   string          `json:"scheme"`
	Services []serviceConfig `json:"services"`
}

type serviceConfig struct {
	Name    string          `json:"name"`
	Addrs   []string        `json:"addrs"`
	Tls     utils.TlsConfig `json:"tls"`
	Balance string          `json:"balance"` // round or first, the default value is round
}

type ClientContext struct {
	ClientConn *grpc.ClientConn
	Ctx        context.Context
	Cancel     context.CancelFunc
}

type connsImp struct {
	conns  map[string]*ClientContext
	config ConnsConfig
}

//Construction component
func newConns(conf []byte) (dot.Dot, error) {
	dconf := &ConnsConfig{}
	err := dot.UnMarshalConfig(conf, dconf)
	if err != nil {
		return nil, err
	}

	d := &connsImp{
		config: *dconf,
	}

	return d, err
}

//Data structure needed when generating newer component
func ConnsTypeLives() *dot.TypeLives {
	return &dot.TypeLives{
		Meta: dot.Metadata{TypeId: ConnsTypeId, NewDoter: func(conf []byte) (dot dot.Dot, err error) {
			return newConns(conf)
		}},
	}
}

//jayce edit
//return config of Conn
func ConnsConfigTypeLives() *dot.ConfigTypeLives {
	slice1 := make([]serviceConfig, 0)
	slice2 := make([]string, 0)
	slice2 = append(slice2, "")
	slice1 = append(slice1, serviceConfig{Addrs: slice2})

	return &dot.ConfigTypeLives{
		TypeIdConfig: ConnsTypeId,
		ConfigInfo: &ConnsConfig{
			Services: slice1,
		},
	}
}

func (c *connsImp) Create(l dot.Line) error {
	logger := dot.Logger()
	var err error = nil
	sa := make(map[string][]string, len(c.config.Services))
	{
		for i := range c.config.Services {
			s := c.config.Services[i]
			sa[s.Name] = s.Addrs
		}
	}
	resolver.Register(lb.NewClientBuilder(c.config.Scheme, sa))
	c.conns = make(map[string]*ClientContext, len(c.config.Services))

	errDo := func(er error) {
		if err != nil {
			logger.Errorln("connsImp", zap.Error(err))
		}
		err = er
	}

ForServices:
	for i := range c.config.Services {
		var e1 error = nil
		s := &c.config.Services[i]
		target := fmt.Sprintf("%s:///%s", c.config.Scheme, s.Name)

		var rpc ClientContext
		funRpc := func(rpc *ClientContext, target string, opts ...grpc.DialOption) error {
			rpc.Ctx, rpc.Cancel = context.WithCancel(context.Background())
			var e error
			rpc.ClientConn, e = grpc.DialContext(rpc.Ctx, target, opts...)
			return e
		}
		{
			switch {
			case len(s.Tls.CaPem) > 0 && len(s.Tls.Key) > 0 && len(s.Tls.Pem) > 0: //both tls
				caPemFile := utils.GetFullPathFile(s.Tls.CaPem)
				if len(caPemFile) < 1 {
					errDo(errors.New("the caPem is not empty, and can not find the file: " + s.Tls.CaPem))
					continue ForServices
				}
				keyFile := utils.GetFullPathFile(s.Tls.Key)
				if len(keyFile) < 1 {
					errDo(errors.New("the Key is not empty, and can not find the file: " + s.Tls.Key))
					continue ForServices
				}

				pemFile := utils.GetFullPathFile(s.Tls.Pem)
				if len(pemFile) < 1 {
					errDo(errors.New("the Pem is not empty, and can not find the file: " + s.Tls.Pem))
					continue ForServices
				}

				var tc credentials.TransportCredentials
				{
					pool := x509.NewCertPool()
					{
						caCrt, err1 := ioutil.ReadFile(caPemFile)
						if err1 != nil {
							errDo(errors.WithStack(err1))
							continue ForServices
						}
						if !pool.AppendCertsFromPEM(caCrt) {
							errDo(errors.New("credentials: failed to append certificates"))
							continue ForServices
						}
					}
					cert, err1 := tls.LoadX509KeyPair(pemFile, keyFile)
					if err1 != nil {
						errDo(errors.WithStack(err1))
						continue ForServices
					}

					tc = credentials.NewTLS(&tls.Config{
						ServerName:   s.Tls.ServerNameOverride,
						Certificates: []tls.Certificate{cert},
						RootCAs:      pool, //ClientConn, use the RootCAs
					})
				}

				logger.Infoln("connsImp", zap.String("", "tls with ca"))
				e1 = funRpc(&rpc, target, lb.Balance(s.Balance), grpc.WithTransportCredentials(tc))
			case len(s.Tls.Pem) > 0: //just server
				pemfile := utils.GetFullPathFile(s.Tls.Pem)
				if len(pemfile) < 1 {
					errDo(errors.New("the Pem is not empty, and can not find the file: " + s.Tls.Pem))
					continue ForServices
				}

				creds, err1 := credentials.NewClientTLSFromFile(pemfile, s.Tls.ServerNameOverride)
				if err1 != nil {
					errDo(errors.WithStack(err1))
					continue ForServices
				}
				logger.Infoln("connsImp", zap.String("", "tls no ca"))
				e1 = funRpc(&rpc, target, lb.Balance(s.Balance), grpc.WithTransportCredentials(creds))
			default: //no tls
				logger.Infoln("connsImp", zap.String("", "no tls"))
				e1 = funRpc(&rpc, target, lb.Balance(s.Balance), grpc.WithInsecure())
			}
		}
		if e1 != nil {
			errDo(errors.WithStack(e1))
		} else {
			c.conns[s.Name] = &rpc
		}
	}
	return err
}

func (c *connsImp) Stop(ignore bool) error {
	var err error = nil
	if len(c.conns) > 0 {
		conns := c.conns
		c.conns = nil
		for _, conn := range conns {
			if conn.ClientConn != nil {
				e1 := conn.ClientConn.Close() //todo Cancel request?
				if e1 != nil {                //do not return , close all connection
					if err != nil { //log the err
						dot.Logger().Errorln(err.Error())
					}
					err = e1
				}
			}
		}
	}

	//todo This function is for testing only,
	resolver.UnregisterForTesting(c.config.Scheme)

	return err
}

func (c *connsImp) DefaultClientConn() *grpc.ClientConn {
	var conn *grpc.ClientConn = nil
	if len(c.conns) == 1 {
		for k := range c.conns {
			conn = c.conns[k].ClientConn
		}
	}

	return conn
}

func (c *connsImp) ClientConn(serviceName string) *grpc.ClientConn {
	var conn *grpc.ClientConn = nil
	if len(c.conns) > 0 {
		if c, ok := c.conns[serviceName]; ok {
			conn = c.ClientConn
		}
	}
	return conn
}

func (c *connsImp) ClientContext(serviceName string) *ClientContext {
	var conn *ClientContext = nil
	if len(c.conns) > 0 {
		if c, ok := c.conns[serviceName]; ok {
			conn = c
		}
	}
	return conn
}

func (c *connsImp) ServiceName() []string {
	var sn []string = nil
	if len(c.config.Services) > 0 {
		sn = make([]string, 0, len(c.config.Services))
		for i := range c.config.Services {
			sn = append(sn, c.config.Services[i].Name)
		}
	}
	return sn
}

func (c *connsImp) SchemeName() string {
	return c.config.Scheme
}

func NewTestConns(conf *ConnsConfig) Conns {
	re := &connsImp{}
	re.config = *conf
	re.Create(nil)
	return re
}
