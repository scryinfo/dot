package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/scryinfo/dot/dots/db/pgs"
	"github.com/scryinfo/scryg/sutils/uuid"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const templateStr = `
package {{$.Package}}
import "github.com/scryinfo/dot/dot"
const {{$.Name}}TypeId = "{{$.Id}}"
type config{{$.Name}} struct {
	//todo add
}
type {{$.Name}} struct {
	conf config{{$.Name}}
	//todo add
}
//func (c *{{$.Name}}) Create(l dot.Line) error {
//	//todo add
//}
//func (c *{{$.Name}}) Injected(l dot.Line) error {
//	//todo add
//}
//func (c *{{$.Name}}) AfterAllInject(l dot.Line) {
//	//todo add
//}
//
//func (c *{{$.Name}}) Start(ignore bool) error {
//	//todo add
//}
//
//func (c *{{$.Name}}) Stop(ignore bool) error {
//	//todo add
//}
//
//func (c *{{$.Name}}) Destroy(ignore bool) error {
//	//todo add
//}

//construct dot
func new{{$.Name}}(conf []byte) (dot.Dot, error) {
	dconf := &config{{$.Name}}{}
	//todo
	//err := dot.UnMarshalConfig(conf, dconf)
	//if err != nil {
	//	return nil, err
	//}

	d := &{{$.Name}}{conf: *dconf}

	return d, nil
}

//{{$.Name}}TypeLives
func {{$.Name}}TypeLives() []*dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{TypeId: {{$.Name}}TypeId, NewDoter: func(conf []byte) (dot.Dot, error) {
			return new{{$.Name}}(conf)
		}},
		//Lives: []dot.Live{
		//	{
		//		LiveId:    {{$.Name}}TypeId,
		//		RelyLives: map[string]dot.LiveId{"some field": "some id"},
		//	},
		//},
	}

	lives := []*dot.TypeLives{tl}

	return lives
}

//{{$.Name}}ConfigTypeLive
func {{$.Name}}ConfigTypeLive() *dot.ConfigTypeLives {
	paths := make([]string, 0)
	paths = append(paths, "")
	return &dot.ConfigTypeLives{
		TypeIdConfig: {{$.Name}}TypeId,
		ConfigInfo: &config{{$.Name}}{
			//todo
		},
	}
}
`

type tData struct {
	Name    string
	Id      string
	Config  bool //default true
	Package string

	BackQuote string
}

func parms(data *tData) {
	flag.StringVar(&data.Name, "name", "AnyName", "")
	flag.StringVar(&data.Id, "id", "", "")
	flag.BoolVar(&data.Config, "config", true, "")
	flag.StringVar(&data.Package, "package", "dot", "")

	flag.Parse()
	if len(data.Id) < 1 {
		data.Id = uuid.GetUuid()
	}
}

func main() {
	log.Println("run dotcli")
	data := &tData{BackQuote: "`"}
	parms(data)

	src := gmodel(data)

	outputName := ""
	{
		types := pgs.Underscore(data.Name)
		baseName := fmt.Sprintf("%s.go", types)
		outputName = filepath.Join(".", strings.ToLower(baseName))
	}

	if _, err := os.Stat(outputName); os.IsNotExist(err) {
		err := ioutil.WriteFile(outputName, src, 0644)
		if err != nil {
			log.Fatalf("writing output: %s", err)
		}
	} else {
		log.Println("exist the file: " + outputName)
	}

	log.Println("finished dotcli")
}

func gmodel(data *tData) []byte {

	temp := templateStr
	var src []byte = nil
	{
		t, err := template.New("").Parse(temp)
		if err != nil {
			log.Fatal(err)
		}
		buff := bytes.NewBufferString("")
		err = t.Execute(buff, data)
		if err != nil {
			log.Fatal(err)
		}

		src, err = format.Source(buff.Bytes())
		if err != nil {
			log.Fatal(err)
		}
	}
	return src
}
