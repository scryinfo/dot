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
	ForceAttemptHTTP2   bool `json:"forceAttemptHTTP2" toml:"forceAttemptHTTP2" yaml:"forceAttemptHTTP2"`
	DisableCompression  bool `json:"disableCompression" toml:"disableCompression" yaml:"disableCompression"`
	MaxIdleConns        int  `json:"maxIdleConns" toml:"maxIdleConns" yaml:"maxIdleConns"`
	MaxIdleConnsPerHost int  `json:"maxIdleConnsPerHost" toml:"maxIdleConnsPerHost" yaml:"maxIdleConnsPerHost"`
	MaxConnsPerHost     int  `json:"maxConnsPerHost" toml:"maxConnsPerHost" yaml:"maxConnsPerHost"`
	// sample "http://localhost:8089"
	ServerAddress string `json:"serverAddress" toml:"serverAddress" yaml:"serverAddress"`

	Name                     string `json:"name" toml:"name" yaml:"name"`
	WithDefaultServiceConfig string `json:"withDefaultServiceConfig" toml:"withDefaultServiceConfig" yaml:"withDefaultServiceConfig"`
	Tls                      TlsConfig
}

func NewHttpClientEtcd(config *HttpClientConfigEtcd, sconf dot.SConfig, baseCert *certificate.BaseCertificate, etcdClient *etcddot.Client, logger *dot.LoggerType) (*HttpClientEtcd, error) {
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
	tlsConfig, err := conf.Tls.MakeTlsConfig(sconf, baseCert)
	if err != nil {
		return nil, err
	}
	tr.TLSClientConfig = tlsConfig

	return &HttpClientEtcd{
		client: http.Client{
			Transport: tr,
		},
		logger: logger,
		conf:   conf,
	}, nil
}

type HttpClientEtcd struct {
	client http.Client
	logger *dot.LoggerType
	conf   HttpClientConfigEtcd
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
