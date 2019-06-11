// Scry Info.  All rights reserved.
// license that can be found in the license file.

package gserver

import (
	"go.uber.org/zap"
	"net"
	"os"
	"path/filepath"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/scryg/sutils/sfile"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
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
	//sample:  keys/s1.pem,   Comparing with current exe path
	PemPath string `json:"pemPath"`
	//sample:  keys/s1.pem,   Comparing with current exe path
	KeyPath string `json:"keyPath"`
}

//grpc server component, without bl; one server can monitor in multi address or API at the same time,support tls
type ServerNoblImp struct {
	conf      ConfigNobl
	server    *grpc.Server
	listeners []net.Listener
}

//Construct component
func newServerNobl(conf interface{}) (dot.Dot, error) {
	var err error = nil
	var bs []byte = nil
	if bt, ok := conf.([]byte); ok {
		bs = bt
	} else {
		return nil, dot.SError.Parameter
	}
	dconf := &ConfigNobl{}
	err = dot.UnMarshalConfig(bs, dconf)
	if err != nil {
		return nil, err
	}

	d := &ServerNoblImp{
		conf: *dconf,
	}

	return d, err
}

//Data structure needed when generating newer component
func TypeLiveConns() *dot.TypeLives {
	return &dot.TypeLives{
		Meta: dot.Metadata{TypeId: ServerNoblTypeId, NewDoter: func(conf interface{}) (dot dot.Dot, err error) {
			return newServerNobl(conf)
		}},
	}
}

func (c *ServerNoblImp) Create(l dot.Line) error {
	var err error = nil
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
	if len(c.conf.KeyPath) > 0 && len(c.conf.PemPath) > 0 {
		pem := ""
		key := ""
		{
			if ex, err2 := os.Executable(); err2 == nil {
				exPath := filepath.Dir(ex)
				pem = filepath.Join(exPath, c.conf.PemPath)
				if !sfile.ExistFile(pem) && sfile.ExistFile(c.conf.PemPath) {
					pem = c.conf.PemPath
				}

				key = filepath.Join(exPath, c.conf.KeyPath)
				if !sfile.ExistFile(key) && sfile.ExistFile(c.conf.KeyPath) {
					key = c.conf.KeyPath
				}
			} else {
				dot.Logger().Errorln(err2.Error()) //just log it,  do not return and set err, the code credentials.NewServerTLSFromFile will make it
				pem = c.conf.PemPath
				key = c.conf.KeyPath
			}
		}
		creds, err2 := credentials.NewServerTLSFromFile(testdata.Path(pem), testdata.Path(key))
		if err2 != nil {
			if err != nil {
				dot.Logger().Errorln(err.Error())
			}
			err = err2
		}
		c.server = grpc.NewServer(grpc.Creds(creds), grpc.StreamInterceptor(StreamServerInterceptor()), grpc.UnaryInterceptor(UnaryServerInterceptor()))
	} else {
		c.server = grpc.NewServer(grpc.StreamInterceptor(StreamServerInterceptor()), grpc.UnaryInterceptor(UnaryServerInterceptor()))
	}

	return err
}

//Run after every component finished start, this can ensure all service has been registered on grpc server
func (c *ServerNoblImp) AfterAllStart(l dot.Line) {
	c.startServer()
}

//Stop stop dot
func (c *ServerNoblImp) Stop(ignore bool) error {
	c.server.GracefulStop()
	c.server = nil
	return nil
}

func (c *ServerNoblImp) Server() *grpc.Server {
	return c.server
}

func (c *ServerNoblImp) startServer() {
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
