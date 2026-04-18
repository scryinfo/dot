// Scry Info.  All rights reserved.
// license that can be found in the license file.

package conns

import (
	"google.golang.org/grpc"
)

type ConnName struct {
	serviceName string
	conn        *grpc.ClientConn
	Conns_      *ConnsImp
	//EtcdConns   *EtcdConns `dot:"?"` //可选组件，优先使用etcd
}

type ConfigName struct {
	Name string `json:"name"`
}

func NewConnName(conf *ConfigName, conns *ConnsImp) *ConnName {
	d := &ConnName{
		serviceName: conf.Name,
		Conns_:      conns,
	}
	d.conn = conns.ClientConn(d.serviceName)
	return d
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

func NewTestConnName(conn *ConnsImp, name string) *ConnName {
	return NewConnName(&ConfigName{Name: name}, conn)
}
