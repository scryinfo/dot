package rpcdot

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"github.com/scryinfo/dot/dot"
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
}

func NewHttpClientEx(config *HttpClientConfig, logger *dot.LoggerType) (*HttpClientEx, error) {
	tr := &http.Transport{
		ForceAttemptHTTP2:   config.ForceAttemptHTTP2,
		DisableCompression:  config.DisableCompression,
		MaxIdleConns:        config.MaxIdleConns,
		MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
		MaxConnsPerHost:     config.MaxConnsPerHost,
	}
	err := http2.ConfigureTransport(tr)
	if err != nil {
		return nil, err
	}

	return &HttpClientEx{
		client: http.Client{
			Transport: tr,
		},
		logger: logger,
		conf:   *config,
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
