// Scry Info.  All rights reserved.
// license that can be found in the license file.

package nobl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/grpc/gserver"
	"github.com/scryinfo/dot/tools/config/data/go_out"
	"github.com/scryinfo/dot/tools/config/data/nobl/tool"
	"gopkg.in/yaml.v2"
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

func (serv *HiServer) Hi(ctx context.Context, req *go_out.ReqData) (*go_out.ResData, error) {
	log.Println("hi:", "name:", req.Name)
	res := &go_out.ResData{Test: "hi, i am serve"}
	return res, nil
}

func (serv *HiServer) FindDot(ctx context.Context, in *go_out.ReqDirs) (*go_out.ResDots, error) {
	dirs := in.Dirs
	bytes, strings2, e := tool.FindDots(dirs)
	fmt.Println(string(bytes), strings2, e)
	//删除中间文件
	/*	del := os.Remove("./callMethod.go")
		del = os.Remove("./result.json")
		if del != nil {
			fmt.Println(del)
		}*/
	resDots := go_out.ResDots{
		DotsInfo:    string(bytes),
		NoExistDirs: strings2,
	}
	if e != nil {
		resDots.Error = e.Error()
	}
	return &resDots, nil
}

func (serv *HiServer) LoadByConfig(context.Context, *go_out.ReqLoad) (*go_out.ResConfig, error) {

	panic("implement me")
}

//tim
//根据配置文件导入信息
//支持三种格式json toml yaml
func (serv *HiServer) ImportByConfig(context.Context, *go_out.ReqImport) (*go_out.ResImport, error) {
	panic("implement me")
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
		for key, value := range fileFormat {

			if key == "json" {
				if err := json.Unmarshal([]byte(data), &target); err != nil {
					log.Panic(err)
				}
				file, _ := os.OpenFile(value, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0)
				defer file.Close()
				enc := json.NewEncoder(file)
				err := enc.Encode(target)
				if err != nil {
					log.Println("Error in encoding json")
				}
			}
			if key == "yaml" {
				if err := yaml.Unmarshal([]byte(data), &target); err != nil {
					log.Panic(err)
				}
				file, _ := os.OpenFile(value, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0)
				defer file.Close()
				enc := yaml.NewEncoder(file)
				err := enc.Encode(target)
				if err != nil {
					log.Println("Error in encoding yaml")
				}
			}
			if key == "toml" {

				/*if _, err := toml.DecodeFile(value, &target); err != nil {
					log.Fatal(err)
				}*/
				/*if err := toml.Unmarshal([]byte(data), &target); err != nil {
					log.Panic(err)
				}*/

				/*file, _ := os.OpenFile(value, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0)
				defer file.Close()
				enc:=toml.NewEncoder(file)
				err := enc.Encode(target)
				if err != nil {
					log.Println("Error in encoding toml")
				}*/
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
	file, err := os.OpenFile(in.Filename[0], os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		fmt.Printf("An error occurred with file opening or creation\n")
		panic("An error occurred with file opening or creation\n")
	}
	defer file.Close()
	_, err = file.WriteString(data)
	if err != nil {
		panic("writeString err")
	}
	return &go_out.ResExport{}, nil
}
