// Scry Info.  All rights reserved.
// license that can be found in the license file.

package lb

import (
	"google.golang.org/grpc/resolver"
)

//客户端的负载均衡
type clientBuilder struct {
	scheme       string
	serviceAddrs map[string][]string //key 服务名， value 服务所对应的地址（如 12.23.23.23：909）
}

func NewClientBuilder(schema string, serviceAddrs map[string][]string) resolver.Builder {
	return &clientBuilder{
		scheme:       schema,
		serviceAddrs: serviceAddrs,
	}
}

func (c *clientBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	r := &clientResolver{
		target:     target,
		conn:       cc,
		addrsStore: c.serviceAddrs,
	}
	r.start()
	return r, nil
}
func (c *clientBuilder) Scheme() string { return c.scheme }

type clientResolver struct {
	target     resolver.Target
	conn       resolver.ClientConn
	addrsStore map[string][]string
}

func (c clientResolver) ResolveNow(resolver.ResolveNowOption) {
}

func (c clientResolver) Close() {
}

func (c *clientResolver) start() {
	addrStrs := c.addrsStore[c.target.Endpoint]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	c.conn.UpdateState(resolver.State{Addresses: addrs})
}
