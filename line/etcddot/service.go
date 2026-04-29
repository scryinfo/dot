package etcddot

import (
	"github.com/scryinfo/dot/dot"
	"go.etcd.io/etcd/server/v3/embed"
)

type ServerConfig struct {
	embed.Config
}

func NewServer(conf *ServerConfig) (*Server, func(), error) {
	e, err := embed.StartEtcd(&conf.Config)
	if err != nil {
		dot.Logger.Error().Err(err).Send()
		return nil, nil, err
	}
	d := Server{
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

type Server struct {
	conf ServerConfig
	etct *embed.Etcd
}
