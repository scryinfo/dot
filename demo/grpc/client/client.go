package main

import (
	"fmt"
	"github.com/scryInfo/dot"
	"github.com/scryInfo/dot/dots/grpc/client"
	"github.com/scryInfo/dot/dots/line"
	"log"
)

func main() {
	l := line.New()
	l.ToLifer().Create(nil)

	gclient.Add(l, "dd05cbec-e3d0-4be3-a7df-87b0522ac46b")

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

	c := dot.NewTestClient(conn)

	c1, err := c.SayHello(ctx, &dot.TestRequest{Name: "shrimpliao"})

	fmt.Println("err", err)

	log.Printf("@@@c1: %s", c1.Message)

	//用户实现 end

	f.Stop(false)
	f.Destroy(false)
}
