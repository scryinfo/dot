// Scry Info.  All rights reserved.
// license that can be found in the license file.

package gserver

import (
	"github.com/gin-gonic/gin"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/gindot"
	"google.golang.org/grpc"
)

const (
	GinNoblTypeId = "afbeac47-e5fd-4bf3-8fb1-f0fb8ec79bd0"
)

//support the http and tcp
type ginNobl struct {
	ServerNobl ServerNobl `dot:""`
	GinRouter  *gindot.Router
	wrapserver *grpcweb.WrappedGrpcServer
}

//GinNoblTypeLives Data structure needed when generating newer component
func GinNoblTypeLives() []*dot.TypeLives {

	tl := &dot.TypeLives{
		Meta: dot.Metadata{TypeId: GinNoblTypeId, NewDoter: func(conf interface{}) (dot dot.Dot, err error) {
			return newHttpNobl(conf)
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
	//options.OptionsPassthrough
	c.wrapserver = grpcweb.WrapServer(c.Server(), grpcweb.WithAllowedRequestHeaders([]string{"Access-Control-Allow-Origin:*", "Access-Control-Allow-Methods:*"}))

	c.GinRouter.Router().Use(func(g *gin.Context) {

		if c.wrapserver.IsGrpcWebRequest(g.Request) {

			resp := g.Writer
			resp.Header().Set("Access-Control-Allow-Origin", "*")  //
			resp.Header().Set("Access-Control-Allow-Methods", "*") //
			resp.Header().Add("Access-Control-Allow-Headers", "content-type,x-grpc-web,x-user-agent")
			c.wrapserver.ServeHTTP(resp, g.Request)
			return
		}
		g.Next()

	})
}
