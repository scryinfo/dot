package etcddot

import (
	"net/url"
	"time"

	"github.com/scryinfo/dot/dot"
	"go.etcd.io/etcd/server/v3/embed"
)

type ServerConfig struct {
	Name                string   `json:"name" toml:"name" yaml:"name"`
	Dir                 string   `json:"dir" toml:"dir" yaml:"dir"`
	ListenClientUrls    []string `json:"listenClientUrls" toml:"listenClientUrls" yaml:"listenClientUrls"`
	AdvertiseClientUrls []string `json:"advertiseClientUrls" toml:"advertiseClientUrls" yaml:"advertiseClientUrls"`
	ListenPeerUrls      []string `json:"listenPeerUrls" toml:"listenPeerUrls" yaml:"listenPeerUrls"`
	AdvertisePeerUrls   []string `json:"advertisePeerUrls" toml:"advertisePeerUrls" yaml:"advertisePeerUrls"`
	// the unique token for the cluster
	InitialClusterToken string `json:"initialClusterToken" toml:"initialClusterToken" yaml:"initialClusterToken"`
	InitialCluster      string `json:"initialCluster" toml:"initialCluster" yaml:"initialCluster"`
	// debug, info, warn, error, panic, or fatal. Default 'info'
	LogLevel string `json:"logLevel" toml:"logLevel" yaml:"logLevel"`
}

func NewServer(conf *ServerConfig, logger *dot.LoggerType) (*Server, func(), error) {
	cfgEtcd := embed.NewConfig()
	cfgEtcd.Name = conf.Name
	if len(conf.Dir) < 1 {
		conf.Dir = "default.etcd"
	}
	cfgEtcd.Dir = conf.Dir
	cfgEtcd.ListenClientUrls = make([]url.URL, 0, len(conf.ListenClientUrls))
	for _, u := range conf.ListenClientUrls {
		clientUrl, err := url.Parse(u)
		if err != nil {
			logger.Error().Err(err).Send()
			return nil, nil, err
		}
		cfgEtcd.ListenClientUrls = append(cfgEtcd.ListenClientUrls, *clientUrl)
	}
	cfgEtcd.AdvertiseClientUrls = make([]url.URL, 0, len(conf.AdvertiseClientUrls))
	for _, u := range conf.AdvertiseClientUrls {
		clientUrl, err := url.Parse(u)
		if err != nil {
			logger.Error().Err(err).Send()
			return nil, nil, err
		}
		cfgEtcd.AdvertiseClientUrls = append(cfgEtcd.AdvertiseClientUrls, *clientUrl)
	}
	cfgEtcd.ListenPeerUrls = make([]url.URL, 0, len(conf.ListenPeerUrls))
	for _, u := range conf.ListenPeerUrls {
		peerUrl, err := url.Parse(u)
		if err != nil {
			logger.Error().Err(err).Send()
			return nil, nil, err
		}
		cfgEtcd.ListenPeerUrls = append(cfgEtcd.ListenPeerUrls, *peerUrl)
	}
	cfgEtcd.AdvertisePeerUrls = make([]url.URL, 0, len(conf.AdvertisePeerUrls))
	for _, u := range conf.AdvertisePeerUrls {
		peerUrl, err := url.Parse(u)
		if err != nil {
			logger.Error().Err(err).Send()
			return nil, nil, err
		}
		cfgEtcd.AdvertisePeerUrls = append(cfgEtcd.AdvertisePeerUrls, *peerUrl)
	}
	cfgEtcd.InitialClusterToken = conf.InitialClusterToken
	cfgEtcd.InitialCluster = conf.InitialCluster
	if len(conf.LogLevel) > 0 {
		cfgEtcd.LogLevel = conf.LogLevel
	}

	etcdServer, err := embed.StartEtcd(cfgEtcd)
	if err != nil {
		logger.Error().Err(err).Send()
		return nil, nil, err
	}

	select {
	case <-etcdServer.Server.ReadyNotify():
		logger.Info().Msg("etcd server is ready")
	case <-time.After(10 * time.Second):
		logger.Error().Msg("etcd server did not become ready within 10 seconds")
		etcdServer.Server.Stop()
		return nil, nil, nil
	}

	d := Server{
		conf:    *conf,
		cfgEtcd: cfgEtcd,
		etct:    etcdServer,
		logger:  logger,
	}
	return &d, func() {
		if d.etct != nil {
			d.etct.Close()
			d.etct = nil
		}
	}, nil
}

type Server struct {
	conf    ServerConfig
	cfgEtcd *embed.Config
	etct    *embed.Etcd
	logger  *dot.LoggerType
}
