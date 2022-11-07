package model

import (
	"github.com/scryinfo/dot/dots/db"
	"github.com/uptrace/bun"
)

type ModelBase struct {
	ID                    string `pg:",pk" json:"id" bun:",pk"`                //id
	UpdateTime            int64  `json:"updateTime"`                           //更新时间
	CreateTime            int64  `json:"createTime"`                           //创建时间
	OptimisticLockVersion int64  `pg:",use_zero" json:"optimisticLockVersion"` //default 0 ，not null
}

type DataType struct {
	Name  string
	Count int64
}

//go:generate gmodel -typeName=Notice -tableName=notices
//go:generate gdao -typeName=Notice -tableName=notices -daoPackage=dao
type Notice struct {
	bun.BaseModel
	ModelBase
	Data   DataType `pg:"composite:data"`
	Status int
	No     int `pg:"-" bun:"-"`
}

//go:generate gmodel -typeName=Sub -tableName=subs
//go:generate gdao -typeName=Sub -tableName=subs -daoPackage=dao
type Sub struct {
	bun.BaseModel
	ModelBase
	Name string
}

//go:generate gmodel -typeName=HasSub -tableName=has_subs
//go:generate gdao -typeName=HasSub -tableName=has_subs -daoPackage=dao
type HasSub struct {
	bun.BaseModel
	ModelBase
	SubId   string
	Count   int
	SubData *Sub `pg:"rel:has-one"`
}

//go:generate gmodel -typeName=AutoData -tableName=auto_datas
//go:generate gdao -typeName=AutoData -tableName=auto_datas -daoPackage=dao -useGorm=true
type AutoData struct {
	db.AutoModelBase
	Name string
	Age  int8
}
