package rpcdot

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"slices"
	"strings"
	"sync/atomic"
	"time"

	"github.com/scryinfo/dot/dot"
	httptools "github.com/scryinfo/dot/line/rpcdot/http_tools"
)

type ConnectHttpServerMux struct {
	http.ServeMux
}

func NewConnectHttpServerMux() *ConnectHttpServerMux {
	return &ConnectHttpServerMux{http.ServeMux{}}
}

type HandlerMiddle func(w http.ResponseWriter, r *http.Request) error

type ConnectServerConfig struct {
	// true : auto run, false: manual run
	AutoRun bool `json:"auto_run" toml:"auto_run" yaml:"auto_run" mapstructure:"auto_run"`
	// sample ":8080"
	Addr                 string        `json:"addr" toml:"addr" yaml:"addr" mapstructure:"addr"`
	ReadTimeout          time.Duration `json:"read_timeout" toml:"read_timeout" yaml:"read_timeout" mapstructure:"read_timeout"`
	WriteTimeout         time.Duration `json:"write_timeout" toml:"write_timeout" yaml:"write_timeout" mapstructure:"write_timeout"`
	MaxConcurrentStreams int           `json:"max_concurrent_streams" toml:"max_concurrent_streams" yaml:"max_concurrent_streams" mapstructure:"max_concurrent_streams"`
	AllowedOrigins       []string      `json:"allowed_origins" toml:"allowed_origins" yaml:"allowed_origins" mapstructure:"allowed_origins"`
	AllowHeaders         []string      `json:"allow_headers" toml:"allow_headers" yaml:"allow_headers" mapstructure:"allow_headers"`
	AllowMethods         []string      `json:"allow_methods" toml:"allow_methods" yaml:"allow_methods" mapstructure:"allow_methods"`
	AllowCredentials     bool          `json:"allow_credentials" toml:"allow_credentials" yaml:"allow_credentials" mapstructure:"allow_credentials"`
	// dont auth these urls,sample: ["/login", "/rpc/test"]
	// ["*"] dont auth all urls
	UnAuthUrls []string `json:"un_auth_urls" toml:"un_auth_urls" yaml:"un_auth_urls" mapstructure:"un_auth_urls"`
	// if it is true, all OPTIONS requests will be returned ok
	OptionMethods bool `json:"option_methods" toml:"option_methods" yaml:"option_methods" mapstructure:"option_methods"`
	// shutdown timeout
	ShutdownTimeout  time.Duration `json:"shutdown_timeout" toml:"shutdown_timeout" yaml:"shutdown_timeout" mapstructure:"shutdown_timeout"`
	HTTP1            bool          `json:"http1" toml:"http1" yaml:"http1" mapstructure:"http1"`
	HTTP2            bool          `json:"http2" toml:"http2" yaml:"http2" mapstructure:"http2"`
	UnencryptedHTTP2 bool          `json:"unencrypted_http2" toml:"unencrypted_http2" yaml:"unencrypted_http2" mapstructure:"unencrypted_http2"`

	Tls TlsConfig `json:"tls" toml:"tls" yaml:"tls" mapstructure:"tls"`
}

type ConnectServer struct {
	HTTPServer *http.Server
	conf       *ConnectServerConfig
	logger     *dot.LoggerType
	started    atomic.Bool
}

