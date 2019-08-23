// Scry Info.  All rights reserved.
// license that can be found in the license file.

package gserver

import (
	"github.com/gin-gonic/gin"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/gindot"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
	"strings"
)

const (
	GinNoblTypeId = "3c9e8119-3d42-45bd-98f9-32939c895c6d"
)

//support the http and tcp
type ginNobl struct {
	ServerNobl ServerNobl     `dot:""`
	GinRouter  *gindot.Router `dot:""`
	wrapserver *grpcweb.WrappedGrpcServer
	preUrl     string
}

//GinNoblTypeLives Data structure needed when generating newer component
func GinNoblTypeLives() []*dot.TypeLives {

	tl := &dot.TypeLives{
		Meta: dot.Metadata{TypeId: GinNoblTypeId, NewDoter: func(conf interface{}) (dot dot.Dot, err error) {
			return &ginNobl{}, nil
		}},
		Lives: []dot.Live{
			dot.Live{
				LiveId:    GinNoblTypeId,
				RelyLives: map[string]dot.LiveId{"GinRouter": gindot.RouterTypeId, "ServerNobl": ServerNoblTypeId},
			},
		},
	}

	lives := []*dot.TypeLives{
		tl, ServerNoblTypeLive(),
	}

	lives = append(lives, gindot.TypeLiveRouter()...)
	return lives
}

//Run after every component finished start, this can ensure all service has been registered on grpc server
func (c *ginNobl) AfterAllStart(l dot.Line) {
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
func (c *ginNobl) Stop(ignore bool) error {
	if c.wrapserver != nil {
		c.wrapserver = nil
	}
	return nil
}

func (c *ginNobl) Server() *grpc.Server {
	return c.ServerNobl.Server()
}

func (c *ginNobl) startServer() {

	logger := dot.Logger()
	c.wrapserver = grpcweb.WrapServer(c.Server(), grpcweb.WithAllowedRequestHeaders([]string{"Access-Control-Allow-Origin:*", "Access-Control-Allow-Methods:*"}))



	handle := func(ctx *gin.Context) {
		logger.Debugln("ginNobl", zap.String("", ctx.Request.RequestURI))
		if c.wrapserver.IsGrpcWebRequest(ctx.Request) {

			if len(c.preUrl) > 0 { // because can not set the "endpointFunc" of WrapServer, do this so so
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
			c.wrapserver.ServeHTTP(resp, ctx.Request)
		} else {
			ctx.String(http.StatusOK, "no rpc")
		}
	}

	url := "/*rpc"

	c.GinRouter.Router().POST(url, handle)
	c.GinRouter.Router().OPTIONS(url, handle)
}
