package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/scryinfo/dot/dots/db/pgs"
	"github.com/scryinfo/scryg/sutils/uuid"
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
	DaoName      string
	TypeName     string
	TableName    string
	ModelPkgName string
	DaoPkgName   string
	BackQuote    string
	Id           string
	//Fields       []DbField
	//StringFields []DbField
}

var params_ struct {
	typeName   string
	tableName  string
	daoPackage string
	suffix     string
}

func parms(data *tData) {
	flag.StringVar(&params_.typeName, "typeName", "", "")
	flag.StringVar(&params_.tableName, "tableName", "", "")
	flag.StringVar(&params_.daoPackage, "daoPackage", "", "")
	flag.StringVar(&params_.suffix, "suffix", "Dao", "")
	flag.Parse()

	if len(params_.tableName) < 1 {
		params_.tableName = pgs.Underscore(params_.typeName)
	}

	data.TypeName = params_.typeName
	data.TableName = params_.tableName
	data.DaoPkgName = params_.daoPackage
	if len(data.DaoPkgName) < 1 {
		data.DaoPkgName = "dao"
	}
	if len(params_.suffix) > 0 {
		data.DaoName = data.TypeName + params_.suffix
	} else {
		data.DaoName = data.TypeName
	}
	data.Id = uuid.GetUuid()
}

func main() {
	log.Println("run gdao")
	data := &tData{}
	data.BackQuote = "`"
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
		types := pgs.Underscore(data.DaoName)
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

	log.Println("finished gdao")
}

func makeData(data *tData) {
	data.ModelPkgName = os.Getenv("GOPACKAGE")
	file := os.Getenv("GOFILE")
	{
		f, err := parser.ParseFile(token.NewFileSet(), file, nil, 0)
		if err != nil {
			log.Fatal(err)
		}
		{ //todo
			dir, err := filepath.Abs(file)
			if err != nil {
				log.Fatal(err)
			}
			dir = filepath.Dir(dir)
			dir = strings.Replace(dir, "\\", "/", -1)
			index := strings.Index(dir, "github.com/scryinfo")
			if index >= 0 {
				data.ModelPkgName = dir[index:]
			} else {
				log.Println("not find the model")
			}
		}
		find := false
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
						find = true
						return false
					}
				}
			}
			return true
		})

		if !find {
			log.Fatal("not find: " + data.TypeName)
		}
	}
	//data.Fields = fields
}

func gmodel(data *tData) []byte {

	temp := `
package {{$.DaoPkgName}}
import (
	"github.com/go-pg/pg/v9"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/db/pgs"
	"{{$.ModelPkgName}}"
)

const {{$.DaoName}}TypeId = "{{$.Id}}"

type {{$.DaoName}} struct {
	*pgs.DaoBase {{$.BackQuote}}dot:""{{$.BackQuote}}
}

//{{$.DaoName}}TypeLives
func {{$.DaoName}}TypeLives() []*dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{TypeId: {{$.DaoName}}TypeId, NewDoter: func(conf []byte) (dot.Dot, error) {
			return &{{$.DaoName}}{}, nil
		}},
	}

	lives := pgs.DaoBaseTypeLives()
	lives = append(lives, tl)
	return lives
}
func (c *{{$.DaoName}}) Query(conn *pg.Conn, condition string, params ...interface{}) (ms []model.{{$.TypeName}}, err error) {
	err = conn.Model(&ms).Where(condition, params...).Select()
	return
}

func (c *{{$.DaoName}}) QueryOne(conn *pg.Conn, condition string, params ...interface{}) (m *model.{{$.TypeName}}, err error) {
	m = &model.{{$.TypeName}}{}
	err = conn.Model(m).Where(condition, params...).First()
	return
}

func (c *{{$.DaoName}}) Insert(conn *pg.Conn, m *model.{{$.TypeName}}) (err error) {
	err = conn.Insert(m)
	return
}

func (c *{{$.DaoName}}) Upsert(conn *pg.Conn, m *model.{{$.TypeName}}) (err error) {
	om := conn.Model(m).OnConflict("(id) DO UPDATE")
	for _, it := range m.ToUpsertSet() {
		om.Set(it)
	}
	_, err = om.Insert()
	return err
}

func (c *{{$.DaoName}}) Delete(conn *pg.Conn, m *model.{{$.TypeName}}) (err error) {
	err = conn.Delete(m)
	return
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
