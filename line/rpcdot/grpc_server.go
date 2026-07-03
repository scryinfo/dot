package rpcdot

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/scryinfo/dot/dot"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GrpcServerConfig struct {
	// true : auto run, false: manual run
	AutoRun bool `json:"auto_run" toml:"auto_run" yaml:"auto_run" mapstructure:"auto_run"`
	// sample ":8080"
	Addr string `json:"addr" toml:"addr" yaml:"addr" mapstructure:"addr"`
	// ReadTimeout          time.Duration `json:"read_timeout" toml:"read_timeout" yaml:"read_timeout" mapstructure:"read_timeout"`
	// WriteTimeout         time.Duration `json:"write_timeout" toml:"write_timeout" yaml:"write_timeout" mapstructure:"write_timeout"`
	// MaxConcurrentStreams int           `json:"max_concurrent_streams" toml:"max_concurrent_streams" yaml:"max_concurrent_streams" mapstructure:"max_concurrent_streams"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout" toml:"shutdown_timeout" yaml:"shutdown_timeout" mapstructure:"shutdown_timeout"`
	Tls             TlsConfig     `json:"tls" toml:"tls" yaml:"tls" mapstructure:"tls"`
}
type GrpcServer struct {
	conf         *GrpcServerConfig
	logger       *dot.LoggerType
	started      atomic.Bool
	server       *grpc.Server
	lock         sync.Mutex
	lazyResister []func(server *grpc.Server)
}

func NewGrpcServer(config *GrpcServerConfig, sconf dot.SConfig, logger *dot.LoggerType) (*GrpcServer, func(), error) {
	err := config.Tls.FullPath(sconf)
	if err != nil {
		return nil, nil, err
	}
	if config.ShutdownTimeout < 0 {
		config.ShutdownTimeout = 10 * time.Second
	}
	server := &GrpcServer{
		server: nil, // create when start with tls or not
		conf:   config,
		logger: logger,
	}

	if config.AutoRun {
		err := server.StartNoListner()
		if err != nil {
			return nil, nil, err
		}
	}
	return server, func() {
		server.Shoutdown()
	}, nil
}

func (p *GrpcServer) StartNoListner() error {
	p.logger.Info().Msg("grpc init without listener")
	listner, err := net.Listen("tcp", p.conf.Addr)
	if err != nil {
		p.logger.Error().Err(err).Send()
		return err
	}
	return p._startWithListener(listner)
}

func (p *GrpcServer) StartWithListener(listner net.Listener) error {
	p.logger.Info().Msg("grpc api init with listener")
	return p._startWithListener(listner)
}

func (p *GrpcServer) _startWithListener(listner net.Listener) error {
	if p.started.Swap(true) {
		p.logger.Info().Msg("had started, just return")
		return nil
	}
	//check tls cert and key
	if (p.conf.Tls.Cert != "" && p.conf.Tls.Key == "") || (p.conf.Tls.Cert == "" && p.conf.Tls.Key != "") {
		err := fmt.Errorf("tls cert and key must be both set or both empty")
		p.logger.Error().Err(err).Send()
		return err
	}

	go func() {
		if p.conf.Tls.NeedsTls() {
			p.logger.Info().Msgf("grpc tls listen(%s)", p.conf.Addr)
			if p.conf.Tls.RootCert != "" {
				// double tls
				serverCert, err := tls.LoadX509KeyPair(p.conf.Tls.Cert, p.conf.Tls.Key)
				if err != nil {
					p.logger.Error().Err(err).Send()
					return
				}
				certPool := x509.NewCertPool()
				{
					caCert, err := os.ReadFile(p.conf.Tls.RootCert)
					if err != nil {
						p.logger.Error().Err(err).Send()
						return
					}
					if ok := certPool.AppendCertsFromPEM(caCert); !ok {
						p.logger.Error().Msg("failed to append ca cert")
						return
					}
				}
				{
					caCert, err := os.ReadFile(p.conf.Tls.PeerCert)
					if err != nil {
						p.logger.Error().Err(err).Send()
						return
					}
					if ok := certPool.AppendCertsFromPEM(caCert); !ok {
						p.logger.Error().Msg("failed to append ca cert")
						return
					}
				}
				tlsConfig := &tls.Config{
					Certificates: []tls.Certificate{serverCert},
					ClientAuth:   tls.RequireAndVerifyClientCert,
					ClientCAs:    certPool,
				}
				p.setServer(grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig))))
			} else {
				// single tls
				serverCert, err := tls.LoadX509KeyPair(p.conf.Tls.Cert, p.conf.Tls.Key)
				if err != nil {
					p.logger.Error().Err(err).Send()
					return
				}

				tlsConfig := &tls.Config{
					Certificates: []tls.Certificate{serverCert},
				}
				creds := credentials.NewTLS(tlsConfig)
				p.setServer(grpc.NewServer(grpc.Creds(creds)))
			}
		} else {
			p.logger.Info().Msgf("grpc listen(%s)", p.conf.Addr)
			p.setServer(grpc.NewServer())
		}
		if p.server != nil {
			if err := p.server.Serve(listner); err != nil {
				p.logger.Error().Err(err).Send()
			} else {
				p.logger.Info().Msg("grpc api done")
			}
		} else {
			p.logger.Error().Msg("cant create the grpc server by tls config")
			return
		}
	}()
	return nil
}

func (p *GrpcServer) Shoutdown() {
	if !p.started.Swap(false) {
		return
	}
	if p.server != nil {
		stopped := make(chan struct{})
		go func() {
			p.server.GracefulStop()
			close(stopped)
		}()

		select {
		case <-stopped:
			p.logger.Info().Msg("gRPC server gracefully stopped.")
		case <-time.After(p.conf.ShutdownTimeout):
			p.logger.Info().Msg("Graceful stop timed out. Forcing sharp stop...")
			p.server.Stop()
		}

		p.server = nil
	}
	p.logger = nil
}

// resister the service, if the grpc server is nil , push it into the lazy task
// when the gprc server is not nil, lazy resister it
func (p *GrpcServer) ResisterOrLazy(f func(server *grpc.Server)) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.server != nil {
		f(p.server)
	} else {
		p.lazyResister = append(p.lazyResister, f)
	}
}

func (p *GrpcServer) setServer(server *grpc.Server) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.server = server
	for _, it := range p.lazyResister {
		it(server)
	}
}
