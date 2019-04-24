# dot
Componentized development framework for Scryinfo
主要功能是实例管理及依赖注入，是简化了的Spring
dot: 组件，没有类型或接口要求，任何一个都可以
line: 存放的组件的，对组件进行增删除改查
newer: 组件的创建方法，line通过newer来创建组件，如果没有指定，那么默认使用 refect.new来创建
lifer: 是组件生命周期管理的接口，如果实现这个接口或它的父接口，对应的方法会运行，下面是这四个接口
```
Creater  // 
Srater
Stopper
Destroyer
```
组件的运行过程如下：
Newer:  创建组件
Creater: 创建完成后调用
Starter: 所有有Create之后调用
程序运行
Stopper: 程序将要退出时调用
Destroyer: 所有的stop调用后再调用

可以通过配置文件或代码，来说明组件之间的关系， 这时line会计算组件之前的依赖关系，使用者不用关系它们的创建顺序

Injecter ：主要是组件注入，增删改查组件， 通过这个接口增加的组件，创建过程是自己完成

