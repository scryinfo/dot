package kits

import (
	"errors"
	"reflect"
	"strings"

	"github.com/scryinfo/scryg/sutils/sfile"
	"github.com/spf13/viper"
)

type _config struct {
}

var Config = _config{}

func (c _config) WriteConfig(file string, pconf any, newFile string) error {
	val := reflect.ValueOf(pconf)
	if val.Kind() != reflect.Ptr {
		return errors.New("the parameter pconf must be a pointer type")
	}
	config := viper.New()
	if sfile.ExistFile(file) {
		config.SetConfigFile(file)
		if err := config.ReadInConfig(); err != nil {
			return err
		}
		if err := config.Unmarshal(pconf); err != nil {
			return err
		}
	}
	tagName := "toml"
	if strings.HasSuffix(file, ".json") {
		tagName = "json"
	} else if strings.HasSuffix(file, ".yaml") {
		tagName = "yaml"
	}
	err := structFields(tagName, config.AllSettings(), reflect.ValueOf(pconf))
	if err != nil {
		return err
	}
	config.SetConfigFile(newFile)
	if err := config.WriteConfig(); err != nil {
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
			val := reflect.ValueOf(v)
			if val.Kind() == reflect.Pointer {
				val = val.Elem()
			}
			if val.Kind() == reflect.Struct {
				if vs, ok := v.(map[string]any); ok {
					if err := structFields(tagName, vs, val); err != nil {
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
