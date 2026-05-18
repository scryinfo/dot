package rpcdot

import (
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/certificate"
	"github.com/scryinfo/dot/line/etcddot"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClientEtcdConfig struct {
	Name                     string `json:"name" toml:"name" yaml:"name"`
	WithDefaultServiceConfig string `json:"withDefaultServiceConfig" toml:"withDefaultServiceConfig" yaml:"withDefaultServiceConfig"`
	Tls                      TlsConfig
}

type GrpcClientEtcd struct {
	config         GrpcClientEtcdConfig
	etcdClient     *etcddot.Client
	grpcClientConn *grpc.ClientConn
	logger         *dot.LoggerType
}

func NewGrpcClientEtcd(config *GrpcClientEtcdConfig, sconf dot.SConfig, baseCert *certificate.BaseCertificate, etcdClient *etcddot.Client, logger *dot.LoggerType) (*GrpcClientEtcd, error) {
	d := &GrpcClientEtcd{
		config:     *config,
		etcdClient: etcdClient,
		logger:     logger,
	}
	if len(config.WithDefaultServiceConfig) < 1 {
		config.WithDefaultServiceConfig = `{"loadBalancingConfig": [{"round_robin": {}}]}`
	}
	err := d.config.Tls.FullPath(sconf)
	if err != nil {
		return nil, err
	}

	err = d.makeClientConn(sconf, baseCert)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (p *GrpcClientEtcd) makeClientConn(sconf dot.SConfig, baseCert *certificate.BaseCertificate) error {
	etcdResolver, err := resolver.NewBuilder(p.etcdClient.EtcdClient())
	if err != nil {
		p.logger.Error().Err(err).Send()
		return err
	}
	tlsConfig, err := p.config.Tls.MakeTlsConfig(sconf, baseCert)
	if err != nil {
		p.logger.Error().Err(err).Send()
		return err
	}
	cred := insecure.NewCredentials()
	if tlsConfig != nil {
		cred = credentials.NewTLS(tlsConfig)
	}
	conn, err := grpc.NewClient(
		"etcd:///"+p.config.Name,
		grpc.WithTransportCredentials(cred),
		grpc.WithDefaultServiceConfig(p.config.WithDefaultServiceConfig),
		grpc.WithResolvers(etcdResolver),
	)
	if err != nil {
		p.logger.Error().Err(err).Send()
		return err
	}
	p.grpcClientConn = conn
	return nil
}

func (p *GrpcClientEtcd) Client() *grpc.ClientConn {
	return p.grpcClientConn
}
