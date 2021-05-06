package data

import (
	"fmt"
	"github.com/gookit/config/json"
	"github.com/gookit/config/toml"
	"github.com/gookit/config/yaml"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

//By default, the configuration file is loaded from the path
// where the executable is located.
//It is also possible to load the configuration
// according to the configPaths parameter passed in by the user.
//The configPaths passed by the user can be an absolute path
// or a path relative to the path of the executable file.
//Support for incoming multiple paths.
//The priority of the configuration files loading is json > toml > yaml
func (sc *Config) LoadConfigFile(configPaths ...string) (map[string]interface{}, error) {
	allFilesAbs := make([]string, 0)
	{
		for _, configPath := range configPaths {
			configPath, err := pathDeal(configPath)
			if err != nil {
				return nil, err
			}
			s, err := os.Stat(configPath)
			if err != nil {
				return nil, err
			}
			if s.IsDir() {
				allFiles, err := ioutil.ReadDir(configPath)
				if err != nil {
					return nil, err
				}
				for _, allFile := range allFiles {
					if !allFile.IsDir() {
						allFilesAbs = append(allFilesAbs, filepath.Join(configPath, allFile.Name()))
					}
				}
			} else {
				allFilesAbs = append(allFilesAbs, configPath)
			}
		}
	}

	var jsonfile, tomlfile, yamlfile []string
	{
		jsonfile = make([]string, 0)
		tomlfile = make([]string, 0)
		yamlfile = make([]string, 0)
		var config_exist bool = false
		for _, fi := range allFilesAbs {
			fileExt := path.Ext(fi)
			if fileExt == ".json" {
				jsonfile = append(jsonfile, fi)
				config_exist = true
			} else if fileExt == ".toml" {
				tomlfile = append(tomlfile, fi)
				config_exist = true
			} else if fileExt == ".yml" || fileExt == ".yaml" {
				yamlfile = append(yamlfile, fi)
				config_exist = true
			} else {
				continue
			}
		}
		if !config_exist {
			return nil, fmt.Errorf("%s", "configfiles are not exist")
		}
	}

	var configdata map[string]interface{}
	{
		sc.Dc.AddDriver(json.Driver)
		sc.Dc.AddDriver(toml.Driver)
		sc.Dc.AddDriver(yaml.Driver)
		sc.Dc.LoadFiles(yamlfile...)
		sc.Dc.LoadFiles(tomlfile...)
		sc.Dc.LoadFiles(jsonfile...)
		configdata = sc.Dc.Data()
	}
	return configdata, nil
}

//Process the path entered by the user.
// If it is a relative path,
// convert it to an absolute path
// corresponding to the relative path
// based on the path where the executable file is located.
func pathDeal(userpath string) (string, error) {
	if filepath.IsAbs(userpath) {
		return userpath, nil
	}
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	exPath := filepath.Dir(ex)
	return filepath.Join(exPath, userpath), nil
}
