// Scry Info.  All rights reserved.
// license that can be found in the license file.

package gserver

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/pkg/errors"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/grpc/shared"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	HttpTypeId = "afbeac47-e5fd-4bf3-8fb1-f0fb8ec79bd0"
)

type httpNoblConf struct {
	//sample :  1.1.1.1:568
	PreUrl string           `json:"preUrl"`
	Addr   string           `json:"addr"`
	Tls    shared.TlsConfig `json:"tls"`
}

//support the http and tcp
type httpNobl struct {
	conf       httpNoblConf
	ServerNobl ServerNobl `dot:""`
	httpServer *http.Server
}

//Construct component
func newHttpNobl(conf interface{}) (dot.Dot, error) {
	var err error = nil
	var bs []byte = nil
	if bt, ok := conf.([]byte); ok {
		bs = bt
	} else {
		return nil, dot.SError.Parameter
	}
	dconf := &httpNoblConf{}
	err = dot.UnMarshalConfig(bs, dconf)
	if err != nil {
		return nil, err
	}
	if len(dconf.PreUrl) > 0 {
		if !strings.HasPrefix(dconf.PreUrl, "/") {
			dconf.PreUrl = "/" + dconf.PreUrl
		}
		if !strings.HasSuffix(dconf.PreUrl, "/") {
			dconf.PreUrl += "/"
		}
	}

	d := &httpNobl{
		conf: *dconf,
	}

	return d, err
}

//HttpNoblTypeLives Data structure needed when generating newer component
func HttpNoblTypeLives() []*dot.TypeLives {

	tl := &dot.TypeLives{
		Meta: dot.Metadata{TypeId: HttpTypeId, NewDoter: func(conf interface{}) (dot dot.Dot, err error) {
			return newHttpNobl(conf)
		}},
		Lives: []dot.Live{
			{
				LiveId:    HttpTypeId,
				RelyLives: map[string]dot.LiveId{"ServerNobl": ServerNoblTypeId},
			},
		},
	}

	return []*dot.TypeLives{
		tl, ServerNoblTypeLive(),
	}
}

//Run after every component finished start, this can ensure all service has been registered on grpc server
func (c *httpNobl) AfterAllStart(l dot.Line) {
	c.startServer()
}

//Stop stop dot
func (c *httpNobl) Stop(ignore bool) error {
	if c.httpServer != nil {
		_ = c.httpServer.Shutdown(context.Background())
		c.httpServer = nil
	}
	return nil
}

func (c *httpNobl) Server() *grpc.Server {
	return c.ServerNobl.Server()
}

func (c *httpNobl) startServer() {
	logger := dot.Logger()
	//options.OptionsPassthrough
	wrappedGrpc := grpcweb.WrapServer(c.Server(), grpcweb.WithAllowedRequestHeaders([]string{"Access-Control-Allow-Origin:*", "Access-Control-Allow-Methods:*"}))

	//start http grpc
	c.httpServer = &http.Server{Addr: c.conf.Addr}
	c.httpServer.Handler = http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		logger.Debugln("httpNobl", zap.String("", req.RequestURI))
		//if wrappedGrpc.IsGrpcWebRequest(req) {
		resp.Header().Set("Access-Control-Allow-Origin", "*")  //
		resp.Header().Set("Access-Control-Allow-Methods", "*") //
		resp.Header().Add("Access-Control-Allow-Headers", "content-type,x-grpc-web,x-user-agent")
		if len(c.conf.PreUrl) > 0 { // because can not set the "endpointFunc" of WrapServer, do this so so
			old := req.URL.Path
			if strings.HasPrefix(old, c.conf.PreUrl) {
				index := len(c.conf.PreUrl) - 1
				req.URL.Path = old[index:]
			}
		}
		wrappedGrpc.ServeHTTP(resp, req)
		//}
		//dot.Logger().Infoln("httpNobl", zap.String("", "it is not grpc request from the http"))
	})

	go func() {
		switch {
		case len(c.conf.Tls.CaPem) > 0 && len(c.conf.Tls.Key) > 0 && len(c.conf.Tls.Pem) > 0: //both tls
			caPem := shared.GetFullPathFile(c.conf.Tls.CaPem)
			if len(caPem) < 1 {
				logger.Errorln("httpNobl", zap.Error(errors.New("the caPem is not empty, and can not find the file: "+c.conf.Tls.CaPem)))
				return
			}
			key := shared.GetFullPathFile(c.conf.Tls.Key)
			if len(key) < 1 {
				logger.Errorln("httpNobl", zap.Error(errors.New("the Key is not empty, and can not find the file: "+c.conf.Tls.Key)))
				return
			}

			pem := shared.GetFullPathFile(c.conf.Tls.Pem)
			if len(pem) < 1 {
				logger.Errorln("httpNobl", zap.Error(errors.New("the Pem is not empty, and can not find the file: "+c.conf.Tls.Pem)))
				return
			}

			pool := x509.NewCertPool()
			{
				caCrt, err1 := ioutil.ReadFile(caPem)
				if err1 != nil {
					logger.Errorln("httpNobl", zap.Error(errors.WithStack(err1)))
					return
				}
				if !pool.AppendCertsFromPEM(caCrt) {
					logger.Errorln("httpNobl", zap.Error(errors.New("credentials: failed to append certificates")))
					return
				}
			}
			cert, err1 := tls.LoadX509KeyPair(pem, key)
			if err1 != nil {
				logger.Errorln("httpNobl", zap.Error(errors.WithStack(err1)))
				return
			}

			c.httpServer.TLSConfig = &tls.Config{
				Certificates: []tls.Certificate{cert},
				ClientCAs:    pool,
				ClientAuth:   tls.RequireAndVerifyClientCert,
			}
			logger.Infoln("httpNobl", zap.String("", "grpc-web server(with ca) will start: "+c.conf.Addr))
			err := c.httpServer.ListenAndServeTLS("", "")
			if err != nil {
				logger.Errorln("httpNobl", zap.Error(errors.WithStack(err)))
			}
		case len(c.conf.Tls.Key) > 0 && len(c.conf.Tls.Pem) > 0: //just server
			pem := shared.GetFullPathFile(c.conf.Tls.Pem)
			if len(pem) < 1 {
				logger.Errorln("httpNobl", zap.Error(errors.New("the pem is not empty, and can not find the file: "+c.conf.Tls.Pem)))
				return
			}
			key := shared.GetFullPathFile(c.conf.Tls.Key)
			if len(key) < 1 {
				logger.Errorln("httpNobl", zap.Error(errors.New("the key is not empty, and can not find the file: "+c.conf.Tls.Key)))
				return
			}

			c.httpServer.TLSConfig = &tls.Config{
				ClientAuth: tls.NoClientCert,
			}

			logger.Infoln("httpNobl", zap.String("", "grpc-web server(no ca) will start: "+c.conf.Addr))
			err := c.httpServer.ListenAndServeTLS(pem, key)
			if err != nil {
				logger.Errorln("httpNobl", zap.Error(errors.WithStack(err)))
			}
		default: //no tls
			logger.Infoln("httpNobl", zap.String("", "grpc-web server(no https) will start: "+c.conf.Addr))
			err := c.httpServer.ListenAndServe()
			if err != nil {
				logger.Errorln("httpNobl", zap.Error(errors.WithStack(err)))
			}
		}

	}()

}
