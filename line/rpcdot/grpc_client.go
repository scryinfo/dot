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

type RpcTls string

const (
	RpcTlsNone     RpcTls = "none"
	RpcTlsInsecure RpcTls = "insecure" //InsecureSkipVerify = true
	RpcTlsSecure   RpcTls = "secure"   //InsecureSkipVerify = false, read name from server cert
	RpcTlsBoth     RpcTls = "both"     // client and server cert are verified
)

type TlsConfig struct {
	Mode       RpcTls
	Cert       string
	Key        string
	RootCert   string
	ServerCert string
}

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
			serverCertFile, err := sconf.FullPath(config.Tls.ServerCert)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get server cert path: %w", err)
			}
			if serverCertFile == "" {
				return nil, nil, fmt.Errorf("server cert is required")
			}
			serverCert, err := baseCert.LoadCertificate(serverCertFile)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to load server cert: %w", err)
			}
			pool.AddCert(serverCert)

			tlsConfig.RootCAs = pool
			tlsConfig.ServerName = baseCert.ServerName(serverCert)
			if tlsConfig.ServerName == "" {
				return nil, nil, fmt.Errorf("cant get server name from server certificate")
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
