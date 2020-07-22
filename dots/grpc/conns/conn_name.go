// Scry Info.  All rights reserved.
// license that can be found in the license file.

package conns

import (
	"google.golang.org/grpc"

	"github.com/scryinfo/dot/dot"
)

const (
	ConnNameTypeID = "2d281e77-3133-4718-a5b7-9eea069591ca"
)

type ConnName struct {
	serviceName string
	conn        *grpc.ClientConn
	Conns_      Conns      `dot:"?"` //可选组件，优先使用etcd
	EtcdConns   *EtcdConns `dot:"?"` //可选组件，优先使用etcd
}

type configName struct {
	Name string `json:"name"`
}

func newConnName(conf []byte) (dot.Dot, error) {
	dconf := &configName{}
	err := dot.UnMarshalConfig(conf, dconf)
	if err != nil {
		return nil, err
	}

	d := &ConnName{
		serviceName: dconf.Name,
	}

	return d, err
}

func ConnNameTypeLives() []*dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{TypeID: ConnNameTypeID, NewDoter: func(conf []byte) (dot dot.Dot, err error) {
			return newConnName(conf)
		}},
		Lives: []dot.Live{
			{
				LiveID:    ConnNameTypeID,
				RelyLives: map[string]dot.LiveID{"Conns_": ConnsTypeID, "EtcdConns": EtcdConnsTypeID},
			},
		},
	}

	lives := EtcdConnsTypeLives()
	lives = append(lives, tl)
	lives = append(lives, ConnsTypeLives())

	return lives
}

//jayce edit
//return config of ConnName
func ConnNameConfigTypeLives() *dot.ConfigTypeLives {
	return &dot.ConfigTypeLives{
		TypeIDConfig: ConnNameTypeID,
		ConfigInfo:   &configName{},
	}
}

func (c *ConnName) AfterAllInject(l dot.Line) {
	if c.conn == nil {
		if c.EtcdConns != nil { //优先使用 etcd中的名字
			c.conn = c.EtcdConns.ClientConn(c.serviceName)
		}
		if c.conn == nil {
			c.conn = c.Conns_.ClientConn(c.serviceName)
		}
	}
}

func (c *ConnName) Conn() *grpc.ClientConn {
	return c.conn
}

func (c *ConnName) ClientContext() *ClientContext {
	var cc *ClientContext
	if c.Conns_ != nil {
		cc = c.Conns_.ClientContext(c.serviceName)
	}
	return cc
}

func (c *ConnName) ServerName() string {
	return c.serviceName
}

func NewTestConnName(conn Conns, name string) *ConnName {
	re := &ConnName{
		serviceName: name,
		conn:        nil,
		Conns_:      conn,
	}

	re.AfterAllInject(nil)

	return re
}
