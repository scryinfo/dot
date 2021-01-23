package model

import "github.com/scryinfo/dot/dots/db/pgs"

type Data struct {
	Name  string
	Count int64
}

//go:generate gmodel -typeName=Notice -tableName=notices
//go:generate gdao -typeName=Notice -tableName=notices -daoPackage=dao
type Notice struct {
	pgs.ModelBase
	Data   Data `pg:"composite:data"`
	Status int
}
