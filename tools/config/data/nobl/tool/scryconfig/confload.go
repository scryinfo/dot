package scryconfig

import (
	"flag"
	"fmt"
	"github.com/gookit/config"
	"github.com/gookit/config/json"
	"sync"
)
//Used to define command line flags
const(
	Json string = "confJson"
	Toml string = "confToml"
	Yaml string = "confYaml"
)
//For storage configuration
type ScryConfig struct {
	Dc                   *config.Config
	FlagJson, FlagToml, FlagYaml string
}
//Configuring storage object instantiation
func New() *ScryConfig {
	return &ScryConfig{config.New("default"),"","",""}
}
//When the user uses go flag package,
//Configure flags bind before parse
var once sync.Once
func (sc *ScryConfig)BindFlag()  {
	var one = func(){
		flag.StringVar(&sc.FlagJson,Json,"","use for json config")//flag is confJson
		flag.StringVar(&sc.FlagToml,Toml, "", "use for toml config")//flag is confToml
		flag.StringVar(&sc.FlagYaml,Yaml, "", "use for yaml config")//flag is confYaml
	}
	once.Do(one)
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
//Get the Json bytes configuration information of this part by key
//Key uses `.` as an interval to represent hierarchical relationships
func (sc *ScryConfig)GetJsonByte(key string) ([]byte, error) {
	var subConf interface{}
	var ok bool
	if key != ""{
		subConf,ok = sc.Dc.Get(key)
		if !ok {
			return nil, fmt.Errorf("can't get data by key")
		}
	}else {
		subConf = sc.Dc.Data()
	}
	jsonByte, err := json.Encoder(subConf)
	if err != nil {
		return nil, err
	}
	return jsonByte,nil
}
