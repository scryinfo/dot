package gserver

import (
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
	//sample:  keys/s1.pem,   相对于当前exe的路径
	PemPath string `json:"pemPath"`
	//sample:  keys/s1.pem,   相对于当前exe的路径
	KeyPath string `json:"keyPath"`
}

//grpc 的 server组件，不带 bl 的；一个server可以同时在多个地址或端口上监听；支持tls
type ServerNoblImp struct {
	conf      ConfigNobl
	server    *grpc.Server
	listeners []net.Listener
}

//构造组件
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

//产生newer组件时需要的数据结构
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
				if !sfile.ExitFile(pem) && sfile.ExitFile(c.conf.PemPath) {
					pem = c.conf.PemPath
				}

				key = filepath.Join(exPath, c.conf.KeyPath)
				if !sfile.ExitFile(key) && sfile.ExitFile(c.conf.KeyPath) {
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
		c.server = grpc.NewServer(grpc.Creds(creds))
	} else {
		c.server = grpc.NewServer()
	}

	return err
}

//在所有的组件完成 start后运行，这样能可以确保所有的 服务都已注册到grpc server上
func (c *ServerNoblImp) AfterStart(l dot.Line) {
	c.startServer()
}

func (c *ServerNoblImp) Server() *grpc.Server {
	return c.server
}

func (c *ServerNoblImp) startServer() {
	for _, lis := range c.listeners {
		go func(li net.Listener) {
			err := c.server.Serve(li)
			if err != nil {
				dot.Logger().Errorln(err.Error())
			}
		}(lis)
	}
}
