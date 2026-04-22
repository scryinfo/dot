package etcddot

import (
	"github.com/scryinfo/dot/dot"
	"go.etcd.io/etcd/embed"
)

type ServiceConfig struct {
	embed.Config
}

func NewService(conf *ServiceConfig) (*Service, func(), error) {
	e, err := embed.StartEtcd(&conf.Config)
	if err != nil {
		dot.Logger.Error().Err(err).Send()
		return nil, nil, err
	}
	d := Service{
		conf: *conf,
		etct: e,
	}
	return &d, func() {
		if d.etct != nil {
			d.etct.Close()
			d.etct = nil
		}
	}, nil
}

type Service struct {
	conf ServiceConfig
	etct *embed.Etcd
}
