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
	DaoName            string
	TypeName           string
	TableName          string
	ModelPkgName       string
	ImportModelPkgName string
	DaoPkgName         string
	BackQuote          string
	Id                 string
	//Fields       []DbField
	//StringFields []DbField
}

var params_ struct {
	typeName   string
	tableName  string
	daoPackage string
	suffix     string
	useLock    bool
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
		{
			dir, err := filepath.Abs(file)
			if err != nil {
				log.Fatal(err)
			}
			dir = filepath.Dir(dir)
			dir = strings.Replace(dir, "\\", "/", -1)
			index := strings.Index(dir, "github.com/scryinfo")
			if index >= 0 {
				data.ImportModelPkgName = dir[index:]
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
	"time"
	
	"github.com/go-pg/pg/v9"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/db/pgs"
	"github.com/scryinfo/scryg/sutils/uuid"
	"{{$.ImportModelPkgName}}"
)

const {{$.DaoName}}TypeId = "{{$.Id}}"

type {{$.DaoName}} struct {
	*pgs.DaoBase {{$.BackQuote}}dot:""{{$.BackQuote}}
}

//{{$.DaoName}}TypeLives
func {{$.DaoName}}TypeLives() []*dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{
			Name: "{{$.DaoName}}",
			TypeId: {{$.DaoName}}TypeId, 
			NewDoter: func(conf []byte) (dot.Dot, error) {
				return &{{$.DaoName}}{}, nil
			},
		},
		Lives: []dot.Live{
			{
				LiveId: {{$.DaoName}}TypeId,
				RelyLives: map[string]dot.LiveId{
					"DaoBase": pgs.DaoBaseTypeId,
				},
			},
		},
	}

	lives := pgs.DaoBaseTypeLives()
	lives = append(lives, tl)
	return lives
}

// if find nothing, return pg.ErrNoRows
func (c *{{$.DaoName}}) GetByIdWithLock(conn *pg.Conn, id string) (m *{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	m = &{{$.ModelPkgName}}.{{$.TypeName}}{Id: id,}
	err = conn.Model(m).WherePK().For("UPDATE").Select()
	if err != nil {
		m = nil
	}
	return
}
func (c *{{$.DaoName}}) GetById(conn *pg.Conn, id string) (m *{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	m = &{{$.ModelPkgName}}.{{$.TypeName}}{Id: id,}
	err = conn.Model(m).WherePK().Select()
	if err != nil {
		m = nil
	}
	return
}


// if find nothing, return pg.ErrNoRows
func (c *{{$.DaoName}}) QueryWithLock(conn *pg.Conn, condition string, params ...interface{}) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	if len(condition) < 1 {
		err = conn.Model(&ms).For("UPDATE").Select()
	} else {
		err = conn.Model(&ms).Where(condition, params...).For("UPDATE").Select()
	}
	if err != nil { //be sure
		ms = nil
	}
	return
}
func (c *{{$.DaoName}}) Query(conn *pg.Conn, condition string, params ...interface{}) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	if len(condition) < 1 {
		err = conn.Model(&ms).Select()
	} else {
		err = conn.Model(&ms).Where(condition, params...).Select()
	}
	if err != nil { //be sure
		ms = nil
	}
	return
}


// if find nothing, return pg.ErrNoRows
func (c *{{$.DaoName}}) ListWithLock(conn *pg.Conn) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	err = conn.Model(&ms).For("UPDATE").Select()
	if err != nil {//be sure
		ms = nil
	}
	return
}
func (c *{{$.DaoName}}) List(conn *pg.Conn) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	err = conn.Model(&ms).Select()
	if err != nil {//be sure
		ms = nil
	}
	return
}

func (c *{{$.DaoName}}) Count(conn *pg.Conn, condition string, params ...interface{}) (count int, err error) {
	if len(condition) < 1 {
		count, err = conn.Model(&{{$.ModelPkgName}}.{{$.TypeName}}{}).Count()
	}else {
		count, err = conn.Model(&{{$.ModelPkgName}}.{{$.TypeName}}{}).Where(condition, params...).Count()
	}
	return
}

