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

//DbField do not use the map, we need the order
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
	ID                 string
	//Fields       []DbField
	//StringFields []DbField
}

var params struct {
	typeName   string
	tableName  string
	daoPackage string
	suffix     string
	useLock    bool
}

func parms(data *tData) {
	flag.StringVar(&params.typeName, "typeName", "", "")
	flag.StringVar(&params.tableName, "tableName", "", "")
	flag.StringVar(&params.daoPackage, "daoPackage", "", "")
	flag.StringVar(&params.suffix, "suffix", "Dao", "")
	flag.Parse()

	if len(params.tableName) < 1 {
		params.tableName = pgs.Underscore(params.typeName)
	}

	data.TypeName = params.typeName
	data.TableName = params.tableName
	data.DaoPkgName = params.daoPackage
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

func main() {
	log.Println("run gdao")
	data := &tData{}
	data.BackQuote = "`"
	parms(data)
	if len(params.typeName) < 1 {
		log.Fatal("type name is null")
	}

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

const {{$.DaoName}}TypeID = "{{$.ID}}"

type {{$.DaoName}} struct {
	*pgs.DaoBase {{$.BackQuote}}dot:""{{$.BackQuote}}
}

//{{$.DaoName}}TypeLives
func {{$.DaoName}}TypeLives() []*dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{
			Name: "{{$.DaoName}}",
			TypeID: {{$.DaoName}}TypeID, 
			NewDoter: func(conf []byte) (dot.Dot, error) {
				return &{{$.DaoName}}{}, nil
			},
		},
		Lives: []dot.Live{
			{
				LiveID: {{$.DaoName}}TypeID,
				RelyLives: map[string]dot.LiveID{
					"DaoBase": pgs.DaoBaseTypeID,
				},
			},
		},
	}

	lives := pgs.DaoBaseTypeLives()
	lives = append(lives, tl)
	return lives
}

