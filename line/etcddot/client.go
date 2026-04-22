package etcddot

import (
	"time"

	"github.com/scryinfo/dot/dot"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func NewClient(conf *ClientConfig) (*Client, func(), error) {

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   conf.Endpoints,
		DialTimeout: conf.DialTimeout,
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
	Endpoints   []string
	DialTimeout time.Duration
}

type Client struct {
	conf   ClientConfig
	client *clientv3.Client
}
