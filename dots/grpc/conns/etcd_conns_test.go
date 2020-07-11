package conns

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/naming"
	"testing"
)

func TestEtcdConns_ClientConn(t *testing.T) {
	//first run the etcd 3.4.9 on default port 2379 , (2380 is for peer)
	const addr = "127.0.0.1:100"
	etcdConns := NewEtcd([]string{"127.0.0.1:2379"}, []string{"test"})
	etcdConns.RegisterServer("test", addr)
	re, err := etcdConns.etcdClient.Get(context.TODO(), "test/"+addr)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, len(re.Kvs) > 0)
	v := &naming.Update{}
	err = json.Unmarshal(re.Kvs[0].Value, v)
	assert.Equal(t, nil, err)
	assert.Equal(t, addr, v.Addr)
	_ = re
	cl := etcdConns.ClientConn("test")
	assert.NotEqual(t, nil, cl)
}
