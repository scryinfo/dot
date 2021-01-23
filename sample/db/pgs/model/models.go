package model

import "github.com/scryinfo/dot/dots/db/pgs"

type DataType struct {
	Name  string
	Count int64
}

//go:generate gmodel -typeName=Notice -tableName=notices
//go:generate gdao -typeName=Notice -tableName=notices -daoPackage=dao
type Notice struct {
	pgs.ModelBase
	Data   DataType `pg:"composite:data"`
	Status int
}

//go:generate gmodel -typeName=Sub -tableName=subs
//go:generate gdao -typeName=Sub -tableName=subs -daoPackage=dao
type Sub struct {
	pgs.ModelBase
	Name string
}

//go:generate gmodel -typeName=HasSub -tableName=has_subs
//go:generate gdao -typeName=HasSub -tableName=has_subs -daoPackage=dao
type HasSub struct {
	pgs.ModelBase
	SubId   string
	Count   int
	SubData *Sub `pg:"rel:has-one"`
}
