package main

import (
	"fmt"
	pb"github.com/scryinfo/dot-0/line/demo/pb"
	"github.com/scryinfo/dot-0/dots/grpc/server"
	"github.com/scryinfo/dot-0/line/lineimp"
	"golang.org/x/net/context"
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

	//llll := l.SLogger()

	//ccc := l.SConfig()

	//ccc.Start(true)

	//fmt.Println(llll.GetLevel())

	l.SLogger().Debugln("123123123123")

	f := &gserver.GrpcServer{}

	l.ToInjecter().Inject(f)

	//用户实现 start
	s := f.Grpcs.GetServer()
	pb.RegisterTestServer(s, &servers{})
	pb.RegisterGreeterServer(s, &servers{})

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