package rpcdot

import (
	"github.com/scryinfo/dot/line/etcddot"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClientEtcdConfig struct {
	Name                     string `json:"name" toml:"name" yaml:"name"`
	WithDefaultServiceConfig string `json:"withDefaultServiceConfig" toml:"withDefaultServiceConfig" yaml:"withDefaultServiceConfig"`
}

type GrpcClientEtcd struct {
	config         *GrpcClientEtcdConfig
	etcdClient     *etcddot.Client
	grpcClientConn *grpc.ClientConn
}

func NewGrpcClientEtcd(config *GrpcClientEtcdConfig, etcdClient *etcddot.Client) *GrpcClientEtcd {
	if len(config.WithDefaultServiceConfig) < 1 {
		config.WithDefaultServiceConfig = `{"loadBalancingConfig": [{"round_robin": {}}]}`
	}
	d := &GrpcClientEtcd{
		config:     config,
		etcdClient: etcdClient,
	}
	d.grpcClientConn = d.makeClientConn()

	return d
}

func (p *GrpcClientEtcd) makeClientConn() *grpc.ClientConn {
	conn, err := grpc.NewClient(
		p.config.Name,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(p.config.WithDefaultServiceConfig),
	)
	if err != nil {
		panic(err)
	}
	return conn
}

func (p *GrpcClientEtcd) Client() *grpc.ClientConn {
	return p.grpcClientConn
}
