package conns

import (
	"github.com/scryinfo/dot/dot"
	"google.golang.org/grpc"
)

const (
	ConnNameTypeId = "2d281e77-3133-4718-a5b7-9eea069591ca"
)

type ConnName struct {
	serviceName string
	conn        *grpc.ClientConn
	Conns_      Conns `dot:""`
}

type configName struct {
	Name string `json:"name"`
}

func NewConnName(conf interface{}) (dot.Dot, error) {
	var err error = nil
	var bs []byte = nil
	if bt, ok := conf.([]byte); ok {
		bs = bt
	} else {
		return nil, dot.SError.Parameter
	}
	dconf := &configName{}
	err = dot.UnMarshalConfig(bs, dconf)
	if err != nil {
		return nil, err
	}

	d := &ConnName{
		serviceName: dconf.Name,
	}

	return d, err
}

func TypeLiveConnName() *dot.TypeLives {
	return &dot.TypeLives{
		Meta: dot.Metadata{TypeId: ConnNameTypeId, NewDoter: func(conf interface{}) (dot dot.Dot, err error) {
			return NewConnName(conf)
		}},
	}
}

func (c *ConnName) Conn() *grpc.ClientConn {
	if c.Conns_ == nil && c.Conns_ != nil {
		c.conn = c.Conns_.ClientConn(c.serviceName)
	}
	return c.conn
}
