package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/scryinfo/dot/dots/db/buns"
	"github.com/scryinfo/dot/dots/db/tools"
	"github.com/scryinfo/scryg/sutils/sfile"
	"github.com/scryinfo/scryg/sutils/uuid"
)

// DbField do not use the map, we need the order
//type DbField struct {
//	Name   string
//	DbName string
//}

type tData struct {
	DaoName            string
	TypeName           string
	TableName          string
	ModelFile          string
	ModelPkgName       string
	ImportModelPkgName string
	DaoPkgName         string
	BackQuote          string
	ID                 string
	//Fields       []DbField
	//StringFields []DbField
}

var params struct {
	typeName   string
	tableName  string
	model      string
	daoPackage string
	suffix     string
	useLock    bool
}

func parms(data *tData) {
	flag.StringVar(&params.typeName, "typeName", "", "")
	flag.StringVar(&params.tableName, "tableName", "", "")
	flag.StringVar(&params.daoPackage, "daoPackage", "", "")
	flag.StringVar(&params.suffix, "suffix", "Dao", "")
	flag.StringVar(&params.model, "model", "models.go", "")
	flag.Parse()

	if len(params.tableName) < 1 {
		params.tableName = buns.Underscore(params.typeName)
	}

	data.TypeName = params.typeName
	data.TableName = params.tableName
	data.DaoPkgName = params.daoPackage
	data.ModelFile = params.model
	if len(data.DaoPkgName) < 1 {
		data.DaoPkgName = "dao"
	}
	if len(params.suffix) > 0 {
		data.DaoName = data.TypeName + params.suffix
	} else {
		data.DaoName = data.TypeName
	}
	data.ID = uuid.GetUuid()
}

// go run dots/db/tools/gdao/gdao.go -typeName Notice -model models.go -daoPackage pgs
func main() {
	log.Println("run gdao")
	data := &tData{}
	data.BackQuote = "`"
	parms(data)
	if len(params.typeName) < 1 {
		log.Fatal("type name is null")
	}
	log.Println(params.typeName)

	if len(params.tableName) < 1 {
		log.Fatal("table name is null")
	}

	var src []byte = nil
	{
		makeData(data)
		src = gmodel(data)
	}

	outputName := ""
	{
		//check the package dir, if not exist, then create it
		daoDir := "../" + data.DaoPkgName
		{
			// dapDir = filepath.Re
			if !sfile.ExistFile(daoDir) {
				err := os.Mkdir(daoDir, os.ModePerm)
				if err != nil {
					log.Fatal(err)
					return
				}
			}
		}
		types := tools.Underscore(data.DaoName)
		baseName := fmt.Sprintf("%s.go", types)
		outputName = filepath.Join(daoDir, strings.ToLower(baseName))
	}

	if _, err := os.Stat(outputName); os.IsNotExist(err) {
		err := os.WriteFile(outputName, src, 0644)
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
		{
			dir, err := filepath.Abs(file)
			if err != nil {
				log.Fatal(err)
			}
			dir = filepath.Dir(dir)
			dir = strings.Replace(dir, "\\", "/", -1)
			data.ImportModelPkgName = strings.TrimLeft(dir, os.Getenv("GOPATH"))
			data.ImportModelPkgName = strings.Trim(data.ImportModelPkgName, "src/")
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
			log.Fatal("type: <" + data.TypeName + "> not found in: " + file)
		}
	}
	//data.Fields = fields
}

func gmodel(data *tData) []byte {
	temp := buns.GetDaoData()
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
