// Scry Info.  All rights reserved.
// license that can be found in the license file.

package gserver

import (
	"net/http"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/gindot"
)

const (
	WebSocketTypeId = "27709862-0a7b-48e9-8427-5bbc86760c41"
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
			Meta: dot.Metadata{TypeId: WebSocketTypeId, NewDoter: func(conf []byte) (dot.Dot, error) {
				return &WebSocket{}, nil
			}},
			//Lives: []dot.Live{
			//	{
			//		TypeId: WebSocketTypeId,
			//		RelyLives: map[string]dot.LiveId{
			//			"GinEngine": gindot.EngineTypeId,
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

// GET wraps the given grpcServer to allow for handling grpc-web requests of websockets - enabling bidirectional requests.
//
// Under the hood, GET takes a HTTP request from gin.Context and if it is a gRPC-WebSocket request wraps it with a compatibility layer
// to transform it to a standard gRPC request for the wrapped gRPC server and transforms the request to comply with
// the gRPC-Web protocol.
//
// Through this mechanism, client (e.g. browser) is able to fully take advantage of WebSocket communication with remote
// gRPC service server, initially routed by the HTTP GET method and the given path.
//
// Note: this GET method can only be called before the underlying gin.Engine starts running, besides the caller must insure that
// all the grpc service servers get appropriately registered with the standard gRPC server.
func (s *WebSocket) GET(relativePath string, grpcServer *grpc.Server) {
	// wrap grpcServer to allow for handling grpc-web requests of websockets - enabling bidirectional requests.
	options := []grpcweb.Option{
		grpcweb.WithWebsockets(true),
		grpcweb.WithWebsocketOriginFunc(func(req *http.Request) bool {
			return true
		}),
		grpcweb.WithOriginFunc(func(origin string) bool {
			return true
		}),
		grpcweb.WithCorsForRegisteredEndpointsOnly(false),
	}
	s.wrappedServer = grpcweb.WrapServer(grpcServer, options...)

	// handlerFunc takes a HTTP request from gin.Context and if it is a gRPC-WebSocket request wraps it with a compatibility layer
	// to transform it to a standard gRPC request for the wrapped gRPC server and transforms the request to comply with
	// the gRPC-Web protocol.
	handlerFunc := func(ctx *gin.Context) {
		dot.Logger().Debugln("WebSocket", zap.String("", ctx.Request.RequestURI))

		if s.wrappedServer.IsGrpcWebSocketRequest(ctx.Request) {
			s.wrappedServer.ServeHTTP(ctx.Writer, ctx.Request)
			return
		}

		ctx.String(http.StatusForbidden, "Not valid gRPC-WebSocket request!")
	}

	// registers the new request handle (without any middleware) with the given path and GET method with gin.Engine.
	s.GinEngine.GinEngine().GET(relativePath, handlerFunc)
}
