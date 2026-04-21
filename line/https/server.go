package https

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/utils"
)

const TypeID = "a998357a-3d92-4390-bd5c-ef6fd6e2ba41"

type ConfigTls struct {
	Addr string          `json:"addr"` // addr smaple:  ":8080"
	Tls  utils.TlsConfig `json:"tls"`
}

type HttpSeverTls struct {
	conf       ConfigTls
	httpServer *http.Server
	mux        *http.ServeMux
}

func (c *HttpSeverTls) Stop() (err error) {
	if c.httpServer != nil {
		err = c.httpServer.Close()
		if err != nil {
			dot.Logger.Error().AnErr("Server", err).Send()
		}
	}
	return
}

// construct dot
func NewHttpSeverTls(conf *ConfigTls, server *http.Server, mux *http.ServeMux) (*HttpSeverTls, func(), error) {
	d := &HttpSeverTls{conf: *conf, httpServer: server, mux: mux}
	d.startServer()
	return d, func() {
		d.Stop()
	}, nil
}

func Test(addr string, tls *utils.TlsConfig) *HttpSeverTls {
	s, _, _ := NewHttpSeverTls(&ConfigTls{
		Addr: addr,
		Tls:  *tls,
	}, nil, nil)

	return s
}

func (c *HttpSeverTls) RegisterHandle(pattern string, handle http.Handler) error {
	if handle == nil {
		return errors.New("service is nil")
	}
	c.mux.Handle(pattern, handle)
	return nil
}

func (c *HttpSeverTls) startServer() error {
	logger := &dot.Logger

	c.httpServer = &http.Server{Handler: c.mux}
	conn, err := net.Listen("tcp", c.conf.Addr)
	if err != nil {
		return err
	}

	var serverRun func() error

	switch {
	case len(c.conf.Tls.CaPem) > 0 && len(c.conf.Tls.Key) > 0 && len(c.conf.Tls.Pem) > 0: //both tls
		caPem := utils.GetFullPathFile(c.conf.Tls.CaPem)
		if len(caPem) < 1 {
			err = errors.New("the caPem is not empty, and can not find the file: " + c.conf.Tls.CaPem)
			logger.Error().AnErr("Server", err).Send()
			return err
		}
		key := utils.GetFullPathFile(c.conf.Tls.Key)
		if len(key) < 1 {
			err = errors.New("the Key is not empty, and can not find the file: " + c.conf.Tls.Key)
			logger.Error().AnErr("Server", err).Send()
			return err
		}

		pem := utils.GetFullPathFile(c.conf.Tls.Pem)
		if len(pem) < 1 {
			err = errors.New("the Pem is not empty, and can not find the file: " + c.conf.Tls.Pem)
			logger.Error().AnErr("Server", err).Send()
			return err
		}

		pool := x509.NewCertPool()
		{
			caCrt, err := os.ReadFile(caPem)
			if err != nil {
				logger.Error().AnErr("Server", err).Send()
				return err
			}
			if !pool.AppendCertsFromPEM(caCrt) {
				err = errors.New("credentials: failed to append certificates")
				logger.Error().AnErr("Server", err).Send()
				return err
			}
		}
		cert, err := tls.LoadX509KeyPair(pem, key)
		if err != nil {
			logger.Error().AnErr("Server", err).Send()
			return err
		}

		c.httpServer.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientCAs:    pool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
		}
		logger.Info().Msg("ApiService server(with ca) will start: " + c.conf.Addr)
		serverRun = func() error {
			return c.httpServer.ServeTLS(conn, "", "")
		}
	case len(c.conf.Tls.Key) > 0 && len(c.conf.Tls.Pem) > 0: //just server
		pem := utils.GetFullPathFile(c.conf.Tls.Pem)
		if len(pem) < 1 {
			err = errors.New("the pem is not empty, and can not find the file: " + c.conf.Tls.Pem)
			logger.Error().AnErr("Server", err).Send()
			return err
		}
		key := utils.GetFullPathFile(c.conf.Tls.Key)
		if len(key) < 1 {
			err = errors.New("the key is not empty, and can not find the file: " + c.conf.Tls.Key)
			logger.Error().AnErr("Server", err).Send()
			return err
		}

		c.httpServer.TLSConfig = &tls.Config{
			ClientAuth: tls.NoClientCert,
		}

		logger.Info().Msg("ApiService server(no ca) will start: " + c.conf.Addr)
		serverRun = func() error {
			return c.httpServer.ServeTLS(conn, pem, key)
		}
	default: //no tls
		logger.Info().Msg("ApiService server(no https) will start: " + c.conf.Addr)
		serverRun = func() error {
			return c.httpServer.Serve(conn)
		}
	}

	go func() {
		e := serverRun()
		if e != nil {
			logger.Info().Err(e).Send()
		}
		conn.Close()
	}()

	return nil
}
