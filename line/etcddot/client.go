package etcddot

import (
	"time"

	"github.com/scryinfo/dot/dot"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func NewClient(conf *ClientConfig) (*Client, func(), error) {

	if conf.DialTimeout < 0 {
		conf.DialTimeout = 10
	}
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   conf.Endpoints,
		DialTimeout: conf.DialTimeoutDuration(),
	})
	if err != nil {
		dot.Logger.Error().Err(err).Send()
		return nil, nil, err
	}
	d := Client{
		conf:   *conf,
		client: client,
	}
	return &d, func() {
		if d.client != nil {
			err := d.client.Close()
			if err != nil {
				dot.Logger.Error().Err(err).Send()
			}
			d.client = nil
		}
	}, nil
}

type ClientConfig struct {
	Endpoints []string
	// milli second count
	DialTimeout int32
}

type Client struct {
	conf   ClientConfig
	client *clientv3.Client
}

func (p *Client) EtcdClient() *clientv3.Client {
	return p.client
}

func (p *ClientConfig) DialTimeoutDuration() time.Duration {
	return time.Duration(p.DialTimeout) * time.Millisecond
}
