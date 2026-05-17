package rpcdot

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/certificate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClientConfig struct {
	ServerAddress string
	Tls           TlsConfig
}

type GrpcClientEx struct {
	conn *grpc.ClientConn
}

func NewGrpcClientEx(config *GrpcClientConfig, sconf dot.SConfig, logger *dot.LoggerType, baseCert *certificate.BaseCertificate) (*GrpcClientEx, func(), error) {
	d := &GrpcClientEx{}
	switch config.Tls.Mode {
	case RpcTlsNone:
		conn, err := grpc.NewClient(config.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, nil, err
		}
		d.conn = conn
	case RpcTlsInsecure:
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		conn, err := grpc.NewClient(config.ServerAddress, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
		if err != nil {
			return nil, nil, err
		}
		d.conn = conn
	case RpcTlsSecure:
		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
		}
		{

			pool := x509.NewCertPool()
			if config.Tls.RootCert != "" {
				rootCertFile, err := sconf.FullPath(config.Tls.RootCert)
				if err != nil {
					return nil, nil, fmt.Errorf("failed to get root cert path: %w", err)
				}
				rootCert, err := baseCert.LoadCertificate(rootCertFile)
				if err != nil {
					return nil, nil, fmt.Errorf("failed to load root cert: %w", err)
				}
				pool.AddCert(rootCert)
			}
			peerCertFile, err := sconf.FullPath(config.Tls.PeerCert)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get peer cert path: %w", err)
			}
			if peerCertFile == "" {
				return nil, nil, fmt.Errorf("peer cert is required")
			}
			peerCert, err := baseCert.LoadCertificate(peerCertFile)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to load peer cert: %w", err)
			}
			pool.AddCert(peerCert)

			tlsConfig.RootCAs = pool
			tlsConfig.ServerName = baseCert.ServerName(peerCert)
			if tlsConfig.ServerName == "" {
				return nil, nil, fmt.Errorf("cant get server name from peer certificate")
			}
		}
		conn, err := grpc.NewClient(config.ServerAddress, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
		if err != nil {
			return nil, nil, err
		}
		d.conn = conn
	case RpcTlsBoth:
		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
		}
		conn, err := grpc.NewClient(config.ServerAddress, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
		if err != nil {
			return nil, nil, err
		}
		d.conn = conn
		//todo
		return nil, nil, fmt.Errorf("dont implement the tls both")
	default:
		conn, err := grpc.NewClient(config.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
