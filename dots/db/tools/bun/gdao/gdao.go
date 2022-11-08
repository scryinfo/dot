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

	"github.com/scryinfo/dot/dots/db/bun"
	"github.com/scryinfo/scryg/sutils/uuid"
)

// DbField do not use the map, we need the order
type DbField struct {
	Name   string
	DbName string
}

type tData struct {
	DaoName            string
	TypeName           string
	TableName          string
	OrmMode            string
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
	ormMode    string // bun|gorm|...
	useLock    bool
	useGorm    bool
}

func parms(data *tData) {
	flag.StringVar(&params.typeName, "typeName", "", "")
	flag.StringVar(&params.tableName, "tableName", "", "")
	flag.StringVar(&params.daoPackage, "daoPackage", "", "")
	flag.StringVar(&params.suffix, "suffix", "Dao", "")
	flag.StringVar(&params.ormMode, "ormMode", "bun", "")
	flag.StringVar(&params.model, "model", "models.go", "")
	flag.BoolVar(&params.useGorm, "useGorm", false, "")
	flag.Parse()

	if len(params.tableName) < 1 {
		params.tableName = bun.Underscore(params.typeName)
	}

	data.TypeName = params.typeName
	data.TableName = params.tableName
	data.DaoPkgName = params.daoPackage
	data.OrmMode = params.ormMode
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

	if len(params.tableName) < 1 {
		log.Fatal("table name is null")
	}

	os.Setenv("GOPACKAGE", "model")
	os.Setenv("GOFILE", data.ModelFile)

	var src []byte = nil
	{
		makeData(data)
		src = gmodel(data, params.useGorm)
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

func gmodel(data *tData, supportGorm bool) []byte {

	temp := ""
	if supportGorm {
		temp = `
package dao

import (
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/db/gorms"
	"github.com/scryinfo/dot/dots/db/pgs"
	"{{$.ImportModelPkgName}}"
	"gorm.io/gorm/clause"
)

const {{$.DaoName}}TypeID = "{{$.ID}}"

type {{$.DaoName}} struct {
	*gorms.DaoBase {{$.BackQuote}}dot:""{{$.BackQuote}}
}

//{{$.DaoName}}TypeLives
func {{$.DaoName}}TypeLives() []*dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{TypeID: {{$.DaoName}}TypeID, NewDoter: func(conf []byte) (dot.Dot, error) {
			return &{{$.DaoName}}{}, nil
		}},
		Lives: []dot.Live{
			{
				LiveID: {{$.DaoName}}TypeID,
				RelyLives: map[string]dot.LiveID{
					"DaoBase": pgs.DaoBaseTypeID,
				},
			},
		},
	}

	lives := gorms.DaoBaseTypeLives()
	lives = append(lives, tl)

	return lives
}

func (c *{{$.DaoName}}) GetByID(id uint) (m *{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	result := c.Wrapper.GetDb().First(&m, id)
	if result.RowsAffected == 1 {
		return m, nil
	}
	return nil, result.Error
}

func (c *{{$.DaoName}}) Query(condition string, params ...interface{}) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {

	db := c.Wrapper.GetDb()
	if len(condition) < 1 {
		db = db.Find(&ms)
	} else {
		db = db.Where(condition, params...).Find(&ms)
	}
	err = db.Error
	if err != nil {
		ms = nil
	}
	return
}

func (c *{{$.DaoName}}) List() (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	db := c.Wrapper.GetDb().Find(&ms)
	if db.Error != nil {
		ms = nil
	}
	return ms, db.Error
}

func (c *{{$.DaoName}}) Count(condition string, params ...interface{}) (count int64, err error) {

	var ms []*{{$.ModelPkgName}}.{{$.TypeName}}
	db := c.Wrapper.GetDb()
	if len(condition) < 1 {
		db = db.Find(&ms)
	} else {
		db = db.Where(condition, params...).Find(&ms)
	}
	err = db.Error
	if err != nil {
		count = 0
	}
	count=db.RowsAffected
	return
}

func (c *{{$.DaoName}}) QueryPage(pageSize int, page int, condition string, params ...interface{}) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	db := c.Wrapper.GetDb()
	if len(condition) < 1 {
		db = db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&ms)
	} else {
		db = db.Limit(pageSize).Offset((page-1)*pageSize).Where(condition, params...).Find(&ms)
	}
	if db.Error != nil { //be sure
		ms = nil
	}
	return ms, db.Error
}

func (c *{{$.DaoName}}) QueryPageWithCount(pageSize int, page int, condition string, params ...interface{}) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, count int64, err error) {
	db := c.Wrapper.GetDb()
	if len(condition) < 1 {
		db = db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&ms)
	} else {
		db = db.Limit(pageSize).Offset((page-1)*pageSize).Where(condition, params...).Find(&ms)
	}
	if db.Error != nil { //be sure
		ms = nil
	}
	return ms, db.RowsAffected, db.Error
}

func (c *{{$.DaoName}}) QueryOne(condition string, params ...interface{}) (m *{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	db := c.Wrapper.GetDb()
	if len(condition) < 1 {
		db = db.First(&m)
	} else {
		db = db.Where(condition, params...).First(&m)
	}
	err = db.Error
	if err != nil {
		m = nil
	}
	return
}

//insert = insertReturn
func (c *{{$.DaoName}}) Insert(m *{{$.ModelPkgName}}.{{$.TypeName}}) (err error) {
	//because id auto increment
	//m.ID = 0
	err = c.Wrapper.GetDb().Create(&m).Error
	return
}

func (c *{{$.DaoName}}) Inserts(ms []*{{$.ModelPkgName}}.{{$.TypeName}}) (err error) {
	//because id auto increment
	//m.ID = 0
	err = c.Wrapper.GetDb().Create(&ms).Error
	return
}

//update everything except ID
func (c *{{$.DaoName}}) Upsert(m *{{$.ModelPkgName}}.{{$.TypeName}}) (err error) {
	//because id auto increment
	//m.ID = 0
	err = c.Wrapper.GetDb().Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&m).Error
	return
}

func (c *{{$.DaoName}}) Upserts(ms []*{{$.ModelPkgName}}.{{$.TypeName}}) (err error) {
	//because id auto increment
	//m.ID = 0
	err = c.Wrapper.GetDb().Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&ms).Error
	return
}

//Updates 方法支持 struct 和 map[string]interface{} 参数
//当使用 struct 更新时，默认情况下，GORM 只会更新非零值的字段
func (c *{{$.DaoName}}) Update(m *{{$.ModelPkgName}}.{{$.TypeName}}) (err error) {
	err = c.Wrapper.GetDb().Updates(&m).Error
	return
}

//Save 会保存所有的字段，即使字段是零值
//todo CreateTime=0
func (c *{{$.DaoName}}) Save(m *{{$.ModelPkgName}}.{{$.TypeName}}) (err error) {
	err = c.Wrapper.GetDb().Save(&m).Error
	return
}

//default where id=m.ID
func (c *{{$.DaoName}}) UpdateColumn(m *{{$.ModelPkgName}}.{{$.TypeName}}, columnName string, value interface{}) (err error) {
	err = c.Wrapper.GetDb().Model(&m).Update(columnName, value).Error
	return
}
//soft delete
func (c *{{$.DaoName}}) DeleteById(id uint) error {
	return c.Wrapper.GetDb().Delete(&{{$.ModelPkgName}}.{{$.TypeName}}{}, id).Error
}
//soft delete
func (c *{{$.DaoName}}) DeleteByIds(ids []uint) error {
	return c.Wrapper.GetDb().Delete(&{{$.ModelPkgName}}.{{$.TypeName}}{}, ids).Error
}
//soft delete
func (c *{{$.DaoName}}) Delete(condition string, params ...interface{}) error {
	return c.Wrapper.GetDb().Where(condition, params...).Delete(&{{$.ModelPkgName}}.{{$.TypeName}}{}).Error
}

//Delete permanently
func (c *{{$.DaoName}}) DeleteByIdUnscoped(id uint) error {
	return c.Wrapper.GetDb().Unscoped().Delete(&{{$.ModelPkgName}}.{{$.TypeName}}{}, id).Error
}

`
	} else {
		temp = `
package {{$.DaoPkgName}}
import (
	"time"
	
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
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
func (c *{{$.DaoName}}) GetByIDWithLock(conn orm.DB, id string) (m *{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	m = &{{$.ModelPkgName}}.{{$.TypeName}}{}
	m.ID = id
	err = conn.Model(m).WherePK().For("UPDATE").Select()
	if err != nil {
		m = nil
	}
	return
}
func (c *{{$.DaoName}}) GetByID(conn orm.DB, id string) (m *{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	m = &{{$.ModelPkgName}}.{{$.TypeName}}{}
	m.ID = id
	err = conn.Model(m).WherePK().Select()
	if err != nil {
		m = nil
	}
	return
}

//update before
//you must get OptimisticLockVersion value
func (c *{{$.DaoName}}) GetLockByID(conn orm.DB, ids ...string) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	for i,_ := range ids {
		m := &{{$.ModelPkgName}}.{{$.TypeName}}{}
		m.ID = ids[i]
		ms=append(ms,m)
	}
	err=conn.Model(&ms).WherePK().Column({{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion,{{$.ModelPkgName}}.{{$.TypeName}}_Struct+"."+{{$.ModelPkgName}}.{{$.TypeName}}_ID).For("UPDATE").Select()
	if err != nil {
		ms=nil
	}
	return
}
func (c *{{$.DaoName}}) GetLockByModelID(conn orm.DB, ms ...*{{$.ModelPkgName}}.{{$.TypeName}}) error {
	return conn.Model(&ms).WherePK().Column({{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion,{{$.ModelPkgName}}.{{$.TypeName}}_Struct+"."+{{$.ModelPkgName}}.{{$.TypeName}}_ID).For("UPDATE").Select()
}

// if find nothing, return pg.ErrNoRows
func (c *{{$.DaoName}}) QueryWithLock(conn orm.DB, condition string, params ...interface{}) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
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
func (c *{{$.DaoName}}) Query(conn orm.DB, condition string, params ...interface{}) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
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
func (c *{{$.DaoName}}) ListWithLock(conn orm.DB) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	err = conn.Model(&ms).For("UPDATE").Select()
	if err != nil {//be sure
		ms = nil
	}
	return
}
func (c *{{$.DaoName}}) List(conn orm.DB) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	err = conn.Model(&ms).Select()
	if err != nil {//be sure
		ms = nil
	}
	return
}

func (c *{{$.DaoName}}) Count(conn orm.DB, condition string, params ...interface{}) (count int, err error) {
	if len(condition) < 1 {
		count, err = conn.Model(&{{$.ModelPkgName}}.{{$.TypeName}}{}).Count()
	}else {
		count, err = conn.Model(&{{$.ModelPkgName}}.{{$.TypeName}}{}).Where(condition, params...).Count()
	}
	return
}

// if find nothing, return pg.ErrNoRows
func (c *{{$.DaoName}}) QueryPageWithLock(conn orm.DB, pageSize int, page int, condition string, params ...interface{}) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
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
func (c *{{$.DaoName}}) QueryPage(conn orm.DB, pageSize int, page int, condition string, params ...interface{}) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
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

// count counts valid records which after conditions filter, rather than whole table's count
func (c *{{$.DaoName}}) QueryPageWithCount(
	conn orm.DB,
	pageSize,
	pageNum int,
	condition string,
	params ...interface{},
) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, count int, err error) {
	if len(condition) < 1 {
		count, err = conn.Model(&ms).Limit(pageSize).Offset((pageNum - 1) * pageSize).SelectAndCount()
	} else {
		count, err = conn.Model(&ms).Where(condition, params...).Limit(pageSize).Offset((pageNum - 1) * pageSize).SelectAndCount()
	}

	if err != nil { //be sure
		ms = nil
	}
	return
}

// if find nothing, return pg.ErrNoRows
func (c *{{$.DaoName}}) QueryOneWithLock(conn orm.DB, condition string, params ...interface{}) (m *{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
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
func (c *{{$.DaoName}}) QueryOne(conn orm.DB, condition string, params ...interface{}) (m *{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
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
func (c *{{$.DaoName}}) Insert(conn orm.DB, m *{{$.ModelPkgName}}.{{$.TypeName}}) (err error) {
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuid()
	}
	m.CreateTime = time.Now().Unix()
	m.UpdateTime = m.CreateTime
	_, err = conn.Model(m).Insert()
	return
}
//if insert nothing, then return pg.ErrNoRows
func (c *{{$.DaoName}}) InsertReturn(conn orm.DB, m *{{$.ModelPkgName}}.{{$.TypeName}}) ( mnew *{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
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
func (c *{{$.DaoName}}) Upsert(conn orm.DB, m *{{$.ModelPkgName}}.{{$.TypeName}}) (err error) {
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
func (c *{{$.DaoName}}) UpsertReturn(conn orm.DB, m *{{$.ModelPkgName}}.{{$.TypeName}}) ( mnew *{{$.ModelPkgName}}.{{$.TypeName}},err error) {
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
func (c *{{$.DaoName}}) Update(conn orm.DB, m *{{$.ModelPkgName}}.{{$.TypeName}}) (err error) {
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
func (c *{{$.DaoName}}) UpdateReturn(conn orm.DB, m *{{$.ModelPkgName}}.{{$.TypeName}}) (mnew *{{$.ModelPkgName}}.{{$.TypeName}},  err error) {
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
func (c *{{$.DaoName}}) Delete(conn orm.DB, m *{{$.ModelPkgName}}.{{$.TypeName}}) (err error) {
	_, err = conn.Model(m).Delete()
	return
}

//if delete nothing, then return pg.ErrNoRows
func (c *{{$.DaoName}}) DeleteByID(conn orm.DB, id string) (err error) {
	_, err = conn.Model((*{{$.ModelPkgName}}.{{$.TypeName}})(nil)).Where({{$.ModelPkgName}}.{{$.TypeName}}_ID+" = ?", id).Delete()
	return
}

//if delete nothing, then return pg.ErrNoRows
func (c *{{$.DaoName}}) DeleteByIDs(conn orm.DB, ids []string, oneMax int) (err error) {
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
func (c *{{$.DaoName}}) DeleteReturn(conn orm.DB, m *{{$.ModelPkgName}}.{{$.TypeName}}) (mnew *{{$.ModelPkgName}}.{{$.TypeName}},err error) {
	mnew = &{{$.ModelPkgName}}.{{$.TypeName}}{}
	_, err = conn.Model(m).WherePK().Returning("*").Delete(mnew)
	if err != nil {
		mnew = nil
	}
	return
}

//example,please edit it
//update designated column with Optimistic Lock
func (c *{{$.DaoName}}) Update{{$.TypeName}}SomeColumn(conn orm.DB, ids []string,/*todo: update parameters*/) (err error) {

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
	}

	if data.OrmMode == "bun" {
		temp = getBunDaoTpl()
	}

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

// returns orm uptrace/bun dao template.
func getBunDaoTpl() (temp string) {
	temp = `
package {{$.DaoPkgName}}
import (
	"context"
	"database/sql"
	"time"
	
    "github.com/scryinfo/dot/dot"
    "github.com/scryinfo/dot/dots/db/bun/pgd"
	"github.com/scryinfo/scryg/sutils/uuid"
	"github.com/uptrace/bun"
	"{{$.ImportModelPkgName}}"
)

const {{$.DaoName}}TypeID = "{{$.ID}}"

type {{$.DaoName}} struct {
	*pgd.DaoBase {{$.BackQuote}}dot:""{{$.BackQuote}}
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

// if find|update|insert nothing, sql.ErrNoRows error may returned

func (c *{{$.DaoName}}) GetByIDWithLock(conn *bun.DB, id string) (m *{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	m = &{{$.ModelPkgName}}.{{$.TypeName}}{}
	m.ID = id
	err = conn.NewSelect().Model(m).WherePK().For("UPDATE").Scan(context.TODO())
	if err != nil {
		m = nil
	}
	return
}
func (c *{{$.DaoName}}) GetByID(conn *bun.DB, id string) (m *{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	m = &{{$.ModelPkgName}}.{{$.TypeName}}{}
	m.ID = id
	err = conn.NewSelect().Model(m).WherePK().Scan(context.TODO())
	if err != nil {
		m = nil
	}
	return
}

//update before
//you must get OptimisticLockVersion value
func (c *{{$.DaoName}}) GetLockByID(conn *bun.DB, ids ...string) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	for i,_ := range ids {
		m := &{{$.ModelPkgName}}.{{$.TypeName}}{}
		m.ID = ids[i]
		ms = append(ms,m)
	}
	err = conn.NewSelect().Model(&ms).WherePK().Column({{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion,{{$.ModelPkgName}}.{{$.TypeName}}_Struct+"."+{{$.ModelPkgName}}.{{$.TypeName}}_ID).For("UPDATE").Scan(context.TODO())
	if err != nil {
		ms = nil
	}
	return
}
func (c *{{$.DaoName}}) GetLockByModelID(conn *bun.DB, ms ...*{{$.ModelPkgName}}.{{$.TypeName}}) error {
	return conn.NewSelect().Model(&ms).WherePK().Column({{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion,{{$.ModelPkgName}}.{{$.TypeName}}_Struct+"."+{{$.ModelPkgName}}.{{$.TypeName}}_ID).For("UPDATE").Scan(context.TODO())
}

func (c *{{$.DaoName}}) QueryWithLock(conn *bun.DB, condition string, params ...interface{}) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	if len(condition) < 1 {
		err = conn.NewSelect().Model(&ms).For("UPDATE").Scan(context.TODO())
	} else {
		err = conn.NewSelect().Model(&ms).Where(condition, params...).For("UPDATE").Scan(context.TODO())
	}
	if err != nil { //be sure
		ms = nil
	}
	return
}
func (c *{{$.DaoName}}) Query(conn *bun.DB, condition string, params ...interface{}) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	if len(condition) < 1 {
		err = conn.NewSelect().Model(&ms).Scan(context.TODO())
	} else {
		err = conn.NewSelect().Model(&ms).Where(condition, params...).Scan(context.TODO())
	}
	if err != nil { //be sure
		ms = nil
	}
	return
}

func (c *{{$.DaoName}}) ListWithLock(conn *bun.DB) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	err = conn.NewSelect().Model(&ms).For("UPDATE").Scan(context.TODO())
	if err != nil {//be sure
		ms = nil
	}
	return
}
func (c *{{$.DaoName}}) List(conn *bun.DB) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	err = conn.NewSelect().Model(&ms).Scan(context.TODO())
	if err != nil {//be sure
		ms = nil
	}
	return
}

func (c *{{$.DaoName}}) Count(conn *bun.DB, condition string, params ...interface{}) (count int, err error) {
	if len(condition) < 1 {
		count, err = conn.NewSelect().Model(&{{$.ModelPkgName}}.{{$.TypeName}}{}).Count(context.TODO())
	} else {
		count, err = conn.NewSelect().Model(&{{$.ModelPkgName}}.{{$.TypeName}}{}).Where(condition, params...).Count(context.TODO())
	}
	return
}

func (c *{{$.DaoName}}) QueryPageWithLock(conn *bun.DB, pageSize int, page int, condition string, params ...interface{}) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	if len(condition) < 1 {
		err = conn.NewSelect().Model(&ms).Limit(pageSize).Offset((page - 1) * pageSize).For("UPDATE").Scan(context.TODO())
	}else {
		err = conn.NewSelect().Model(&ms).Where(condition, params...).Limit(pageSize).Offset((page - 1) * pageSize).For("UPDATE").Scan(context.TODO())
	}
	if err != nil { //be sure
		ms = nil
	}
	return
}
func (c *{{$.DaoName}}) QueryPage(conn *bun.DB, pageSize int, page int, condition string, params ...interface{}) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	if len(condition) < 1 {
		err = conn.NewSelect().Model(&ms).Limit(pageSize).Offset((page - 1) * pageSize).Scan(context.TODO())
	}else {
		err = conn.NewSelect().Model(&ms).Where(condition, params...).Limit(pageSize).Offset((page - 1) * pageSize).Scan(context.TODO())
	}
	if err != nil { //be sure
		ms = nil
	}
	return
}

// count counts valid records which after conditions filter, rather than whole table's count
func (c *{{$.DaoName}}) QueryPageWithCount(
	conn *bun.DB,
	pageSize,
	pageNum int,
	condition string,
	params ...interface{},
) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, count int, err error) {
	if len(condition) < 1 {
		count, err = conn.NewSelect().Model(&ms).Limit(pageSize).Offset((pageNum - 1) * pageSize).ScanAndCount(context.TODO())
	} else {
		count, err = conn.NewSelect().Model(&ms).Where(condition, params...).Limit(pageSize).Offset((pageNum - 1) * pageSize).ScanAndCount(context.TODO())
	}

	if err != nil { //be sure
		ms = nil
	}
	return
}

func (c *{{$.DaoName}}) QueryOneWithLock(conn *bun.DB, condition string, params ...interface{}) (m *{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	m = &{{$.ModelPkgName}}.{{$.TypeName}}{}
	if len(condition) < 1 {
		err = conn.NewSelect().Model(m).For("UPDATE").Limit(1).Scan(context.TODO())
	} else {
		err = conn.NewSelect().Model(m).Where(condition, params...).For("UPDATE").Limit(1).Scan(context.TODO())
	}
	if err != nil {//be sure
		m = nil
	}
	return
}

func (c *{{$.DaoName}}) QueryOne(conn *bun.DB, condition string, params ...interface{}) (m *{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	m = &{{$.ModelPkgName}}.{{$.TypeName}}{}
	if len(condition) < 1 {
		err = conn.NewSelect().Model(m).Scan(context.TODO())
	} else {
		err = conn.NewSelect().Model(m).Where(condition, params...).Scan(context.TODO())
	}
	if err != nil {//be sure
		m = nil
	}
	return
}

func (c *{{$.DaoName}}) Insert(conn *bun.DB, m *{{$.ModelPkgName}}.{{$.TypeName}}) (err error) {
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuid()
	}
	m.CreateTime = time.Now().Unix()
	m.UpdateTime = m.CreateTime
	_, err = conn.NewInsert().Model(m).Exec(context.TODO())
	return
}

func (c *{{$.DaoName}}) InsertReturn(conn *bun.DB, m *{{$.ModelPkgName}}.{{$.TypeName}}) ( mnew *{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuid()
	}
	m.CreateTime = time.Now().Unix()
	m.UpdateTime = m.CreateTime

	mnew = &{{$.ModelPkgName}}.{{$.TypeName}}{}
	_, err = conn.NewInsert().Model(m).Returning("*").Exec(context.TODO(), mnew)
	if err != nil{
		mnew = nil
	}
	return
}

func (c *{{$.DaoName}}) Upsert(conn *bun.DB, m *{{$.ModelPkgName}}.{{$.TypeName}}) (err error) {
    m.UpdateTime = time.Now().Unix()
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuid()
		m.CreateTime = m.UpdateTime
	} else if m.CreateTime == 0 {
		m.CreateTime = m.UpdateTime
	}
	m.OptimisticLockVersion++
	om := conn.NewInsert().Model(m).On("CONFLICT (id) DO UPDATE").Where({{$.ModelPkgName}}.{{$.TypeName}}_Struct+"."+{{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion+" = ?",m.OptimisticLockVersion-1)
	for _, it := range m.ToUpsertSet() {
		om.Set(it)
	}
	res, err := om.Exec(context.TODO())
	if n, _ := res.RowsAffected(); n == 0 {
		newm, err := c.GetLockByID(conn, m.ID)
		if err != nil {
			return err
		}
		m.OptimisticLockVersion = newm[0].OptimisticLockVersion
		err = c.Update(conn, m)
	}
	return err
}

func (c *{{$.DaoName}}) UpsertReturn(conn *bun.DB, m *{{$.ModelPkgName}}.{{$.TypeName}}) ( mnew *{{$.ModelPkgName}}.{{$.TypeName}},err error) {
	m.UpdateTime = time.Now().Unix()
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuid()
		m.CreateTime = m.UpdateTime
	} else if m.CreateTime == 0 {
		m.CreateTime = m.UpdateTime
	}
	m.OptimisticLockVersion++
	om := conn.NewInsert().Model(m).On("CONFLICT (id) DO UPDATE").Where({{$.ModelPkgName}}.{{$.TypeName}}_Struct+"."+{{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion+" = ?",m.OptimisticLockVersion-1)
	for _, it := range m.ToUpsertSet() {
		om.Set(it)
	}
	mnew = &{{$.ModelPkgName}}.{{$.TypeName}}{}
	res, err := om.Returning("*").Exec(context.TODO(), mnew)
	if n, _ := res.RowsAffected(); n == 0 {
		ms, err := c.GetLockByID(conn, m.ID)
		if err != nil {
			return nil, err
		}
		m.OptimisticLockVersion = ms[0].OptimisticLockVersion
		mnew, err = c.UpdateReturn(conn, m)
	}
	return
}

func (c *{{$.DaoName}}) Update(conn *bun.DB, m *{{$.ModelPkgName}}.{{$.TypeName}}) (err error) {
	m.UpdateTime = time.Now().Unix()
	m.OptimisticLockVersion++
	res, err := conn.NewUpdate().Model(m).Where({{$.ModelPkgName}}.{{$.TypeName}}_ID+" = ? and "+{{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion+" = ?", m.ID, m.OptimisticLockVersion-1).Exec(context.TODO())
	if n, _ := res.RowsAffected(); n == 0 {
		err = sql.ErrNoRows
	}
	return
}

func (c *{{$.DaoName}}) UpdateReturn(conn *bun.DB, m *{{$.ModelPkgName}}.{{$.TypeName}}) (mnew *{{$.ModelPkgName}}.{{$.TypeName}},  err error) {
	m.UpdateTime = time.Now().Unix()
	m.OptimisticLockVersion++
	mnew = &{{$.ModelPkgName}}.{{$.TypeName}}{}
	res, err := conn.NewUpdate().Model(m).Where({{$.ModelPkgName}}.{{$.TypeName}}_ID+" = ? and "+{{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion+" = ?", m.ID, m.OptimisticLockVersion-1).Returning("*").Exec(context.TODO(), mnew)
	if err != nil {
		mnew = nil
	}
	if n, _ := res.RowsAffected(); n == 0 {
		err = sql.ErrNoRows
	}
	return
}

func (c *{{$.DaoName}}) Delete(conn *bun.DB, m *{{$.ModelPkgName}}.{{$.TypeName}}) (err error) {
	return c.DeleteByID(conn, m.ID)
}

func (c *{{$.DaoName}}) DeleteByID(conn *bun.DB, id string) (err error) {
	_, err = conn.NewDelete().Model((*{{$.ModelPkgName}}.{{$.TypeName}})(nil)).Where({{$.ModelPkgName}}.{{$.TypeName}}_ID+" = ?", id).Exec(context.TODO())
	return
}

func (c *{{$.DaoName}}) DeleteByIDs(conn *bun.DB, ids []string, oneMax int) (err error) {
	m := (*{{$.ModelPkgName}}.{{$.TypeName}})(nil)
	max := oneMax
	times := len(ids)/max;
	for i := 1; i < times; i++ {
		oneIDs := ids[(i-1) * max:i * max -1]
		_, err = conn.NewDelete().Model(m).Where({{$.ModelPkgName}}.{{$.TypeName}}_ID+" in (?)", bun.In(oneIDs)).Exec(context.TODO())
		if err != nil {
			return
		}
	}

	if max * times < len(ids) {
		oneIDs := ids[max * times:]
		_, err = conn.NewDelete().Model(m).Where({{$.ModelPkgName}}.{{$.TypeName}}_ID+" in (?)", bun.In(oneIDs)).Exec(context.TODO())
	}
	return 
}

func (c *{{$.DaoName}}) DeleteReturn(conn *bun.DB, m *{{$.ModelPkgName}}.{{$.TypeName}}) (mnew *{{$.ModelPkgName}}.{{$.TypeName}},err error) {
	mnew = &{{$.ModelPkgName}}.{{$.TypeName}}{}
	_, err = conn.NewDelete().Model(m).WherePK().Returning("*").Exec(context.TODO(), mnew)
	if err != nil {
		mnew = nil
	}
	return
}

//example,please edit it
//update designated column with Optimistic Lock
func (c *{{$.DaoName}}) Update{{$.TypeName}}SomeColumn(conn *bun.DB, ids []string,/*todo: update parameters*/) (err error) {

	ms, err := c.GetLockByID(conn, ids...)
	if err != nil {
		return
	}
	ctx := context.TODO()
	condition := {{$.ModelPkgName}}.{{$.TypeName}}_ID+" = ? and "+{{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion+" = ?"
	for i, _ := range ms {
		ms[i].UpdateTime = time.Now().Unix()
		ms[i].OptimisticLockVersion++
		_, err = conn.NewUpdate().Model(ms[i]).Where(condition, ms[i].ID, ms[i].OptimisticLockVersion-1).Column(/*{{$.ModelPkgName}}.{{$.TypeName}}_xx,*/ {{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion, {{$.ModelPkgName}}.{{$.TypeName}}_UpdateTime).Exec(ctx)
		if err != nil {
			dot.Logger().Debugln(err.Error())
			return
		}
	}
	return
}
`
	return
}
