// Scry Info.  All rights reserved.
// license that can be found in the license file.

package gserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"go.uber.org/zap"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/gindot"
)

const (
	WebSocketTypeId = "27709862-0a7b-48e9-8427-5bbc86760c41"
)

// WebSocket component makes it possible the communication between a gRPC-WebSocket request and a standard gRPC server.
type WebSocket struct {
	ServerNobl ServerNobl     `dot:""`
	GinRouter  *gindot.Router `dot:""`

	wrappedServer *grpcweb.WrappedGrpcServer
}

func WebSocketTypeLive() *dot.TypeLives {
	newWebSocket := func(conf []byte) (dot.Dot, error) {
		return &WebSocket{}, nil
	}

	return &dot.TypeLives{
		Meta: dot.Metadata{TypeId: WebSocketTypeId, NewDoter: newWebSocket},
		Lives: []dot.Live{
			{
				TypeId: WebSocketTypeId,
				RelyLives: map[string]dot.LiveId{
					"ServerNobl": ServerNoblTypeId,
					"GinRouter":  gindot.RouterTypeId,
				},
			},
		},
	}
}

func WebSocketAndRelyTypeLives() []*dot.TypeLives {
	typeLives := []*dot.TypeLives{WebSocketTypeLive(), ServerNoblTypeLive()}
	typeLives = append(typeLives, gindot.TypeLiveRouter()...)
	return typeLives
}

func (s *WebSocket) AfterAllStart(l dot.Line) {
	s.serve()
}

func (s *WebSocket) Stop(ignore bool) error {
	if s.wrappedServer != nil {
		s.wrappedServer = nil
	}
	return nil
}

func (s *WebSocket) serve() {
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
	s.wrappedServer = grpcweb.WrapServer(s.ServerNobl.Server(), options...)

	handlerFunc := func(ctx *gin.Context) {
		dot.Logger().Debugln("WebSocket", zap.String("", ctx.Request.RequestURI))

		if s.wrappedServer.IsGrpcWebSocketRequest(ctx.Request) {
			s.wrappedServer.ServeHTTP(ctx.Writer, ctx.Request)
			return
		}

		ctx.String(http.StatusForbidden, "Not valid gRPC-WebSocket request!")
	}

	s.GinRouter.Router().GET("/*grpc-websocket", handlerFunc)
}
