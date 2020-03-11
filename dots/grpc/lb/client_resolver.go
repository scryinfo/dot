// Scry Info.  All rights reserved.
// license that can be found in the license file.

package lb

import (
	"google.golang.org/grpc/resolver"
)

//Client load balancing
type clientBuilder struct {
	scheme       string
	serviceAddrs map[string][]string //key service name, value service corresponding address(such as 12.23.23.23：909）
}

func NewClientBuilder(schema string, serviceAddrs map[string][]string) resolver.Builder {
	return &clientBuilder{
		scheme:       schema,
		serviceAddrs: serviceAddrs,
	}
}

func (c *clientBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
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

func (c clientResolver) ResolveNow(resolver.ResolveNowOptions) {
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
