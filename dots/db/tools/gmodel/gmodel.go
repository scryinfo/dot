package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/scryinfo/dot/dots/db/bun"
	"go/ast"
	"go/format"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/jinzhu/inflection"
	"golang.org/x/tools/go/packages"
)

//todo bug说明 类似如下情况
/*type Token struct {
	ID string `pg:",pk,type:varchar(36)"`
	DefaultToken
}

type DefaultToken struct {
	ID 	 string
	Name string
}
gmodel result
const (
	Lock_Table      = "locks"
	Lock_ID         = "id"
	Lock_ID         = "id"
	Lock_Name       = "name"
pgs.CreateSchema result
	DefaultToken.ID	字段缺失
)*/

// DbField do not use the map, we need the order
type DbField struct {
	Name        string
	DbName      string
	HasRelation bool //has one, has many, belong to
}

type tData struct {
	TypeName     string
	TableName    string
	PkgName      string
	DbObjectName string
	MapExcludes  map[string]bool
	Fields       []DbField
	StringFields []DbField
	ModelFile    string
}

var params struct {
	typeName    string
	mapExcludes string
	model       string
}

func parms(data *tData) {
	flag.StringVar(&params.typeName, "typeName", "", "")
	flag.StringVar(&params.mapExcludes, "mapExcludes", "", "split ','")
	flag.StringVar(&params.model, "model", "models.go", "")
	flag.Parse()

	if len(params.mapExcludes) > 0 {
		exes := strings.Split(params.mapExcludes, ",")
		data.MapExcludes = make(map[string]bool, len(exes))
		for i := range exes {
			it := bun.CamelCased(exes[i])
			data.MapExcludes[it] = true
		}
	} else {
		data.MapExcludes = make(map[string]bool, 0)
	}

	data.TypeName = params.typeName
	data.DbObjectName = bun.Underscore(params.typeName)
	data.TableName = tableNameInflector(data.DbObjectName)
	data.ModelFile = params.model
}

var tableNameInflector = inflection.Plural

// env:   GOPACKAGE=model;GOFILE=D:\gopath\src\github.com\scryinfo\dot\sample\db\pgs\model\models.go
func main() {

	log.Println("run gmodel")
	data := &tData{}
	parms(data)
	if len(params.typeName) < 1 {
		log.Fatal("type name is null")
	}
	_ = os.Setenv("GOPACKAGE", "model")
	_ = os.Setenv("GOFILE", data.ModelFile)

	var src []byte = nil
	{
		makeData(data)
		src = gmodel(data)
	}

	outputName := ""
	{
		types := bun.Underscore(data.TypeName)
		baseName := fmt.Sprintf("%s_model.go", types)
		outputName = filepath.Join(".", strings.ToLower(baseName))
	}

	if _, err := os.Stat(outputName); os.IsNotExist(err) {
		err := os.WriteFile(outputName, src, 0644)
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
	var pkg *packages.Package
	{
		dir := filepath.Dir(file)
		cfg := &packages.Config{
			Mode: packages.NeedName |
				//packages.NeedFiles |
				//packages.NeedCompiledGoFiles |
				packages.NeedImports |
				packages.NeedDeps |
				//packages.NeedExportsFile |
				//packages.NeedTypes |
				packages.NeedSyntax,
			//packages.NeedTypesInfo |
			//packages.NeedTypesSizes ,
			Dir:   dir,
			Tests: true,
			Env:   append(os.Environ(), "GO111MODULE=off", "GOPROXY=off"), //"GOPATH="+dir,
		}
		file, err := filepath.Abs(file)
		if err != nil {
			log.Fatal(err)
		}
		pkgs, err := packages.Load(cfg, file)
		if err != nil {
			log.Fatal(err)
		}
		pkg = pkgs[0]
		if len(pkg.Syntax) != 1 {
			log.Fatal("parse file not " + file)
		}
	}
	{
		f := pkg.Syntax[0]
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
						fields, _ = listFields(data, structT, fields, pkg)
					}
				}
			}

			return true
		})

	}
	data.Fields = fields
}
func listFields(data *tData, st *ast.StructType, fields []DbField, pkg *packages.Package) ([]DbField, error) {
	for _, field := range st.Fields.List {
		var subT *ast.StructType
		//|| (field.Tag != nil && strings.Contains(field.Tag.Value,"composite"))
		if field.Names == nil { //这个字段直接是类型，或者是组合
			if ident, ok := field.Type.(*ast.Ident); ok {
				if decl, ok := ident.Obj.Decl.(*ast.TypeSpec); ok {
					if subT, ok = decl.Type.(*ast.StructType); ok {
						//fields, _ = listFields(data, subT, fields, pkg)
					}
				}
			} else if expr, ok := field.Type.(*ast.SelectorExpr); ok {
				//这个字段的类型是从import中导入的，
				obj := expr.Sel.Obj
				if obj == nil {
					packageName := expr.X.(*ast.Ident).Name
					objName := expr.Sel.Name
				objLoop:
					for _, importPackage := range pkg.Imports {
						if importPackage.Name == packageName {
							for _, importFile := range importPackage.Syntax {
								if importFile.Scope != nil {
									obj = importFile.Scope.Lookup(objName)
									if obj != nil {
										break objLoop
									}
								}
							}
						}
					}
				}
				if obj == nil {
					//todo 内部错误
				} else {
					if decl, ok := obj.Decl.(*ast.TypeSpec); ok {
						if subT, ok = decl.Type.(*ast.StructType); ok {
							//fields, _ = listFields(data, subT, fields, pkg)
						}
					}
				}

			} else { //todo 内部错误

			}
		} else {
			name := field.Names[0].Name
			dbField := DbField{Name: name, DbName: bun.Underscore(name)}
			tag := ""
			if field.Tag != nil {
				tag = field.Tag.Value
			}
			if strings.Contains(tag, `pg:"-"`) || strings.Contains(tag, "rel:has-one") || strings.Contains(tag, "rel:has-many") || strings.Contains(tag, "rel:belongs-to") {
				dbField.HasRelation = true
			}
			fields = append(fields, dbField)
			ftype, ok := field.Type.(*ast.Ident)
			if ok && ftype.Name == "string" {
				data.StringFields = append(data.StringFields, dbField)
			}
		}
		if subT != nil {
			fields, _ = listFields(data, subT, fields, pkg)
		}
	}
	return fields, nil
}

func gmodel(data *tData) []byte {

	temp := `
package {{.PkgName}}
import (
	"fmt"
	"github.com/scryinfo/dot/dots/db/bun/pgd"
)
	const (
		{{$.TypeName}}_Table       = "{{$.TableName}}"
		{{$.TypeName}}_Struct      = "{{$.DbObjectName}}"
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
			{{- range $.Fields}}
			{{if .HasRelation}}{{else}}fmt.Sprintf("%s = EXCLUDED.%s", {{$.TypeName}}_{{.Name}}, {{$.TypeName}}_{{.Name}}),
			{{- end}}
			{{- end}}
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
