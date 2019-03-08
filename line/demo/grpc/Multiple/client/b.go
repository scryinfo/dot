package client

import (
	"fmt"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/grpc/client"
	pb "github.com/scryinfo/dot/line/demo/pb"
	"github.com/scryinfo/dot/line/lineimp"
	"log"
)

func B()  {
	l := lineimp.New()
	l.ToLifer().Create(nil)

	gclient.Add(l,dot.LiveId("dd05cbec-e3d0-4be3-a7df-87b0522ac46b"))

	err := l.Rely()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = l.CreateDots()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = l.ToLifer().Start(false)
	if err != nil {
		fmt.Println(err)
		return
	}

	var f gclient.GrpcClienter
	{
		d, _:= l.ToInjecter().GetByLiveId("dd05cbec-e3d0-4be3-a7df-87b0522ac46b")
		f = d.(gclient.GrpcClienter)
	}



	conn := f.GetConn()

	ctx := f.GetCtx()


	//用户实现 start

	c := pb.NewTestClient(conn)

	c1, err := c.SayHello(ctx, &pb.TestRequest{Name: "shrimpliao"})

	fmt.Println("err",err)

	log.Printf("@@@c1: %s", c1.Message)

	//用户实现 end

	f.Stop(false)
	f.Destroy(false)
}