package import_config

import (
	"fmt"
	"github.com/gookit/config"
	"github.com/gookit/config/json"
)

//For storage configuration
type Config struct {
	Dc *config.Config
}

//Configuring storage object instantiation
func New() *Config {
	return &Config{config.New("default")}
}

//Return loaded data
func (sc *Config) ConfLoad(configPaths ...string) (map[string]interface{}, error) {
	_, err := sc.LoadConfigFile(configPaths...)
	if err != nil {
		return nil, err
	}
	return sc.Dc.Data(), nil
}

//Get the Json bytes configuration information of this part by key
//Key uses `.` as an interval to represent hierarchical relationships
func (sc *Config) GetJsonByte(key string) ([]byte, error) {
	var subConf interface{}
	var ok bool
	if key != "" {
		subConf, ok = sc.Dc.Get(key)
		if !ok {
			return nil, fmt.Errorf("can't get data by key")
		}
	} else {
		subConf = sc.Dc.Data()
	}
	jsonByte, err := json.Encoder(subConf)
	if err != nil {
		return nil, err
	}
	return jsonByte, nil
}