// if find nothing, return pg.ErrNoRows
func (c *{{$.DaoName}}) QueryPageWithLock(conn *pg.Conn, pageSize int, page int, condition string, params ...interface{}) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	if len(condition) < 1 {
		err = conn.Model(&ms).Limit(pageSize).Offset((page - 1) * pageSize).For("UPDATE").Select()
	}else {
		err = conn.Model(&ms).Where(condition, params...).Limit(pageSize).Offset((page - 1) * pageSize).For("UPDATE").Select()
	}
	if err != nil { //be sure
		ms = nil
	}
	return
}
func (c *{{$.DaoName}}) QueryPage(conn *pg.Conn, pageSize int, page int, condition string, params ...interface{}) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	if len(condition) < 1 {
		err = conn.Model(&ms).Limit(pageSize).Offset((page - 1) * pageSize).Select()
	}else {
		err = conn.Model(&ms).Where(condition, params...).Limit(pageSize).Offset((page - 1) * pageSize).Select()
	}
	if err != nil { //be sure
		ms = nil
	}
	return
}

// if find nothing, return pg.ErrNoRows
func (c *{{$.DaoName}}) QueryOneWithLock(conn *pg.Conn, condition string, params ...interface{}) (m *{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	m = &{{$.ModelPkgName}}.{{$.TypeName}}{}
	if len(condition) < 1 {
		err = conn.Model(m).For("UPDATE").First()
	} else {
		err = conn.Model(m).Where(condition, params...).For("UPDATE").First()
	}
	if err != nil {//be sure
		m = nil
	}
	return
}
func (c *{{$.DaoName}}) QueryOne(conn *pg.Conn, condition string, params ...interface{}) (m *{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	m = &{{$.ModelPkgName}}.{{$.TypeName}}{}
	if len(condition) < 1 {
		err = conn.Model(m).First()
	} else {
		err = conn.Model(m).Where(condition, params...).First()
	}
	if err != nil {//be sure
		m = nil
	}
	return
}

//if insert nothing, then return pg.ErrNoRows
func (c *{{$.DaoName}}) Insert(conn *pg.Conn, m *{{$.ModelPkgName}}.{{$.TypeName}}) (err error) {
	if len(m.Id) < 1 {
		m.Id = uuid.GetUuid()
	}
	m.CreateTime = time.Now().Unix()
	m.UpdateTime = m.CreateTime
	err = conn.Insert(m)
	return
}
//if insert nothing, then return pg.ErrNoRows
func (c *{{$.DaoName}}) InsertReturn(conn *pg.Conn, m *{{$.ModelPkgName}}.{{$.TypeName}}) ( mnew *{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	if len(m.Id) < 1 {
		m.Id = uuid.GetUuid()
	}
	m.CreateTime = time.Now().Unix()
	m.UpdateTime = m.CreateTime

	mnew = &{{$.ModelPkgName}}.{{$.TypeName}}{}
	_, err = conn.Model(m).Returning("*").Insert(mnew)
	if err != nil{
		mnew = nil
	}
	return
}


//if update nothing, then return pg.ErrNoRows
func (c *{{$.DaoName}}) Upsert(conn *pg.Conn, m *{{$.ModelPkgName}}.{{$.TypeName}}) (err error) {
    m.UpdateTime = time.Now().Unix()
	if len(m.Id) < 1 {
		m.Id = uuid.GetUuid()
		m.CreateTime = m.UpdateTime
	} else if m.CreateTime == 0 {
		m.CreateTime = m.UpdateTime
	}
	
	om := conn.Model(m).OnConflict("(id) DO UPDATE").Where("{{$.TypeName}}."+{{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion+" = ?",m.OptimisticLockVersion)
	for _, it := range m.ToUpsertSet() {
		om.Set(it)
	}
	res, err := om.Insert()
	if res.RowsAffected() == 0 {
		err = pg.ErrNoRows
	}
	return err
}

//if update nothing, then return pg.ErrNoRows
func (c *{{$.DaoName}}) UpsertReturn(conn *pg.Conn, m *{{$.ModelPkgName}}.{{$.TypeName}}) ( mnew *{{$.ModelPkgName}}.{{$.TypeName}},err error) {
	m.UpdateTime = time.Now().Unix()
	if len(m.Id) < 1 {
		m.Id = uuid.GetUuid()
		m.CreateTime = m.UpdateTime
	} else if m.CreateTime == 0 {
		m.CreateTime = m.UpdateTime
	}

	om := conn.Model(m).OnConflict("(id) DO UPDATE").Where("{{$.TypeName}}."+{{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion+" = ?",m.OptimisticLockVersion)
	for _, it := range m.ToUpsertSet() {
		om.Set(it)
	}
	mnew = &{{$.ModelPkgName}}.{{$.TypeName}}{}
	_, err = om.Returning("*").Insert(mnew)
	if err != nil {
		mnew = nil
	}
	return
}
//if update nothing, then return pg.ErrNoRows
func (c *{{$.DaoName}}) Update(conn *pg.Conn, m *{{$.ModelPkgName}}.{{$.TypeName}}) (err error) {
	m.UpdateTime = time.Now().Unix()
	m.OptimisticLockVersion++
	//err = conn.Update(m)
	res, err := conn.Model(m).Where({{$.ModelPkgName}}.{{$.TypeName}}_Id+" = ? and "+{{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion+" = ?", m.Id, m.OptimisticLockVersion-1).Update()
	if res.RowsAffected() == 0 {
		err = pg.ErrNoRows
	}
	return
}

//if update nothing, then return pg.ErrNoRows
func (c *{{$.DaoName}}) UpdateReturn(conn *pg.Conn, m *{{$.ModelPkgName}}.{{$.TypeName}}) (mnew *{{$.ModelPkgName}}.{{$.TypeName}},  err error) {
	m.UpdateTime = time.Now().Unix()
	m.OptimisticLockVersion++
	mnew = &{{$.ModelPkgName}}.{{$.TypeName}}{}
	res, err := conn.Model(m).Where({{$.ModelPkgName}}.{{$.TypeName}}_Id+" = ? and "+{{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion+" = ?", m.Id, m.OptimisticLockVersion-1).Returning("*").Update(mnew)
	if err != nil {
		mnew = nil
	}
	if res.RowsAffected() == 0 {
		err = pg.ErrNoRows
	}
	return
}

//if delete nothing, then return pg.ErrNoRows
func (c *{{$.DaoName}}) Delete(conn *pg.Conn, m *{{$.ModelPkgName}}.{{$.TypeName}}) (err error) {
	err = conn.Delete(m)
	return
}

//if delete nothing, then return pg.ErrNoRows
func (c *{{$.DaoName}}) DeleteById(conn *pg.Conn, id string) (err error) {
	_, err = conn.Model((*{{$.ModelPkgName}}.{{$.TypeName}})(nil)).Where({{$.ModelPkgName}}.{{$.TypeName}}_Id+" = ?", id).Delete()
	return
}

//if delete nothing, then return pg.ErrNoRows
func (c *{{$.DaoName}}) DeleteByIds(conn *pg.Conn, ids []string, oneMax int) (err error) {
	m := (*{{$.ModelPkgName}}.{{$.TypeName}})(nil)
	max := oneMax
	times := len(ids)/max;
	for i := 1; i < times; i++ {
		oneIds := ids[(i-1) * max:i * max -1]
		_, err = conn.Model(m).Where({{$.ModelPkgName}}.{{$.TypeName}}_Id+" in (?)", pg.In(oneIds)).Delete()
		if err != nil {
			return
		}
	}

	if max * times < len(ids) {
		oneIds := ids[max * times:]
		_, err = conn.Model(m).Where({{$.ModelPkgName}}.{{$.TypeName}}_Id+" in (?)", pg.In(oneIds)).Delete()
	}
	return 
}

//if delete nothing, then return pg.ErrNoRows
func (c *{{$.DaoName}}) DeleteReturn(conn *pg.Conn, m *{{$.ModelPkgName}}.{{$.TypeName}}) (mnew *{{$.ModelPkgName}}.{{$.TypeName}},err error) {
	mnew = &{{$.ModelPkgName}}.{{$.TypeName}}{}
	_, err = conn.Model(m).WherePK().Returning("*").Delete(mnew)
	if err != nil {
		mnew = nil
	}
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
