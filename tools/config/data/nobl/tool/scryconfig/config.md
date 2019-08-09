#Component configuration
## Features
 - Load from `configuration files` which supports multiple formats `Json`, `Toml`, `Yaml`.
   - By default, the configuration file is loaded from the path where the executable is located.
   - It is also possible to load the configuration according to the configPaths parameter passed in by the user.
   - The configPaths passed by the user can be an absolute path or a path relative to the path of the executable file.
   - Support for incoming multiple paths.
   - The priority of the configuration files loading is `Json` > `Toml` > `Yaml`
 - Load from `cmd flags` which supports multiple formats `Json`, `Toml`, `Yaml`. `Cmd flags` has a higher priority than the `configuration files`.
   - Get three types of configuration information yaml, toml, json through function parameters and load by the priority json > toml > yaml
   - When the user uses go flag package, user passes the method ConFlag() as a ConfLoad() first parameter.
   - When the user uses other flag packages, user passes custom method which `signature` is func ()(flagJson, flagToml, flagYaml string).
   - If there is no data of this type, please return "".
 - Get `partial configuration information` in `Json` format.
   - Get the Json bytes configuration information of this part by key.
   - Key uses `.` as an interval to represent hierarchical relationships.
## Usage
 - First, instantiate configuration information object.
 ```go
    //Configuring storage object instantiation
    func New() *ScryConfig {
    	return &ScryConfig{config.New("default"),"","",""}
    }
 ```
 - Sencond, divided into two cases:
   - One case, use go flag package.
     - First, use `BindFlag()` to bind the configuration flag.
       - Guarantee `BindFlag()` to be called before `flag.Parse()`
     - Second, call `ConfLoad()` and use `BindFlag` as the first parameter. 
   - Another case, use other flag packages.
     - First, define a function signature as `() (flagJson, flagToml, flagYaml string)` .
       - Ensure that the returned parameters correspond to the format.
       - if don't have the format, return "".
     - Second, use such a function as the first argument of `ConfLoad()` and call `ConfLoad()`.
```go
    //When the user uses go flag package,
    //Configure flags bind before parse
    func (sc *ScryConfig)BindFlag()  {
    	flag.StringVar(&sc.FlagJson,Json,"","use for json config")   //flag is confJson
    	flag.StringVar(&sc.FlagToml,Toml, "", "use for toml config") //flag is confToml
    	flag.StringVar(&sc.FlagYaml,Yaml, "", "use for yaml config") //flag is confYaml
    }
    //When the user uses go flag package,
    //user passes this method as a ConfLoad() first parameter
    //When the user uses other flag packages,
    //user passes custom method
    // which signature is func ()(flagJson, flagToml, flagYaml string)
    //If there is no data of this type, return ""
    func (sc *ScryConfig)ConFlag() (flagJson, flagToml, flagYaml string) {
    	if !flag.Parsed(){
    		flag.Parse()
    	}
    	flagJson = sc.FlagJson
    	flagToml = sc.FlagToml
    	flagYaml = sc.FlagYaml
    	return
    }
    //Complete the entire load configuration process
    //as followed:
    //
    //By default, the configuration file is loaded from the path
    // where the executable is located.
    //It is also possible to load the configuration
    // according to the configPaths parameter passed in by the user.
    //The configPaths passed by the user can be an absolute path
    // or a path relative to the path of the executable file.
    //Support for incoming multiple paths.
    //The priority of the configuration files loading is json > toml > yaml
    //
    //Get three types of configuration information yaml, toml, json through function parameters
    //and load by the priority json > toml > yaml
    //
    //Return loaded data
    func (sc *ScryConfig)ConfLoad(cmd func()(flagJson, flagToml, flagYaml string),configPaths ...string) (map[string]interface{}, error){
    	_, err := sc.LoadConfigFile(configPaths...)
    	if err != nil {
    		return nil,err
    	}
    	conJson, conToml, conYaml := cmd()
    	if conYaml != ""{
    		err = sc.Dc.LoadSources(config.Yaml,[]byte(conYaml))
    		if err != nil{
    			return nil, err
    		}
    	}
    	if conToml != ""{
    		err = sc.Dc.LoadSources(config.Toml,[]byte(conToml))
    		if err != nil{
    			return nil, err
    		}
    	}
    	if conJson != ""{
    		err = sc.Dc.LoadSources(config.JSON,[]byte(conJson))
    		if err != nil{
    			return nil, err
    		}
    	}
    	return sc.Dc.Data(), nil
    }

```
- At last, Get the Json bytes configuration information of this part by key.
   	- Key uses `.` as an interval to represent hierarchical relationships.
```go
    //Get the Json bytes configuration information of this part by key
    //Key uses `.` as an interval to represent hierarchical relationships
    func (sc *ScryConfig)GetJsonByte(key string) ([]byte, error) {
    	subConf,ok := sc.Dc.Get(key)
    	if !ok {
    		return nil, fmt.Errorf("can't get data by key")
    	}
    	jsonByte, err := json.Encoder(subConf)
    	if err != nil {
    		return nil, err
    	}
    	return jsonByte,nil
    }
```
Here is a vary simple example:

```go
    scry := scryconfig.New()
	scry.BindFlag()
	_, err :=scry.ConfLoad(scry.ConFlag,"testdata")
	if err != nil {
		fmt.Println(err)
	}
	value, err :=scry.GetJsonByte("lang.allowed")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(value))
```

