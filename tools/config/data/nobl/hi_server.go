// Scry Info.  All rights reserved.
// license that can be found in the license file.

package nobl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/grpc/gserver"
	"github.com/scryinfo/dot/tools/config/data/go_out"
	"github.com/scryinfo/dot/tools/config/data/nobl/tool"
	"github.com/scryinfo/dot/tools/config/data/nobl/tool/scryconfig"
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
//todo 命名
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
func (serv *HiServer) Hi(ctx context.Context, req *go_out.ReqData) (*go_out.ResData, error) {
	log.Println("hi:", "name:", req.Name)
	res := &go_out.ResData{Test: "hi, i am serve"}
	return res, nil
}

func (serv *HiServer) FindDot(ctx context.Context, in *go_out.ReqDirs) (*go_out.ResDots, error) {
	dirs := in.Dirs
	bytes, invalidDirectory, e := tool.FindDots(dirs)
	//删除中间文件
	{
		del := os.Remove("./callMethod.go")
		del = os.Remove("./result.json")
		if del != nil {
			fmt.Println(del)
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

//从配置文件中加载某一个typeId对应的实例
//暂时只支持json格式
//copyPaste或者文件路径都可
type allinfo struct {
	Dots []dotinfo
}
type dotinfo struct {
	MetaData meta          `json:"metaData"`
	Lives    []interface{} `json:"lives"`
}
type meta struct {
	TypeId string `json:"typeId"`
}

func (serv *HiServer) LoadByConfig(ctx context.Context, in *go_out.ReqLoad) (*go_out.ResConfig, error) {

	var result string
	var err string
	var configinfo allinfo
	//读文件
	{
		data, err := ioutil.ReadFile(in.DataFilepath)
		if err != nil {
			log.Panic("File reading error", err)
		}
		if err := json.Unmarshal(data, &configinfo); err != nil {
			log.Panic("Unmarshal file,", err)
		}
	}
	//筛选
	for _, value := range configinfo.Dots {
		if value.MetaData.TypeId == in.TypeId {
			data, err := json.Marshal(value)
			if err != nil {
				log.Panic("Json marshaling failed：%s", err)
			}
			result = string(data)
		}
	}
	res := go_out.ResConfig{
		ConfigJson: result,
		ErrInfo:    err,
	}
	return &res, nil
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
	scry := scryconfig.New()
	scry.BindFlag()
	_, err := scry.ConfLoad(scry.ConFlag, im.Filepath)
	if err != nil {
		resConfig := go_out.ResImport{
			Error: err.Error(),
		}
		return &resConfig, nil
	}
	value, err := scry.GetJsonByte("")
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
