# 使用方法（已放弃）
> go redis orm，符合我们要求的，经筛选仅有zoom库，但该库不支持redis分片与集群。  
> 所以目前放弃使用orm

> 这个东西是用来验证redis作为组件，调用它的可行性的，如无特殊需要，不建议使用。  
> 也就是说你可以关掉这个文件了 ^_^

步骤：
 - clone dot代码到本地
 - 修改**main.json**配置文件，主要是数据库地址
 - 像其他dot组件一样启动（主要指配置文件路径）
 - 启动之后，即可在**Start**部分的日志找到demo执行记录
 
说明：
 - 将ModelImp匿名组合进你的自定义结构体即可：（`redis:"-"`表示该字段不保存在redis）
    ```go 
    type Demo struct {
    Str                string `redis:"key_str"`
    I                  int    `redis:"key_int"`
    dot_redis.ModelImp `redis:"-"`
    }
    ```
 - 结构体采用hash表存储，表名（hash表的key）形如：`demo:3qVkFK00046Hmy2yRjj7gk`，其结构被`:`分为两部分。  
   前面一部分为`collection`的`name`，可以在注册`collection`时指定：
   ```go 
	demoInsWithName := &Demo{}
	demoInsWithName.SetModelName("demo")

	c.Collections = make([]*dot_redis.Collection, 0)
	c.Collections = append(c.Collections, c.Redis.RegisterCollections([]dot_redis.Model{demoInsWithName})...)
   ```
   后面一部分是`model`的`id`，如未初始化，则会自动生成一段随机`id`
# see
http://doc.redisfans.com
# redis 可视化工具
[fastonosql](https://github.com/fastogt/fastonosql/releases)
[RedisStudio window](https://github.com/cinience/RedisStudio/releases)