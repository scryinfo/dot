[中文](./README-cn.md)  
[EN](./README.md)  
# dot 
组件开发规范，主要有组件定义、组件依赖关系、组件生命周期、依赖注入、及常用的基础组件  
* Dot: 组件，它没有类型或接口要求，都可以成为一个组件  
* Line: 存放组件的容器，对组件进行增删除改查，依赖注入  
* Newer: 组件的创建方法，line通过newer来创建组件，如果没有指定，那么默认使用 refect.new来创建  
* Lifer: 是组件生命周期管理接口，实现接口后方法会被Line自动运行，下面是这四个接口 
```
Creator  // 
Starter
Stopper
Destroyer
```
* Injecter ：主要是组件依赖注入，增删改查组件， 通过这个接口增加的组件，创建过程是自己完成，它是Line的一部分   

组件的运行过程如下：  
***
Create Config and Log  
1. Make Default Log
2. Make Config 
3. Make Log of config
***
Create  
1. Builder.BeforeCreate 
2. dot.Newer
3. dot.SetterLine
4. dot.SetterTypeAndLiveId
5. Events.BeforeCreate //for type id
6. Events.BeforeCreate //for live id
7. dot.Creator
8. Events.AfterCreate //for live id
9. Events.AfterCreate //for type id, go to "2. Newer", until all done  
10. Inject all dependent dots  
11. AfterAllInjecter  
12. Builder.AfterCreate  
***
Start  
1. Builder.BeforeStart 
2. Events.BeforeStart
3. Events.BeforeStart //for live id
4. dot.Starter
5. Events.AfterStart //for live id
6. Events.AfterStart //go to "2. Events.BeforeStart", until all done
7. dot.AfterAllStart
8. Builder.AfterStart  
***
Stop  
1. Builder.BeforeStop
2. dot.BeforeAllStopper
3. Events.BeforeStop //for type id
4. Events.BeforeStop //for live id
5. dot.Stopper
6. Events.AfterStop //for live id
7. Events.AfterStop //for type id go to "2. Events.BeforeStop", until all done
8. Builder.AfterStop  
***
Destroy  
1. Builder.BeforeDestroy 
2. Events.BeforeDestroy //for type id
3. Events.BeforeDestroy //for live id
4. dot.Destroyer
5. Events.AfterStop //for live id
6. Events.AfterStop //for type id go to "2. Events.BeforeDestroy", until all done
7. Builder.AfterDestroy  

可以通过配置文件或代码，来说明组件之间的关系， 这时line会计算组件之间的依赖关系，使用者不用管它们的创建顺序  

# 默认组件
## 配置 dots/sconfig
现在配置支持json格式，以后会支持toml、yaml、命令行及环境变量
## 日志 dots/slog
基于zap的日志
## grpc组件
dots/grpc/conns： 客户端负载均衡组件， 支持服务端tls, 双向tls认证  
dots/grpc/gserver/http_nobl: 进程内的grpc-web支持，支持https， sample/grpc/http是使用例子  
* 运行build_go.bat生成 js及go代码
* 编译sample/grpc/http/server中的代码
* 运行编译出来的server.exe (server.exe --configfile="server_http.json"), 注：如果可执行文件名与配置文件同名，可以不指定配置文件 
* cd 进入 sample/grpc/http/client
* npm install
* npm run start
* 在浏览器中输入 “http://localhost:9000”, 查看控制台的输出，注：如果这里使用了https，是自己签名的证书时需要加入允许证书中或参考网上关于证书的说明  

dots/grpc/gserver/server_nobl： grpc server组件，支持服务端tls, 双向tls认证， sample/grpc/nobl是使用例子  
提供注：与tls或https相关的例子在 sample/grpc/tls文件夹下面  
## 证书生成组件 dots/certificate
生成根证书及子证书， sample/certificate 是一个使用的例子
## gin组件
简化gin的使用，且把slog与gin的日志整合在一起， sample/gindot 是一个使用例子

# [Code Style -- Go](https://github.com/scryinfo/scryg/blob/master/codestyle_go-cn.md)
