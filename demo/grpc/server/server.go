package main

import (
	proto "github.com/scryinfo/dot/demo/pb"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/grpc/server"
	"github.com/scryinfo/dot/dots/line"
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

func main() {

	l, _ := line.BuildAndStart(func(l dot.Line) error {
		gserver.Add(l)
		return nil
	})

	l.SLogger().SetLogFile("out1.log")

	l.SLogger().Debugln("123123123123")

	//d := l.Config().Dots

	f := &gserver.GrpcServer{}

	l.ToInjecter().Inject(f)

	//pp := l.Config().Dots

	//fmt.Println(pp)

	//dddddddd:= f.Grpcs.GetD()
	//fmt.Println(dddddddd)

	//用户实现 start
	s := f.Grpcs.GetServer()
	proto.RegisterTestServer(s, &servers{})
	proto.RegisterGreeterServer(s, &servers{})

	//用户实现 end

	f.Grpcs.StartService()

	//关闭
	line.StopAndDestroy(l, true)

}

type servers struct{}

func (s *servers) SayHello(ctx context.Context, in *proto.TestRequest) (*proto.TestReply, error) {
	return &proto.TestReply{Message: "SayHello" + in.Name}, nil
}

func (s *servers) SayHello1(ctx context.Context, in *proto.HelloRequest) (*proto.HelloReply, error) {
	return &proto.HelloReply{Message: "SayHello1" + in.Name}, nil
}
