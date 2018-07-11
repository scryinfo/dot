package dot

import "errors"

var (
	_ Errorer = (*sError)(nil)
)

//Errorer dot error interface
type Errorer interface {
	error
	Code() string
}

type sError struct {
	err error
	id  string
}

//Code error id
func (c *sError) Code() string {
	return c.id
}

//SError return error info
func (c *sError) Error() string {
	return c.err.Error()
}

//NewError new Errorer
func NewError(id string, info string) Errorer {
	err := sError{err: errors.New(info), id: id}
	return &err
}

//Error dot error
type Error struct {
	ErrNullParameter    Errorer
	ErrExisted          Errorer
	ErrNotExisted       Errorer
	ErrParameter        Errorer
	ErrRelyTypeNotMatch Errorer
}

//SError dot的全局常用 error对象
var SError = &Error{}

func init() {
	SError.ErrNullParameter = NewError("dot_null_parameter", "the parameter is null")
	SError.ErrExisted = NewError("dot_existed", "the value exited:")
	SError.ErrNotExisted = NewError("dot_not_existed", "the not exited:")
	SError.ErrParameter = NewError("dot_error_parameter", "the parameter error ")
	SError.ErrRelyTypeNotMatch = NewError("dot_rely_type","rely type not match ")
}
