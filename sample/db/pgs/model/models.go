package model

//go:generate gmodel -typeName=Notice -tableName=notice
//go:generate gdao -typeName=Notice -tableName=notice -daoPackage=dao
type Notice struct {
	ID     string `pk`
	Status int

	CreateTime            int64
	UpdateTime            int64
	OptimisticLockVersion int64 `pg:",use_zero"` //default 0 ，not null
}
