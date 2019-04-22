package main

import (
	"fmt"
	"github.com/scryInfo/dot/dots/grpc/server"
	pb "github.com/scryInfo/dot/line/demo/pb"
	"github.com/scryInfo/dot/line/lineimp"
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

	l := lineimp.New()
	//l.
	l.ToLifer().Create(nil)

	//l.(*lineimp).CreateLog(nil,"ut.log")

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

	//ccc := l.SCon
	//
	//
	// fig()

	//ccc.Start(true)

	//fmt.Println(llll.GetLevel())#243447
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
	pb.RegisterTestServer(s, &servers{})

	//用户实现 end

	f.Grpcs.StartService()

	//关闭
	f.Grpcs.Stop(false)
	f.Grpcs.Destroy(false)

}

type servers struct{}

func (s *servers) SayHello(ctx context.Context, in *pb.TestRequest) (*pb.TestReply, error) {
	return &pb.TestReply{Message: "SayHelloBBBB" + in.Name}, nil
}