func NewConnetServer(config *ConnectServerConfig, sconf dot.SConfig, connetMux *ConnectHttpServerMux, logger *dot.LoggerType, middle HandlerMiddle) (*ConnectServer, func(), error) {
	if config.ShutdownTimeout < 0 {
		config.ShutdownTimeout = 10 * time.Second
	}
	if len(config.UnAuthUrls) > 0 {
		unauthUrls = config.UnAuthUrls
	}
	if len(config.AllowMethods) < 1 {
		config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	}
	if len(config.AllowHeaders) < 1 {
		config.AllowHeaders = []string{
			"Origin", "Upgrade", "Connection", "X-Requested-With", "X-HTTP-Protocol", "Content-Type",
			"Accept", "Cookie", "connect-protocol-version", "connect-timeout-ms",
			httptools.TokenName, httptools.TokenGame, httptools.Authorization, httptools.AuthorizationGame}
	}
	if len(config.AllowedOrigins) < 1 {
		config.AllowedOrigins = []string{"http://localhost", "http://127.0.0.1"}
	}

	allowMethods := strings.Join(config.AllowMethods, ",")
	allowHeaders := strings.Join(config.AllowHeaders, ",")
	muxEx := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		{
			originHeader := r.Header.Get("Origin")
			if len(config.AllowedOrigins) == 1 && config.AllowedOrigins[0] == "*" {
				header.Set("Access-Control-Allow-Origin", "*")
			} else if slices.Contains(config.AllowedOrigins, originHeader) {
				header.Set("Access-Control-Allow-Origin", originHeader)
			} else if slices.ContainsFunc(config.AllowedOrigins, func(origin string) bool {
				if strings.HasPrefix(origin, "http://localhost") || strings.HasPrefix(origin, "http://127.0.0.1") ||
					strings.HasPrefix(origin, "https://localhost") || strings.HasPrefix(origin, "https://127.0.0.1") {
					return true
				}
				return false
			}) {
				header.Set("Access-Control-Allow-Origin", originHeader)
			}
		}
		header.Set("Access-Control-Allow-Methods", allowMethods)
		header.Set("Access-Control-Allow-Headers", allowHeaders)
		header.Set("Access-Control-Expose-Headers", "X-Protocol, X-Response-Time")
		if config.AllowCredentials {
			header.Set("Access-Control-Allow-Credentials", "true")
		}
		if config.OptionMethods && r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if middle != nil {
			err := middle(w, r)
			if err != nil {
				//dont log
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		connetMux.ServeHTTP(w, r)
	})
	server := &http.Server{
		Addr:         config.Addr,
		Handler:      muxEx,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		Protocols:    &http.Protocols{},
	}
	if config.HTTP2 {
		server.HTTP2 = &http.HTTP2Config{
			MaxConcurrentStreams: config.MaxConcurrentStreams,
		}
		server.Protocols.SetHTTP2(true)
		server.Protocols.SetUnencryptedHTTP2(config.UnencryptedHTTP2)
	}
	server.Protocols.SetHTTP1(config.HTTP1)
	d := &ConnectServer{
		HTTPServer: server,
		conf:       config,
		logger:     logger,
		started:    atomic.Bool{},
	}
	if config.AutoRun {
		err := d.StartNoListner(sconf)
		if err != nil {
			return nil, nil, err
		}
	}
	return d, func() {
		d.Shoutdown()
	}, nil
}

func (p *ConnectServer) StartNoListner(sconf dot.SConfig) error {
	p.logger.Info().Msg("rpc api init without listener")
	if p.started.Swap(true) {
		return nil
	}
	err := p.conf.Tls.FullPath(sconf)
	if err != nil {
		return err
	}

	//check tls cert and key
	if (p.conf.Tls.Cert != "" && p.conf.Tls.Key == "") || (p.conf.Tls.Cert == "" && p.conf.Tls.Key != "") {
		return fmt.Errorf("tls cert and key must be both set or both empty")
	}

	go func() {
		if p.conf.Tls.NeedsTls() {
			p.logger.Info().Msgf("rpc tls listen(%s)", p.conf.Addr)
			if err := p.HTTPServer.ListenAndServeTLS(p.conf.Tls.Cert, p.conf.Tls.Key); err != nil {
				p.logger.Error().Err(err).Send()
			} else {
				p.logger.Info().Msg("rpc api done")
			}
		} else {
			p.logger.Info().Msgf("rpc listen(%s)", p.conf.Addr)
			if err := p.HTTPServer.ListenAndServe(); err != nil {
				p.logger.Error().Err(err).Send()
			} else {
				p.logger.Info().Msg("rpc api done")
			}
		}
	}()
	return nil
}

func (p *ConnectServer) StartWithListener(sconf dot.SConfig, listner net.Listener) error {
	p.logger.Info().Msg("rpc api init with listener")
	if p.started.Swap(true) {
		return nil
	}
	err := p.conf.Tls.FullPath(sconf)
	if err != nil {
		return err
	}

	//check tls cert and key
	if (p.conf.Tls.Cert != "" && p.conf.Tls.Key == "") || (p.conf.Tls.Cert == "" && p.conf.Tls.Key != "") {
		return fmt.Errorf("tls cert and key must be both set or both empty")
	}

	go func() {
		if p.conf.Tls.NeedsTls() {
			p.logger.Info().Msgf("rpc tls listen(%s)", p.conf.Addr)
			if err := p.HTTPServer.ServeTLS(listner, p.conf.Tls.Cert, p.conf.Tls.Key); err != nil {
				p.logger.Error().Err(err).Send()
			} else {
				p.logger.Info().Msg("rpc api done")
			}
		} else {
			p.logger.Info().Msgf("rpc listen(%s)", p.conf.Addr)
			if err := p.HTTPServer.Serve(listner); err != nil {
				p.logger.Error().Err(err).Send()
			} else {
				p.logger.Info().Msg("rpc api done")
			}
		}
	}()
	return nil
}

func (p *ConnectServer) Shoutdown() {
	if !p.started.Swap(false) {
		return
	}
	if p.HTTPServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), p.conf.ShutdownTimeout)
		defer cancel()
		if err := p.HTTPServer.Shutdown(ctx); err != nil {
			p.logger.Error().Err(err).Send()
		} else {
			p.logger.Info().Msg("rpc api shutdown successfully")
		}
		p.HTTPServer = nil
	}
	p.logger = nil
}
