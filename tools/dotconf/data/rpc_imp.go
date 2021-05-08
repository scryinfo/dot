// Scry Info.  All rights reserved.
// license that can be found in the license file.

package data

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/grpc/gserver"
	"github.com/scryinfo/dot/tools/dotconf/data/rpc"
)

const (
	ServerTypeID = "rpcImplement"
)

type RpcImplement struct {
	ServerNobl gserver.ServerNobl `dot:""`
}

func (c *RpcImplement) Start(ignore bool) error {
	rpc.RegisterDotConfigFaceServer(c.ServerNobl.Server(), c)
	return nil
}

//RpcImplementTypeLives make all type lives
func RpcImplementTypeLives() []*dot.TypeLives {
	lives := []*dot.TypeLives{{
		Meta: dot.Metadata{TypeID: ServerTypeID, NewDoter: func(conf []byte) (dot.Dot, error) {
			return &RpcImplement{}, nil
		}},
		Lives: []dot.Live{
			dot.Live{
				LiveID:    ServerTypeID,
				RelyLives: map[string]dot.LiveID{"ServerNobl": gserver.ServerNoblTypeID},
			},
		},
	}}
	lives = append(lives, gserver.ServerNoblTypeLives()...)
	return lives
}

//rpc implement

func (c *RpcImplement) FindDot(_ context.Context, in *rpc.FindReq) (*rpc.FindRes, error) {
	res := &rpc.FindRes{}

	dirs := in.Dirs
	bytes, invalidDirectory, err := FindDots(dirs)
	//删除运行时产生的中间文件
	{
		err := os.Remove("./run_out/callMethod.go")
		err = os.Remove("./run_out/result.json")
		if err != nil {
			dot.Logger().Error(err.Error)
		}
	}
	if err != nil {
		res.Error = err.Error()
	} else {
		res.DotsInfo = string(bytes)
		res.NoExistDirs = invalidDirectory
	}
	return res, nil
}

func (c *RpcImplement) ImportByDot(_ context.Context, in *rpc.ImportReq) (*rpc.ImportRes, error) {

	res := &rpc.ImportRes{}
	data, err := ioutil.ReadFile(in.Filepath)
	if err != nil {
		res.Error = err.Error()
		dot.Logger().Errorln("File reading error", zap.Error(err))
	} else {
		res.Json = string(data)
	}
	return res, nil
}

//支持三种格式json toml yaml
func (c *RpcImplement) ImportByConfig(_ context.Context, im *rpc.ImportReq) (*rpc.ImportRes, error) {
	res := &rpc.ImportRes{}
	config := NewConfig()
	_, err := config.ConfLoad(im.Filepath)
	if err != nil {
		res.Error = err.Error()
		return res, nil
	}
	value, err := config.GetJsonByte("")
	if err != nil {
		res.Error = err.Error()
		return res, nil
	}
	res.Json = string(value)
	return res, nil
}

//获取预置组件
func (c *RpcImplement) InitImport(_ context.Context, im *rpc.ImportReq) (*rpc.ImportRes, error) {
	res := &rpc.ImportRes{}
	data, err := ioutil.ReadFile("./dots.json")
	if err != nil {
		res.Error = err.Error()
		dot.Logger().Errorln("File reading error", zap.Error(err))
	} else {
		res.Json = string(data)
	}
	return res, nil
}

//导出配置信息
//支持三种格式json toml yaml
//由文件名来区分不同格式
func (c *RpcImplement) ExportConfig(_ context.Context, in *rpc.ExportReq) (*rpc.ExportRes, error) {
	var data = in.Configdata
	var target interface{}

	var fileFormat = make(map[string]string)
	{
		for _, value := range in.Filename {
			sliceF := strings.Split(value, ".")
			fileFormat[sliceF[len(sliceF)-1]] = value //[json]filename.json
		}
	}
	{
		//data->map
		if err := json.Unmarshal([]byte(data), &target); err != nil {
			log.Panic(err)
		}
		//map->file
		for key, value := range fileFormat {
			value = "./run_out/" + value
			file, err := os.OpenFile(value, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
			if err != nil {
				dot.Logger().Errorln("an error occurred with file opening or creation")
			}
			switch key {
			case "json":
				enc := json.NewEncoder(file)
				err := enc.Encode(target)
				if err != nil {
					log.Println("Error in encoding json")
				}
			case "yaml":
				enc := yaml.NewEncoder(file)
				err := enc.Encode(target)
				if err != nil {
					log.Println("Error in encoding yaml")
				}
			case "toml":
				enc := toml.NewEncoder(file)
				err := enc.Encode(target)
				if err != nil {
					log.Println("Error in encoding toml")
				}
			}
			if file != nil {
				file.Close()
			}
		}
	}

	return &rpc.ExportRes{}, nil
}

//导出组件信息
//json
func (c *RpcImplement) ExportDot(_ context.Context, in *rpc.ExportReq) (*rpc.ExportRes, error) {
	var data = in.Dotdata
	name := ""
	if in.Filename == nil || len(in.Filename) < 1 {
		//export init dot
		name = "./dots.json"
	} else {
		//normal export
		name = "./run_out/" + in.Filename[0]
	}
	file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic("An error occurred with file opening or creation\n")
	}
	defer file.Close()
	_, err = file.WriteString(data)
	if err != nil {
		panic("writeString err")
	}
	return &rpc.ExportRes{}, nil
}
