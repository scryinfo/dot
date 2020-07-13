# 使用方法
> 这个东西是用来验证redis作为组件，调用它的可行性的，如无特殊需要，不建议使用。  
> 也就是说你可以关掉这个文件了 ^_^

步骤：
 - clone dot代码到本地
 - 将下面的**test code**代码复制到**dot/dots/db/redis/redis.go**（放到文件的任意位置即可，不排除后续代码结构调整、路径改变的可能）
 - 修改**main.json**配置文件，主要是数据库地址
 - 像其他dot组件一样使用（主要指配置文件路径）

# test code
```go 
// For test only
func (c *Redis) Start(_ bool) error {
	dot.Logger().Debugln("Node: show redis param", 
		zap.Any("conf", c.conf),
		zap.Bool("db init success?", c.DB != nil))

	c.DB.FlushAll(c.DB.Context())

	// simulate query in cache first (no result)
	v1, err := c.DB.Get(c.DB.Context(), "demo").Result()
	if err != redis.Nil {
		dot.Logger().Infoln("Example: get value not run as suppose", zap.NamedError("error", err))
		return nil
	}

	// skip query in db, only simulate update cache
	if err = c.DB.Set(c.DB.Context(), "demo", "basic process demo", 0).Err(); err != nil {
		dot.Logger().Infoln("Example: set value failed, error:" + err.Error())
		return nil
	}

	// suppose a request comes now, query in cache (has result)
	v2, err := c.DB.Get(c.DB.Context(), "demo").Result()
	if err != nil {
		dot.Logger().Infoln("Example: get value failed, error:" + err.Error())
		return nil
	}

	dot.Logger().Infoln(`Node: value suppose:{"v1 (nil) ": "", "v2": "basic process demo"}`)
	dot.Logger().Infoln("Node: show values", zap.String("v1 (nil) ", v1), zap.String("v2", v2))

    c.DB.FlushAll(c.DB.Context())

	return nil
}
```
