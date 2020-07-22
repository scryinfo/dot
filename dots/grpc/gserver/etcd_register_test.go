package gserver

import (
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/scryinfo/dot/dots/grpc/conns"
	"github.com/scryinfo/dot/dots/line"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/naming"
	"testing"
)

var controller *gomock.Controller

// generate the mock
// cd to folder gserver, run "mockgen -source ./server_nobl.go -destination ./server_nobl_mock.go -package gserver"

func TestEtcdConns_EtcdRegister(t *testing.T) {
	//first run the etcd 3.4.9 on default port 2379 , (2380 is for peer)

	const (
		addr = "127.0.0.1:100"
		name = "test"
	)
	etcdConns := conns.NewEtcd([]string{"127.0.0.1:2379"}, []string{name})
	{
		l, err := line.BuildAndStart(nil)
		assert.Equal(t, nil, err)
		l.ToInjecter().ReplaceOrAddByLiveID(etcdConns, conns.EtcdConnsTypeID)
		{
			controller := gomock.NewController(t)
			s := NewMockServerNobl(controller)
			s.EXPECT().ServerItem().Return(ServerItem{
				Name:  name,
				Addrs: []string{addr},
			})

			l.ToInjecter().ReplaceOrAddByLiveID(s, GinNoblTypeID)
		}

		etcdRegister := &EtcdRegister{EtcdConns: etcdConns}
		l.ToInjecter().ReplaceOrAddByLiveID(etcdRegister, EtcdRegisterTypeID)
		etcdRegister.AfterAllInject(l)
		etcdRegister.AfterAllStart()
	}

	{
		re, err := etcdConns.EtcdClient().Get(context.TODO(), name+"/"+addr)
		assert.Equal(t, nil, err)
		assert.Equal(t, true, len(re.Kvs) > 0)
		v := &naming.Update{}
		err = json.Unmarshal(re.Kvs[0].Value, v)
		assert.Equal(t, nil, err)
		assert.Equal(t, addr, v.Addr)
		cl := etcdConns.ClientConn(name)
		assert.NotEqual(t, nil, cl)
	}

	{
		err := etcdConns.UnRegisterServer(name, addr)
		assert.Equal(t, nil, err)

		etcdConns.RegisterServer(name, addr)
		re, err := etcdConns.EtcdClient().Get(context.TODO(), name+"/"+addr)
		assert.Equal(t, nil, err)
		assert.Equal(t, true, len(re.Kvs) > 0)
		v := &naming.Update{}
		err = json.Unmarshal(re.Kvs[0].Value, v)
		assert.Equal(t, nil, err)
		assert.Equal(t, addr, v.Addr)
	}

	if controller != nil {
		controller.Finish()
		controller = nil
	}
}
