package main

import (
	"fmt"
	pb"github.com/scryinfo/dot/line/demo/pb"
	"github.com/scryinfo/dot/dots/grpc/server"
	"github.com/scryinfo/dot/line/lineimp"
	"golang.org/x/net/context"
	"time"
)

//使用注释
// 一 流程
/*
1 创建组件模块 l := lineimp.New()
2 初始化组件 l.ToLifer().Create(nil)
2.1 添加组件 server.Add(l)
3 判断依赖 err := l.Rely()
4 创建dots err = l.CreateDots()
5 组件模块启动 err = l.ToLifer().Start(false)
6 用户初始化组件对象 f := &server.GrpcServer{}
7 注入组件 l.ToInjecter().Inject(f)
...
停止，关闭，销毁
	f.Grpcs.Stop(false)
	f.Grpcs.Destroy(false)
*/
// 二 用户准备
/*
用户需要实现 grpc server 具体方法内容
*/



func main()  {

	l := lineimp.New()
	l.ToLifer().Create(nil)

	//l.(*lineimp).CreateLog(nil,"out.log")

	gserver.Add(l)

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

	//l.SLogger().SetLevel()

	//llll := l.SLogger()

	//ccc := l.SConfig()

	//ccc.Start(true)

	//fmt.Println(llll.GetLevel())#243447

	l.SLogger().Debugln("123123123123")

	//d := l.Config().Dots

	f := &gserver.GrpcServer{}

	l.ToInjecter().Inject(f)

	//pp := l.Config().Dots

	//fmt.Println(pp)

	dddddddd:= f.Grpcs.GetD()
	fmt.Println(dddddddd)

	//用户实现 start
	s := f.Grpcs.GetServer()
	pb.RegisterTestServer(s, &servers{})
	pb.RegisterGreeterServer(s, &servers{})
	pb.RegisterRechargeServer(s,&servers{})

	//用户实现 end

	f.Grpcs.StartService()

	//关闭
	f.Grpcs.Stop(false)
	f.Grpcs.Destroy(false)

}

type servers struct{}

func (s *servers) SayHello(ctx context.Context, in *pb.TestRequest) (*pb.TestReply, error) {
	return &pb.TestReply{Message: "SayHello" + in.Name}, nil
}

func (s *servers) SayHello1(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "SayHello1" + in.Name}, nil
}

func (s *servers) GetRechargeAddr(ctx context.Context,in *pb.GetRechargeAddrRequest)(*pb.GetRechargeAddrReply,error){
	return &pb.GetRechargeAddrReply{EeeAddr:"123",Ts:time.Now().Unix(),MetaData:"哈哈哈",PubKey:[]byte("123456"),SignedTxHash:[]byte("456789")},nil
}