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
	d.register()
	return d, func() {
		d.unregister()
	}, nil
}

type ConnectServerEtcdConfig struct {
	// the name of the server, used for etcd discovery
	Name string `json:"name" toml:"name" yaml:"name"`
	// the address of the server, used for etcd discovery
	Address string `json:"address" toml:"address" yaml:"address"`
	// the ttl of the lease, in seconds,
	// default is 10 seconds, if the ttl is <= 0
	Ttl int64 `json:"ttl" toml:"ttl" yaml:"ttl"`
}

type ConnectServerEtcd struct {
	config     *ConnectServerEtcdConfig
	etcdClient *etcddot.Client
	server     *ConnectServer
	logger     *dot.LoggerType
	leaseId    *clientv3.LeaseID
}

func (p *ConnectServerEtcd) register() {
	lease := clientv3.NewLease(p.etcdClient.EtcdClient())
	leaseResp, err := lease.Grant(context.Background(), p.config.Ttl)
	if err != nil {
		p.logger.Error().Err(err).Send()
		return

	}
	p.leaseId = new(leaseResp.ID)
	key := fmt.Sprintf("%s/%s", p.config.Name, p.config.Address)
	_, err = p.etcdClient.EtcdClient().Put(context.Background(), key, p.config.Address, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		p.logger.Error().Err(err).Send()
		return
	}
	p.logger.Info().Msgf("register server: %s, addr: %s, ttl: %d", p.config.Name, p.config.Address, p.config.Ttl)
	_, err = lease.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		p.logger.Error().Err(err).Send()
		return
	}

}

func (p *ConnectServerEtcd) unregister() {
	if p.leaseId != nil {
		_, err := p.etcdClient.EtcdClient().Revoke(context.Background(), *p.leaseId)
		if err != nil {
			p.logger.Error().Err(err).Send()
		} else {
			p.logger.Info().Msgf("unregister server: %s, addr: %s", p.config.Name, p.config.Address)
		}
	}
}
