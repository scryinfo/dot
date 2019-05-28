[中文](./README-cn.md)  
[EN](./README.md)  
# dot  
Component development specification, including component definition, component dependencies, component life cycle, dependency injection, and common basic components  
* Dot: A component which has no type or interface requirements, anything can be a component  
* Line: A container that holds components, adds, deletes, modifies, and injects dependencies into components  
* Newer:  Construct component, the Newer is used to construct the component, and if it is not specified, then construct it by default "refect.New"
* Lifer: Is the component life cycle management interface, the implementation of the interface and the method will be automatically run by Line, the following are the four interfaces 
```go
Creator
Starter
Stopper
Destroyer
```
* Injecter ：It is component dependency injection, adding, deleting and checking components. The creation process of components added through this interface is completed by ourselves, which is part of Line   

The process that component runs as follows：  
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
5. Events.BeforeCreate
6. dot.Creator
7. Events.AfterCreate //go to "2. Newer", untill all done  
8. Inject all dependentes of dots  
9. Builder.AfterCreate  
***
Start  
1. Builder.BeforeStart 
2. Events.BeforeStart
3. dot.Starter
4. Events.AfterStart //go to "2. Events.BeforeStart", untill all done
5. dot.AfterAllStart
6. Builder.AfterStart  
***
Stop  
1. Builder.BeforeStop
2. dot.BeforeAllStopper
3. Events.BeforeStop
4. dot.Stopper
5. Events.AfterStop //go to "2. Events.BeforeStop", until all done
6. Builder.AfterStop  
***
Destroy  
1. Builder.BeforeDestroy 
2. Events.BeforeDestroy
3. dot.Destroyer
4. Events.AfterStop //go to "2. Events.BeforeDestroy", until all done
5. Builder.AfterDestroy  

The relationships between components can be set by configuration files or code, Line computes the dependencies between components, regardless of the order in which they are created .

# Default components 
## Config: dots/sconfig
Now use the json format,  later will support toml, yaml, command line, and environment variables.
## Log: dots/slog
High performance logs based on zap.

## GRPC client balance:  dots/grpc/conns
 Client load balancing for GRPC. "sample /grpc_conns" is an example.
## Certificate generated: dots/certificate
Generate root and sub certificates. "sample/certificate" is an example.

# [Code Style -- Go](https://github.com/scryinfo/scryg/blob/master/codestyle_go.md)

