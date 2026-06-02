package rpcdot

// todo
import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/certificate"
	"github.com/scryinfo/dot/line/etcddot"
	"golang.org/x/net/http2"
)

type HttpClientConfigEtcd struct {
	ForceAttemptHTTP2   bool `json:"force_attempt_http2" toml:"force_attempt_http2" yaml:"force_attempt_http2"`
	DisableCompression  bool `json:"disable_compression" toml:"disable_compression" yaml:"disable_compression"`
	MaxIdleConns        int  `json:"max_idle_conns" toml:"max_idle_conns" yaml:"max_idle_conns"`
	MaxIdleConnsPerHost int  `json:"max_idle_conns_per_host" toml:"max_idle_conns_per_host" yaml:"max_idle_conns_per_host"`
	MaxConnsPerHost     int  `json:"max_conns_per_host" toml:"max_conns_per_host" yaml:"max_conns_per_host"`
	// sample "http://localhost:8089"
	ServerAddress string `json:"server_address" toml:"server_address" yaml:"server_address"`

	Name                     string    `json:"name" toml:"name" yaml:"name"`
	WithDefaultServiceConfig string    `json:"with_default_service_config" toml:"with_default_service_config" yaml:"with_default_service_config"`
	Tls                      TlsConfig `json:"tls" toml:"tls" yaml:"tls"`
}

func NewHttpClientEtcd(config *HttpClientConfigEtcd, sconf dot.SConfig, baseCert *certificate.BaseCertificate, etcdClient *etcddot.Client, logger *dot.LoggerType) (*HttpClientEtcd, error) {
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

	return &HttpClientEtcd{
		client: http.Client{
			Transport: tr,
		},
		logger: logger,
		conf:   config,
	}, nil
}

type HttpClientEtcd struct {
	client http.Client
	logger *dot.LoggerType
	conf   *HttpClientConfigEtcd
}

func (p *HttpClientEtcd) NotCompressOptions() []connect.ClientOption {
	return []connect.ClientOption{
		connect.WithGRPC(),
		connect.WithSendCompression(""),
		connect.WithAcceptCompression("", nil, nil),
		connect.WithInterceptors(disableGRPCCompressionInterceptor()),
	}
}

func (p *HttpClientEtcd) ServerAddress() string {
	return p.conf.ServerAddress
}

func (p *HttpClientEtcd) Client() *http.Client {
	return &p.client
}
