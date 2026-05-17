package rpcdot

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/certificate"
	"golang.org/x/net/http2"
)

type HttpClientConfig struct {
	ForceAttemptHTTP2   bool `json:"forceAttemptHTTP2" toml:"forceAttemptHTTP2" yaml:"forceAttemptHTTP2"`
	DisableCompression  bool `json:"disableCompression" toml:"disableCompression" yaml:"disableCompression"`
	MaxIdleConns        int  `json:"maxIdleConns" toml:"maxIdleConns" yaml:"maxIdleConns"`
	MaxIdleConnsPerHost int  `json:"maxIdleConnsPerHost" toml:"maxIdleConnsPerHost" yaml:"maxIdleConnsPerHost"`
	MaxConnsPerHost     int  `json:"maxConnsPerHost" toml:"maxConnsPerHost" yaml:"maxConnsPerHost"`
	// sample "http://localhost:8089"
	ServerAddress string `json:"serverAddress" toml:"serverAddress" yaml:"serverAddress"`

	Tls TlsConfig
}

func NewHttpClientEx(config *HttpClientConfig, sconf dot.SConfig, baseCert *certificate.BaseCertificate, logger *dot.LoggerType) (*HttpClientEx, error) {
	conf := *config
	err := conf.Tls.FullPath(sconf)
	if err != nil {
		return nil, err
	}

	tr := &http.Transport{
		ForceAttemptHTTP2:   conf.ForceAttemptHTTP2,
		DisableCompression:  conf.DisableCompression,
		MaxIdleConns:        conf.MaxIdleConns,
		MaxIdleConnsPerHost: conf.MaxIdleConnsPerHost,
		MaxConnsPerHost:     conf.MaxConnsPerHost,
	}
	if !conf.Tls.NeedsTls() {
		// add support http2 not tls
		err := http2.ConfigureTransport(tr)
		if err != nil {
			return nil, err
		}
	}

	switch conf.Tls.Mode {
	case RpcTlsNone:
		tr.TLSClientConfig = nil
	case RpcTlsInsecure:
		tr.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	case RpcTlsSecure:
		tr.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: false,
		}
		{
			pool := x509.NewCertPool()
			if conf.Tls.RootCert != "" {
				rootCertFile, err := sconf.FullPath(conf.Tls.RootCert)
				if err != nil {
					return nil, fmt.Errorf("failed to get root cert path: %w", err)
				}
				rootCert, err := baseCert.LoadCertificate(rootCertFile)
				if err != nil {
					return nil, fmt.Errorf("failed to load root cert: %w", err)
				}
				pool.AddCert(rootCert)
			}
			peerCertFile, err := sconf.FullPath(conf.Tls.PeerCert)
			if err != nil {
				return nil, fmt.Errorf("failed to get peer cert path: %w", err)
			}
			if peerCertFile == "" {
				return nil, fmt.Errorf("peer cert is required")
			}
			peerCert, err := baseCert.LoadCertificate(peerCertFile)
			if err != nil {
				return nil, fmt.Errorf("failed to load peer cert: %w", err)
			}
			pool.AddCert(peerCert)

			tr.TLSClientConfig.RootCAs = pool
			tr.TLSClientConfig.ServerName = baseCert.ServerName(peerCert)
			if tr.TLSClientConfig.ServerName == "" {
				return nil, fmt.Errorf("cant get server name from peer certificate")
			}
		}
	case RpcTlsBoth:
		tr.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: false,
		}
		//todo
		return nil, fmt.Errorf("dont implement the tls both")
	default:
		tr.TLSClientConfig = nil
	}

	return &HttpClientEx{
		client: http.Client{
			Transport: tr,
		},
		logger: logger,
		conf:   conf,
	}, nil
}

type HttpClientEx struct {
	client http.Client
	logger *dot.LoggerType
	conf   HttpClientConfig
}

func (p *HttpClientEx) NotCompressOptions() []connect.ClientOption {
	return []connect.ClientOption{
		connect.WithGRPC(),
		connect.WithSendCompression(""),
		connect.WithAcceptCompression("", nil, nil),
		connect.WithInterceptors(disableGRPCCompressionInterceptor()),
	}
}

func (p *HttpClientEx) ServerAddress() string {
	return p.conf.ServerAddress
}

func (p *HttpClientEx) Client() *http.Client {
	return &p.client
}

func disableGRPCCompressionInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			req.Header().Del("Grpc-Encoding")
			req.Header().Del("Grpc-Accept-Encoding")
			return next(ctx, req)
		}
	}
}
