package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/scryinfo/dot/dots/db/pgs"
	"go/ast"
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

//env:   GOPACKAGE=model;GOFILE=D:\peace\gopath\src\github.com\scryinfo\cashbox_site\shared\db\model\model_api.go
func main() {
	log.Println("run gmodel")
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
		log.Println("exist the file: " + outputName)
	}
	log.Println("finished gmodel")
}

func makeData(data *tData) {
	data.PkgName = os.Getenv("GOPACKAGE")
	file := os.Getenv("GOFILE")
	fields := make([]DbField, 0)
	{
		f, err := parser.ParseFile(token.NewFileSet(), file, nil, 0)
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
						fields, _ = listFields(data, structT, fields)
					}
				}
			}

			return true
		})

	}
	data.Fields = fields
}
func listFields(data *tData, st *ast.StructType, fields []DbField) ([]DbField, error) {
	for _, f := range st.Fields.List {
		if f.Names == nil { //这个字段直接是类型，或者是组合
			if ident, ok := f.Type.(*ast.Ident); ok {
				if decl, ok := ident.Obj.Decl.(*ast.TypeSpec); ok {
					if subT, ok := decl.Type.(*ast.StructType); ok {
						fields, _ = listFields(data, subT, fields)
					}
				}
			} else { //todo 内部错误

			}
		} else {
			n := f.Names[0].Name
			fields = append(fields, DbField{Name: n, DbName: pgs.Underscore(n)})
			ftype, ok := f.Type.(*ast.Ident)
			if ok && ftype.Name == "string" {
				data.StringFields = append(data.StringFields, DbField{Name: n, DbName: pgs.Underscore(n)})
			}
		}
	}
	return fields, nil
}

func gmodel(data *tData) []byte {

	temp := `
package {{.PkgName}}
import (
	"fmt"
	"github.com/scryinfo/dot/dots/db/pgs"
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

	//todo Please modify with lock
	//fmt.Sprintf("%s = EXCLUDED.%s+1", {{$.TypeName}}_Version, {{$.TypeName}}_Version),
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

		src, err = format.Source(buff.Bytes())
		if err != nil {
			log.Fatal(err)
		}
	}
	return src
}
