// Scry Info.  All rights reserved.
// license that can be found in the license file.

package gclient

import (
	"encoding/json"
	"log"
	"reflect"
	"time"

	"github.com/scryinfo/dot/dot"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

)

const (
	DotTypeId = "24b8f11f-8600-4786-91e1-0d0dc7bc8969"
)

// Deprecated: Use grpc/conns instead
type GrpcClienter interface {
	GetCtx() context.Context
	GetConn() *grpc.ClientConn
	Stop(ignore bool) error
	Destroy(ignore bool) error
}

type Grpc struct {
	Address string
	ctx     context.Context
	conn    *grpc.ClientConn
	cancel  context.CancelFunc
}

//Add newer to gclient, it will be always used unless specially specified
func AddType(l dot.Line) error {
	err := l.AddNewerByTypeId(DotTypeId, func(conf interface{}) (d dot.Dot, err error) {
		d = &Grpc{}
		err = nil
		t := reflect.ValueOf(conf)
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

func (g *Grpc) GetCtx() context.Context {
	return g.ctx
}

func (g *Grpc) GetConn() *grpc.ClientConn {
	return g.conn
}

func (g *Grpc) Create(conf dot.SConfig) error {
	conn, err := grpc.Dial(g.Address, grpc.WithInsecure())
	g.conn = conn
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	g.ctx = ctx
	g.cancel = cancel
	//defer cancel()
	return err
}

//Start connection
func (g *Grpc) Start(ignore bool) error {
	return nil
}

//Stop
//ignore When calling other Lifer, true erred will continue, false erred will return directly
func (g *Grpc) Stop(ignore bool) error {
	return nil
}

//Destroy Destroy Dot
//ignore When calling other Lifer, true erred will continue, false erred will return directly
func (g *Grpc) Destroy(ignore bool) error {
	g.conn.Close()
	g.cancel()
	return nil
}
