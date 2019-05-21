# dot  
Component development specification, including component definition, component dependencies, component life cycle, dependency injection, and common basic components  
* Dot: A component which has no type or interface requirements, anything can be a component  
* Line: A container that holds components, adds, deletes, modifies, and injects dependencies into components  
* Newer:  Construct component, the Newer is used to construct the component, and if it is not specified, then construct it by default "refect.New"
* Lifer: Is the component life cycle management interface, the implementation of the interface and the method will be automatically run by Line, the following are the four interfaces 
```go
Creator 
Srater
Stopper
Destroyer
```
* Injecter ：It is component dependency injection, adding, deleting and checking components. The creation process of components added through this interface is completed by ourselves, which is part of Line   

The process that component runs as follows：  
* Newer:  Construct component  
* Creator: Call after construct   
* Starter: Call after Create  
* Run programs  
* Stopper: Call before exit programs 
* Destroyer: Call after stop  

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

