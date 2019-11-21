// Scry Info.  All rights reserved.
// license that can be found in the license file.

package nobl

import (
	"context"
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/grpc/gserver"
	"github.com/scryinfo/dot/tools/config/data/go_out"
	"github.com/scryinfo/dot/tools/config/data/nobl/tool/findDot"
	"github.com/scryinfo/dot/tools/config/data/nobl/tool/importConfig"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	ServerTypeId = "rpcImplement"
)

type config struct {
	Name string `json:"name"`
}

//todo 命名
type RpcImplement struct {
	ServerNobl gserver.ServerNobl `dot:""`
	conf       config
}

func newRpcImplement(conf []byte) (dot.Dot, error) {
	dconf := &config{}
	err := dot.UnMarshalConfig(conf, dconf)
	if err != nil {
		return nil, err
	}

	d := &RpcImplement{
		conf: *dconf,
	}

	return d, err
}

func (serv *RpcImplement) Start(ignore bool) error {
	go_out.RegisterDotConfigServer(serv.ServerNobl.Server(), serv)
	return nil
}

//RpcImplementTypeLives make all type lives
func RpcImplementTypeLives() []*dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{TypeId: ServerTypeId, NewDoter: func(conf []byte) (dot.Dot, error) {
			return newRpcImplement(conf)
		}},
		Lives: []dot.Live{
			dot.Live{
				LiveId:    ServerTypeId,
				RelyLives: map[string]dot.LiveId{"ServerNobl": gserver.ServerNoblTypeId},
			},
		},
	}

	lives := []*dot.TypeLives{gserver.ServerNoblTypeLive(), tl}

	return lives
}

//rpc implement

func (serv *RpcImplement) FindDot(ctx context.Context, in *go_out.ReqDirs) (*go_out.ResDots, error) {
	dirs := in.Dirs
	bytes, invalidDirectory, e := findDot.FindDots(dirs)
	//删除运行时产生的中间文件
	{
		del := os.Remove("./run_out/callMethod.go")
		del = os.Remove("./run_out/result.json")
		if del != nil {
			log.Println(del)
		}
	}
	resDots := go_out.ResDots{
		DotsInfo:    string(bytes),
		NoExistDirs: invalidDirectory,
	}
	if e != nil {
		resDots.Error = e.Error()
	}
	return &resDots, nil
}

func (serv *RpcImplement) ImportByDot(ctx context.Context, in *go_out.ReqImport) (*go_out.ResImport, error) {
	var errStr string
	data, err := ioutil.ReadFile(in.Filepath)
	if err != nil {
		errStr = err.Error()
		log.Println("File reading error", err)
	}
	res := go_out.ResImport{
		Json:  string(data),
		Error: errStr,
	}
	return &res, nil

}

//支持三种格式json toml yaml
func (serv *RpcImplement) ImportByConfig(con context.Context, im *go_out.ReqImport) (*go_out.ResImport, error) {
	config := importConfig.New()
	_, err := config.ConfLoad(im.Filepath)
	if err != nil {
		resConfig := go_out.ResImport{
			Error: err.Error(),
		}
		return &resConfig, nil
	}
	value, err := config.GetJsonByte("")
	if err != nil {
		resConfig := go_out.ResImport{
			Error: err.Error(),
		}
		return &resConfig, nil
	}
	resConfig := go_out.ResImport{
		Json: string(value),
	}
	return &resConfig, nil
}

//导出配置信息
//支持三种格式json toml yaml
//由文件名来区分不同格式
func (serv *RpcImplement) ExportConfig(ctx context.Context, in *go_out.ReqExport) (*go_out.ResExport, error) {
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
				log.Println("An error occurred with file opening or creation\n")
			}
			defer file.Close()
			if key == "json" {
				enc := json.NewEncoder(file)
				err := enc.Encode(target)
				if err != nil {
					log.Println("Error in encoding json")
				}
			}
			if key == "yaml" {
				enc := yaml.NewEncoder(file)
				err := enc.Encode(target)
				if err != nil {
					log.Println("Error in encoding yaml")
				}
			}
			if key == "toml" {
				enc := toml.NewEncoder(file)
				err := enc.Encode(target)
				if err != nil {
					log.Println("Error in encoding toml")
				}
			}
		}
	}

	return &go_out.ResExport{}, nil
}

//导出组件信息
//json
func (serv *RpcImplement) ExportDot(ctx context.Context, in *go_out.ReqExport) (*go_out.ResExport, error) {
	var data = in.Dotdata
	if in.Filename == nil {
		in.Filename[0] = "./run_out/dots.json"
	}
	name := "./run_out/" + in.Filename[0]
	file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic("An error occurred with file opening or creation\n")
	}
	defer file.Close()
	_, err = file.WriteString(data)
	if err != nil {
		panic("writeString err")
	}
	return &go_out.ResExport{}, nil
}
