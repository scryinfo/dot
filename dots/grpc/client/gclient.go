package gclient

import (
	"encoding/json"
	"github.com/scryinfo/dot-0/dot"
	"github.com/scryinfo/dot-0/line"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"reflect"
	"time"
)


type GrpcClienter interface {
	GetCtx() context.Context
	GetConn() *grpc.ClientConn
	Stop(ignore bool) error
	Destroy(ignore bool) error
}

type Grpc struct {
	Address string
	ctx context.Context
	conn *grpc.ClientConn
	cancel context.CancelFunc
}

type GrpcClient struct {
	Grpcs GrpcClienter `dot:"dd05cbec-e3d0-4be3-a7df-87b0522ac46a"`
}

func Add(l line.Line)  {
	l.AddNewerByLiveId(dot.LiveId("dd05cbec-e3d0-4be3-a7df-87b0522ac46a"), func(conf interface{}) (d dot.Dot, err error) {
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
}


func (g *Grpc) GetCtx () context.Context {
	return g.ctx
}

func (g *Grpc) GetConn() *grpc.ClientConn {
	return g.conn
}

func (g *Grpc) Create (conf dot.SConfig) error {
	conn,err := grpc.Dial(g.Address,grpc.WithInsecure())
	g.conn = conn
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	g.ctx=ctx
	g.cancel=cancel
	//defer cancel()
	return err
}

//启动连接
func (g *Grpc) Start (ignore bool) error {
	return nil
}

//Stop
//ignore 在调用其它Lifer时，true 出错出后继续，false 出现一个错误直接返回
func (g *Grpc) Stop(ignore bool) error {
	return nil
}

//Destroy 销毁 Dot
//ignore 在调用其它Lifer时，true 出错出后继续，false 出现一个错误直接返回
func (g *Grpc) Destroy(ignore bool) error{
	g.conn.Close()
	g.cancel()
	return nil
}