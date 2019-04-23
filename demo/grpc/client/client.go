package main

import (
	"fmt"
	proto "github.com/scryInfo/dot/demo/pb"
	"github.com/scryInfo/dot/dot"
	"github.com/scryInfo/dot/dots/grpc/client"
	"github.com/scryInfo/dot/dots/line"
	"log"
)

func main() {
	l, _ := line.BuildAndStart(func(l dot.Line) error {
		gclient.Add(l, "dd05cbec-e3d0-4be3-a7df-87b0522ac46b")
		return nil
	})

	//f := &gclient.GrpcClient{}
	//l.ToInjecter().Inject(f)

	var f gclient.GrpcClienter
	{
		d, _ := l.ToInjecter().GetByLiveId("dd05cbec-e3d0-4be3-a7df-87b0522ac46b")
		f = d.(gclient.GrpcClienter)
	}

	conn := f.GetConn()

	ctx := f.GetCtx()

	//用户实现 start

	c := proto.NewTestClient(conn)

	c1, err := c.SayHello(ctx, &proto.TestRequest{Name: "shrimpliao"})

	fmt.Println("err", err)

	log.Printf("@@@c1: %s", c1.Message)

	//用户实现 end

	line.StopAndDestroy(l, true)
}
