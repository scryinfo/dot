package connectdot

import (
	"context"
	"net/http"
	"slices"
	"strings"
	"sync/atomic"
	"time"

	"github.com/scryinfo/dot/dot"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func NewHttpServerMux() *http.ServeMux {
	mux := http.NewServeMux()
	return mux
}

type HandlerMiddle func(w http.ResponseWriter, r *http.Request) error

type HttpServerConfig struct {
	// sample ":8080"
	Addr                 string
	ReadTimeout          time.Duration
	WriteTimeout         time.Duration
	MaxConcurrentStreams uint32
	AllowedOrigins       []string
	AllowHeaders         []string
	AllowMethods         []string
	AllowCredentials     bool
	// dont auth these urls,sample: ["/login", "/rpc/test"]
	UnAuthUrls []string
	// if it is true, all OPTIONS requests will be returned ok
	OptionMethods bool

	ShutdownTimeout time.Duration
}

type HttpServerEx struct {
	HTTPServer *http.Server
	conf       HttpServerConfig
	rawMux     *http.ServeMux
	logger     *dot.LoggerType
	started    atomic.Bool
}

func NewHttpServerEx(conf *HttpServerConfig, mux *http.ServeMux, logger *dot.LoggerType, middle HandlerMiddle) (*HttpServerEx, func(), error) {
	if conf.ShutdownTimeout < 0 {
		conf.ShutdownTimeout = 10 * time.Second
	}
	if len(conf.UnAuthUrls) < 1 {
		unauthUrls = conf.UnAuthUrls
	}
	allowMethods := strings.Join(conf.AllowMethods, ",")
	allowHeaders := strings.Join(conf.AllowHeaders, ",")
	muxEx := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		{
			origin := r.Header.Get("Origin")
			if slices.Contains(conf.AllowedOrigins, origin) {
				header.Set("Access-Control-Allow-Origin", origin)
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
		mux.ServeHTTP(w, r)
	})
	server := &http.Server{
		Addr: conf.Addr,
		Handler: h2c.NewHandler(muxEx, &http2.Server{
			MaxConcurrentStreams: conf.MaxConcurrentStreams,
		}),
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
	}
	d := &HttpServerEx{
		HTTPServer: server,
		conf:       *conf,
		rawMux:     mux,
		logger:     logger,
		started:    atomic.Bool{},
	}
	d.Start()

	return d, func() {
		d.Shoutdown()
	}, nil
}

func (p *HttpServerEx) Start() {
	p.logger.Info().Msg("rpc api init")
	if p.started.Swap(true) {
		return
	}

	go func() {
		p.logger.Info().Msgf("rpc listen(%s)", p.conf.Addr)
		if err := p.HTTPServer.ListenAndServe(); err != nil {
			p.logger.Error().Err(err).Send()
		} else {
			p.logger.Info().Msg("rpc api done")
		}
	}()
}
func (p *HttpServerEx) Shoutdown() {
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
	p.rawMux = nil
	p.logger = nil
}