// if find nothing, return pg.ErrNoRows
func (c *{{$.DaoName}}) GetByIDWithLock(conn *pg.Conn, id string) (m *{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	m = &{{$.ModelPkgName}}.{{$.TypeName}}{ID: id,}
	err = conn.Model(m).WherePK().For("UPDATE").Select()
	if err != nil {
		m = nil
	}
	return
}
func (c *{{$.DaoName}}) GetByID(conn *pg.Conn, id string) (m *{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	m = &{{$.ModelPkgName}}.{{$.TypeName}}{ID: id,}
	err = conn.Model(m).WherePK().Select()
	if err != nil {
		m = nil
	}
	return
}

//update before
//you must get OptimisticLockVersion value
func (c *{{$.DaoName}}) GetLockByID(conn *pg.Conn, ids ...string) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	for i,_ := range ids {
		m := &{{$.ModelPkgName}}.{{$.TypeName}}{ID: ids[i],}
		ms=append(ms,m)
	}
	err=conn.Model(&ms).WherePK().Column({{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion,{{$.ModelPkgName}}.{{$.TypeName}}_ID).For("UPDATE").Select()
	if err != nil {
		ms=nil
	}
	return
}
func (c *{{$.DaoName}}) GetLockByModelID(conn *pg.Conn, ms ...*{{$.ModelPkgName}}.{{$.TypeName}}) error {
	return conn.Model(&ms).WherePK().Column({{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion,{{$.ModelPkgName}}.{{$.TypeName}}_ID).For("UPDATE").Select()
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
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuid()
	}
	m.CreateTime = time.Now().Unix()
	m.UpdateTime = m.CreateTime
	err = conn.Insert(m)
	return
}
//if insert nothing, then return pg.ErrNoRows
func (c *{{$.DaoName}}) InsertReturn(conn *pg.Conn, m *{{$.ModelPkgName}}.{{$.TypeName}}) ( mnew *{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuid()
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
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuid()
		m.CreateTime = m.UpdateTime
	} else if m.CreateTime == 0 {
		m.CreateTime = m.UpdateTime
	}
	m.OptimisticLockVersion++
	om := conn.Model(m).OnConflict("(id) DO UPDATE").Where({{$.ModelPkgName}}.{{$.TypeName}}_Struct+"."+{{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion+" = ?",m.OptimisticLockVersion-1)
	for _, it := range m.ToUpsertSet() {
		om.Set(it)
	}
	res, err := om.Insert()
	if res.RowsAffected() == 0 {
		//err = pg.ErrNoRows
		newm ,err:=c.GetLockByID(conn,m.ID)
		if err != nil {
			return err
		}
		m.OptimisticLockVersion =newm[0].OptimisticLockVersion
		err =c.Update(conn,m)
		return err
	}
	return err
}

//if update nothing, then return pg.ErrNoRows
func (c *{{$.DaoName}}) UpsertReturn(conn *pg.Conn, m *{{$.ModelPkgName}}.{{$.TypeName}}) ( mnew *{{$.ModelPkgName}}.{{$.TypeName}},err error) {
	m.UpdateTime = time.Now().Unix()
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuid()
		m.CreateTime = m.UpdateTime
	} else if m.CreateTime == 0 {
		m.CreateTime = m.UpdateTime
	}
	m.OptimisticLockVersion++
	om := conn.Model(m).OnConflict("(id) DO UPDATE").Where({{$.ModelPkgName}}.{{$.TypeName}}_Struct+"."+{{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion+" = ?",m.OptimisticLockVersion-1)
	for _, it := range m.ToUpsertSet() {
		om.Set(it)
	}
	mnew = &{{$.ModelPkgName}}.{{$.TypeName}}{}
	_, err = om.Returning("*").Insert(mnew)
	if err ==pg.ErrNoRows {
		new_m ,err:=c.GetLockByID(conn,m.ID)
		if err != nil {
			return nil,err
		}
		m.OptimisticLockVersion =new_m[0].OptimisticLockVersion
		mnew,err =c.UpdateReturn(conn,m)
		return mnew,err
	}
	return
}
//if update nothing, then return pg.ErrNoRows
func (c *{{$.DaoName}}) Update(conn *pg.Conn, m *{{$.ModelPkgName}}.{{$.TypeName}}) (err error) {
	m.UpdateTime = time.Now().Unix()
	m.OptimisticLockVersion++
	//err = conn.Update(m)
	res, err := conn.Model(m).Where({{$.ModelPkgName}}.{{$.TypeName}}_ID+" = ? and "+{{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion+" = ?", m.ID, m.OptimisticLockVersion-1).Update()
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
	res, err := conn.Model(m).Where({{$.ModelPkgName}}.{{$.TypeName}}_ID+" = ? and "+{{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion+" = ?", m.ID, m.OptimisticLockVersion-1).Returning("*").Update(mnew)
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
func (c *{{$.DaoName}}) DeleteByID(conn *pg.Conn, id string) (err error) {
	_, err = conn.Model((*{{$.ModelPkgName}}.{{$.TypeName}})(nil)).Where({{$.ModelPkgName}}.{{$.TypeName}}_ID+" = ?", id).Delete()
	return
}

//if delete nothing, then return pg.ErrNoRows
func (c *{{$.DaoName}}) DeleteByIDs(conn *pg.Conn, ids []string, oneMax int) (err error) {
	m := (*{{$.ModelPkgName}}.{{$.TypeName}})(nil)
	max := oneMax
	times := len(ids)/max;
	for i := 1; i < times; i++ {
		oneIDs := ids[(i-1) * max:i * max -1]
		_, err = conn.Model(m).Where({{$.ModelPkgName}}.{{$.TypeName}}_ID+" in (?)", pg.In(oneIDs)).Delete()
		if err != nil {
			return
		}
	}

	if max * times < len(ids) {
		oneIDs := ids[max * times:]
		_, err = conn.Model(m).Where({{$.ModelPkgName}}.{{$.TypeName}}_ID+" in (?)", pg.In(oneIDs)).Delete()
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

//example,please edit it
//update designated column with Optimistic Lock
func (c *{{$.DaoName}}) Update{{$.TypeName}}SomeColumn(conn *pg.Conn, ids []string,/*todo: update parameters*/) (err error) {

	ms, err := c.GetLockByID(conn, ids...)
	if err != nil {
		return
	}

	for i, _ := range ms {
		ms[i].UpdateTime = time.Now().Unix()
		ms[i].OptimisticLockVersion++
		//todo ms[i].xx=parameter
		_, err = conn.Model(ms[i]).Where({{$.ModelPkgName}}.{{$.TypeName}}_ID+" = ? and "+{{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion+" = ?", ms[i].ID, ms[i].OptimisticLockVersion-1).Column(/*{{$.ModelPkgName}}.{{$.TypeName}}_xx,*/ {{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion, {{$.ModelPkgName}}.{{$.TypeName}}_UpdateTime).Update()
		if err != nil {
			dot.Logger().Debugln(err.Error())
			return
		}
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
