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
Srater
Stopper
Destroyer
```
* Injecter ：主要是组件依赖注入，增删改查组件， 通过这个接口增加的组件，创建过程是自己完成，它是Line的一部分   

组件的运行过程如下：  
* Newer:  创建组件  
* Creator: 创建完成后调用  
* Starter: 所有有Create之后调用  
* 程序运行  
* Stopper: 程序将要退出时调用  
* Destroyer: 所有的stop调用后再调用  

可以通过配置文件或代码，来说明组件之间的关系， 这时line会计算组件之间的依赖关系，使用者不用管它们的创建顺序  

# 默认组件
## 配置 dots/sconfig
现在配置支持json格式，以后会支持toml、yaml、命令行及环境变量
## 日志 dots/slog
基于zap的日志
## grpc的客户端负载均衡 dots/grpc/conns
提供grpc的客户端负载均衡， sample/grpc_conns是一个使用的例子
## 证书生成组件 dots/certificate
生成根证书及子证书， sample/certificate 是一个使用的例子

