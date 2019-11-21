// Scry Info.  All rights reserved.
// license that can be found in the license file.

package gserver

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/pkg/errors"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"net"
)

const (
	ServerNoblTypeId = "77a766e7-c288-413f-946b-bc9de6df3d70"
)

type ServerNobl interface {
	Server() *grpc.Server
}

type ConfigNobl struct {
	//sample :  1.1.1.1:568
	Addrs []string `json:"addrs"`

	Tls utils.TlsConfig `json:"tls"`
}

//grpc server component, without bl; one server can monitor in multi address or API at the same time,support tls
type serverNoblImp struct {
	conf      ConfigNobl
	server    *grpc.Server
	listeners []net.Listener
}

//Construct component
func newServerNobl(conf []byte) (dot.Dot, error) {
	dconf := &ConfigNobl{}
	err := dot.UnMarshalConfig(conf, dconf)
	if err != nil {
		return nil, err
	}

	d := &serverNoblImp{
		conf: *dconf,
	}

	return d, err
}

//Data structure needed when generating newer component
func ServerNoblTypeLive() *dot.TypeLives {
	return &dot.TypeLives{
		Meta: dot.Metadata{TypeId: ServerNoblTypeId, NewDoter: func(conf []byte) (dot dot.Dot, err error) {
			return newServerNobl(conf)
		}},
	}
}

//jayce edit
//return config of ServerNobl
func ServerNoblConfigTypeLive() *dot.ConfigTypeLives {
	addrs := make([]string, 0)
	addrs = append(addrs, "")
	return &dot.ConfigTypeLives{
		TypeIdConfig: ServerNoblTypeId,
		ConfigInfo: &ConfigNobl{
			Addrs: addrs,
		},
	}
}

func (c *serverNoblImp) Create(l dot.Line) error {
	logger := dot.Logger()
	var err error = nil
	errDo := func(er error) {
		if err != nil {
			logger.Errorln("serverNoblImp", zap.Error(err))
		}
		err = er
	}
	{
		c.listeners = make([]net.Listener, 0, len(c.conf.Addrs))
		for i := range c.conf.Addrs {
			addr := c.conf.Addrs[i]
			var err2 error = nil
			lis, err2 := net.Listen("tcp", addr)
			if err2 != nil {
				if err != nil {
					dot.Logger().Errorln(err.Error())
				}
				err = err2
			} else {
				c.listeners = append(c.listeners, lis)
			}
		}
	}
	//see https://bbengfort.github.io/programmer/2017/03/03/secure-grpc.html
	switch {
	case len(c.conf.Tls.CaPem) > 0 && len(c.conf.Tls.Pem) > 0 && len(c.conf.Tls.Key) > 0:
		pem := utils.GetFullPathFile(c.conf.Tls.Pem)
		if len(pem) < 1 {
			errDo(errors.New("the pem is not empty, and can not find the file: " + c.conf.Tls.Pem))
			return err
		}
		key := utils.GetFullPathFile(c.conf.Tls.Key)
		if len(key) < 1 {
			errDo(errors.New("the key is not empty, and can not find the file: " + c.conf.Tls.Key))
			return err
		}
		ca := utils.GetFullPathFile(c.conf.Tls.CaPem)
		if len(ca) < 1 {
			errDo(errors.New("the ca pem is not empty, and can not find the file: " + c.conf.Tls.CaPem))
			return err
		}

		var tc credentials.TransportCredentials
		{
			creds, err2 := tls.LoadX509KeyPair(pem, key)
			if err2 != nil {
				errDo(errors.WithStack(err2))
				return err
			}

			caCert, err2 := ioutil.ReadFile(ca)
			if err2 != nil {
				errDo(errors.WithStack(err2))
				return err
			}
			caCertPool := x509.NewCertPool()
			if !caCertPool.AppendCertsFromPEM(caCert) {
				errDo(errors.New("credentials: failed to append certificates"))
				return err
			}

			tc = credentials.NewTLS(&tls.Config{
				Certificates: []tls.Certificate{creds},
				ClientCAs:    caCertPool, //server, use the ClientCAs
				ClientAuth:   tls.RequireAndVerifyClientCert,
			})
		}
		logger.Infoln("serverNoblImp", zap.String("", "tls with ca"))

		c.server = grpc.NewServer(grpc.Creds(tc), grpc.StreamInterceptor(StreamServerInterceptor()), grpc.UnaryInterceptor(UnaryServerInterceptor()))
	case len(c.conf.Tls.Pem) > 0 && len(c.conf.Tls.Key) > 0:
		pem := utils.GetFullPathFile(c.conf.Tls.Pem)
		if len(pem) < 1 {
			errDo(errors.New("the pem is not empty, and can not find the file: " + c.conf.Tls.Pem))
			return err
		}
		key := utils.GetFullPathFile(c.conf.Tls.Key)
		if len(key) < 1 {
			errDo(errors.New("the key is not empty, and can not find the file: " + c.conf.Tls.Key))
			return err
		}

		//var tc credentials.TransportCredentials
		//{
		//	creds, err2 := tls.LoadX509KeyPair(pem, key)
		//	if err2 != nil {
		//		errDo(errors.WithStack(err2))
		//		return err
		//	}
		//
		//	tc = credentials.NewTLS(&tls.Config{
		//		Certificates: []tls.Certificate{creds},
		//	})
		//}

		tc, err2 := credentials.NewServerTLSFromFile(pem, key)
		if err2 != nil {
			errDo(errors.WithStack(err2))
			return err
		}
		logger.Infoln("serverNoblImp", zap.String("", "tls no ca"))
		c.server = grpc.NewServer(grpc.Creds(tc), grpc.StreamInterceptor(StreamServerInterceptor()), grpc.UnaryInterceptor(UnaryServerInterceptor()))

	default:
		logger.Infoln("serverNoblImp", zap.String("", "no tls"))
		c.server = grpc.NewServer(grpc.StreamInterceptor(StreamServerInterceptor()), grpc.UnaryInterceptor(UnaryServerInterceptor()))
	}

	return err
}

//Run after every component finished start, this can ensure all service has been registered on grpc server
func (c *serverNoblImp) AfterAllStart(l dot.Line) {
	c.startServer()
}

//Stop stop dot
func (c *serverNoblImp) Stop(ignore bool) error {
	if c.server != nil {
		c.server.GracefulStop()
		c.server = nil
	}

	if len(c.listeners) > 0 {
		for _, lis := range c.listeners {
			if lis != nil {
				lis.Close()
			}
		}
		c.listeners = nil
	}

	return nil
}

func (c *serverNoblImp) Server() *grpc.Server {
	return c.server
}

func (c *serverNoblImp) startServer() {
	for _, lis := range c.listeners {
		go func(li net.Listener) {
			logger := dot.Logger()
			logger.Infoln("ServerNobl", zap.String("", li.Addr().String()))
			err := c.server.Serve(li)
			if err != nil {
				logger.Errorln(err.Error())
			}
		}(lis)
	}
}
