// Scry Info.  All rights reserved.
// license that can be found in the license file.

package gserver

import (
	"net/http"
	"os"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/gindot"
)

const (
	WebSocketTypeID = "27709862-0a7b-48e9-8427-5bbc86760c41"
)

// WebSocket component makes it possible the full-duplex communication with low latency
// between a gRPC-WebSocket request and a remote standard gRPC server.
type WebSocket struct {
	GinEngine *gindot.Engine `dot:""`

	wrappedServer *grpcweb.WrappedGrpcServer
}

func WebSocketTypeLives() []*dot.TypeLives {
	websocktTypeLive := func() *dot.TypeLives {
		return &dot.TypeLives{
			Meta: dot.Metadata{TypeID: WebSocketTypeID, NewDoter: func(conf []byte) (dot.Dot, error) {
				return &WebSocket{}, nil
			}},
			//Lives: []dot.Live{
			//	{
			//		TypeID: WebSocketTypeID,
			//		RelyLives: map[string]dot.LiveID{
			//			"GinEngine": gindot.EngineTypeID,
			//		},
			//	},
			//},
		}
	}
	return []*dot.TypeLives{websocktTypeLive(), gindot.TypeLiveGinDot()}
}

func (s *WebSocket) Stop(ignore bool) error {
	if s.wrappedServer != nil {
		s.wrappedServer = nil
	}
	return nil
}

// Wrap wraps the given grpcServer to allow for handling grpc-web requests of websockets - enabling bidirectional requests.
//
// Under the hood, Wrap takes a HTTP request from gin.Context and if it is a gRPC-WebSocket request wraps it with a compatibility layer
// to transform it to a standard gRPC request for the wrapped gRPC server and transforms the request to comply with
// the gRPC-Web protocol.
//
// Through this mechanism, client (e.g. browser) is able to fully take advantage of WebSocket communication with remote
// gRPC service server, initially routed by the HTTP GET method and the URLs of resources that are registered on gRPC server
//
// Note: this Wrap method can only be called before the underlying gin.Engine starts running, besides the caller must insure that
// all the grpc service servers get appropriately registered with the standard gRPC server.
func (s *WebSocket) Wrap(grpcServer *grpc.Server) {
	// Argument validation of the passed grpcServer
	if grpcServer == nil {
		dot.Logger().Errorln("nil argument of grpcServer not allowed")
		os.Exit(1)
	}
	if len(grpcServer.GetServiceInfo()) == 0 {
		dot.Logger().Errorln("no service registered with grpcServer")
		os.Exit(1)
	}

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
	}
	s.wrappedServer = grpcweb.WrapServer(grpcServer, options...)

	// HandlerFunc takes a HTTP request from gin.Context and if it is a gRPC-WebSocket request wraps it with a compatibility layer
	// to transform it to a standard gRPC request for the wrapped gRPC server and transforms the request to comply with
	// the gRPC-Web protocol.
	handlerFunc := func(ctx *gin.Context) {
		dot.Logger().Debugln("WebSocket", zap.String("", ctx.Request.RequestURI))

		if s.wrappedServer.IsGrpcWebSocketRequest(ctx.Request) {
			s.wrappedServer.ServeHTTP(ctx.Writer, ctx.Request)
			return
		}

		ctx.String(http.StatusForbidden, "Not a valid gRPC-WebSocket request!")
	}

	// registers the handlerFunc (without any middleware) with all URLs of resources that are registered on gRPC server.
	urls := grpcweb.ListGRPCResources(grpcServer)
	for _, url := range urls {
		s.GinEngine.GinEngine().GET(url, handlerFunc)
	}
}
