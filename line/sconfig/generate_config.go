package sconfig

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/scryg/sutils/sfile"
	"github.com/spf13/viper"
)

var OkGenerateConfig = fmt.Errorf("generate config ok")

// GenerateConfig generates a config file from the config struct
// config: the config struct
// rootConfig: the root config struct
// args: whether to parse command line arguments
// returns: whether to make config, error
func GenerateConfigWithArgs[T any](config dot.SConfig, rootConfig *T) (*T, error) {
	generateConfig := false
	{
		// parse command line arguments "GenerateConfig", dont use the flag package
		if len(os.Args) > 1 {
			name := "GenerateConfig"
			for i, it := range os.Args[1:] {
				it = strings.TrimSpace(it)
				suf := strings.TrimPrefix(it, "-"+name+"=")
				if suf == it {
					suf = strings.TrimPrefix(it, "--"+name+"=")
				}
				if suf == it {
					if it == "-"+name || it == "--"+name {
						if i < len(os.Args)-1 {
							itNext := strings.TrimSpace(os.Args[i+1])
							if !strings.HasPrefix(itNext, "-") {
								suf = itNext
							} else {
								suf = ""
							}
						} else {
							// the last argument is "GenerateConfig",
							suf = ""
						}
					}
				}

				if suf == "" {
					// the argument is "GenerateConfig",
					// no value, -GenerateConfig=true or --GenerateConfig=true or -GenerateConfig true or --GenerateConfig true)
					generateConfig = true
					break
				} else if suf != it {
					// the argument is "GenerateConfig",
					// has value, parse it as bool
					v, err := strconv.ParseBool(suf)
					if err == nil {
						generateConfig = v
					} else {
						fmt.Printf("cant parse GenerateConfig value: %s, set the GenerateConfig to false, err: %v\n", suf, err)
						generateConfig = false
					}
					break
				} else {
					// dont find the argument GenerateConfig, go next
				}
			}
		}
	}

	if generateConfig {
		ok, err := GenerateConfigGo(config, rootConfig)
		if err != nil {
			return nil, err
		}
		if ok {
			return rootConfig, OkGenerateConfig
		}
	}
	return rootConfig, nil
}

// GenerateConfigGo generates a config file from the config struct
// config: the config struct
// rootConfig: the root config struct
// returns: whether to make config, error
func GenerateConfigGo[T any](config dot.SConfig, rootConfig *T) (bool, error) {
	err := MergeConfig(config, rootConfig)
	if err != nil {
		// dont use the logger here, the logger is not initialized yet
		fmt.Printf("make config err: %v\n", err)
		return false, err
	} else {
		fmt.Println("make config success")
		return true, nil
	}
}

func MergeConfig[T any](sconf dot.SConfig, confg *T) error {
	name := filepath.Join(sconf.ConfigPath(), sconf.ConfigFile())
	ext := filepath.Ext(name)
	newName := name[:len(name)-len(ext)] + "_gen" + ext
	return MergeConfigToNew(name, confg, newName)
}

func MergeConfigToNew[T any](file string, pconf *T, newFile string) error {
	val := reflect.ValueOf(pconf)
	if val.Kind() != reflect.Pointer {
		return errors.New("the parameter pconf must be a pointer type")
	}
	tagName := "toml"
	if strings.HasSuffix(file, ".json") {
		tagName = "json"
	} else if strings.HasSuffix(file, ".yaml") {
		tagName = "yaml"
	}
	config := viper.New()
	if sfile.ExistFile(file) {
		config.SetConfigFile(file)
		config.SetConfigType(tagName)
		if err := config.ReadInConfig(); err != nil {
			return err
		}
		err := config.Unmarshal(pconf, func(dc *mapstructure.DecoderConfig) {
			// dc.ErrorUnused = true
			// dc.ErrorUnset = true
			dc.TagName = tagName
		})
		if err != nil {
			return err
		}
	}
	kv := config.AllSettings()
	err := structFields(tagName, kv, reflect.ValueOf(pconf))
	if err != nil {
		return err
	}
	err = config.MergeConfigMap(kv)
	if err != nil {
		return err
	}
	err = config.WriteConfigAs(newFile)
	if err != nil {
		return err
	}
	return nil
}

func structFields(tagName string, kv map[string]any, s reflect.Value) error {
	if s.Kind() == reflect.Pointer {
		s = s.Elem()
	}
	typ := s.Type()
	if typ.Kind() != reflect.Struct {
		return nil
	}
	for i := 0; i < typ.NumField(); i++ {
		fieldType := typ.Field(i)
		fieldValue := s.Field(i)
		if fieldType.Type.Kind() == reflect.Pointer && fieldType.Type.Elem().Kind() == reflect.Struct {
			if fieldValue.IsNil() {
				fieldValue.Set(reflect.New(fieldType.Type.Elem()))
			}
		}

		key := fieldType.Tag.Get(tagName)
		if key == "" {
			key = fieldType.Tag.Get("mapstructure")
			if key == "" {

			}
			key = fieldType.Name
		}
		if v, ok := kv[key]; ok {
			fType := fieldType.Type
			if fieldType.Type.Kind() == reflect.Pointer {
				fType = fType.Elem()
			}
			if fType.Kind() == reflect.Struct {
				if vs, ok := v.(map[string]any); ok {
					if err := structFields(tagName, vs, fieldValue); err != nil {
						return err
					}
				} else {
					return errors.New("the value of key " + key + " must be a map[string]any")
				}

			} else {
				// the field is exist in the config, do nothing
			}
		} else {

			if fieldType.Type.Kind() == reflect.Struct || (fieldType.Type.Kind() == reflect.Pointer && fieldType.Type.Elem().Kind() == reflect.Struct) {
				subStruct := make(map[string]any)
				kv[key] = subStruct
				if err := structFields(tagName, subStruct, fieldValue); err != nil {
					return err
				}
			} else {
				kv[key] = fieldValue.Interface()
			}
		}

	}
	return nil
}
