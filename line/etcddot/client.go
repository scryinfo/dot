package etcddot

import (
	"context"
	"time"

	"github.com/scryinfo/dot/dot"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func NewClient(conf *ClientConfig, logger *dot.LoggerType) (*Client, func(), error) {

	if conf.DialTimeout < 0 {
		conf.DialTimeout = 10
	}
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   conf.Endpoints,
		DialTimeout: conf.DialTimeoutDuration(),
		// DialOptions: []grpc.DialOption{
		// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
		// 	grpc.WithDefaultServiceConfig(`{}`),
		// },
	})
	if err != nil {
		logger.Error().Err(err).Send()
		return nil, nil, err
	}
	d := Client{
		conf:   conf,
		client: client,
		logger: logger,
	}
	return &d, func() {
		if d.client != nil {
			err := d.client.Close()
			if err != nil {
				logger.Error().Err(err).Send()
			}
			d.client = nil
		}
	}, nil
}

type ClientConfig struct {
	Endpoints []string `json:"endpoints" toml:"endpoints" yaml:"endpoints"`
	// milli second count
	DialTimeout int32 `json:"dial_timeout" toml:"dial_timeout" yaml:"dial_timeout"`
}

type Client struct {
	conf   *ClientConfig
	client *clientv3.Client
	logger *dot.LoggerType
}

func (p *Client) EtcdClient() *clientv3.Client {
	return p.client
}

func (p *ClientConfig) DialTimeoutDuration() time.Duration {
	return time.Duration(p.DialTimeout) * time.Millisecond
}

// ping server
// returns nil if ping is successful, error otherwise
func (p *Client) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), p.conf.DialTimeoutDuration())
	defer cancel()
	_, err := p.client.Status(ctx, p.client.Endpoints()[0])
	if err != nil {
		p.logger.Error().Err(err).Send()
		return err
	}
	return nil
}
