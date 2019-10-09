package model

//go:generate gmodel -typeName=Notice -tableName=notice
type Notice struct {
	Id         string `pk`
	Status     int
	CreateTime int64
	UpdateTime int64
}
