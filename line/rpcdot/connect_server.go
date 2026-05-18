package rpcdot

import (
	"context"
	"fmt"
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
	// sample ":8080"
	Addr                 string        `json:"addr" toml:"addr" yaml:"addr"`
	ReadTimeout          time.Duration `json:"readTimeout" toml:"readTimeout" yaml:"readTimeout"`
	WriteTimeout         time.Duration `json:"writeTimeout" toml:"writeTimeout" yaml:"writeTimeout"`
	MaxConcurrentStreams int           `json:"maxConcurrentStreams" toml:"maxConcurrentStreams" yaml:"maxConcurrentStreams"`
	AllowedOrigins       []string      `json:"allowedOrigins" toml:"allowedOrigins" yaml:"allowedOrigins"`
	AllowHeaders         []string      `json:"allowHeaders" toml:"allowHeaders" yaml:"allowHeaders"`
	AllowMethods         []string      `json:"allowMethods" toml:"allowMethods" yaml:"allowMethods"`
	AllowCredentials     bool          `json:"allowCredentials" toml:"allowCredentials" yaml:"allowCredentials"`
	// dont auth these urls,sample: ["/login", "/rpc/test"]
	// ["*"] dont auth all urls
	UnAuthUrls []string `json:"unAuthUrls" toml:"unAuthUrls" yaml:"unAuthUrls"`
	// if it is true, all OPTIONS requests will be returned ok
	OptionMethods bool `json:"optionMethods" toml:"optionMethods" yaml:"optionMethods"`
	// shutdown timeout
	ShutdownTimeout  time.Duration `json:"shutdownTimeout" toml:"shutdownTimeout" yaml:"shutdownTimeout"`
	HTTP1            bool          `json:"http1" toml:"http1" yaml:"http1"`
	HTTP2            bool          `json:"http2" toml:"http2" yaml:"http2"`
	UnencryptedHTTP2 bool          `json:"unencryptedHTTP2" toml:"unencryptedHTTP2" yaml:"unencryptedHTTP2"`

	Tls TlsConfig `json:"tls" toml:"tls" yaml:"tls"`
}

type ConnectServer struct {
	HTTPServer *http.Server
	conf       ConnectServerConfig
	logger     *dot.LoggerType
	started    atomic.Bool
}

func NewConnetServer(conf *ConnectServerConfig, sconf dot.SConfig, connetMux *ConnectHttpServerMux, logger *dot.LoggerType, middle HandlerMiddle) (*ConnectServer, func(), error) {
	if conf.ShutdownTimeout < 0 {
		conf.ShutdownTimeout = 10 * time.Second
	}
	if len(conf.UnAuthUrls) > 0 {
		unauthUrls = conf.UnAuthUrls
	}
	if len(conf.AllowMethods) < 1 {
		conf.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	}
	if len(conf.AllowHeaders) < 1 {
		conf.AllowHeaders = []string{
			"Origin", "Upgrade", "Connection", "X-Requested-With", "X-HTTP-Protocol", "Content-Type",
			"Accept", "Cookie", "connect-protocol-version", "connect-timeout-ms",
			httptools.TokenName, httptools.TokenGame, httptools.Authorization, httptools.AuthorizationGame}
	}
	if len(conf.AllowedOrigins) < 1 {
		conf.AllowedOrigins = []string{"http://localhost", "http://127.0.0.1"}
	}

	allowMethods := strings.Join(conf.AllowMethods, ",")
	allowHeaders := strings.Join(conf.AllowHeaders, ",")
	muxEx := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		{
			originHeader := r.Header.Get("Origin")
			if len(conf.AllowedOrigins) == 1 && conf.AllowedOrigins[0] == "*" {
				header.Set("Access-Control-Allow-Origin", "*")
			} else if slices.Contains(conf.AllowedOrigins, originHeader) {
				header.Set("Access-Control-Allow-Origin", originHeader)
			} else if slices.ContainsFunc(conf.AllowedOrigins, func(origin string) bool {
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
		if conf.AllowCredentials {
			header.Set("Access-Control-Allow-Credentials", "true")
		}
		if conf.OptionMethods && r.Method == "OPTIONS" {
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
		Addr:         conf.Addr,
		Handler:      muxEx,
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
		Protocols:    &http.Protocols{},
	}
	if conf.HTTP2 {
		server.HTTP2 = &http.HTTP2Config{
			MaxConcurrentStreams: conf.MaxConcurrentStreams,
		}
		server.Protocols.SetHTTP2(true)
		server.Protocols.SetUnencryptedHTTP2(conf.UnencryptedHTTP2)
	}
	server.Protocols.SetHTTP1(conf.HTTP1)
	d := &ConnectServer{
		HTTPServer: server,
		conf:       *conf,
		logger:     logger,
		started:    atomic.Bool{},
	}
	err := d.start(sconf)
	return d, func() {
		d.Shoutdown()
	}, err
}

func (p *ConnectServer) start(sconf dot.SConfig) error {
	p.logger.Info().Msg("rpc api init")
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
