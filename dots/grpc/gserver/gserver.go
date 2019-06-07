// Scry Info.  All rights reserved.
// license that can be found in the license file.

package gserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"reflect"

	"github.com/scryinfo/dot/dot"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	DotId = "634658a8-7598-4785-b4ac-bb201ff0010f"
)

// Deprecated: Use grpc/server/ServerNoblImp instead
type GrpcInterface interface {
	GetServer() *grpc.Server
	StartService()
	Stop(ignore bool) error
	Destroy(ignore bool) error
}

type Grpc struct {
	Port string
	S    *grpc.Server
	lis  net.Listener
}

type GrpcServer struct {
	Grpcs GrpcInterface `dot:"634658a8-7598-4785-b4ac-bb201ff0010f"`
}

func Add(l dot.Line) error {
	err := l.AddNewerByLiveId(dot.LiveId(DotId), func(conf interface{}) (d dot.Dot, err error) {
		d = &Grpc{}
		err = nil
		t := reflect.ValueOf(conf)
		fmt.Println(t.Kind())
		if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
			if t.Len() > 0 && t.Index(0).Kind() == reflect.Uint8 {
				v := t.Slice(0, t.Len())
				json.Unmarshal(v.Bytes(), d)
			}
		} else {
			err = dot.SError.Parameter
		}
		return
	})
	return err
}

func (g *Grpc) GetServer() (s *grpc.Server) {
	return g.S
}

func (g *Grpc) Create(conf dot.SConfig) (err error) {
	g.lis, err = net.Listen("tcp", g.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	g.S = grpc.NewServer()
	return err
}

//Start connection
func (g *Grpc) Start(ignore bool) error {
	return nil
}

//Stop
//ignore When calling other Lifer, if true erred then continue, false erred then return directly
func (g *Grpc) Stop(ignore bool) error {
	g.S.Stop()
	return nil
}

//Destroy Destroy Dot
//ignore When calling other Lifer, if true erred then continue, false erred then return directly
func (g *Grpc) Destroy(ignore bool) error {
	g.lis.Close()
	return nil
}

//Register service and operate
func (g Grpc) StartService() {
	reflection.Register(g.S)
	if err := g.S.Serve(g.lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return
}
