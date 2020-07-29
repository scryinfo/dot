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
const {{$.Name}}TypeID = "{{$.ID}}"
type config{{$.Name}} struct {
	//todo add
}
type {{$.Name}} struct {
	conf config{{$.Name}}
	//todo add
}
//func (c *{{$.Name}}) Create(l dot.Line) error {
//	
//}
//func (c *{{$.Name}}) Injected(l dot.Line) error {
//	
//}
//func (c *{{$.Name}}) AfterAllInject(l dot.Line) {
//	
//}
//
//func (c *{{$.Name}}) Start(ignore bool) error {
//	
//}
//
//func (c *{{$.Name}}) Stop(ignore bool) error {
//	
//}
//
//func (c *{{$.Name}}) Destroy(ignore bool) error {
//	
//}

//construct dot
func new{{$.Name}}(conf []byte) (dot.Dot, error) {
	dconf := &config{{$.Name}}{}
	
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
		Meta: dot.Metadata{TypeID: {{$.Name}}TypeID, NewDoter: func(conf []byte) (dot.Dot, error) {
			return new{{$.Name}}(conf)
		}},
		//Lives: []dot.Live{
		//	{
		//		LiveID:    {{$.Name}}TypeID,
		//		RelyLives: map[string]dot.LiveID{"some field": "some id"},
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
		TypeIDConfig: {{$.Name}}TypeID,
		ConfigInfo: &config{{$.Name}}{
			//todo
		},
	}
}
`

type tData struct {
	Name      string
	ID        string
	Config    bool //default true
	Package   string
	BackQuote string
}

var help bool = false

func parms(data *tData) {
	flag.StringVar(&data.Name, "name", "AnyName", "struct name")
	flag.StringVar(&data.ID, "id", "", "dot id, if not set, will make a new")
	//flag.BoolVar(&data.Config, "config", true, "")
	flag.StringVar(&data.Package, "package", "", "package name")

	flag.BoolVar(&help, "h", false, "")

	flag.Parse()
	if len(data.ID) < 1 {
		data.ID = uuid.GetUuid()
	}
	if len(data.Package) < 1 {
		curPath, err := os.Getwd()
		if err == nil {
			data.Package = filepath.Base(curPath)
		}
	}

	if len(data.Package) < 1 {
		data.Package = "dot"
	}
}

func main() {
	log.Println("run dotcli")
	data := &tData{BackQuote: "`"}
	parms(data)

	if help {
		flag.PrintDefaults()
		return
	}

	src := gdao(data)

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

func gdao(data *tData) []byte {

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
