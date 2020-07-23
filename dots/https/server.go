package https

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/utils"
	"go.uber.org/zap"
	"io/ioutil"
	"net"
	"net/http"
)

const TypeID = "a998357a-3d92-4390-bd5c-ef6fd6e2ba41"

type config struct {
	Addr string          `json:"addr"` // addr smaple:  ":8080"
	Tls  utils.TlsConfig `json:"tls"`
}

type Server struct {
	conf          config
	httpServer    *http.Server
	defaultHandle http.ServeMux
}

func (c *Server) Stop(ignore bool) (err error) {
	if c.httpServer != nil {
		err = c.httpServer.Close()
		if err != nil {
			dot.Logger().Errorln("Server", zap.Error(err))
		}
	}
	return
}

//construct dot
func newServer(conf []byte) (dot.Dot, error) {
	dconf := &config{}
	err := json.Unmarshal(conf, dconf)
	if err != nil {
		return nil, err
	}

	d := &Server{conf: *dconf}

	return d, err
}

//ServerTypeLives
func ServerTypeLives() []*dot.TypeLives {
	lives := []*dot.TypeLives{
		{
			Meta: dot.Metadata{TypeID: TypeID, NewDoter: func(conf []byte) (dot dot.Dot, err error) {
				return newServer(conf)
			}},
		},
	}
	return lives
}

func ServerConfigTypeLive() *dot.ConfigTypeLive {
	paths := make([]string, 0)
	paths = append(paths, "")
	return &dot.ConfigTypeLive{
		TypeIDConfig: TypeID,
		ConfigInfo:   &config{},
	}
}

func Test(addr string, tls *utils.TlsConfig) *Server {
	s := &Server{conf: config{
		Addr: addr,
	}}
	if tls != nil {
		s.conf.Tls = *tls
	}

	s.AfterAllStart(nil)

	return s
}

func (c *Server) RegisterHandle(pattern string, handle http.Handler) error {
	if handle == nil {
		return errors.New("service is nil")
	}
	c.defaultHandle.Handle(pattern, handle)
	return nil
}

//AfterAllStart run the function after start
func (c *Server) AfterAllStart(l dot.Line) {
	go c.startServer()
}

func (c *Server) startServer() error {
	logger := dot.Logger()

	c.httpServer = &http.Server{Handler: &c.defaultHandle}
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
			logger.Errorln("Server", zap.Error(err))
			return err
		}
		key := utils.GetFullPathFile(c.conf.Tls.Key)
		if len(key) < 1 {
			err = errors.New("the Key is not empty, and can not find the file: " + c.conf.Tls.Key)
			logger.Errorln("Server", zap.Error(err))
			return err
		}

		pem := utils.GetFullPathFile(c.conf.Tls.Pem)
		if len(pem) < 1 {
			err = errors.New("the Pem is not empty, and can not find the file: " + c.conf.Tls.Pem)
			logger.Errorln("Server", zap.Error(err))
			return err
		}

		pool := x509.NewCertPool()
		{
			caCrt, err := ioutil.ReadFile(caPem)
			if err != nil {
				logger.Errorln("Server", zap.Error(err))
				return err
			}
			if !pool.AppendCertsFromPEM(caCrt) {
				err = errors.New("credentials: failed to append certificates")
				logger.Errorln("Server", zap.Error(err))
				return err
			}
		}
		cert, err := tls.LoadX509KeyPair(pem, key)
		if err != nil {
			logger.Errorln("Server", zap.Error(err))
			return err
		}

		c.httpServer.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientCAs:    pool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
		}
		logger.Infoln("Server", zap.String("", "ApiService server(with ca) will start: "+c.conf.Addr))
		serverRun = func() error {
			return c.httpServer.ServeTLS(conn, "", "")
		}
	case len(c.conf.Tls.Key) > 0 && len(c.conf.Tls.Pem) > 0: //just server
		pem := utils.GetFullPathFile(c.conf.Tls.Pem)
		if len(pem) < 1 {
			err = errors.New("the pem is not empty, and can not find the file: " + c.conf.Tls.Pem)
			logger.Errorln("Server", zap.Error(err))
			return err
		}
		key := utils.GetFullPathFile(c.conf.Tls.Key)
		if len(key) < 1 {
			err = errors.New("the key is not empty, and can not find the file: " + c.conf.Tls.Key)
			logger.Errorln("Server", zap.Error(err))
			return err
		}

		c.httpServer.TLSConfig = &tls.Config{
			ClientAuth: tls.NoClientCert,
		}

		logger.Infoln("Server", zap.String("", "ApiService server(no ca) will start: "+c.conf.Addr))
		serverRun = func() error {
			return c.httpServer.ServeTLS(conn, pem, key)
		}
	default: //no tls
		logger.Infoln("Server", zap.String("", "ApiService server(no https) will start: "+c.conf.Addr))
		serverRun = func() error {
			return c.httpServer.Serve(conn)
		}
	}

	go func() {
		e := serverRun()
		if e != nil {
			logger.Infoln("Server", zap.Error(e))
		}
		conn.Close()
	}()

	return nil
}
