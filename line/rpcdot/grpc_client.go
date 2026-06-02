package rpcdot

import (
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/certificate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClientConfig struct {
	// sample: localhost:8090,
	// no https/http:
	ServerAddress string    `toml:"server_address" yaml:"server_address" json:"server_address"`
	Tls           TlsConfig `toml:"tls" yaml:"tls" json:"tls"`
}

type GrpcClientEx struct {
	conn *grpc.ClientConn
}

func NewGrpcClientEx(config *GrpcClientConfig, sconf dot.SConfig, logger *dot.LoggerType, baseCert *certificate.BaseCertificate) (*GrpcClientEx, func(), error) {
	d := &GrpcClientEx{}
	err := config.Tls.FullPath(sconf)
	if err != nil {
		return nil, nil, err
	}
	tlsConfig, err := config.Tls.MakeTlsConfig(sconf, baseCert)
	if err != nil {
		return nil, nil, err
	}
	if tlsConfig == nil {
		conn, err := grpc.NewClient(config.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, nil, err
		}
		d.conn = conn
	} else {
		conn, err := grpc.NewClient(config.ServerAddress, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
		if err != nil {
			return nil, nil, err
		}
		d.conn = conn
	}

	return d, func() {
		if d.conn != nil {
			err := d.conn.Close()
			if err != nil {
				logger.Error().Err(err).Send()
			}
		}
	}, nil
}

func (g *GrpcClientEx) Client() *grpc.ClientConn {
	return g.conn
}
