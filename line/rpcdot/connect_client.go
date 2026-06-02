package rpcdot

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/certificate"
	"golang.org/x/net/http2"
)

type HttpClientConfig struct {
	ForceAttemptHTTP2   bool `json:"force_attempt_http2" toml:"force_attempt_http2" yaml:"force_attempt_http2"`
	DisableCompression  bool `json:"disable_compression" toml:"disable_compression" yaml:"disable_compression"`
	MaxIdleConns        int  `json:"max_idle_conns" toml:"max_idle_conns" yaml:"max_idle_conns"`
	MaxIdleConnsPerHost int  `json:"max_idle_conns_per_host" toml:"max_idle_conns_per_host" yaml:"max_idle_conns_per_host"`
	MaxConnsPerHost     int  `json:"max_conns_per_host" toml:"max_conns_per_host" yaml:"max_conns_per_host"`
	// sample "http://localhost:8089"
	ServerAddress string `json:"server_address" toml:"server_address" yaml:"server_address"`

	Tls TlsConfig `json:"tls" toml:"tls" yaml:"tls"`
}

func NewHttpClientEx(config *HttpClientConfig, sconf dot.SConfig, baseCert *certificate.BaseCertificate, logger *dot.LoggerType) (*HttpClientEx, error) {
	err := config.Tls.FullPath(sconf)
	if err != nil {
		return nil, err
	}

	tr := &http.Transport{
		ForceAttemptHTTP2:   config.ForceAttemptHTTP2,
		DisableCompression:  config.DisableCompression,
		MaxIdleConns:        config.MaxIdleConns,
		MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
		MaxConnsPerHost:     config.MaxConnsPerHost,
	}
	if !config.Tls.NeedsTls() {
		// add support http2 not tls
		err := http2.ConfigureTransport(tr)
		if err != nil {
			return nil, err
		}
	}
	tlsConfig, err := config.Tls.MakeTlsConfig(sconf, baseCert)
	if err != nil {
		return nil, err
	}
	tr.TLSClientConfig = tlsConfig

	return &HttpClientEx{
		client: http.Client{
			Transport: tr,
		},
		logger: logger,
		conf:   config,
	}, nil
}

type HttpClientEx struct {
	client http.Client
	logger *dot.LoggerType
	conf   *HttpClientConfig
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
