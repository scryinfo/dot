package model

//go:generate gmodel -typeName=Notice -tableName=notice
//go:generate gdao -typeName=Notice -tableName=notice -daoPackage=dao
type Notice struct {
	Id         string `pk`
	Status     int
	CreateTime int64
	UpdateTime int64
	Version    int64
}
