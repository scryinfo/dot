package kits

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/knadh/koanf"
	kjson "github.com/knadh/koanf/parsers/json"
	ktoml "github.com/knadh/koanf/parsers/toml"
	kyaml "github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	kfile "github.com/knadh/koanf/providers/file"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/scryg/sutils/sfile"
)

type _config struct {
}

var Config = _config{}

func (c _config) MakeConfig(sconf dot.SConfig, confg any) error {
	name := filepath.Join(sconf.ConfigPath(), sconf.ConfigFile())
	ext := filepath.Ext(name)
	newName := name[:len(name)-len(ext)] + "_new" + ext
	return c.WriteConfig(name, confg, newName)
}

func (c _config) WriteConfig(file string, pconf any, newFile string) error {
	val := reflect.ValueOf(pconf)
	if val.Kind() != reflect.Pointer {
		return errors.New("the parameter pconf must be a pointer type")
	}
	tagName := "toml"
	var parse koanf.Parser = ktoml.Parser()
	if strings.HasSuffix(file, ".json") {
		tagName = "json"
		parse = kjson.Parser()
	} else if strings.HasSuffix(file, ".yaml") {
		tagName = "yaml"
		parse = kyaml.Parser()
	}
	config := koanf.New(".")
	if sfile.ExistFile(file) {
		if err := config.Load(kfile.Provider(file), parse); err != nil {
			return err
		}
		if err := config.Unmarshal("", pconf); err != nil {
			return err
		}
	}
	kv := config.Raw()
	err := structFields(tagName, kv, reflect.ValueOf(pconf))
	if err != nil {
		return err
	}
	err = config.Load(confmap.Provider(kv, "."), nil)
	if err != nil {
		return err
	}
	data, err := config.Marshal(parse)
	if err != nil {
		return err
	}
	err = os.WriteFile(newFile, data, 0644)
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
