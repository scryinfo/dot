package model

import (
	"github.com/scryinfo/dot/dots/db"
)

type DataType struct {
	Name  string
	Count int64
}

//go:generate gmodel -typeName=Notice -tableName=notices
//go:generate gdao -typeName=Notice -tableName=notices -daoPackage=dao
type Notice struct {
	db.ModelBase
	Data   DataType `pg:"composite:data"`
	Status int
	No     int `pg:"-"`
}

//go:generate gmodel -typeName=Sub -tableName=subs
//go:generate gdao -typeName=Sub -tableName=subs -daoPackage=dao
type Sub struct {
	db.ModelBase
	Name string
}

//go:generate gmodel -typeName=HasSub -tableName=has_subs
//go:generate gdao -typeName=HasSub -tableName=has_subs -daoPackage=dao
type HasSub struct {
	db.ModelBase
	SubId   string
	Count   int
	SubData *Sub `pg:"rel:has-one"`
}
