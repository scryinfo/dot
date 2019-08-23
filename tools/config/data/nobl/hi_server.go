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
	HiServerTypeId = "hiserver"
)

type config struct {
	Name string `json:"name"`
}

type HiServer struct {
	ServerNobl gserver.ServerNobl `dot:""`
	conf       config
}

func newHiServer(conf interface{}) (dot.Dot, error) {
	var err error = nil
	var bs []byte = nil
	if bt, ok := conf.([]byte); ok {
		bs = bt
	} else {
		return nil, dot.SError.Parameter
	}
	dconf := &config{}
	err = dot.UnMarshalConfig(bs, dconf)
	if err != nil {
		return nil, err
	}

	d := &HiServer{
		conf: *dconf,
	}

	return d, err
}

func (serv *HiServer) Start(ignore bool) error {
	go_out.RegisterHiDotServer(serv.ServerNobl.Server(), serv)
	return nil
}

//HiServerTypeLives make all type lives
func HiServerTypeLives() []*dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{TypeId: HiServerTypeId, NewDoter: func(conf interface{}) (dot.Dot, error) {
			return newHiServer(conf)
		}},
		Lives: []dot.Live{
			dot.Live{
				LiveId:    HiServerTypeId,
				RelyLives: map[string]dot.LiveId{"ServerNobl": gserver.ServerNoblTypeId},
			},
		},
	}

	lives := []*dot.TypeLives{gserver.ServerNoblTypeLive(), tl}

	return lives
}

//rpc implement

func (serv *HiServer) FindDot(ctx context.Context, in *go_out.ReqDirs) (*go_out.ResDots, error) {
	dirs := in.Dirs
	bytes, invalidDirectory, e := findDot.FindDots(dirs)
	//删除中间文件
	/*{
		del := os.Remove("./callMethod.go")
		del = os.Remove("./result.json")
		if del != nil {
			fmt.Println(del)
		}
	}*/
	resDots := go_out.ResDots{
		DotsInfo:    string(bytes),
		NoExistDirs: invalidDirectory,
	}
	if e != nil {
		resDots.Error = e.Error()
	}
	return &resDots, nil
}

func (serv *HiServer) ImportByDot(ctx context.Context, in *go_out.ReqImport) (*go_out.ResImport, error) {
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
func (serv *HiServer) ImportByConfig(con context.Context, im *go_out.ReqImport) (*go_out.ResImport, error) {
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
func (serv *HiServer) ExportConfig(ctx context.Context, in *go_out.ReqExport) (*go_out.ResExport, error) {
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
func (serv *HiServer) ExportDot(ctx context.Context, in *go_out.ReqExport) (*go_out.ResExport, error) {
	var data = in.Dotdata
	if in.Filename == nil {
		in.Filename[0] = "dots.json"
	}
	file, err := os.OpenFile(in.Filename[0], os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
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
