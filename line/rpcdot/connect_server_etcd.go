package rpcdot

import (
	"context"
	"fmt"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/etcddot"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func NewConnectServerEtcd(config *ConnectServerEtcdConfig, etcdClient *etcddot.Client, server *ConnectServer, logger *dot.LoggerType) (*ConnectServerEtcd, func(), error) {
	if config.Ttl <= 0 {
		config.Ttl = 10
	}
	d := &ConnectServerEtcd{
		config:     config,
		etcdClient: etcdClient,
		server:     server,
		logger:     logger,
	}
	err := d.register()
	if err != nil {
		return nil, nil, err
	}
	return d, func() {
		d.unregister()
	}, nil
}

type ConnectServerEtcdConfig struct {
	// the name of the server, used for etcd discovery
	Name string `json:"name" toml:"name" yaml:"name" mapstructure:"name"`
	// the address of the server, used for etcd discovery
	Addr string `json:"addr" toml:"addr" yaml:"addr" mapstructure:"addr"`
	// the ttl of the lease, in seconds,
	// default is 10 seconds, if the ttl is <= 0
	Ttl int64 `json:"ttl" toml:"ttl" yaml:"ttl" mapstructure:"ttl"`
}

type ConnectServerEtcd struct {
	config     *ConnectServerEtcdConfig
	etcdClient *etcddot.Client
	server     *ConnectServer
	logger     *dot.LoggerType
	leaseId    *clientv3.LeaseID
}

func (p *ConnectServerEtcd) register() error {
	lease := clientv3.NewLease(p.etcdClient.EtcdClient())
	leaseResp, err := lease.Grant(context.Background(), p.config.Ttl)
	if err != nil {
		p.logger.Error().Err(err).Send()
		return err

	}
	p.leaseId = new(leaseResp.ID)
	key := fmt.Sprintf("%s/%s", p.config.Name, p.config.Addr)
	addrValue := fmt.Sprintf(`{"Addr":"%s"}`, p.config.Addr)
	_, err = p.etcdClient.EtcdClient().Put(context.Background(), key, addrValue, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		p.logger.Error().Err(err).Send()
		return err
	}
	p.logger.Info().Msgf("register server: %s, addr: %s, ttl: %d", p.config.Name, p.config.Addr, p.config.Ttl)
	_, err = lease.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		p.logger.Error().Err(err).Send()
		return err
	}

	return nil
}

func (p *ConnectServerEtcd) unregister() {
	if p.leaseId != nil {
		_, err := p.etcdClient.EtcdClient().Revoke(context.Background(), *p.leaseId)
		if err != nil {
			p.logger.Error().Err(err).Send()
		} else {
			p.logger.Info().Msgf("unregister server: %s, addr: %s", p.config.Name, p.config.Addr)
		}
	}
}
