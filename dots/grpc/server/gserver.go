package gserver

import (
	"encoding/json"
	"fmt"
	"github.com/scryInfo/dot/dot"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"reflect"
)

const (
	DotId = "634658a8-7598-4785-b4ac-bb201ff0010f"
)

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

//启动连接
func (g *Grpc) Start(ignore bool) error {
	return nil
}

//Stop
//ignore 在调用其它Lifer时，true 出错出后继续，false 出现一个错误直接返回
func (g *Grpc) Stop(ignore bool) error {
	g.S.Stop()
	return nil
}

//Destroy 销毁 Dot
//ignore 在调用其它Lifer时，true 出错出后继续，false 出现一个错误直接返回
func (g *Grpc) Destroy(ignore bool) error {
	g.lis.Close()
	return nil
}

//注册服务，并执行
func (g Grpc) StartService() {
	reflection.Register(g.S)
	if err := g.S.Serve(g.lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return
}
