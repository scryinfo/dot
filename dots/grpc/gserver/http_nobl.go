package gserver

import (
	"context"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/scryinfo/dot/dot"
	"google.golang.org/grpc"
	"net/http"
)

const (
	HttpTypeId = "afbeac47-e5fd-4bf3-8fb1-f0fb8ec79bd0"
)

type httpNoblConf struct {
	//sample :  1.1.1.1:568
	Addr string `json:"addr"`
	//sample:  keys/s1.pem,   Comparing with current exe path
	PemPath string `json:"pemPath"`
	//sample:  keys/s1.pem,   Comparing with current exe path
	KeyPath string `json:"keyPath"`
}

//support the http and tcp
type httpNobl struct {
	conf       httpNoblConf
	ServerNobl     ServerNobl `dot:""`
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

	d := &httpNobl{
		conf: *dconf,
	}

	return d, err
}

//Data structure needed when generating newer component
func HttpNoblTypeLives() []*dot.TypeLives {

	tl := &dot.TypeLives{
		Meta: dot.Metadata{TypeId: HttpTypeId, NewDoter: func(conf interface{}) (dot dot.Dot, err error) {
			return newHttpNobl(conf)
		}},
		Lives:[]dot.Live{
			dot.Live{
				LiveId: HttpTypeId,
				RelyLives: map[string]dot.LiveId{"ServerNobl": ServerNoblTypeId},
			},
		},
	}

	return []*dot.TypeLives{
		tl, ServerNoblTypeLive(),
	}
}

func (c *httpNobl) Create(l dot.Line) error {
	var err error = nil


	return err
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

	//options.OptionsPassthrough
	wrappedGrpc := grpcweb.WrapServer(c.Server(), grpcweb.WithAllowedRequestHeaders([]string{"Access-Control-Allow-Origin:*", "Access-Control-Allow-Methods:*"}) )

	//start http grpc
	c.httpServer = &http.Server{Addr:c.conf.Addr}
	c.httpServer.Handler = http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		//if wrappedGrpc.IsGrpcWebRequest(req) {
			wrappedGrpc.ServeHTTP(resp, req)
		//}
		//dot.Logger().Infoln("httpNobl", zap.String("", "it is not grpc request from the http"))
	})
	go func() {
		_ = c.httpServer.ListenAndServe()
	}()

}
