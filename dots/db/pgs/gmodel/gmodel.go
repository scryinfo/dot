package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/scryinfo/dot/dots/db/pgs"
	"github.com/scryinfo/dot/dots/db/tools"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//do not use the map, we need the order
type DbField struct {
	Name   string
	DbName string
}

type tData struct {
	TypeName     string
	TableName    string
	PkgName      string
	MapExcludes  map[string]bool
	Fields       []DbField
	StringFields []DbField
}

var params_ struct {
	typeName    string
	tableName   string
	mapExcludes string
}

func parms(data *tData) {
	flag.StringVar(&params_.typeName, "typeName", "", "")
	flag.StringVar(&params_.tableName, "tableName", "", "")
	flag.StringVar(&params_.mapExcludes, "mapExcludes", "", "split ','")
	flag.Parse()

	if len(params_.mapExcludes) > 0 {
		exes := strings.Split(params_.mapExcludes, ",")
		data.MapExcludes = make(map[string]bool, len(exes))
		for i, _ := range exes {
			it := pgs.CamelCased(exes[i])
			data.MapExcludes[it] = true
		}
	} else {
		data.MapExcludes = make(map[string]bool, 0)
	}

	if len(params_.tableName) < 1 {
		params_.tableName = pgs.Underscore(params_.typeName)
	}

	data.TypeName = params_.typeName
	data.TableName = params_.tableName
}

func main() {
	data := &tData{}
	parms(data)
	if len(params_.typeName) < 1 {
		log.Fatal("type name is null")
	}

	if len(params_.tableName) < 1 {
		log.Fatal("table name is null")
	}

	var src []byte = nil
	{
		makeData(data)
		src = gmodel(data)
		//fmt.Println(string(src))
	}

	outputName := ""
	{
		types := strings.Split(params_.typeName, ",")
		baseName := fmt.Sprintf("%s_model.go", types[0])
		outputName = filepath.Join(".", strings.ToLower(baseName))
	}

	if _, err := os.Stat(outputName); os.IsNotExist(err) {
		err := ioutil.WriteFile(outputName, src, 0644)
		if err != nil {
			log.Fatalf("writing output: %s", err)
		}
	} else {
		fmt.Println("exist the file: " + outputName)
	}

}

func makeData(data *tData) {
	data.PkgName = os.Getenv("GOPACKAGE")
	fields := make([]DbField, 0)
	{
		var pkgInfo *build.Package = nil
		var err error
		{
			pkgInfo, err = build.ImportDir(".", 0)
			if err != nil {
				log.Fatal(err)
			}
		}
		if data.PkgName == "" {
			data.PkgName = pkgInfo.Name
		}

		fset := token.NewFileSet()
		for _, file := range pkgInfo.GoFiles {
			f, err := parser.ParseFile(fset, file, nil, 0)
			if err != nil {
				log.Fatal(err)
			}
			//typ := ""
			ast.Inspect(f, func(n ast.Node) bool {
				if n != nil {
					decl, ok := n.(*ast.GenDecl)
					// just STRUCT
					if !ok || decl.Tok != token.TYPE {
						return true
					}
					for _, spec := range decl.Specs {
						typeS, ok := spec.(*ast.TypeSpec)
						if ok && typeS.Name.Name == data.TypeName {
							structT := typeS.Type.(*ast.StructType)
							for _, f := range structT.Fields.List {
								n := f.Names[0].Name
								fields = append(fields, DbField{Name: n, DbName: tools.Underscore(n)})
								ftype, ok := f.Type.(*ast.Ident)
								if ok && ftype.Name == "string" {
									data.StringFields = append(data.StringFields, DbField{Name: n, DbName: tools.Underscore(n)})
								}
							}
						}
					}
				}

				return true
			})
		}

	}
	data.Fields = fields
}

func gmodel(data *tData) []byte {

	temp := `
package {{.PkgName}}
import (
	"fmt"
	"github.com/scryinfo/cashbox-backend/kits"
)
	const (
		{{$.TypeName}}_Table       = "{{$.TableName}}"
		{{range $.Fields}} {{$.TypeName}}_{{.Name}} = "{{.DbName}}"
		{{end}}
	)

	func (m *{{$.TypeName}}) String() string {
		//todo please change the format string
		//{{range $.Fields}}m.{{.Name}}, {{end}}
		str := fmt.Sprintf("{{$.TypeName}}<{{range $.StringFields}}%s {{end}}>",
		{{range $.StringFields}}m.{{.Name}}, {{end}}
		)
		return str
	}

	func (m *{{$.TypeName}}) ToMap() map[string]string {
		res := pgs.ToMap(m, map[string]bool{ {{range $k,$v := $.MapExcludes}}"{{$k}}":{{$v}} {{end}} })
		return res
	}

	func (m *{{$.TypeName}}) ToUpsertSet() []string {
		res := []string{
		{{range $.Fields}}
			fmt.Sprintf("%s = EXCLUDED.%s", {{$.TypeName}}_{{.Name}}, {{$.TypeName}}_{{.Name}}), {{end}}
		}
		return res
	}
`

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

		//test
		//fmt.Println(string(buff.Bytes()))

		src, err = format.Source(buff.Bytes())
		if err != nil {
			log.Fatal(err)
		}
	}
	return src
}
