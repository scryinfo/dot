// Scry Info.  All rights reserved.
// license that can be found in the license file.

package gserver

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/gindot"
)

const (
	GinNoblTypeID = "3c9e8119-3d42-45bd-98f9-32939c895c6d"
)

//support the http and tcp
type ginNobl struct {
	ServerNobl        ServerNobl     `dot:""`
	GinRouter         *gindot.Router `dot:""`
	wrappedGrpcServer *grpcweb.WrappedGrpcServer
	preUrl            string
}

//GinNoblTypeLives Data structure needed when generating newer component
func GinNoblTypeLives() []*dot.TypeLives {

	lives := []*dot.TypeLives{{
		Meta: dot.Metadata{TypeID: GinNoblTypeID, RefType: reflect.TypeOf((*ginNobl)(nil)).Elem(), NewDoter: func(conf []byte) (dot dot.Dot, err error) {
			return &ginNobl{}, nil
		}},
		Lives: []dot.Live{
			{
				LiveID:    GinNoblTypeID,
				RelyLives: map[string]dot.LiveID{"GinRouter": gindot.RouterTypeID, "ServerNobl": ServerNoblTypeID},
			},
		},
	}}

	lives = append(lives, ServerNoblTypeLives()...)
	lives = append(lives, gindot.RouterTypeLives()...)
	return lives
}

//Run after every component finished start, this can ensure all service has been registered on grpc server
func (c *ginNobl) AfterAllStart(dot.Line) {
	if rp := c.GinRouter.RelativePath(); len(rp) > 0 && rp != "/" {
		if !strings.HasPrefix(rp, "/") {
			rp = "/" + rp
		}
		if !strings.HasSuffix(rp, "/") {
			rp += "/"
		}
		c.preUrl = rp
	} else {
		c.preUrl = ""
	}
	c.startServer()
}

//Stop stop dot
func (c *ginNobl) Stop(bool) error {
	if c.wrappedGrpcServer != nil {
		c.wrappedGrpcServer = nil
	}
	return nil
}

func (c *ginNobl) Server() *grpc.Server {
	return c.ServerNobl.Server()
}

func (c *ginNobl) startServer() {

	logger := dot.Logger()
	// Control the behaviour of the gRPC-WebSocket wrapper (e.g. modifying CORS behaviour) using `With*` options.
	options := []grpcweb.Option{
		// Allows for handling grpc-web requests of websockets - enabling bidirectional requests.
		grpcweb.WithWebsockets(true),
		// Accept all requests from remote origins,
		// don't check whether the origin of the request matches the host of the request.
		grpcweb.WithWebsocketOriginFunc(func(req *http.Request) bool {
			return true
		}),
		// Not allow requests incoming with a path prefix added to the URL,
		// exposing the endpoint as the root resource, to avoid
		// the performance cost of path processing for every request.
		grpcweb.WithAllowNonRootResource(false),
		// Only allow CORS requests for registered endpoints,
		// not allow handling gRPC requests for unknown endpoints (e.g. for proxying).
		grpcweb.WithCorsForRegisteredEndpointsOnly(true),
		grpcweb.WithAllowedRequestHeaders([]string{"Access-Control-Allow-Origin:*", "Access-Control-Allow-Methods:*"}),
	}//todo #49

	c.wrappedGrpcServer = grpcweb.WrapServer(c.Server(), options...)

	handle := func(ctx *gin.Context) {
		logger.Debugln("ginNobl", zap.String("", ctx.Request.RequestURI))

		if c.wrappedGrpcServer.IsGrpcWebRequest(ctx.Request) {
			if len(c.preUrl) > 0 { //todo #49
				old := ctx.Request.URL.Path
				if strings.HasPrefix(old, c.preUrl) {
					index := len(c.preUrl) - 1
					ctx.Request.URL.Path = old[index:]
				}
			}

			resp := ctx.Writer
			resp.Header().Set("Access-Control-Allow-Origin", "*")  //
			resp.Header().Set("Access-Control-Allow-Methods", "*") //
			resp.Header().Add("Access-Control-Allow-Headers", "content-type,x-grpc-web,x-user-agent")
			c.wrappedGrpcServer.ServeHTTP(resp, ctx.Request)
			return
		} else if c.wrappedGrpcServer.IsGrpcWebSocketRequest(ctx.Request) {
			if len(c.preUrl) > 0 { // because can not set the "endpointFunc" of WrapServer, do this so so
				old := ctx.Request.URL.Path
				if strings.HasPrefix(old, c.preUrl) {
					index := len(c.preUrl) - 1
					ctx.Request.URL.Path = old[index:]
				}
			}
			c.wrappedGrpcServer.HandleGrpcWebsocketRequest(ctx.Writer, ctx.Request)
			return
		} else {
			ctx.String(http.StatusOK, "no rpc")
		}
	}

	url := "/*rpc"

	c.GinRouter.Router().POST(url, handle)
	c.GinRouter.Router().OPTIONS(url, func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")  //
		ctx.Header("Access-Control-Allow-Methods", "*") //
		ctx.Header("Access-Control-Allow-Headers", "content-type,x-grpc-web,x-user-agent")
		ctx.String(http.StatusOK, "ok")
	})
	c.GinRouter.Router().GET(url, handle) //for websocket
}
