package rpcdot

import (
	"net/http"
	"slices"
	"strings"
	"sync/atomic"
	"time"

	"github.com/scryinfo/dot/dot"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

func NewBothHttpServer(conf *ConnectServerConfig, connectMux *http.ServeMux, grpcServer *grpc.Server, logger *dot.LoggerType, middle HandlerMiddle) (*BothHttpServer, func(), error) {
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
		if r.ProtoMajor == 1 {
			connectMux.ServeHTTP(w, r)
			return
		}

		contentType := r.Header.Get("Content-Type")
		if strings.Contains(contentType, "grpc-web") || strings.Contains(contentType, "/connect+") {
			connectMux.ServeHTTP(w, r)
			return
		}
		grpcServer.ServeHTTP(w, r)
	})
	server := &http.Server{
		Addr: conf.Addr,
		Handler: h2c.NewHandler(muxEx, &http2.Server{
			MaxConcurrentStreams: conf.MaxConcurrentStreams,
		}),
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
	}
	d := &BothHttpServer{
		ConnectServer: ConnectServer{
			HTTPServer: server,
			conf:       *conf,
			logger:     logger,
			started:    atomic.Bool{},
		},
	}
	d.Start()

	return d, func() {
		d.Shoutdown()
	}, nil
}

type BothHttpServer struct {
	ConnectServer
}
